package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSharedVariableResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSharedVariableResourceConfigDefault("1234567890"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFixtureAttr("railway_shared_variable.test", "id", "{project}:{environment}:API_KEY"),
					resource.TestCheckResourceAttr("railway_shared_variable.test", "name", "API_KEY"),
					resource.TestCheckResourceAttr("railway_shared_variable.test", "value", "1234567890"),
					resource.TestCheckResourceAttrPair("railway_shared_variable.test", "environment_id", "railway_project.fixture", "default_environment.id"),
					resource.TestCheckResourceAttrPair("railway_shared_variable.test", "project_id", "railway_project.fixture", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_shared_variable.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{project}:staging:API_KEY"),
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccSharedVariableResourceConfigDefault("1234567890"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFixtureAttr("railway_shared_variable.test", "id", "{project}:{environment}:API_KEY"),
					resource.TestCheckResourceAttr("railway_shared_variable.test", "name", "API_KEY"),
					resource.TestCheckResourceAttr("railway_shared_variable.test", "value", "1234567890"),
					resource.TestCheckResourceAttrPair("railway_shared_variable.test", "environment_id", "railway_project.fixture", "default_environment.id"),
					resource.TestCheckResourceAttrPair("railway_shared_variable.test", "project_id", "railway_project.fixture", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccSharedVariableResourceConfigDefault("nice"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFixtureAttr("railway_shared_variable.test", "id", "{project}:{environment}:API_KEY"),
					resource.TestCheckResourceAttr("railway_shared_variable.test", "name", "API_KEY"),
					resource.TestCheckResourceAttr("railway_shared_variable.test", "value", "nice"),
					resource.TestCheckResourceAttrPair("railway_shared_variable.test", "environment_id", "railway_project.fixture", "default_environment.id"),
					resource.TestCheckResourceAttrPair("railway_shared_variable.test", "project_id", "railway_project.fixture", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_shared_variable.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{project}:staging:API_KEY"),
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSharedVariableResourceConfigDefault(value string) string {
	return testAccFixtureConfig() + fmt.Sprintf(`
resource "railway_shared_variable" "test" {
  name = "API_KEY"
  value = "%s"
  environment_id = railway_project.fixture.default_environment.id
  project_id = railway_project.fixture.id
}
`, value)
}
