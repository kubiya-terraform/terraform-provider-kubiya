---
page_title: "kubiya_scheduled_task Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_scheduled_task resource manages scheduled automation tasks in the Kubiya platform.
---

# kubiya_scheduled_task (Resource)

The `kubiya_scheduled_task` resource allows you to create and manage scheduled tasks in the Kubiya platform. Scheduled tasks execute Kubiya agents on defined schedules using cron expressions, enabling automated recurring operations.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. At least one configured agent to execute
4. Understanding of cron expression syntax

## Example Usage

### 1. Basic Daily Task

Create a simple daily scheduled task:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # API key is automatically read from KUBIYA_API_KEY environment variable
}

resource "kubiya_agent" "daily_reporter" {
  name         = "daily-reporter"
  runner       = "kubiya-hosted"
  description  = "Daily reporting agent"
  instructions = "You generate daily reports and summaries."
}

resource "kubiya_scheduled_task" "daily_report" {
  name        = "daily-status-report"
  description = "Generate daily status report"
  agent       = kubiya_agent.daily_reporter.name
  schedule    = "0 9 * * *"  # 9 AM daily
  prompt      = "Generate a comprehensive daily status report for all systems"
  enabled     = true
}
```

**Expected Outcome**: Creates a task that runs daily at 9 AM to generate status reports.

### 2. Backup Verification Task

Schedule regular backup verifications:

```hcl
resource "kubiya_agent" "backup_agent" {
  name         = "backup-verifier"
  runner       = "kubiya-hosted"
  description  = "Backup verification agent"
  instructions = "You verify backup completion and integrity."
  
  integrations = ["aws"]
}

resource "kubiya_scheduled_task" "backup_check" {
  name        = "nightly-backup-verification"
  description = "Verify nightly backup completion"
  agent       = kubiya_agent.backup_agent.name
  schedule    = "30 2 * * *"  # 2:30 AM daily
  prompt      = "Check if all nightly backups completed successfully and verify their integrity"
  enabled     = true
  
  notification {
    method      = "Slack"
    destination = "#backups"
  }
}
```

**Expected Outcome**: Creates a nightly task that verifies backups and sends notifications to Slack.

### 3. Weekly Report Generation

Generate weekly reports:

```hcl
resource "kubiya_scheduled_task" "weekly_report" {
  name        = "weekly-infrastructure-report"
  description = "Generate weekly infrastructure report"
  agent       = "infrastructure-agent"
  schedule    = "0 10 * * MON"  # 10 AM every Monday
  prompt      = <<-EOT
    Generate a comprehensive weekly infrastructure report including:
    - Resource utilization trends
    - Incident summary
    - Performance metrics
    - Cost analysis
    - Recommendations for optimization
  EOT
  enabled     = true
  
  notification {
    method      = "Email"
    destination = "team@example.com"
  }
}
```

**Expected Outcome**: Creates a weekly task that generates infrastructure reports every Monday.

### 4. Frequent Health Checks

Schedule frequent health monitoring:

```hcl
resource "kubiya_agent" "health_monitor" {
  name         = "health-monitor"
  runner       = "kubiya-hosted"
  description  = "System health monitoring agent"
  instructions = "You monitor system health and alert on issues."
  
  integrations = ["kubernetes", "datadog"]
}

resource "kubiya_scheduled_task" "health_check" {
  name        = "frequent-health-check"
  description = "Check critical service health"
  agent       = kubiya_agent.health_monitor.name
  schedule    = "*/15 * * * *"  # Every 15 minutes
  prompt      = "Check health status of critical services and alert if any issues are detected"
  enabled     = true
  
  notification {
    method      = "Http"
    destination = "https://monitoring.example.com/webhook"
  }
}
```

**Expected Outcome**: Creates a task that runs every 15 minutes to check system health.

### 5. Business Hours Monitoring

Monitor only during business hours:

```hcl
resource "kubiya_scheduled_task" "business_hours_monitor" {
  name        = "business-hours-monitoring"
  description = "Monitor application performance during business hours"
  agent       = "performance-agent"
  schedule    = "*/30 9-17 * * MON-FRI"  # Every 30 minutes, 9 AM-5 PM, Monday-Friday
  prompt      = "Monitor application performance metrics and report any degradation"
  enabled     = true
  
  notification {
    method      = "Slack"
    destination = "#performance"
  }
}
```

**Expected Outcome**: Creates a task that monitors performance during business hours only.

### 6. Monthly Maintenance Tasks

Schedule monthly maintenance:

```hcl
resource "kubiya_agent" "maintenance_agent" {
  name         = "maintenance-agent"
  runner       = "kubiya-hosted"
  description  = "System maintenance agent"
  instructions = <<-EOT
    You perform system maintenance tasks including:
    - Log rotation and cleanup
    - Database optimization
    - Cache clearing
    - Temporary file cleanup
  EOT
  
  integrations = ["kubernetes", "aws"]
}

