terraform {
  required_providers {
    kubiya = {
      source = "local/provider/example"
    }
  }
}

provider "kubiya" {
  //Your Kubiya API Key will be taken from the
  //environment variable KUBIYA_API_KEY
  //To set the key, please use export KUBIYA_API_KEY="YOUR_API_KEY"
}

resource "kubiya_scheduled_task" "example" {
  repeat         = "" // Optional. Allowed values: hourly, daily, weekly, monthly. Leaving this value empty or omitting it will cause the task to be executed only once.
  channel_id    = "C082X4R0FL0"
  agent          = "Which Agent should perform the task?"
  scheduled_time = "2024-12-01T05:00:00"
  description    = "Describe the task"
}

output "output" {
  value = kubiya_scheduled_task.example
}