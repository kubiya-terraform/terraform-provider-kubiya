terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
      # source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_inline_source" "hello_world_tool" {
  name   = "integration_test_tool"
  runner = ""

  # tools = jsonencode([])

  workflows = jsonencode([
    {
      name        = "maayanbauer2"
      description = "Example DAG demonstrating data flow between tool steps"
      steps = [
        {
          name        = "generate-data"
          description = "First tool generates data for the second step 1"
          executor = {
            type = "tool"
            config = {
              tool_def = {
                name        = "data-generator"
                description = "Generates sample data for the next step"
                type        = "docker"
                image       = "python:3.12-slim-bullseye"
                with_files = [
                  {
                    destination = "/tmp/ascript.py"
                    content     = "#!/usr/bin/env python3\nimport json\nimport random\n\n# Generate some random data\ndata = {\n    \"id\": random.randint(1000, 9999),\n    \"values\": [random.randint(1, 100) for _ in range(5)],\n    \"name\": f\"Sample-{random.choice(['A','B', 'C'])}\"\n}\n\n# Output the data as JSON\nprint(json.dumps(data))"
                  }
                ]
                content = "set -e\npython /tmp/ascript.py"
              }
            }
          }
          output = "GENERATED_DATA"
        },
        {
          name        = "process-data"
          description = "Second tool processes data from first tool"
          depends = [
            "generate-data"
          ]
          executor = {
            type = "tool"
            config = {
              tool_def = {
                name        = "data-processor"
                description = "Processes data from previous step"
                type        = "docker"
                image       = "python:3.12-slim-bullseye"
                with_files = [
                  {
                    destination = "/tmp/ascript.py"
                    content     = "#!/usr/bin/env python3\nimport os\nimport json\n\n# Get the data from the previous step\ninput_data = os.environ.get('data')\n\ntry:\n    # Parse the JSON data\ndata = json.loads(input_data)\n\n    # Process the data\ntotal = sum(data.get('values', [0]))\navg = total / len(data.get('values', [1]))\n\n    # Output the results as a single compact line to avoid truncation issues\nresult = {\n    \"source_id\": data.get('id'),\n    \"source_name\": data.get('name'),\n    \"processed\": {\n        \"total\": total,\n        \"average\": avg,\n        \"count\": len(data.get('values', []))\n    }\n}\n\n# Output as a single line with no formatting to avoid truncation\nprint(f\"RESULT:{json.dumps(result)}\")\n\nexcept json.JSONDecodeError as e:\n    print(f\"Error parsing data: {input_data}\")\n    print(f\"Error details: {str(e)}\")\n    raise"
                  }
                ]
                content = "set -e\npython /tmp/ascript.py"
                args = [
                  {
                    name        = "data"
                    type        = "string"
                    description = "JSON data from previous step"
                    required    = true
                  }
                ]
              }
              args = {
                data = "$${GENERATED_DATA}"
              }
            }
          }
          output = "PROCESSED_DATA"
        },
        {
          name = "send-to-slack"
          executor = {
            type = "agent"
            config = {
              teammate_name = "demo_teammate"
              message       = "Send a Slack msg to channel #tf-test saying $${PROCESSED_DATA}. and run it don't ask for questions"
            }
          }
          output = "SLACK_RESPONSE"
          depends = [
            "process-data",
            "generate-data"
          ]
        }
      ]
    }
  ])
}

resource "kubiya_agent" "helloworld_teammate" {
  name          = "integration_test_teammate"
  description   = "HelloWorldTeammate is designed to interact with users and execute a simple 'Hello World' tool. This teammate requires minimal permissions and serves as a basic example of Kubiya's capabilities."
  image         = "ghcr.io/kubiyabot/kubiya-agent:stable"
  model         = "azure/gpt-4o"
  runner        = "core-testing-1"
  is_debug_mode = false

  # environment_variables = {}
  # integrations          = []
  # links                 = []
  # groups                = []
  # secrets               = []
  # sources               = []
  # tool_sources          = []
  # users                 = []
  instructions = ""
}

output "output" {
  value = kubiya_inline_source.hello_world_tool
}