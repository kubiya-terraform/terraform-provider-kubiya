terraform {
  required_providers {
    kubiya = {
      # source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_inline_source" "with_all" {
  name = "with_all"
  tools = [
    {
      name        = var.tool_name
      type        = var.tool_type
      image       = var.tool_image
      content     = var.tool_content
      description = var.tool_description
    }
  ]
  runner         = var.source_runner
  dynamic_config = var.source_config
}
output "with_all" {
  value = kubiya_inline_source.with_all
}

# resource "kubiya_inline_source" "no_name" {
#   tools = [
#     {
#       name        = var.tool_name
#       type        = var.tool_type
#       image       = var.tool_image
#       content     = var.tool_content
#       description = var.tool_description
#     }
#   ]
#   runner         = var.source_runner
#   dynamic_config = var.source_config
# }
# output "no_name" {
#   value = kubiya_inline_source.no_name
# }
#
# resource "kubiya_inline_source" "no_tools" {
#   name           = "no_tools"
#   runner         = var.source_runner
#   dynamic_config = var.source_config
# }
# output "no_tools" {
#   value = kubiya_inline_source.no_tools
# }