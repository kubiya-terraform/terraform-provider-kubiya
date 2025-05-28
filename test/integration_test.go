package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestKubiyaAgent(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/agents",
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

func TestKubiyaKnowledge(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/knowledge",
		EnvVars: map[string]string{
			"KUBIYA_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "knowledge")
	t.Log(output)

	output = terraform.Destroy(t, terraformOptions)
	t.Log(output)
}

func TestKubiyaIntegrations(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/integrations",
		EnvVars: map[string]string{
			"KUBIYA_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "integrations")
	t.Log(output)

	output = terraform.Destroy(t, terraformOptions)
	t.Log(output)
}

func TestKubiyaSources(t *testing.T) {
	apiKey := os.Getenv("KUBIYA_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_API_KEY environment variable is not set")
	}

	t.Run("git source", func(t *testing.T) {
		t.Parallel()

		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/sources/git_source",
			EnvVars: map[string]string{
				"KUBIYA_API_KEY": apiKey,
			},
		}

		terraform.InitAndApply(t, terraformOptions)

		output := terraform.Output(t, terraformOptions, "output")
		t.Log(output)

		output = terraform.Destroy(t, terraformOptions)
		t.Log(output)
	})

	t.Run("inline source", func(t *testing.T) {
		t.Parallel()

		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/sources/inline_source",
			EnvVars: map[string]string{
				"KUBIYA_API_KEY": apiKey,
			},
		}

		terraform.InitAndApply(t, terraformOptions)

		output := terraform.Output(t, terraformOptions, "inline_source")
		t.Log(output)

		output = terraform.Destroy(t, terraformOptions)
		t.Log(output)
	})
}
