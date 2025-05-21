package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestKubiyaAgentCreate(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "testdata/agent",
		EnvVars: map[string]string{
			"KUBIYA_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "agent")
	t.Log(output)

	output = terraform.Destroy(t, terraformOptions)
	t.Log(output)
}
