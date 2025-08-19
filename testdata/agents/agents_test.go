package agents

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAgentExample(t *testing.T) {
	options := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data/agents",
	})

	defer terraform.Destroy(t, options)

	initResult := terraform.Init(t, options)
	planResult := terraform.Plan(t, options)
	terraform.Get(t, options)

	show := terraform.InitAndPlan(t, options)
	fmt.Println(show)

	output := terraform.Output(t, options, "agent")
	assert.Equal(t, "default-agent", output)
}
