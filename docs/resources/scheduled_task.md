---
page_title: "kubiya_scheduled_task Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Scheduled Task resource to automate agent execution.
---

# kubiya_scheduled_task (Resource)

Provides a Kubiya Scheduled Task resource. This allows scheduled tasks to be created, updated, and deleted on the Kubiya platform. Scheduled tasks execute a Kubiya agent on a defined schedule.

## Example Usage

```hcl
# Daily task example
resource "kubiya_scheduled_task" "daily_backup_check" {
  name        = "daily-backup-check"
  description = "Verify backup completion daily"
  agent       = kubiya_agent.example.name
  schedule    = "0 9 * * *"  # 9 AM daily in cron format
  prompt      = "Check if all daily backups completed successfully"
  enabled     = true
  
  notification {
    method      = "Slack"
    destination = "#backups"
  }
}

# Weekday task example
resource "kubiya_scheduled_task" "weekly_report" {
  name        = "weekly-infrastructure-report"
  description = "Generate weekly infrastructure report"
  agent       = kubiya_agent.example.name
  schedule    = "0 10 * * MON"  # 10 AM every Monday
  prompt      = "Generate a comprehensive infrastructure status report for the past week"
  
  notification {
    method      = "Email"
    destination = "team@example.com"
  }
}

# Every 15 minutes task
resource "kubiya_scheduled_task" "monitoring_check" {
  name        = "service-health-check"
  description = "Check critical service health"
  agent       = kubiya_agent.monitoring.name
  schedule    = "*/15 * * * *"  # Every 15 minutes
  prompt      = "Check health status of critical services"
  
  notification {
    method      = "Teams"
    team_name   = "Operations"
    destination = "Alerts"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the scheduled task.
* `description` - (Required) A description of the scheduled task.
* `agent` - (Required) The name of the agent to execute.
* `schedule` - (Required) A cron expression defining when the task should run. For example:
  * `0 9 * * *` - Daily at 9 AM
  * `0 9 * * MON-FRI` - Weekdays at 9 AM
  * `0 0 1 * *` - Monthly on the 1st at midnight
  * `*/15 * * * *` - Every 15 minutes

### Optional Arguments

* `prompt` - (Optional) The prompt to send to the agent when executing.
* `enabled` - (Optional) Whether the scheduled task is enabled. Defaults to true.
* `notification` - (Optional) Configuration for task notifications.
  * `method` - The notification method. Values: "Slack", "Teams", "Email".
  * `destination` - The destination for notifications. 
  * `team_name` - Required for Teams notifications, the team name.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the scheduled task.
* `last_run` - The timestamp of the last execution.
* `next_run` - The timestamp of the next scheduled execution.
* `status` - The current status of the scheduled task.

## Import

Scheduled tasks can be imported using the `id`:

```
$ terraform import kubiya_scheduled_task.example TASK_ID
``` 