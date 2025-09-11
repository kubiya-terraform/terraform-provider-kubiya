# Resource Instrumentation Template

## How to Add Sentry Error Tracking to Resources

### 1. Add Imports (Already Done by Script)

```go
import (
// ... existing imports
"github.com/getsentry/sentry-go"
kubiyasentry "terraform-provider-kubiya/internal/sentry"
)
```

### 2. Instrument Create Operation

```go
func (r *resourceName) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// Start tracing for the create operation
ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_RESOURCE", "", kubiyasentry.OpResourceCreate)
defer kubiyasentry.FinishSpan(span)

// Add breadcrumb
kubiyasentry.AddBreadcrumb("resource", "Creating RESOURCE resource", sentry.LevelInfo, nil)

var plan entities.RESOURCEModel

diags := req.Plan.Get(ctx, &plan)
resp.Diagnostics.Append(diags...)

if resp.Diagnostics.HasError() {
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
return
}

// Add plan data to span (optional)
if !plan.Name.IsNull() {
kubiyasentry.SetSpanData(span, "RESOURCE.name", plan.Name.ValueString())
}

state, err := r.client.CreateRESOURCE(ctx, &plan)
if err != nil {
// Record error in span and capture to Sentry
kubiyasentry.RecordError(ctx, err)
CaptureResourceError(ctx, "kubiya_RESOURCE", "", "create", err)

resp.Diagnostics.AddError(
resourceActionError(createAction, r.name, err.Error()),
)
return
}

// Update span with created resource ID
if state != nil && !state.Id.IsNull() {
id := state.Id.ValueString()
kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
kubiyasentry.SetSpanData(span, "RESOURCE.created_id", id)
}

// Success
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
```

### 3. Instrument Read Operation

```go
func (r *resourceName) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// Start tracing for the read operation
ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_RESOURCE", "", kubiyasentry.OpResourceRead)
defer kubiyasentry.FinishSpan(span)

// Add breadcrumb
kubiyasentry.AddBreadcrumb("resource", "Reading RESOURCE resource", sentry.LevelInfo, nil)

var state entities.RESOURCEModel

diags := req.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)

if resp.Diagnostics.HasError() {
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
return
}

// Update span with resource ID if available
if !state.Id.IsNull() {
id := state.Id.ValueString()
kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
kubiyasentry.SetSpanData(span, "RESOURCE.id", id)
}

if err := r.client.ReadRESOURCE(ctx, &state); err != nil {
// Record error in span and capture to Sentry
kubiyasentry.RecordError(ctx, err)
CaptureResourceError(ctx, "kubiya_RESOURCE", state.Id.ValueString(), "read", err)

resp.Diagnostics.AddError(
resourceActionError(readAction, r.name, err.Error()),
)
return
}

// Success
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

### 4. Instrument Update Operation

```go
func (r *resourceName) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// Start tracing for the update operation
ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_RESOURCE", "", kubiyasentry.OpResourceUpdate)
defer kubiyasentry.FinishSpan(span)

// Add breadcrumb
kubiyasentry.AddBreadcrumb("resource", "Updating RESOURCE resource", sentry.LevelInfo, nil)

var plan, state entities.RESOURCEModel

resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

if resp.Diagnostics.HasError() {
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
return
}

id := state.Id.ValueString()
kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
kubiyasentry.SetSpanData(span, "RESOURCE.id", id)

// Update logic here...

if err := r.client.UpdateRESOURCE(ctx, &state); err != nil {
// Record error in span and capture to Sentry
kubiyasentry.RecordError(ctx, err)
CaptureResourceError(ctx, "kubiya_RESOURCE", id, "update", err)

resp.Diagnostics.AddError(
resourceActionError(updateAction, r.name, err.Error()),
)
return
}

// Success
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

### 5. Instrument Delete Operation

```go
func (r *resourceName) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// Start tracing for the delete operation
ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_RESOURCE", "", kubiyasentry.OpResourceDelete)
defer kubiyasentry.FinishSpan(span)

// Add breadcrumb
kubiyasentry.AddBreadcrumb("resource", "Deleting RESOURCE resource", sentry.LevelInfo, nil)

var state entities.RESOURCEModel

diags := req.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)

if resp.Diagnostics.HasError() {
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
return
}

// Update span with resource ID
id := state.Id.ValueString()
kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
kubiyasentry.SetSpanData(span, "RESOURCE.id", id)

if err := r.client.DeleteRESOURCE(ctx, &state); err != nil {
// Record error in span and capture to Sentry
kubiyasentry.RecordError(ctx, err)
CaptureResourceError(ctx, "kubiya_RESOURCE", id, "delete", err)

resp.Diagnostics.AddError(
resourceActionError(deleteAction, r.name, err.Error()),
)
return
}

// Success
kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
kubiyasentry.AddBreadcrumb("resource", "Successfully deleted RESOURCE "+id, sentry.LevelInfo, nil)
}
```

## Replace Placeholders

For each resource, replace:

- `RESOURCE` → actual resource name (e.g., `webhook`, `trigger`)
- `resourceName` → struct name (e.g., `webhookResource`, `triggerResource`)
- `RESOURCEModel` → entity model (e.g., `WebhookModel`, `TriggerModel`)
- `CreateRESOURCE`, `ReadRESOURCE`, etc. → actual client methods

## Resource Mapping

| Resource File                    | Resource Type               | Entity Model             |
|----------------------------------|-----------------------------|--------------------------|
| `webhook_resource.go`            | `kubiya_webhook`            | `WebhookModel`           |
| `trigger_resource.go`            | `kubiya_trigger`            | `TriggerModel`           |
| `shceduled_task_resource.go`     | `kubiya_scheduled_task`     | `ScheduledTaskModel`     |
| `knowledge_resource.go`          | `kubiya_knowledge`          | `KnowledgeModel`         |
| `external_knowledge_resource.go` | `kubiya_external_knowledge` | `ExternalKnowledgeModel` |
| `integration_resource.go`        | `kubiya_integration`        | `IntegrationModel`       |
| `source_resource.go`             | `kubiya_source`             | `SourceModel`            |
| `inline_source_resource.go`      | `kubiya_inline_source`      | `InlineSourceModel`      |
| `secret.go`                      | `kubiya_secret`             | `SecretModel`            |

## Quick Implementation Steps

1. ✅ Imports added by script
2. For each resource file:
    - Find each CRUD function (Create, Read, Update, Delete)
    - Add the tracing wrapper at the beginning
    - Add error handling with Sentry capture
    - Set appropriate span status
3. Test build: `go build`
4. Verify instrumentation works