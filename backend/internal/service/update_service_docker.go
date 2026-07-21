package service

import (
	"context"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/sysutil"
)

// performDockerHotUpdate pulls the standby GHCR image. Container recreate is
// deferred to Restart so the HTTP response can complete (same UX as binary update).
func (s *UpdateService) performDockerHotUpdate(ctx context.Context) error {
	if !sysutil.DockerHotUpdateConfigured() {
		return ErrDockerUpdateOnly
	}
	if ctx == nil {
		ctx = context.Background()
	}
	type result struct {
		upToDate bool
		err      error
	}
	ch := make(chan result, 1)
	go func() {
		up, err := sysutil.DockerPullLatest()
		ch <- result{upToDate: up, err: err}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Minute):
		return infraerrors.BadRequest("DOCKER_PULL_TIMEOUT", "docker pull timed out after 10 minutes")
	case r := <-ch:
		if r.err != nil {
			return infraerrors.BadRequest("DOCKER_PULL_FAILED", r.err.Error())
		}
		if r.upToDate {
			return ErrNoUpdateAvailable
		}
		return nil
	}
}

func dockerUpdateConfigured() bool {
	return sysutil.DockerHotUpdateConfigured()
}
