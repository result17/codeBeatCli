package api

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/result17/codeBeatCli/internal/api"
	"github.com/result17/codeBeatCli/pkg/log"

	tz "github.com/gandarez/go-olson-timezone"
)

func NewClient(ctx context.Context, url string) (*api.Client, error) {
	return newClient(ctx, url)
}

func newClient(ctx context.Context, url string) (*api.Client, error) {
	logger := log.Extract(ctx)
	logger.Debugf("Creating client, the baseurl is %s", api.BaseURL)
	return api.NewClient(url), nil
}

func timezone() (name string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked: failed to get timezone: %v. Stack: %s", r, string(debug.Stack()))
		}
	}()

	name, err = tz.Name()

	return name, err
}
