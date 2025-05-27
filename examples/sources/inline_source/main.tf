terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
      # source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_inline_source" "hello_world_tool" {
  name   = "integration_test_tool"
  runner = ""

  tools = jsonencode([
    {
      name        = "hello_world_tool"
      description = "A simple tool that prints 'Hello World' to the console."
      image       = "python:3.9"
      content = "print('Hello World')"
      # long_running = false
      type        = ""
      # workflow     = false
    }
  ])
}

resource "kubiya_agent" "helloworld_teammate" {
  name          = "integration_test_teammate"
  description   = "HelloWorldTeammate is designed to interact with users and execute a simple 'Hello World' tool. This teammate requires minimal permissions and serves as a basic example of Kubiya's capabilities."
  image         = "ghcr.io/kubiyabot/kubiya-agent:stable"
  model         = "azure/gpt-4o"
  runner        = "core-testing-1"
  is_debug_mode = false

  # environment_variables = {}
  # integrations          = []
  # links                 = []
  # groups                = []
  # secrets               = []
  # sources               = []
  # tool_sources          = []
  # users                 = []
  instructions = ""
}

