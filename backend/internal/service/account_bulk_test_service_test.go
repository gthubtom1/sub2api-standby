package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruncateRunes(t *testing.T) {
	require.Equal(t, "hello", truncateRunes("hello", 10))
	require.Equal(t, "hel…", truncateRunes("hello", 3))
	require.Equal(t, "测…", truncateRunes("测试文本", 1))
}

func TestBulkTestAccountsEmpty(t *testing.T) {
	svc := &AccountTestService{}
	results := svc.BulkTestAccounts(context.Background(), BulkTestAccountsRequest{})
	require.Empty(t, results)
}

func TestBulkTestAccountsDedupAndSkipInvalid(t *testing.T) {
	// Without a repo, GetByID is skipped; RunTestBackground will fail fast with
	// "Account not found" style errors. We only assert orchestration shape here.
	svc := &AccountTestService{}
	results := svc.BulkTestAccounts(context.Background(), BulkTestAccountsRequest{
		AccountIDs:  []int64{0, 1, 1, -3, 2},
		Concurrency: 2,
	})
	require.Len(t, results, 2)
	require.Equal(t, int64(1), results[0].AccountID)
	require.Equal(t, int64(2), results[1].AccountID)
	for _, item := range results {
		require.Equal(t, "failed", item.Status)
		require.Equal(t, "account repository unavailable", item.ErrorMessage)
	}
}
