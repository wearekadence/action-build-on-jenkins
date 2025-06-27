package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bndr/gojenkins"
	"kadence.co/build-on-jenkins/internal"
)

func main() {
	ctx := context.Background()

	config := internal.NewConfig()

	jenkins := gojenkins.CreateJenkins(nil, config.JenkinsURL, config.Username, config.ApiToken)

	_, err := jenkins.Init(ctx)

	if err != nil {
		fmt.Printf("Error initializing Jenkins client: %v\n", err)
		panic("Something Went Wrong")
	}

	job, err := jenkins.GetJob(ctx, config.JobName)
	if err != nil {
		fmt.Printf("Failed to get job %s\n", config.JobName)
		panic(err)
	}

	lastBuild, err := job.GetLastBuild(ctx)
	if err != nil {
		fmt.Println("Failed to get last build")
		panic(err)
	}

	var params map[string]string

	err = json.Unmarshal([]byte(config.Parameters), &params)
	if err != nil {
		panic(fmt.Errorf("failed to parse parameters: %v", err))
	}

	_, err = job.InvokeSimple(ctx, params)
	if err != nil {
		fmt.Println("Failed to queue job")
		panic(err)
	}

	currentBuild, err := completeOrDie(
		func() (*gojenkins.Build, error) {
			return getCurrentBuild(job, ctx, lastBuild.GetBuildNumber())
		},
		config.StartTimeout,
	)

	if err != nil {
		fmt.Println("Failed to get current build")
		panic(err)
	}

	fmt.Println("Current Build:", currentBuild.GetBuildNumber())
	fmt.Println("URL:", currentBuild.GetUrl())
	fmt.Println("Console output will be available here after the build has finished")

	_, err = completeOrDie(
		func() (interface{}, error) {
			for currentBuild.IsRunning(ctx) {
				time.Sleep(time.Duration(config.Interval) * time.Millisecond)
				currentBuild.Poll(ctx)
			}
			return nil, nil
		},
		config.Timeout,
	)

	if err != nil {
		fmt.Println("The build failed to finish")
		panic(err)
	}

	fmt.Println(currentBuild.GetConsoleOutput(ctx))

	fmt.Printf(
		"Build number %d with result: %v\nURL: %spipeline-overview/",
		currentBuild.GetBuildNumber(),
		currentBuild.GetResult(),
		currentBuild.GetUrl(),
	)

	if currentBuild.GetResult() != "SUCCESS" {
		panic("Build failed")
	}
}

func getCurrentBuild(job *gojenkins.Job, ctx context.Context, lastBuildNumber int64) (*gojenkins.Build, error) {
	for {
		allBuildIds, err := job.GetAllBuildIds(ctx)
		if err != nil {
			fmt.Println("Failed to get build ids")
			return &gojenkins.Build{}, err
		}
		latestBuildId := allBuildIds[0]
		if latestBuildId.Number == lastBuildNumber {
			fmt.Println("Waiting for build to start")
			time.Sleep(2000 * time.Millisecond)
		} else {
			currentBuild, err := job.GetBuild(ctx, latestBuildId.Number)
			if err != nil {
				fmt.Println("Failed to get build with id")
				return &gojenkins.Build{}, err
			}
			return currentBuild, nil
		}
	}
}

type mustComplete[T any] func() (T, error)

func completeOrDie[T any](fn mustComplete[T], timeout int) (T, error) {
	c := make(chan error, 1)
	r := make(chan T, 1)
	var result T
	go func() {
		res, err := fn()
		if err != nil {
			c <- err
			return
		}
		r <- res
	}()
	select {
	case result = <-r:
		return result, nil
	case err := <-c:
		return result, err
	case <-time.After(time.Duration(timeout) * time.Second):
		return result, fmt.Errorf("operation timed out after %d seconds", timeout)
	}
}
