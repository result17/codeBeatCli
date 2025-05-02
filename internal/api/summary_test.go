package api_test

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/result17/codeBeatCli/internal/api"
	summaryPkg "github.com/result17/codeBeatCli/pkg/summary"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryTodaySummary(t *testing.T) {
	testURL, router, tearDown := setupTestServer()

	var (
		numCalls int
	)
	defer tearDown()
	router.HandleFunc(api.TodaySummaryAPIRouter, func(w http.ResponseWriter, r *http.Request) {
		numCalls++
		assert.Equal(t, []string{"application/json"}, r.Header["Accept"])

		f, err := os.Open("testdata/api_summary_today_response.json")
		require.NoError(t, err)
		defer f.Close()

		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, f)
		require.NoError(t, err)
	})

	v := viper.New()
	v.Set("api-url", testURL)
	v.Set("today-summary", true)

	summary, err := summaryPkg.TodaySummary(t.Context(), v)
	require.NoError(t, err)
	assert.Exactly(t, summary.GrandTotal.Text, "37 mins")
	assert.Exactly(t, len(summary.Timeline), 12)
	assert.Eventually(t, func() bool { return numCalls == 1 }, time.Second, 50*time.Millisecond)
}
