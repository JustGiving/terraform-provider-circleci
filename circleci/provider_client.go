package circleci

import (
	"net/url"

	circleciapi "github.com/jszwedko/go-circleci"
)

// ProviderClient is a thin commodity wrapper on top of circleciapi
type ProviderClient struct {
	client       *circleciapi.Client
	vcsType      string
	organization string
}

// NewConfig initialize circleci API client and returns a new config object
func NewConfig(token, vscType, organization, baseURL string) (*ProviderClient, error) {
	parsedUrl, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &ProviderClient{
		client: &circleciapi.Client{
			BaseURL: parsedUrl,
			Token:   token,
		},
		vcsType:      vscType,
		organization: organization,
	}, nil
}

func (pv *ProviderClient) GetProject(projectName string) (*circleciapi.Project, error) {
	return pv.client.GetProject(pv.organization, projectName)
}

func (pv *ProviderClient) FollowProject(projectName string) (*circleciapi.Project, error) {
	return pv.client.FollowProject(pv.vcsType, pv.organization, projectName)
}

// GetEnvVar get the environment variable with given name
// It returns an empty structure if no environment variable exists with that name
func (pv *ProviderClient) GetEnvVar(projectName, envVarName string) (*circleciapi.EnvVar, error) {
	return pv.client.GetEnvVar(pv.vcsType, pv.organization, projectName, envVarName)
}

// EnvVarExists check if environment variable exists with given name
func (pv *ProviderClient) EnvVarExists(projectName, envVarName string) (bool, error) {
	envVar, err := pv.client.GetEnvVar(pv.vcsType, pv.organization, projectName, envVarName)
	if err != nil {
		return false, err
	}
	return bool(envVar.Name != ""), nil
}

// AddEnvVar create an environment variable with given name and value
func (pv *ProviderClient) AddEnvVar(projectName, envVarName, envVarValue string) (*circleciapi.EnvVar, error) {
	return pv.client.AddEnvVar(pv.vcsType, pv.organization, projectName, envVarName, envVarValue)
}

// DeleteEnvVar delete the environment variable with given name
func (pv *ProviderClient) DeleteEnvVar(projectName, envVarName string) error {
	return pv.client.DeleteEnvVar(pv.vcsType, pv.organization, projectName, envVarName)
}
