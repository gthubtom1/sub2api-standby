package service

import (
	"context"
	"sync"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
)

const (
	// AccountBulkTestMaxBatchSize limits one HTTP request to keep timeouts predictable.
	AccountBulkTestMaxBatchSize = 100
	// AccountBulkTestDefaultConcurrency is used when the client does not specify concurrency.
	AccountBulkTestDefaultConcurrency = 5
	// AccountBulkTestMaxConcurrency hard-caps concurrent upstream probes per request.
	AccountBulkTestMaxConcurrency = 20
	// accountBulkTestResponseLimit truncates response text in bulk results.
	accountBulkTestResponseLimit = 200
)

// BulkTestAccountResult is one account's connection-test outcome.
type BulkTestAccountResult struct {
	AccountID    int64  `json:"account_id"`
	Name         string `json:"name,omitempty"`
	Platform     string `json:"platform,omitempty"`
	Status       string `json:"status"` // success | failed
	ErrorMessage string `json:"error_message,omitempty"`
	LatencyMs    int64  `json:"latency_ms"`
	ResponseText string `json:"response_text,omitempty"`
}

// BulkTestAccountsRequest controls a bulk connection test run.
type BulkTestAccountsRequest struct {
	AccountIDs  []int64
	ModelID     string
	Concurrency int
}

// BulkTestAccounts runs connection tests for many accounts concurrently.
// It reuses RunTestBackground so behavior matches scheduled/single tests.
func (s *AccountTestService) BulkTestAccounts(ctx context.Context, req BulkTestAccountsRequest) []BulkTestAccountResult {
	if len(req.AccountIDs) == 0 {
		return []BulkTestAccountResult{}
	}

	concurrency := req.Concurrency
	if concurrency <= 0 {
		concurrency = AccountBulkTestDefaultConcurrency
	}
	if concurrency > AccountBulkTestMaxConcurrency {
		concurrency = AccountBulkTestMaxConcurrency
	}

	// Deduplicate while preserving first-seen order.
	seen := make(map[int64]struct{}, len(req.AccountIDs))
	ids := make([]int64, 0, len(req.AccountIDs))
	for _, id := range req.AccountIDs {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return []BulkTestAccountResult{}
	}

	// Preload names/platforms for richer result rows.
	nameByID := make(map[int64]string, len(ids))
	platformByID := make(map[int64]string, len(ids))
	for _, id := range ids {
		if s.accountRepo == nil {
			break
		}
		acc, err := s.accountRepo.GetByID(ctx, id)
		if err != nil || acc == nil {
			continue
		}
		nameByID[id] = acc.Name
		platformByID[id] = acc.Platform
	}

	results := make([]BulkTestAccountResult, len(ids))
	var mu sync.Mutex
	var g errgroup.Group
	g.SetLimit(concurrency)

	for i, accountID := range ids {
		idx := i
		id := accountID
		g.Go(func() error {
			if ctx.Err() != nil {
				mu.Lock()
				results[idx] = BulkTestAccountResult{
					AccountID:    id,
					Name:         nameByID[id],
					Platform:     platformByID[id],
					Status:       "failed",
					ErrorMessage: ctx.Err().Error(),
				}
				mu.Unlock()
				return nil
			}

			item := BulkTestAccountResult{
				AccountID: id,
				Name:      nameByID[id],
				Platform:  platformByID[id],
				Status:    "failed",
			}
			if s.accountRepo == nil {
				item.ErrorMessage = "account repository unavailable"
				mu.Lock()
				results[idx] = item
				mu.Unlock()
				return nil
			}

			result, err := s.RunTestBackground(ctx, id, req.ModelID)
			if err != nil {
				item.ErrorMessage = err.Error()
			}
			if result != nil {
				item.Status = result.Status
				item.ErrorMessage = result.ErrorMessage
				item.LatencyMs = result.LatencyMs
				item.ResponseText = truncateRunes(result.ResponseText, accountBulkTestResponseLimit)
				if item.Status == "" {
					if item.ErrorMessage == "" {
						item.Status = "success"
					} else {
						item.Status = "failed"
					}
				}
			}
			if item.Status != "success" && item.ErrorMessage == "" {
				item.ErrorMessage = "test failed"
			}

			mu.Lock()
			results[idx] = item
			mu.Unlock()
			return nil
		})
	}

	_ = g.Wait()
	return results
}

func truncateRunes(s string, limit int) string {
	if limit <= 0 || s == "" {
		return s
	}
	if utf8.RuneCountInString(s) <= limit {
		return s
	}
	runes := []rune(s)
	return string(runes[:limit]) + "…"
}
