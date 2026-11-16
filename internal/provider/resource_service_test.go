package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig("tf-acc-service"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service.test", "name", "tf-acc-service"),
					resource.TestCheckResourceAttr("railway_service.test", "project_id", testAccProjectId),
				),
			},
			{
				ResourceName:      "railway_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceResourceConfig("tf-acc-service-renamed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service.test", "name", "tf-acc-service-renamed"),
					resource.TestCheckResourceAttr("railway_service.test", "project_id", testAccProjectId),
				),
			},
		},
	})
}

func testAccServiceResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "railway_service" "test" {
  name       = "%s"
  project_id = "%s"
}
`, name, testAccProjectId)
}
