package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/result17/codeBeatCli/internal/api"
	"github.com/result17/codeBeatCli/internal/summary"
	"github.com/result17/codeBeatCli/pkg/duration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() (string, *http.ServeMux, func()) {
	router := http.NewServeMux()
	srv := httptest.NewServer(router)
	return srv.URL, router, func() { srv.Close() }
}

func TestQueryTodayDuration(t *testing.T) {
	testURL, router, tearDown := setupTestServer()

	var (
		numCalls int
	)
	defer tearDown()
	router.HandleFunc(api.TodayDurationAPIRouter, func(w http.ResponseWriter, r *http.Request) {
		numCalls++
		assert.Equal(t, []string{"application/json"}, r.Header["Accept"])

		grandTotal, err := summary.NewGrandTotal(0)
		require.NoError(t, err)

		rsp, err := json.Marshal(grandTotal)

		w.WriteHeader(http.StatusOK)
		w.Write(rsp)
	})

	v := viper.New()
	v.Set("api-url", testURL)
	v.Set("today-duration", true)

	text, err := duration.TodayDuration(t.Context(), v)
	require.NoError(t, err)
	assert.Exactly(t, text, "")
}

func TestParseDurationResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/api_duration_today_response.json")
	require.NoError(t, err)

	_, err = api.ParseClientGrandTotalResponse(data)
	require.NoError(t, err)
}

func TestQueryTodayDurationWithLocalServer(t *testing.T) {
	TestSendHeartbeatsToLocalServer(t)
	v := viper.New()
	v.Set("api-url", "http://127.0.0.1:3000")
	v.Set("today-duration", true)

	text, err := duration.TodayDuration(t.Context(), v)
	require.NoError(t, err)
	assert.NotEmpty(t, text)
}
