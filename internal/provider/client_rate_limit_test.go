package provider

import (
	"context"
	"errors"
	"testing"
	"time"
)

const (
	projectRateLimitMessage     = "input:3: projectCreate You are creating projects too quickly. This workspace allows 1 project per 30 seconds. Try again shortly."
	environmentRateLimitMessage = "input:3: environmentCreate You are creating environments too quickly. This workspace allows 1 environment per 30 seconds. Try again shortly."
)

// scriptedCall fails with the given errors in order, then succeeds.
type scriptedCall struct {
	failures []error
	calls    int
}

func (c *scriptedCall) run() (string, error) {
	c.calls++
	if c.calls <= len(c.failures) {
		return "", c.failures[c.calls-1]
	}
	return "created", nil
}

type recordingSleeper struct {
	waits []time.Duration
}

func (s *recordingSleeper) sleep(ctx context.Context, duration time.Duration) error {
	s.waits = append(s.waits, duration)
	return nil
}

func TestRetryWhileRateLimitedRetriesProjectAndEnvironmentLimits(t *testing.T) {
	for _, message := range []string{projectRateLimitMessage, environmentRateLimitMessage} {
		call := &scriptedCall{failures: []error{errors.New(message), errors.New(message)}}
		sleeper := &recordingSleeper{}

		result, err := retryWhileRateLimited(context.Background(), sleeper.sleep, call.run)

		if err != nil || result != "created" {
			t.Fatalf("got (%q, %v), want the created result", result, err)
		}
		if call.calls != 3 || len(sleeper.waits) != 2 {
			t.Fatalf("got %d calls and %d waits, want 3 and 2", call.calls, len(sleeper.waits))
		}
	}
}

func TestRetryWhileRateLimitedDoesNotRetryOtherErrors(t *testing.T) {
	call := &scriptedCall{failures: []error{errors.New("input:3: projectCreate Not Authorized")}}
	sleeper := &recordingSleeper{}

	_, err := retryWhileRateLimited(context.Background(), sleeper.sleep, call.run)

	if err == nil {
		t.Fatal("expected the error to be returned")
	}
	if call.calls != 1 || len(sleeper.waits) != 0 {
		t.Fatalf("got %d calls and %d waits, want 1 and 0", call.calls, len(sleeper.waits))
	}
}

func TestRetryWhileRateLimitedGivesUpAfterTheBudget(t *testing.T) {
	failures := make([]error, 50)
	for i := range failures {
		failures[i] = errors.New(environmentRateLimitMessage)
	}
	call := &scriptedCall{failures: failures}
	sleeper := &recordingSleeper{}

	_, err := retryWhileRateLimited(context.Background(), sleeper.sleep, call.run)

	if err == nil {
		t.Fatal("expected the rate limit error once the budget is spent")
	}
	var total time.Duration
	for _, wait := range sleeper.waits {
		total += wait
	}
	if total > rateLimitRetryBudget {
		t.Fatalf("waited %s, more than the %s budget", total, rateLimitRetryBudget)
	}
	if call.calls != len(sleeper.waits)+1 {
		t.Fatalf("got %d calls for %d waits", call.calls, len(sleeper.waits))
	}
}

func TestRetryWhileRateLimitedStopsWhenTheContextEnds(t *testing.T) {
	call := &scriptedCall{failures: []error{errors.New(projectRateLimitMessage)}}
	cancelled := func(ctx context.Context, duration time.Duration) error {
		return context.Canceled
	}

	_, err := retryWhileRateLimited(context.Background(), cancelled, call.run)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
	if call.calls != 1 {
		t.Fatalf("got %d calls, want 1", call.calls)
	}
}