resource "kubiya_scheduled_task" "monthly_maintenance" {
  name        = "monthly-system-maintenance"
  description = "Perform monthly system maintenance"
  agent       = kubiya_agent.maintenance_agent.name
  schedule    = "0 3 1 * *"  # 3 AM on the 1st of each month
  prompt      = "Execute monthly maintenance tasks: log cleanup, database optimization, cache clearing"
  enabled     = true
  
  notification {
    method      = "Email"
    destination = "ops@example.com"
  }
}
```

**Expected Outcome**: Creates a monthly maintenance task that runs on the first day of each month.

### 7. Multi-Environment Scheduled Tasks

Create scheduled tasks for multiple environments:

```hcl
locals {
  environments = {
    dev = {
      schedule = "0 6 * * *"  # 6 AM daily
      prompt   = "Check development environment health"
    }
    staging = {
      schedule = "0 7 * * *"  # 7 AM daily
      prompt   = "Verify staging environment readiness"
    }
    prod = {
      schedule = "*/10 * * * *"  # Every 10 minutes
      prompt   = "Monitor production environment health and performance"
    }
  }
}

resource "kubiya_agent" "env_monitors" {
  for_each = local.environments
  
  name         = "${each.key}-monitor"
  runner       = "kubiya-hosted"
  description  = "Monitor for ${each.key} environment"
  instructions = "You monitor the ${each.key} environment health and performance."
  
  environment_variables = {
    ENVIRONMENT = each.key
  }
}

resource "kubiya_scheduled_task" "env_tasks" {
  for_each = local.environments
  
  name        = "${each.key}-health-check"
  description = "Health check for ${each.key} environment"
  agent       = kubiya_agent.env_monitors[each.key].name
  schedule    = each.value.schedule
  prompt      = each.value.prompt
  enabled     = each.key == "prod" ? true : false  # Only enable production by default
  
  notification {
    method      = "Slack"
    destination = "#${each.key}-alerts"
  }
}
```

**Expected Outcome**: Creates environment-specific scheduled tasks with different schedules and configurations.

### 8. Conditional Scheduling

Create tasks with conditional execution:

```hcl
resource "kubiya_scheduled_task" "conditional_task" {
  name        = "conditional-deployment-check"
  description = "Check deployments based on conditions"
  agent       = "deployment-agent"
  schedule    = "0 */4 * * *"  # Every 4 hours
  prompt      = <<-EOT
    Check if there are pending deployments:
    1. If pending deployments exist, validate them
    2. If validation passes, approve for next window
    3. If no pending deployments, check system health
    4. Generate appropriate report based on findings
  EOT
  enabled     = true
  
  notification {
    method      = "Slack"
    destination = "#deployments"
  }
}
```

**Expected Outcome**: Creates a task that performs different actions based on conditions.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the scheduled task. Must be unique within your organization.
* `description` - (Required, String) A description of the scheduled task's purpose.
* `agent` - (Required, String) The name of the agent to execute.
* `schedule` - (Required, String) Cron expression defining the schedule. Format: "minute hour day month weekday".
* `prompt` - (Required, String) The prompt to send to the agent when the task executes.

### Optional Arguments

* `enabled` - (Optional, Boolean) Whether the scheduled task is enabled. Defaults to `true`.
* `notification` - (Optional, Block) Notification configuration:
  - `method` - (Required within block, String) Notification method: "Slack", "Email", "Http", "Teams".
  - `destination` - (Required within block, String) Destination for notifications (channel, email, URL).
  - `team_name` - (Optional within block, String) Team name for Microsoft Teams notifications.

## Cron Expression Reference

The schedule uses standard cron expression format:

```
* * * * *
│ │ │ │ │
│ │ │ │ └─── Day of week (0-7, MON-SUN)
│ │ │ └───── Month (1-12, JAN-DEC)
│ │ └─────── Day of month (1-31)
│ └───────── Hour (0-23)
└─────────── Minute (0-59)
```

Common patterns:
- `0 9 * * *` - Daily at 9 AM
- `0 0 * * 0` - Weekly on Sunday at midnight
- `*/15 * * * *` - Every 15 minutes
- `0 9-17 * * MON-FRI` - Hourly during business hours
- `0 0 1 * *` - Monthly on the 1st at midnight

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the scheduled task.
* `created_at` - The timestamp when the scheduled task was created.
* `updated_at` - The timestamp when the scheduled task was last updated.
* `last_run` - The timestamp of the last execution.
* `next_run` - The timestamp of the next scheduled execution.
* `status` - The current status of the scheduled task.

## Import

Scheduled tasks can be imported using their ID:

```shell
terraform import kubiya_scheduled_task.example <scheduled-task-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Cron expressions are evaluated in UTC by default
* Minimum schedule frequency may be limited by platform tier
* Agent must exist and be accessible before task creation

## Best Practices

1. **Time Zones**: Be aware that schedules run in UTC; adjust accordingly
2. **Frequency**: Avoid overly frequent schedules that may overwhelm the system
3. **Error Handling**: Include error handling in agent prompts
4. **Notifications**: Configure appropriate notifications for critical tasks
5. **Testing**: Test schedules in development before production deployment
6. **Monitoring**: Monitor task execution history and success rates
7. **Maintenance Windows**: Schedule intensive tasks during off-peak hours
8. **Documentation**: Document the purpose and expected behavior of each task
9. **Idempotency**: Design tasks to be idempotent to handle potential reruns safely