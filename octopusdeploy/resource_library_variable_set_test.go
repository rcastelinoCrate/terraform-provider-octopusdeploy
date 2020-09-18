package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOctopusDeployLibraryVariableSetBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_library_variable_set.foo"
	const libraryVariableSetName = "Funky Set"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployLibraryVariableSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLibraryVariableSetBasic(libraryVariableSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", libraryVariableSetName),
				),
			},
		},
	})
}

func TestAccOctopusDeployLibraryVariableSetWithUpdate(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_library_variable_set.foo"
	const libraryVariableSetName = "Funky Set"
	const description = "I am a new set description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployLibraryVariableSetDestroy,
		Steps: []resource.TestStep{
			// create variable set with no description
			{
				Config: testAccLibraryVariableSetBasic(libraryVariableSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", libraryVariableSetName),
				),
			},
			// create update it with a description
			{
				Config: testAccLibraryVariableSetWithDescription(libraryVariableSetName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", libraryVariableSetName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", description),
				),
			},
			// update again by remove its description
			{
				Config: testAccLibraryVariableSetBasic(libraryVariableSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", libraryVariableSetName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", ""),
				),
			},
		},
	})
}

func testAccLibraryVariableSetBasic(name string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_library_variable_set" "foo" {
			name           = "%s"
		  }
		`,
		name,
	)
}
func testAccLibraryVariableSetWithDescription(name, description string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_library_variable_set" "foo" {
			name           = "%s"
			description    = "%s"
		  }
		`,
		name, description,
	)
}

func testAccCheckOctopusDeployLibraryVariableSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	if err := destroyHelperLibraryVariableSet(s, client); err != nil {
		return err
	}
	if err := destroyEnvHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployLibraryVariableSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		if err := existsHelperLibraryVariableSet(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyHelperLibraryVariableSet(s *terraform.State, apiClient *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := apiClient.LibraryVariableSets.Get(r.Primary.ID); err != nil {
			if err == client.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving library variable set %s", err)
		}
		return fmt.Errorf("library variable set still exists")
	}
	return nil
}

func existsHelperLibraryVariableSet(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_libraryVariableSet" {
			if _, err := client.LibraryVariableSets.Get(r.Primary.ID); err != nil {
				return fmt.Errorf("received an error retrieving library variable set %s", err)
			}
		}
	}
	return nil
}
