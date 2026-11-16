package provider

import (
	"net/http"
	"os"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"railway": providerserver.NewProtocol6WithError(New("test")()),
}

// Acceptance tests run against a dedicated Railway workspace holding one fixture project, so
// they never touch real projects. The fixture project needs an environment named `staging`
// and one service in it.
var (
	testAccWorkspaceId   = os.Getenv("RAILWAY_TEST_WORKSPACE_ID")
	testAccProjectId     = os.Getenv("RAILWAY_TEST_PROJECT_ID")
	testAccEnvironmentId = os.Getenv("RAILWAY_TEST_ENVIRONMENT_ID")
	testAccServiceId     = os.Getenv("RAILWAY_TEST_SERVICE_ID")
)

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{
		"RAILWAY_TOKEN",
		"RAILWAY_TEST_WORKSPACE_ID",
		"RAILWAY_TEST_PROJECT_ID",
		"RAILWAY_TEST_ENVIRONMENT_ID",
		"RAILWAY_TEST_SERVICE_ID",
	} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s must be set for acceptance tests", name)
		}
	}
}

func testAccClient() graphql.Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			token:   os.Getenv("RAILWAY_TOKEN"),
			wrapped: http.DefaultTransport,
		},
	}
	return graphql.NewClient("https://backboard.railway.app/graphql/v2?source=terraform_provider_railway", &httpClient)
}
