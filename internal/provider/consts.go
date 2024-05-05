package provider

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

func configureProviderError(err error) (string, string) {
	const (
		summary = "Unable to create api client"
		details = "An unexpected error occurred when creating the kubiya api client. If the error is not clear, please contact the provider developers. Kubiya AgentsClient Error: %s"
	)

	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return summary, format(details, msg)
}

func resourceActionError(action, name, err string) (string, string) {
	const (
		summary = "Failed to %s %s resource."
		details = "Could not %s %s data. Error: %s"
	)

	return format(summary, action, name), format(details, action, name, err)
}
