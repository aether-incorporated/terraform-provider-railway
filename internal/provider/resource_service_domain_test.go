package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccServiceDomainResourceDefault(t *testing.T) {
	// Railway subdomains are global, so each run claims fresh ones.
	subdomain := "tf-acc-" + acctest.RandString(8)
	secondSubdomain := subdomain + "-2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServiceDomainResourceConfigDefault(subdomain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", subdomain),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "environment_id", "railway_project.fixture", "default_environment.id"),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "project_id", "railway_project.fixture", "id"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "domain", subdomain+".up.railway.app"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "suffix", "up.railway.app"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_service_domain.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{service}:staging:" + subdomain + ".up.railway.app"),
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccServiceDomainResourceConfigDefault(subdomain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", subdomain),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "environment_id", "railway_project.fixture", "default_environment.id"),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "project_id", "railway_project.fixture", "id"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "domain", subdomain+".up.railway.app"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "suffix", "up.railway.app"),
				),
			},
			// Update with default values
			{
				Config: testAccServiceDomainResourceConfigDefault(secondSubdomain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", secondSubdomain),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "environment_id", "railway_project.fixture", "default_environment.id"),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_service_domain.test", "project_id", "railway_project.fixture", "id"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "domain", secondSubdomain+".up.railway.app"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "suffix", "up.railway.app"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_service_domain.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{service}:staging:" + secondSubdomain + ".up.railway.app"),
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccServiceDomainResourceExternalDeletion(t *testing.T) {
	projectName := "tf-acc-domain-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	subdomain := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	config := testAccServiceDomainResourceExternalDeletionConfig(projectName, subdomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					testAccDeleteServiceDomain,
				),
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", subdomain),
				),
			},
		},
	})
}

func testAccDeleteServiceDomain(state *terraform.State) error {
	domain, ok := state.RootModule().Resources["railway_service_domain.test"]
	if !ok {
		return fmt.Errorf("service domain resource is missing from state")
	}
	_, err := deleteServiceDomain(context.Background(), testAccClient(), domain.Primary.ID)
	return err
}

func testAccServiceDomainResourceConfigDefault(name string) string {
	return testAccFixtureConfig() + fmt.Sprintf(`
resource "railway_service_domain" "test" {
  subdomain = "%s"
  environment_id = railway_project.fixture.default_environment.id
  service_id = railway_service.fixture.id
}
`, name)
}

func testAccServiceDomainResourceExternalDeletionConfig(projectName string, subdomain string) string {
	return fmt.Sprintf(`
resource "railway_project" "test" {
  workspace_id = "`+testAccWorkspaceId+`"
  name = "%s"
}

resource "railway_service" "test" {
  name       = "domain-test"
  project_id = railway_project.test.id
}

resource "railway_service_domain" "test" {
  subdomain     = "%s"
  environment_id = railway_project.test.default_environment.id
  service_id    = railway_service.test.id
}
`, projectName, subdomain)
}
