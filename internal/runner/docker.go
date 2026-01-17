package runner

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerRunner struct {
	Cli *client.Client
}

func NewDockerRunner() (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerRunner{Cli: cli}, nil
}

func (r *DockerRunner) CreateApp(ctx context.Context, name, image string, userID int) error {
	containerName := fmt.Sprintf("poor-exe-%s", name)

	// Pull image if not present (ignore errors, image may exist locally)
	pullResp, err := r.Cli.ImagePull(ctx, image, client.ImagePullOptions{})
	if err == nil {
		pullResp.Wait(ctx)
		pullResp.Close()
	}

	resp, err := r.Cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: containerName,
		Config: &container.Config{
			Image: image,
			Labels: map[string]string{
				"poor-exe": "true",
				"user_id":  fmt.Sprintf("%d", userID),
				"app_name": name,
			},
			Tty: true,
		},
	})
	if err != nil {
		return err
	}

	if _, err := r.Cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (r *DockerRunner) Attach(ctx context.Context, appName string, stdin io.Reader, stdout, stderr io.Writer, sess ssh.Session) error {
	containerName := fmt.Sprintf("poor-exe-%s", appName)

	execConfig := client.ExecCreateOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          true,
		Cmd:          []string{"/bin/sh"},
	}

	execIDResp, err := r.Cli.ExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return err
	}

	execStartConfig := client.ExecAttachOptions{
		TTY: true,
	}

	resp, err := r.Cli.ExecAttach(ctx, execIDResp.ID, execStartConfig)
	if err != nil {
		return err
	}
	defer resp.Close()

	// Handle window resize if it's a TTY
	_, windowChanges, isPty := sess.Pty()
	if isPty {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case win, ok := <-windowChanges:
					if !ok {
						return
					}
					r.Cli.ExecResize(ctx, execIDResp.ID, client.ExecResizeOptions{
						Height: uint(win.Height),
						Width:  uint(win.Width),
					})
				}
			}
		}()
	}

	// Keepalive
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sess.SendRequest("keepalive@openssh.com", true, nil)
			}
		}
	}()

	// Bridge data
	errCh := make(chan error, 2)
	go func() {
		io.Copy(resp.Conn, sess)
		errCh <- nil
	}()
	go func() {
		io.Copy(sess, resp.Conn)
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *DockerRunner) GetAppStatus(ctx context.Context, appName string) (string, error) {
	containerName := fmt.Sprintf("poor-exe-%s", appName)
	inspect, err := r.Cli.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
	if err != nil {
		// If container not found, return that
		return "stopped", nil
	}
	return string(inspect.Container.State.Status), nil
}
