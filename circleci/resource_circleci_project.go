package circleci

import (
	"time"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	circleciapi "github.com/jszwedko/go-circleci"
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
		apiErr := &circleciapi.APIError{}

		if errors.As(err, &apiErr) && apiErr.HTTPStatusCode == 400 && apiErr.Message == "Branch not found" {
			// The API responds with "400: Branch not found" if the GitHub repository for the project is empty.
			// However, it still follows the project and the pipeline will kick off after the first commit with the config file gets pushed.
			//
			// Therefore, we can ignore this error response here.
		} else {
			return err
		}
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
