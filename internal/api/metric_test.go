package api_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/result17/codeBeatCli/pkg/exitcode"
	metricPkg "github.com/result17/codeBeatCli/pkg/metric"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryTodayMetricDuration(t *testing.T) {
	testURL, router, tearDown := setupTestServer()

	var (
		numCalls int
	)
	defer tearDown()
	router.HandleFunc(fmt.Sprintf("/api/metric/duration/today/%s", "project"), func(w http.ResponseWriter, r *http.Request) {
		numCalls++
		assert.Equal(t, []string{"application/json"}, r.Header["Accept"])

		f, err := os.Open("testdata/api_metric_duration_response.json")
		require.NoError(t, err)
		defer f.Close()

		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, f)
		require.NoError(t, err)
	})

	v := viper.New()
	v.Set("api-url", testURL)
	v.Set("today-metric-duration", "project")

	metric, err := metricPkg.TodayMetricDuration[string](t.Context(), v)
	require.NoError(t, err)
	assert.Exactly(t, len(metric.Ratios), 8)
	assert.Exactly(t, metric.GrandTotal.Text, "9 hrs 4 mins")
	assert.Eventually(t, func() bool { return numCalls == 1 }, time.Second, 50*time.Millisecond)
}

func TestQueryTodayMetricDurationWithLocalServer(t *testing.T) {
	v := viper.New()
	v.Set("api-url", "http://127.0.0.1:3000")
	v.Set("today-metric-duration", "project")

	code, err := metricPkg.TypeRun(t.Context(), v, "project")
	require.NoError(t, err)
	assert.Exactly(t, code, exitcode.Success)
}
