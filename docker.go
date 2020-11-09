package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

func BuildImage(tag string, cli client.Client) error {
	contextPath, _ := filepath.Abs("./files/")
	buildOpts := types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		Tags:       []string{tag},
	}

	buildCtx, err := archive.TarWithOptions(contextPath, &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := cli.ImageBuild(context.Background(), buildCtx, buildOpts)
	if err != nil {
		log.Fatalf("Failed to build image - %v", err)
		return err
	}
	defer resp.Body.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stderr, termFd, isTerm, nil)

	return nil
}

func PullImageIfNotExists(cli *client.Client, tag string) {
	list, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatalf("Failed to list images - %v", err)
	}
	found := false
	for _, i := range list {
		for _, t := range i.RepoTags {
			if t == tag {
				found = true
			}
		}
		if found {
			break
		}
	}

	if !found {
		pullResp, err := cli.ImagePull(context.Background(), tag, types.ImagePullOptions{})
		if err != nil {
			log.Fatalf("Failed to pull image %s - %v", tag, err)
		}
		defer pullResp.Close()

		termFd, isTerm := term.GetFdInfo(os.Stderr)
		jsonmessage.DisplayJSONMessagesStream(pullResp, os.Stderr, termFd, isTerm, nil)
	}
}

func EvaluateScript(script string) string {
	tag := "lukaspj/t3deval:4_0Preview"

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	PullImageIfNotExists(cli, tag)

	var containerResp container.ContainerCreateCreatedBody
	containerResp, err = cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        tag,
			AttachStderr: true,
			AttachStdout: true,
			Cmd:          []string{script},
			Tty:          true,
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		nil,
		"t3deval-worker-1",
	)
	if err != nil {
		log.Fatalf("Failed to create container %s - %v", containerResp.ID, err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = cli.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatalf("Failed to run container %s - %v", containerResp.ID, err)
	}

	waitCh, errCh := cli.ContainerWait(context.Background(), containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		log.Fatalf("Failed to wait on container %s - %v", containerResp.ID, err)
	case <-waitCh:

	}

	defer func() {
		err = cli.ContainerRemove(context.Background(), containerResp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Fatalf("Failed to remove container %s - %v", containerResp.ID, err)
		}
	}()

	reader, err := cli.ContainerLogs(context.Background(), containerResp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Fatalf("Failed to get container logs from %s - %v", containerResp.ID, err)
	}

	logs, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to read logs %s - %v", containerResp.ID, err)
	}

	return string(logs)
}
