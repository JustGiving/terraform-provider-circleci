package circleci

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/mrolla/terraform-provider-circleci/circleci/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCLECI_TOKEN", nil),
				Description: "The token key for API operations.",
			},
			"vcs_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCLECI_VCS_TYPE", "github"),
				Description: "The VCS type for the organization.",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCLECI_ORGANIZATION", nil),
				Description: "The CircleCI organization.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCLECI_URL", "https://circleci.com/api/v2/"),
				Description: "The URL of the Circle CI API (v2)",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"circleci_project":                      resourceCircleCIProject(),
			"circleci_environment_variable":         resourceCircleCIEnvironmentVariable(),
			"circleci_context":                      resourceCircleCIContext(),
			"circleci_context_environment_variable": resourceCircleCIContextEnvironmentVariable(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"circleci_context": dataSourceCircleCIContext(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return client.New(client.Config{
		URL:          d.Get("url").(string),
		Token:        d.Get("api_token").(string),
		Organization: d.Get("organization").(string),
		VCS:          d.Get("vcs_type").(string),
	})
}
