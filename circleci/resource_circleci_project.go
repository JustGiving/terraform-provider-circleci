package circleci

import (
	"time"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCircleCIProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCIProjectCreate,
		Read:   resourceCircleCIProjectRead,
		Delete: resourceCircleCIProjectDelete,
		// Exists: resourceCircleCIProjectExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the project",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCircleCIProjectCreate(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	projectName := d.Get("name").(string)

	if _, err := providerClient.FollowProject(projectName); err != nil {
		return err
	}

	d.SetId(projectName)

	return resourceCircleCIProjectRead(d, m)
}

func resourceCircleCIProjectRead(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	projectName := d.Get("name").(string)

	project, err := providerClient.GetProject(projectName)
	if err != nil {
		return err
	}

	if project == nil {
		return errors.New("Invalid response from CircleCI API")
	}

	if err := d.Set("name", project.Reponame); err != nil {
		return err
	}

	return nil
}

func resourceCircleCIProjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// func resourceCircleCIProjectExists(d *schema.ResourceData, m interface{}) (bool, error) {
// 	providerClient := m.(*ProviderClient)

// 	projectName := d.Get("name").(string)

// 	exists, err := providerClient.ProjectExists(projectName)

// 	if err != nil {
// 		return false, err
// 	}

// 	return exists, nil
// }
