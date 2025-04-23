terraform {
  required_providers {
    kubiya = {
      #       source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_scheduled_task" "example" {
  channel_id  = "C08041WFAKT"
  agent       = "mevrat_agent"
  repeat      = "*/15 * * * *"
  description = "mevrat mevrat test task description"
}

output "output" {
  value = kubiya_scheduled_task.example
}

/*

	•	Persistent Storage
Store flags in a durable database (PostgreSQL, MongoDB, etc.) instead of in-memory, so flags survive restarts and can be queried by history or audit.
	•	Web Dashboard & CLI
Provide a simple web UI (and/or CLI) for non-developers to create, edit, toggle, and roll back flags, view flag statuses, and inspect evaluation logs.
	•	Audit Logging & Rollbacks
Record every flag change (who, when, from → to) in an audit log. Expose an API to list past changes and roll back to a previous flag state.
	•	Environment & App Scoping
Support multiple environments (e.g. “staging”, “production”) and allow flags to be scoped per environment and per application.
	•	Percentage Rollouts & Targeting
Allow partial rollouts by defining a percentage of users (via hashing on user ID) or custom user-property rules (country, plan tier, custom attributes).
	•	Real-Time Updates
Push flag changes to SDKs via SSE or WebSocket so clients pick up toggles immediately without polling.
	•	Fallback Defaults & Caching
In the SDK, cache the last fetched flag values locally with TTL and provide a default value if the server is unreachable.
	•	Role-Based Access Control
Restrict who can read vs. modify flags via API tokens with scoped permissions (read, write, audit).
	•	Feature-Driven Metrics & Dashboards
Correlate flag usage events with business metrics—e.g. A/B test results, error rates, performance impact—so you can measure flag impact.
	•	SDK Multi-Language Support
Provide client libraries in other languages (Python, JavaScript, Java, .NET) with the same evaluation and reporting capabilities.
	•	Flag Dependencies & Groups
Define flag hierarchies or make one flag depend on another, and group related flags for bulk operations.
	•	Scheduled Flags
Add the ability to automatically enable/disable flags at specified times or based on cron schedules.
	•	Health & Status Endpoint
Expose a /health endpoint on the server (and in the SDK) to verify service readiness and connectivity.
	•	High Availability & Scaling
Run the server clustered behind a load balancer, with shared storage and a message bus for cache invalidation.
	•	Documentation & Tutorials
Ship clear docs and code samples for setting up the server, integrating the SDK, and using advanced features like targeting and rollbacks.
 */