package main

import (
	"context"
	"flag"
	"io"
	"os"
	"time"

	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"sync"
)

func runContainer(ctx context.Context, cli *client.Client) (string, error) {
	imageName := "nginx"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, err
}

func stopContainer(ctx context.Context, cli *client.Client, id string) {
	cli.ContainerStop(ctx, id, nil)
	cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
}

func PauseAndResume(ctx context.Context, cli *client.Client, id string) {
	var pauseTotal, resumeTotal int64 = 0, 0

	for i := 0; i < 100; i++ {
		start := time.Now()
		cli.ContainerPause(ctx, id)
		end := time.Now()
		pauseTotal += (end.Sub(start).Nanoseconds()) / 1000000

		start = time.Now()
		cli.ContainerUnpause(ctx, id)
		end = time.Now()
		resumeTotal += (end.Sub(start).Nanoseconds()) / 1000000
	}

	fmt.Printf("container %s, Average pause time : %d ms, average resume time : %d ms. \n", id, pauseTotal/100, resumeTotal/100)
}

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	var concurrency int
	flag.IntVar(&concurrency, "c", 1, "concurrency")
	flag.Parse()
	fmt.Printf("Concurrency %d\n", concurrency)

	var wg sync.WaitGroup
	wg.Add(concurrency)
	ids := make([]string, concurrency, concurrency)
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			ids[index], err = runContainer(ctx, cli)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id string) {
			PauseAndResume(ctx, cli, id)

			stopContainer(ctx, cli, id)
			wg.Done()
		}(ids[i])
	}

	wg.Wait()
}
