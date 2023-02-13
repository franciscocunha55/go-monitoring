package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func RunDockerCompose(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	fmt.Println(cmd)
	var stdout, stderr bytes.Buffer
	// & because the value of standard output is in bytes
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Printf("stdout: %s\nstderr: %s\n", stdout.String(), stderr.String())

	return nil
}
