terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
      # source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_source" "git_source" {
  url            = "https://github.com/test-org-for-project/test-repo/blob/main/tools/test.yaml"
}

output "output" {
  value = kubiya_source.git_source
}