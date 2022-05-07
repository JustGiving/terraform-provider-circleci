terraform {
  required_providers {
    circleci = {
      source = "terraform.justgiving.com/justgiving/circleci"
      version = "0.0.1"
    }
  }
}

variable circleci_token {
  type = string
}

provider "circleci" {
  api_token = var.circleci_token
  vcs_type = "github"
  organization = "JustGiving"
}

resource "circleci_project" "service" {
  name = "JG.Scorecards"
}

output project_id {
  value = circleci_project.service.id
}
