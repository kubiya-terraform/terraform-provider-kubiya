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
  name   = "mevrat_raz_tool"
  runner = "core-testing-1"

  workflows = jsondecode([
    {
      file        = "File"
      name        = "Name"
      group       = "Group"
      description = "Description"
      tags = ["Tags"]
      schedule = [
        {
          expression = "Schedule"
        }
      ]
      stop_schedule = [
        {
          expression = "StopSchedule"
        }
      ]
      restart_schedule = [
        {
          expression = "RestartSchedule"
        }
      ]
      skip_if_successful = false
      env = ["Env"]
      dotenv = ["Dotenv"]
      default_params     = "DefaultParams"
      params = [
        {
          key = "Params"
        }
      ]
      secrets = ["Secrets"]
      entrypoint        = "Entrypoint"
      state             = "State"
      log_dir           = "LogDir"
      timeout           = 0
      delay             = 0
      restart_wait      = 0
      max_active_runs   = 0
      max_clean_up_time = 0
      steps = [
        {
          name        = "Name"
          description = "Description"
          icon        = "Icon"
          label       = "Label"
          next_steps = ["step-2"]
          conditions = [
            {
              condition = "2025-05-04T12:32:20.161864+03:00"
            }
          ]
          command = "Command"
          dir     = "Dir"
          env = ["env1=value1", "env2=value2"]
          args = {
            arg1 = "value1"
          }
          depends_on = ["DependsOn-1"]
          timeout      = 0
          retry        = 0
          retry_delay  = 0
          ignore_error = false
        }
      ]
      preconditions = [
        {
          type = "Preconditions"
        }
      ]
      handler_on = {
        failure = {
          name        = "Name"
          description = "Description"
          icon        = "Icon"
          label       = "Label"
          next_steps = ["step-2"]
          conditions = [
            {
              condition = "2025-05-04T12:32:20.161865+03:00"
            }
          ]
          command = "Command"
          dir     = "Dir"
          env = ["env1=value1", "env2=value2"]
          args = {
            arg1 = "value1"
          }
          depends_on = ["DependsOn-1"]
          timeout      = 0
          retry        = 0
          retry_delay  = 0
          ignore_error = false
        }
        success = {
          name        = "Name"
          description = "Description"
          icon        = "Icon"
          label       = "Label"
          next_steps = ["step-2"]
          conditions = [
            {
              condition = "2025-05-04T12:32:20.161865+03:00"
            }
          ]
          command = "Command"
          dir     = "Dir"
          env = ["env1=value1", "env2=value2"]
          args = {
            arg1 = "value1"
          }
          depends_on = ["DependsOn-1"]
          timeout      = 0
          retry        = 0
          retry_delay  = 0
          ignore_error = false
        }
        cancel = {
          name        = "Name"
          description = "Description"
          icon        = "Icon"
          label       = "Label"
          next_steps = ["step-2"]
          conditions = [
            {
              condition = "2025-05-04T12:32:20.161865+03:00"
            }
          ]
          command = "Command"
          dir     = "Dir"
          env = ["env1=value1", "env2=value2"]
          args = {
            arg1 = "value1"
          }
          depends_on = ["DependsOn-1"]
          timeout      = 0
          retry        = 0
          retry_delay  = 0
          ignore_error = false
        }
        exit = {
          name        = "Name"
          description = "Description"
          icon        = "Icon"
          label       = "Label"
          next_steps = ["step-2"]
          conditions = [
            {
              condition = "2025-05-04T12:32:20.161866+03:00"
            }
          ]
          command = "Command"
          dir     = "Dir"
          env = ["env1=value1", "env2=value2"]
          args = {
            arg1 = "value1"
          }
          depends_on = ["DependsOn-1"]
          timeout      = 0
          retry        = 0
          retry_delay  = 0
          ignore_error = false
        }
      }
      smtp = {
        host     = "Host"
        port     = "Port"
        username = "Username"
        password = "Password"
      }
      error_mail = {
        from        = "From"
        to          = "To"
        prefix      = "Prefix"
        attach_logs = false
      }
      info_mail = {
        from        = "From"
        to          = "To"
        prefix      = "Prefix"
        attach_logs = false
      }
      mail_on = {
        failure = true
        success = false
      }
      hist_retention_days = 0
    }
  ])
}

