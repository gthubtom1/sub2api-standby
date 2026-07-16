//go:build unit

package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestApplyTestUpstreamFailure_UsesRateLimitServiceForAPIKey401(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	rateLimitSvc := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	svc := &AccountTestService{
		accountRepo:      repo,
		rateLimitService: rateLimitSvc,
	}
	account := &Account{
		ID:       1,
		Name:     "apikey-1",
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Status:   StatusActive,
		Credentials: map[string]any{
			"api_key": "bad",
		},
	}

	svc.applyTestUpstreamFailure(context.Background(), account, http.StatusUnauthorized, http.Header{}, []byte(`{"error":{"message":"bad key"}}`), "claude-sonnet-4")
	require.Equal(t, 1, repo.setErrorCalls)
	require.Equal(t, int64(1), repo.lastErrorID)
	require.Contains(t, repo.lastErrorMsg, "401")
}

func TestApplyTestUpstreamFailure_FallbackWithoutRateLimitService(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	svc := &AccountTestService{accountRepo: repo}
	account := &Account{
		ID:       2,
		Name:     "apikey-2",
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Status:   StatusActive,
	}

	svc.applyTestUpstreamFailure(context.Background(), account, http.StatusForbidden, nil, []byte("forbidden"), "")
	require.Equal(t, 1, repo.setErrorCalls)
	require.Equal(t, int64(2), repo.lastErrorID)
}
