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

	hearbeatPkg "github.com/result17/codeBeatCli/pkg/entity"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	v.Set("project", "test-cli")
	v.Set("lineno", 19)
	v.Set("plugin", plugin)
	v.Set("time", 1585598059.1)
	v.Set("timeout", 5)

	offlineQueueFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)
	defer offlineQueueFile.Close()

	router.HandleFunc("/cb/heartbeats", func(w http.ResponseWriter, r *http.Request) {
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
