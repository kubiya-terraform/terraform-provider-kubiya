terraform {
  required_providers {
    kubiya = {
      # source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_inline_source" "example" {
  name = "mevrat-inline-tools-1" // Required
  runner = var.source_runner // Optional + Computed
  dynamic_config = var.source_config // Optional + Computed

  tools = [
    // Required
    {
      name = "tool_1" // Required
      type = "" // Optional + Computed
      description = "Example build tool"
      workflow = false // Optional + Computed
      long_running = false // Optional + Computed
      icon = "icon.png" // Optional
      image = "example/image:latest" // Optional
      content = "Example content" // Optional
      mermaid = "Example mermaid content" // Optional
      on_start = "start_script.sh" // Optional
      on_build = "build_script.sh" // Optional
      on_complete = "complete_script.sh" // Optional
      env = ["ENV_VAR2=value2"] // Optional
      secrets = ["SECRET_VAR2"] // Optional
      entrypoint = ["entrypoint.sh"] // Optional

      args = [
        // Optional
        {
          name = "arg1" // Required
          description = "Argument 1" // Required
          required = true // Optional
          default = "default_value" // Optional
          options = ["option2"] // Optional
          type = "" // Optional + Computed
          options_from = {
            // Optional
            image = "example/image:v2" // Required
            script = "options_script_2.sh" // Required
          }
        }
      ]

      files = [
        // Optional
        {
          source = "" // Optional
          destination = "" // Required
          content = "" // Optional:
        }
      ]
    }
  ]
}
output "with_all" {
  value = kubiya_inline_source.example
}