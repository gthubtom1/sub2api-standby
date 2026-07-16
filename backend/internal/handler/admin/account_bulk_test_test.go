package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupBulkTestRouter(handler *AccountHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/admin/accounts/bulk-test", handler.BulkTest)
	return r
}

func TestAccountHandlerBulkTestRequiresIDs(t *testing.T) {
	handler := &AccountHandler{accountTestService: &service.AccountTestService{}}
	router := setupBulkTestRouter(handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/bulk-test", bytes.NewBufferString(`{"account_ids":[]}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountHandlerBulkTestRejectsOversizeBatch(t *testing.T) {
	handler := &AccountHandler{accountTestService: &service.AccountTestService{}}
	router := setupBulkTestRouter(handler)

	ids := make([]int64, service.AccountBulkTestMaxBatchSize+1)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	body, err := json.Marshal(map[string]any{"account_ids": ids})
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/bulk-test", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAccountHandlerBulkTestServiceUnavailable(t *testing.T) {
	handler := &AccountHandler{}
	router := setupBulkTestRouter(handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/bulk-test", bytes.NewBufferString(`{"account_ids":[1]}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}
