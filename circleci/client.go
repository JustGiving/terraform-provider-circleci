package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseEndpoint   = "https://circleci.com/api/v1.1"
	envvarEndpoint = "envvar"
)

// EnvironmentVariable inside a CircleCI project
type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Client for the CircleCI API
type Client struct {
	token        string
	vcsType      string
	organization string
	httpClient   *http.Client
}

// NewClient creates a new CircleCI API client
func NewClient(token, vcsType, organization string) (*Client, error) {
	return &Client{
		token:        token,
		vcsType:      vcsType,
		organization: organization,
		httpClient:   http.DefaultClient,
	}, nil
}

func (c *Client) buildApiURL(projectName string, endpoint string) string {
	return fmt.Sprintf("%s/project/%s/%s/%s/%s", baseEndpoint, c.vcsType, c.organization, projectName, endpoint)
}

// AddEnvironmentVariable creates a new environment variable.
// https://circleci.com/docs/api/#add-environment-variables
func (c *Client) AddEnvironmentVariable(projectName, envName, envValue string) error {
	endpointURL := c.buildApiURL(projectName, envvarEndpoint)

	e := EnvironmentVariable{
		Name:  envName,
		Value: envValue,
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(e); err != nil {
		// TODO(matteo): proper error handling
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpointURL, b)
	if err != nil {
		// TODO(matteo): proper error handling
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(c.token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// TODO(matteo): proper error handling
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("client: create wrong status code %d", resp.StatusCode)
	}

	return nil
}

// GetEnvironmentVariable returns the value of an environment variable of a
// project given its name.
// https://circleci.com/docs/api/#get-single-environment-variable
func (c *Client) GetEnvironmentVariable(projectName, envName string) (*EnvironmentVariable, error) {
	endpointURL := fmt.Sprintf("%s/%s", c.buildApiURL(projectName, envvarEndpoint), envName)

	req, err := http.NewRequest(http.MethodGet, endpointURL, nil)
	if err != nil {
		// TODO(matteo): proper error handling
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// TODO(matteo): proper error handling
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: get wrong status code %d", resp.StatusCode)
	}

	e := new(EnvironmentVariable)
	if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
		// TODO(matteo): proper error handling
		return nil, err
	}

	return e, nil
}

// DeleteEnvironmentVariable deletes an environment variable from a project
// given its name.
// https://circleci.com/docs/api/#delete-environment-variables
func (c *Client) DeleteEnvironmentVariable(projectName, envName string) error {
	endpointURL := fmt.Sprintf("%s/%s", c.buildApiURL(projectName, envvarEndpoint), envName)

	req, err := http.NewRequest(http.MethodDelete, endpointURL, nil)
	if err != nil {
		// TODO(matteo): proper error handling
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// TODO(matteo): proper error handling
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("client: delete wrong status code %d", resp.StatusCode)
	}

	return nil
}
