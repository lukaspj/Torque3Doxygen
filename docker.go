package main

import (
	"context"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "Failed to build image")
	}
	defer resp.Body.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stderr, termFd, isTerm, nil)
}

func PullImageIfNotExists(cli *client.Client, tag string) error {
	list, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Failed to list images")
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
			return errors.Wrapf(err, "Failed to pull image %s", tag)
		}
		defer pullResp.Close()

		termFd, isTerm := term.GetFdInfo(os.Stderr)
		return jsonmessage.DisplayJSONMessagesStream(pullResp, os.Stderr, termFd, isTerm, nil)
	}
	return nil
}

func EvaluateScript(script string) (string, error) {
	tag := "lukaspj/t3deval:4_0Preview"

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	err = PullImageIfNotExists(cli, tag)
	if err != nil {
		return "", err
	}

	var containerResp container.ContainerCreateCreatedBody
	containerResp, err = cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        tag,
			Cmd:          []string{script},
			Tty:          true,
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		nil,
		"t3deval-worker-1",
	)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create container %s", containerResp.ID)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = cli.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to run container %s", containerResp.ID)
	}

	waitCh, errCh := cli.ContainerWait(context.Background(), containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return "", errors.Wrapf(err, "Failed to wait on container %s", containerResp.ID)
	case <-waitCh:

	}

	defer func() {
		err = cli.ContainerRemove(context.Background(), containerResp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Printf("Failed to remove container %s", containerResp.ID)
		}
	}()

	reader, err := cli.ContainerLogs(context.Background(), containerResp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get container logs from %s", containerResp.ID)
	}

	logs, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read logs %s", containerResp.ID)
	}

	return string(logs), nil
}
