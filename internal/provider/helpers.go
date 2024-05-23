package provider

import "fmt"

const (
	readAction   = "read"
	createAction = "create"
	deleteAction = "delete"
	updateAction = "update"
)

func configResourceError(t any) (string, string) {
	const (
		summary = "Unexpected Data Source Configure Type"
		details = "Expected *clients.Client, got: %T. Please report this issue to the provider developers."
	)

	return summary, format(details, t)
}

func format(l string, i ...any) string {
	return fmt.Sprintf(l, i...)
}

func resourceActionError(action, name, err string) (string, string) {
	const (
		summary = "Failed to %s %s resource."
		details = "Could not %s %s data. Error: %s"
	)

	return format(summary, action, name), format(details, action, name, err)
}
