package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	hearbeatAPI "github.com/result17/codeBeatCli/internal/api"
	hearbeatPkg "github.com/result17/codeBeatCli/pkg/entity"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO FIX response json
func TestSendHeartbeats(t *testing.T) {
	testURL, router, tearDown := setupTestServer()
	defer tearDown()

	var (
		plugin   = "codebeat/0.0.1"
		numCalls int
	)

	v := viper.New()
	v.Set("api-url", testURL)
	v.Set("cursorpos", 125)
	v.Set("entity", "testdata/main.go")
	v.Set("language", "Go")
	v.Set("alternate-project", "test-cli")
	v.Set("lineno", 19)
	v.Set("lines-in-file", 38)
	v.Set("plugin", plugin)
	v.Set("time", 1585598059.1)
	v.Set("timeout", 5)
	v.Set("project-path", "/sys/usr/codebeat/")

	offlineQueueFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	defer offlineQueueFile.Close()

	router.HandleFunc(hearbeatAPI.CollectHeartbeatRouter, func(w http.ResponseWriter, r *http.Request) {
		numCalls++
		assert.Equal(t, []string{"application/json"}, r.Header["Accept"])
		assert.Equal(t, []string{"application/json"}, r.Header["Content-Type"])

		fmtStr, err := os.ReadFile("testdata/api_heartbeats_request_template.json")
		require.NoError(t, err)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var entity struct {
			Entity string `json:"entity"`
		}

		err = json.Unmarshal(body, &[]any{&entity})
		require.NoError(t, err)

		expectedBodyStr := fmt.Sprintf(
			string(fmtStr),
			entity.Entity,
			hearbeatPkg.UserAgent(t.Context(), plugin),
		)

		assert.True(t, strings.HasSuffix(entity.Entity, "testdata/main.go"))
		assert.JSONEq(t, expectedBodyStr, string(body))

		// send response
		w.WriteHeader(http.StatusCreated)

		f, err := os.Open("testdata/api_heartbeats_response.json")
		require.NoError(t, err)
		defer f.Close()

		_, err = io.Copy(w, f)
		require.NoError(t, err)
	})

	err = hearbeatPkg.SendHeartbeats(t.Context(), v, offlineQueueFile.Name())
	require.NoError(t, err)

	assert.Eventually(t, func() bool { return numCalls == 1 }, time.Second, 50*time.Millisecond)
}

func TestSendHeartbeatsToLocalServer(t *testing.T) {
	var (
		plugin = "codebeat/0.0.123"
	)

	v := viper.New()
	v.Set("api-url", "http://127.0.0.1:3000")
	v.Set("cursorpos", 125)
	v.Set("entity", "testdata/main.go")
	v.Set("language", "Go")
	v.Set("alternate-project", "test-cli")
	v.Set("lineno", 19)
	v.Set("lines-in-file", 38)
	v.Set("plugin", plugin)
	v.Set("time", 1585598059100)
	v.Set("timeout", 5)
	v.Set("project-path", "/sys/usr/codebeat/")

	offlineQueueFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	defer offlineQueueFile.Close()

	err = hearbeatPkg.SendHeartbeats(t.Context(), v, offlineQueueFile.Name())
	require.NoError(t, err)
}

func TestSendHeartbeatsToLocalServerWithoutTime(t *testing.T) {
	var (
		plugin = "codebeat/0.0.1"
	)

	v := viper.New()
	v.Set("api-url", "http://127.0.0.1:3000")
	v.Set("cursorpos", 125)
	v.Set("entity", "testdata/main.go")
	v.Set("language", "Go")
	v.Set("alternate-project", "test-cli")
	v.Set("lineno", 19)
	v.Set("lines-in-file", 38)
	v.Set("plugin", plugin)
	v.Set("timeout", 5)
	v.Set("project-path", "/sys/usr/codebeat/")

	offlineQueueFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	defer offlineQueueFile.Close()

	err = hearbeatPkg.SendHeartbeats(t.Context(), v, offlineQueueFile.Name())
	require.NoError(t, err)
}

func TestHeartbeatResults(t *testing.T) {
	data, err := os.ReadFile("testdata/api_heartbeat_list_response.json")
	require.NoError(t, err)

	_, err = hearbeatAPI.ParseHeartbeatResponses(t.Context(), data)
	require.NoError(t, err)
}
