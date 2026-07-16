package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// SetRateLimitService injects the live-traffic error policy engine so connection
// tests can move accounts into the same error / temp-unsched / rate-limit states
// that real gateway requests would apply.
func (s *AccountTestService) SetRateLimitService(svc *RateLimitService) {
	if s == nil {
		return
	}
	s.rateLimitService = svc
}

// applyTestUpstreamFailure applies the same upstream error policy used by the
// gateway when a connection test receives a non-success HTTP status.
//
// Without this, bad accounts (401/403/429/503 with rules, custom codes, etc.)
// stay schedulable until a real user request burns client retries (Codex, etc.).
func (s *AccountTestService) applyTestUpstreamFailure(
	ctx context.Context,
	account *Account,
	statusCode int,
	headers http.Header,
	body []byte,
	requestedModel string,
) {
	if s == nil || account == nil || statusCode < 400 {
		return
	}
	if headers == nil {
		headers = make(http.Header)
	}

	// Preserve OpenAI plan-type extraction that test probes historically did on 429.
	if statusCode == http.StatusTooManyRequests && account.IsOpenAI() && s.accountRepo != nil {
		persistOpenAI429PlanType(ctx, s.accountRepo, account, body)
	}

	if s.rateLimitService != nil {
		model := strings.TrimSpace(requestedModel)
		if model != "" {
			s.rateLimitService.HandleUpstreamError(ctx, account, statusCode, headers, body, model)
		} else {
			s.rateLimitService.HandleUpstreamError(ctx, account, statusCode, headers, body)
		}
		return
	}

	// Fallback for unit tests / partially-wired environments: keep historical
	// partial behavior so callers still get something useful.
	if s.accountRepo == nil {
		return
	}
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusPaymentRequired:
		msg := fmt.Sprintf("API returned %d: %s", statusCode, truncateRunes(string(body), 512))
		_ = s.accountRepo.SetError(ctx, account.ID, msg)
	case http.StatusTooManyRequests:
		s.reconcileOpenAI429State(ctx, account, headers, body)
	}
}
