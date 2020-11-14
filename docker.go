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

func BuildImage(tag string, cli client.Client, ctx context.Context) error {
	contextPath, _ := filepath.Abs("./files/")
	buildOpts := types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		Tags:       []string{tag},
	}

	buildCtx, err := archive.TarWithOptions(contextPath, &archive.TarOptions{})
	if err != nil {
		return err
	}

	resp, err := cli.ImageBuild(ctx, buildCtx, buildOpts)
	if err != nil {
		return errors.Wrap(err, "Failed to build image")
	}
	defer resp.Body.Close()

	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stderr, termFd, isTerm, nil)
}

func PullImageIfNotExists(cli *client.Client, tag string, ctx context.Context) error {
	list, err := cli.ImageList(ctx, types.ImageListOptions{})
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
		pullResp, err := cli.ImagePull(ctx, tag, types.ImagePullOptions{})
		if err != nil {
			return errors.Wrapf(err, "Failed to pull image %s", tag)
		}
		defer pullResp.Close()

		termFd, isTerm := term.GetFdInfo(os.Stderr)
		return jsonmessage.DisplayJSONMessagesStream(pullResp, os.Stderr, termFd, isTerm, nil)
	}
	return nil
}

func EvaluateScript(script string, ctx context.Context) (string, error) {
	tag := "lukaspj/t3deval:4_0Preview"

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	err = PullImageIfNotExists(cli, tag, ctx)
	if err != nil {
		return "", err
	}

	var containerResp container.ContainerCreateCreatedBody
	containerResp, err = cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        tag,
			Cmd:          []string{script},
			Tty:          true,
			NetworkDisabled: true,
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		nil,
		"t3deval-worker-1",
	)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create container %s", containerResp.ID)
	}

	err = cli.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to run container %s", containerResp.ID)
	}

	defer func() {
		err = cli.ContainerRemove(context.Background(), containerResp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Printf("Failed to remove container %s", containerResp.ID)
		}
	}()

	timeoutCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	waitCh, errCh := cli.ContainerWait(timeoutCtx, containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return "", errors.Wrapf(err, "Failed to wait on container %s", containerResp.ID)
	case <-waitCh:

	}

	reader, err := cli.ContainerLogs(ctx, containerResp.ID, types.ContainerLogsOptions{
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
