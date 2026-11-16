package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
)

const projectRateLimitMessage = "input:3: projectCreate You are creating projects too quickly. This workspace allows 1 project per 30 seconds. Try again shortly."

type failingGraphQLClient struct {
	failures []error
	calls    int
}

func (c *failingGraphQLClient) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	c.calls++
	if c.calls <= len(c.failures) {
		return c.failures[c.calls-1]
	}
	return nil
}

type recordingSleeper struct {
	waits []time.Duration
}

func (s *recordingSleeper) sleep(ctx context.Context, duration time.Duration) error {
	s.waits = append(s.waits, duration)
	return nil
}

func TestCreateProjectRetriesWhileRateLimited(t *testing.T) {
	client := &failingGraphQLClient{failures: []error{
		errors.New(projectRateLimitMessage),
		errors.New(projectRateLimitMessage),
	}}
	sleeper := &recordingSleeper{}

	_, err := createProjectWithRetry(context.Background(), client, ProjectCreateInput{Name: "p"}, sleeper.sleep)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.calls != 3 {
		t.Fatalf("got %d calls, want 3", client.calls)
	}
	if len(sleeper.waits) != 2 {
		t.Fatalf("got %d waits, want 2", len(sleeper.waits))
	}
}

func TestCreateProjectDoesNotRetryOtherErrors(t *testing.T) {
	client := &failingGraphQLClient{failures: []error{errors.New("input:3: projectCreate Not Authorized")}}
	sleeper := &recordingSleeper{}

	_, err := createProjectWithRetry(context.Background(), client, ProjectCreateInput{Name: "p"}, sleeper.sleep)

	if err == nil {
		t.Fatal("expected the error to be returned")
	}
	if client.calls != 1 || len(sleeper.waits) != 0 {
		t.Fatalf("got %d calls and %d waits, want 1 and 0", client.calls, len(sleeper.waits))
	}
}

func TestCreateProjectGivesUpAfterTheRetryBudget(t *testing.T) {
	failures := make([]error, 50)
	for i := range failures {
		failures[i] = errors.New(projectRateLimitMessage)
	}
	client := &failingGraphQLClient{failures: failures}
	sleeper := &recordingSleeper{}

	_, err := createProjectWithRetry(context.Background(), client, ProjectCreateInput{Name: "p"}, sleeper.sleep)

	if err == nil {
		t.Fatal("expected the rate limit error once the budget is spent")
	}
	var total time.Duration
	for _, wait := range sleeper.waits {
		total += wait
	}
	if total > projectCreateRetryBudget {
		t.Fatalf("waited %s, more than the %s budget", total, projectCreateRetryBudget)
	}
	if client.calls != len(sleeper.waits)+1 {
		t.Fatalf("got %d calls for %d waits", client.calls, len(sleeper.waits))
	}
}

func TestCreateProjectStopsWhenTheContextEnds(t *testing.T) {
	client := &failingGraphQLClient{failures: []error{errors.New(projectRateLimitMessage)}}
	cancelled := func(ctx context.Context, duration time.Duration) error {
		return context.Canceled
	}

	_, err := createProjectWithRetry(context.Background(), client, ProjectCreateInput{Name: "p"}, cancelled)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
	if client.calls != 1 {
		t.Fatalf("got %d calls, want 1", client.calls)
	}
}
