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

// Acceptance tests create every project they use inside a dedicated Railway workspace and
// delete it afterwards, so they never touch real projects. The workspace id is set on each
// project because an account token can see several workspaces.
var testAccWorkspaceId = os.Getenv("RAILWAY_TEST_WORKSPACE_ID")

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"RAILWAY_TOKEN", "RAILWAY_TEST_WORKSPACE_ID"} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s must be set for acceptance tests", name)
		}
	}
}

// testAccFixtureConfig declares a throwaway project, whose default environment is named
// `staging`, and one service in it, for tests of resources that attach to a service. The
// name is fixed because every test step re-renders the configuration, and a changing name
// would replace the project between steps; Railway does not require unique project names.
func testAccFixtureConfig() string {
	return fmt.Sprintf(`
resource "railway_project" "fixture" {
  name         = "tf-acc-fixture"
  workspace_id = "%s"

  default_environment = {
    name = "staging"
  }
}

resource "railway_service" "fixture" {
  name       = "fixture"
  project_id = railway_project.fixture.id
}
`, testAccWorkspaceId)
}

// fillFixtureTemplate replaces `{project}`, `{environment}` and `{service}` with the ids the
// fixture resources received in state.
func fillFixtureTemplate(state *terraform.State, template string) (string, error) {
	project, ok := state.RootModule().Resources["railway_project.fixture"]
	if !ok {
		return "", fmt.Errorf("fixture project is missing from state")
	}
	service, ok := state.RootModule().Resources["railway_service.fixture"]
	if !ok {
		return "", fmt.Errorf("fixture service is missing from state")
	}

	replacer := strings.NewReplacer(
		"{project}", project.Primary.ID,
		"{environment}", project.Primary.Attributes["default_environment.id"],
		"{service}", service.Primary.ID,
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
