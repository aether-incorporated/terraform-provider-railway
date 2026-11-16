package provider

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"railway": providerserver.NewProtocol6WithError(New("test")()),
}

// Acceptance tests run inside a dedicated Railway workspace, so they never touch real projects.
// Tests of resources that attach to a service share one long-lived project in it and create
// their own environment and service there, because Railway allows only one new project per
// workspace every 30 seconds. Tests of project-level behaviour still create their own project,
// pinned to the workspace because an account token can see several.
var (
	testAccWorkspaceId = os.Getenv("RAILWAY_TEST_WORKSPACE_ID")
	testAccProjectId   = os.Getenv("RAILWAY_TEST_PROJECT_ID")

	// Fixed for the whole run, because every test step re-renders its configuration and a
	// changing name would replace the environment between steps. Random across runs, so an
	// environment left behind by an interrupted run never collides with the next one.
	testAccFixtureEnvironmentName = "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
)

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"RAILWAY_TOKEN", "RAILWAY_TEST_WORKSPACE_ID", "RAILWAY_TEST_PROJECT_ID"} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s must be set for acceptance tests", name)
		}
	}
}

// testAccFixtureConfig declares a throwaway environment in the shared project.
func testAccFixtureConfig() string {
	return fmt.Sprintf(`
resource "railway_environment" "fixture" {
  name       = "%s"
  project_id = "%s"
}
`, testAccFixtureEnvironmentName, testAccProjectId)
}

// testAccFixtureServiceConfig adds a service to the fixture environment, for tests of
// resources that attach to one. Railway caps new accounts at 25 new services a day, so tests
// that do not attach to a service use testAccFixtureConfig instead.
func testAccFixtureServiceConfig() string {
	return testAccFixtureConfig() + `
resource "railway_service" "fixture" {
  name       = "fixture"
  project_id = railway_environment.fixture.project_id
}
`
}

// fillFixtureTemplate replaces `{project}`, `{environment}`, `{environment_name}` and
// `{service}` with the values the fixture resources received in state.
func fillFixtureTemplate(state *terraform.State, template string) (string, error) {
	environment, ok := state.RootModule().Resources["railway_environment.fixture"]
	if !ok {
		return "", fmt.Errorf("fixture environment is missing from state")
	}

	serviceId := ""
	if strings.Contains(template, "{service}") {
		service, ok := state.RootModule().Resources["railway_service.fixture"]
		if !ok {
			return "", fmt.Errorf("fixture service is missing from state")
		}
		serviceId = service.Primary.ID
	}

	replacer := strings.NewReplacer(
		"{project}", environment.Primary.Attributes["project_id"],
		"{environment_name}", environment.Primary.Attributes["name"],
		"{environment}", environment.Primary.ID,
		"{service}", serviceId,
	)
	return replacer.Replace(template), nil
}

func testAccFixtureImportStateId(template string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		return fillFixtureTemplate(state, template)
	}
}

func testAccCheckFixtureAttr(name string, key string, template string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		expected, err := fillFixtureTemplate(state, template)
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(name, key, expected)(state)
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
