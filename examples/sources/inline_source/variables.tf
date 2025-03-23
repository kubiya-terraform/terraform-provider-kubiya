variable "tool_name" {
  description = "tool name"
  type        = string
  default     = "python tool"
}
variable "tool_type" {
  description = "tool type"
  type        = string
  default     = "docker"
}
variable "tool_image" {
  description = "tool image"
  type        = string
  default     = "python:latest"
}
variable "tool_content" {
  description = "tool content"
  type        = string
  default     = "echo 'hello world'"
}
variable "tool_description" {
  description = "tool description"
  type        = string
  default     = "python tool that echoes hello world"
}

variable "source_runner" {
  description = "inline source runner"
  type        = string
  default     = "runnerv2-5-vcluster"
}
variable "source_config" {
  description = "List of Kubiya integrations to enable. Supports multiple values. For AWS integration, the main account must be provided."
  type        = string
  default     = "{}"
}

