package circleci

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	client "github.com/mrolla/terraform-provider-circleci/circleci/client"
)

func resourceCircleCIProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCIProjectCreate,
		Read:   resourceCircleCIProjectRead,
		Delete: resourceCircleCIProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the project",
			},
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceCircleCIProjectResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceCircleCIProjectUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceCircleCIProjectResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the project",
			},
		},
	}
}

func resourceCircleCIProjectUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	c := meta.(*client.Client)

	project, err := c.GetProject(rawState["name"].(string))
	if err != nil {
		return nil, err
	}

	rawState["id"] = project.Id

	return rawState, nil
}

func resourceCircleCIProjectCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	projectName := d.Get("name").(string)

	err := c.FollowProject(projectName)

	if err != nil {
		return fmt.Errorf("error following project: %w", err)
	}

	d.SetId(projectName)

	return resourceCircleCIProjectRead(d, m)
}

func resourceCircleCIProjectRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	projectName := d.Get("name").(string)

	project, err := c.GetProject(projectName)
	if err != nil {
		return err
	}

	if project == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", project.Name)
	return nil
}

func resourceCircleCIProjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
