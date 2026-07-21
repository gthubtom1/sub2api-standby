// Package sysutil provides system-level utilities for process management.
package sysutil

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// RestartService triggers a service restart.
//
// Binary/systemd installs: exit 0 and rely on Restart=always.
// Standby Docker installs: recreate the compose service so the new image is used
// (plain exit would restart the same old image layer).
func RestartService() error {
	if runtime.GOOS != "linux" {
		log.Println("Service restart via exit only works on Linux")
		return nil
	}

	if DockerHotUpdateConfigured() {
		log.Println("Initiating Docker compose recreate for hot-update...")
		go func() {
			time.Sleep(500 * time.Millisecond)
			if err := RecreateDockerService(); err != nil {
				log.Printf("Docker recreate failed: %v", err)
				log.Println("Falling back to process exit; container may still run old image")
				os.Exit(0)
			}
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}()
		return nil
	}

	log.Println("Initiating service restart by graceful exit...")
	log.Println("systemd will automatically restart the service (Restart=always)")

	go func() {
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	return nil
}

// RestartServiceAsync is a fire-and-forget version of RestartService.
func RestartServiceAsync() {
	if err := RestartService(); err != nil {
		log.Printf("Service restart failed: %v", err)
		log.Println("Please restart the service manually: sudo systemctl restart sub2api")
	}
}

// DockerHotUpdateConfigured reports whether one-click Docker hot-update can run.
func DockerHotUpdateConfigured() bool {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("UPDATE_DOCKER_ENABLED")), "false") {
		return false
	}
	if _, err := os.Stat("/var/run/docker.sock"); err != nil {
		return false
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	return true
}

func dockerUpdateImage() string {
	if v := strings.TrimSpace(os.Getenv("UPDATE_DOCKER_IMAGE")); v != "" {
		return v
	}
	return "ghcr.io/gthubtom1/sub2api-standby:latest"
}

func dockerComposeDir() string {
	if v := strings.TrimSpace(os.Getenv("UPDATE_DOCKER_COMPOSE_DIR")); v != "" {
		return v
	}
	return "/opt/sub2api-standby"
}

func dockerComposeFile() string {
	if v := strings.TrimSpace(os.Getenv("UPDATE_DOCKER_COMPOSE_FILE")); v != "" {
		return v
	}
	return filepath.Join(dockerComposeDir(), "docker-compose.yml")
}

func dockerComposeService() string {
	if v := strings.TrimSpace(os.Getenv("UPDATE_DOCKER_SERVICE")); v != "" {
		return v
	}
	return "sub2api"
}

// DockerPullLatest pulls the standby GHCR image via the host Docker daemon.
// Returns (upToDate, error).
func DockerPullLatest() (bool, error) {
	if !DockerHotUpdateConfigured() {
		return false, fmt.Errorf("docker hot-update not configured")
	}
	image := dockerUpdateImage()
	cmd := exec.Command("docker", "pull", image)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
	if err != nil {
		return false, fmt.Errorf("docker pull %s failed: %v; output: %s", image, err, truncateOut(out, 800))
	}
	lower := strings.ToLower(out)
	if strings.Contains(lower, "image is up to date") || strings.Contains(lower, "already up to date") {
		return true, nil
	}
	return false, nil
}

// RecreateDockerService force-recreates the compose app service on the new image.
func RecreateDockerService() error {
	if !DockerHotUpdateConfigured() {
		return fmt.Errorf("docker hot-update not configured")
	}
	dir := dockerComposeDir()
	file := dockerComposeFile()
	service := dockerComposeService()
	args := []string{
		"compose",
		"-f", file,
		"--project-directory", dir,
		"up", "-d",
		"--force-recreate",
		"--no-deps",
		service,
	}
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		out := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
		return fmt.Errorf("docker compose recreate failed: %w; output: %s", err, truncateOut(out, 800))
	}
	return nil
}

func truncateOut(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
