package admin

import (
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type bulkTestAccountsRequest struct {
	AccountIDs  []int64 `json:"account_ids" binding:"required"`
	ModelID     string  `json:"model_id"`
	Concurrency int     `json:"concurrency"`
}

// BulkTest handles bulk account connectivity testing.
// POST /api/v1/admin/accounts/bulk-test
func (h *AccountHandler) BulkTest(c *gin.Context) {
	if h.accountTestService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Account test service unavailable")
		return
	}

	var req bulkTestAccountsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if len(req.AccountIDs) == 0 {
		response.BadRequest(c, "account_ids is required")
		return
	}
	if len(req.AccountIDs) > service.AccountBulkTestMaxBatchSize {
		response.BadRequest(c, "account_ids exceeds maximum batch size of 100")
		return
	}

	results := h.accountTestService.BulkTestAccounts(c.Request.Context(), service.BulkTestAccountsRequest{
		AccountIDs:  req.AccountIDs,
		ModelID:     req.ModelID,
		Concurrency: req.Concurrency,
	})

	successCount := 0
	failedCount := 0
	for _, item := range results {
		if item.Status == "success" {
			successCount++
			if h.rateLimitService != nil {
				if _, err := h.rateLimitService.RecoverAccountAfterSuccessfulTest(c.Request.Context(), item.AccountID); err != nil {
					_ = c.Error(err)
				}
			}
			continue
		}
		failedCount++
	}

	response.Success(c, gin.H{
		"total":   len(results),
		"success": successCount,
		"failed":  failedCount,
		"results": results,
	})
}
