package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDomainResourceDefault(t *testing.T) {
	// Custom domains are unique across Railway, so each run attaches a fresh one.
	domain := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum) + ".example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomDomainResourceConfigDefault(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_custom_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "domain", domain),
					resource.TestCheckNoResourceAttr("railway_custom_domain.test", "target_port"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "environment_id", "railway_environment.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "project_id", "railway_environment.fixture", "project_id"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "host_label", "terraform"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "zone", "example.com"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "dns_record_value"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_host_label"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_record_value"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_custom_domain.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{service}:{environment_name}:" + domain),
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccCustomDomainResourceConfigDefault(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_custom_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "domain", domain),
					resource.TestCheckNoResourceAttr("railway_custom_domain.test", "target_port"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "environment_id", "railway_environment.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "project_id", "railway_environment.fixture", "project_id"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "host_label", "terraform"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "zone", "example.com"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "dns_record_value"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_host_label"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_record_value"),
				),
			},
			// Update and Read testing
			{
				Config: testAccCustomDomainResourceConfigNonDefault(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_custom_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "domain", domain),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "target_port", "3000"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "environment_id", "railway_environment.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "project_id", "railway_environment.fixture", "project_id"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "host_label", "terraform"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "zone", "example.com"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "dns_record_value"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_host_label"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_record_value"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_custom_domain.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{service}:{environment_name}:" + domain),
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCustomDomainResourceNonDefault(t *testing.T) {
	// Custom domains are unique across Railway, so each run attaches a fresh one.
	domain := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum) + ".example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomDomainResourceConfigNonDefault(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_custom_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "domain", domain),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "target_port", "3000"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "environment_id", "railway_environment.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "project_id", "railway_environment.fixture", "project_id"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "host_label", "terraform"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "zone", "example.com"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "dns_record_value"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_host_label"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_record_value"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_custom_domain.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{service}:{environment_name}:" + domain),
				ImportStateVerify: true,
			},
			// Update with same values
			{
				Config: testAccCustomDomainResourceConfigNonDefault(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_custom_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "domain", domain),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "target_port", "3000"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "environment_id", "railway_environment.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "project_id", "railway_environment.fixture", "project_id"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "host_label", "terraform"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "zone", "example.com"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "dns_record_value"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_host_label"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_record_value"),
				),
			},
			// Update with default values
			{
				Config: testAccCustomDomainResourceConfigDefault(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_custom_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "domain", domain),
					resource.TestCheckNoResourceAttr("railway_custom_domain.test", "target_port"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "environment_id", "railway_environment.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "service_id", "railway_service.fixture", "id"),
					resource.TestCheckResourceAttrPair("railway_custom_domain.test", "project_id", "railway_environment.fixture", "project_id"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "host_label", "terraform"),
					resource.TestCheckResourceAttr("railway_custom_domain.test", "zone", "example.com"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "dns_record_value"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_host_label"),
					resource.TestCheckResourceAttrSet("railway_custom_domain.test", "verification_record_value"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_custom_domain.test",
				ImportState:       true,
				ImportStateIdFunc: testAccFixtureImportStateId("{service}:{environment_name}:" + domain),
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomDomainResourceConfigDefault(name string) string {
	return testAccFixtureConfig() + fmt.Sprintf(`
resource "railway_custom_domain" "test" {
  domain = "%s"
  environment_id = railway_environment.fixture.id
  service_id = railway_service.fixture.id
}
`, name)
}

func testAccCustomDomainResourceConfigNonDefault(name string) string {
	return testAccFixtureConfig() + fmt.Sprintf(`
resource "railway_custom_domain" "test" {
  domain = "%s"
  environment_id = railway_environment.fixture.id
  service_id = railway_service.fixture.id

  target_port = 3000
}
`, name)
}
