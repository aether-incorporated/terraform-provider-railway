package provider

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)

	return t.wrapped.RoundTrip(req)
}

// Railway allows one new project, and separately one new environment, per workspace every 30
// seconds, and refuses the rest with "You are creating … too quickly" without creating
// anything, so a refused create is safe to repeat.
const (
	rateLimitRetryInterval = 10 * time.Second
	rateLimitRetryBudget   = 90 * time.Second
	rateLimitText          = "too quickly"
)

func retryWhileRateLimited[T any](
	ctx context.Context,
	sleep func(context.Context, time.Duration) error,
	create func() (T, error),
) (T, error) {
	var waited time.Duration
	for {
		result, err := create()
		if err == nil || !strings.Contains(err.Error(), rateLimitText) {
			return result, err
		}
		if waited+rateLimitRetryInterval > rateLimitRetryBudget {
			return result, err
		}
		tflog.Debug(ctx, "creation rate limited by Railway, retrying", map[string]interface{}{"retry_in": rateLimitRetryInterval.String()})
		if sleepErr := sleep(ctx, rateLimitRetryInterval); sleepErr != nil {
			return result, sleepErr
		}
		waited += rateLimitRetryInterval
	}
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
