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


variable "mermaid" {
  type    = string
  default = <<-EOT
    flowchart TD
    A["User Initiates Request"] --> B["AI Analysis Engine"]
    B --> C["Policy Generation"]
    C --> D["Risk Assessment"]
    D --> E["Admin Review Queue"]
    E --> F{"Multi-Level Decision"}
    F -->|"Approved (L1)"| G["Secondary Review"]
    G -->|"Approved (L2)"| H["Policy Activation"]
    F -->|"Rejected"| I["Request Denied"]
    G -->|"Rejected"| I
    H --> J["Active Permission"]
    J --> K["Continuous Monitoring"]
    K --> L["Auto-Cleanup"]
    I --> M["Feedback to User"]
    L --> N["Access Log Updated"]

    style A fill:#4aa1ff,stroke:#333,stroke-width:2px,color:#fff
    style B fill:#4aa1ff,stroke:#333,stroke-width:2px,color:#fff
    style C fill:#3ebd64,stroke:#333,stroke-width:2px,color:#fff
    style D fill:#ff9800,stroke:#333,stroke-width:2px,color:#fff
    style E fill:#ff9800,stroke:#333,stroke-width:2px,color:#fff
    style F fill:#9c27b0,stroke:#333,stroke-width:2px,color:#fff
    style G fill:#9c27b0,stroke:#333,stroke-width:2px,color:#fff
    style H fill:#3ebd64,stroke:#333,stroke-width:2px,color:#fff
    style I fill:#e91e63,stroke:#333,stroke-width:2px,color:#fff
    style J fill:#3ebd64,stroke:#333,stroke-width:2px,color:#fff
    style K fill:#ff9800,stroke:#333,stroke-width:2px,color:#fff
    style L fill:#666666,stroke:#333,stroke-width:2px,color:#fff
    style M fill:#666666,stroke:#333,stroke-width:2px,color:#fff
    style N fill:#666666,stroke:#333,stroke-width:2px,color:#fff
  EOT
}
