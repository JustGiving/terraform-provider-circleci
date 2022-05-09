package client

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/mrolla/terraform-provider-circleci/circleci/client/rest"
)

type project struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Id   string `json:"id"`
}

type projectEnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) FollowProject(name string) error {
	slug, err := c.Slug(c.organization, name)
	if err != nil {
		return err
	}

	u := &url.URL{
		Path: fmt.Sprintf("project/%s/follow", slug),
	}

	req, err := c.v1.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	_, err = c.rest.DoRequest(req, nil)
	if err != nil {
		var httpError *rest.HTTPError

		if errors.As(err, &httpError) && httpError.Code == 400 && httpError.Message == "Branch not found" {
			// The API responds with "400: Branch not found" if the GitHub repository for the project is empty.
			// However, it still follows the project and the pipeline will kick off after the first commit with the config file gets pushed.
			//
			// Therefore, we can ignore this error response here.
			return nil
		}

		return err
	}

	return nil
}

func (c *Client) GetProject(name string) (*project, error) {
	slug, err := c.Slug(c.organization, name)
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Path: fmt.Sprintf("project/%s", slug),
	}

	req, err := c.rest.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	project := &project{}

	_, err = c.rest.DoRequest(req, project)
	if err != nil {
		var httpError *rest.HTTPError
		if errors.As(err, &httpError) && httpError.Code == 404 {
			return nil, nil
		}

		return nil, err
	}

	return project, nil
}

// HasProjectEnvironmentVariable checks for the existence of a matching project environment variable by name
func (c *Client) HasProjectEnvironmentVariable(org, project, name string) (bool, error) {
	slug, err := c.Slug(org, project)
	if err != nil {
		return false, err
	}

	u := &url.URL{
		Path: fmt.Sprintf("project/%s/envvar/%s", slug, name),
	}

	req, err := c.rest.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}

	_, err = c.rest.DoRequest(req, nil)
	if err != nil {
		var httpError *rest.HTTPError
		if errors.As(err, &httpError) && httpError.Code == 404 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateProjectEnvironmentVariable creates a new project environment variable
func (c *Client) CreateProjectEnvironmentVariable(org, project, name, value string) error {
	slug, err := c.Slug(org, project)
	if err != nil {
		return err
	}

	u := &url.URL{
		Path: fmt.Sprintf("project/%s/envvar", slug),
	}

	req, err := c.rest.NewRequest("POST", u, &projectEnvironmentVariable{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return err
	}

	_, err = c.rest.DoRequest(req, nil)
	return err
}

// DeleteProjectEnvironmentVariable deletes an existing project environment variable
func (c *Client) DeleteProjectEnvironmentVariable(org, project, name string) error {
	slug, err := c.Slug(org, project)
	if err != nil {
		return err
	}

	u := &url.URL{
		Path: fmt.Sprintf("project/%s/envvar/%s", slug, name),
	}

	req, err := c.rest.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = c.rest.DoRequest(req, nil)
	return err
}
