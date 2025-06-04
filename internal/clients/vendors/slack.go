package vendors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"terraform-provider-kubiya/internal/entities"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SlackVendor implements the VendorClient interface for Slack
type SlackVendor struct {
	name string
}

// NewSlackVendor creates a new Slack vendor implementation
func NewSlackVendor() VendorClient {
	return &SlackVendor{
		name: "slack",
	}
}

// Slack-specific structures
type slackIntegration struct {
	BaseExternalKnowledge
	ChannelIDs []string `json:"channel_ids"`
}

type slackIntegrationRequest struct {
	ChannelIDs []string `json:"channel_ids"`
}

type slackIntegrationResponse struct {
	Org        string   `json:"org"`
	UserEmail  string   `json:"user_email"`
	UUID       string   `json:"uuid"`
	ChannelIDs []string `json:"channel_ids"`
	StartDate  string   `json:"start_date"`
	Message    string   `json:"message"`
}

// GetVendorName returns the vendor identifier
func (s *SlackVendor) GetVendorName() string {
	return s.name
}

// ValidateConfig validates the Slack configuration
func (s *SlackVendor) ValidateConfig(config types.Map) error {
	// Check if channel_ids exists and is not empty
	if v, ok := config.Elements()["channel_ids"]; ok {
		// Handle different value types
		switch val := v.(type) {
		case types.Dynamic:
			extracted := ExtractDynamicValue(val)
			if channelIDs, ok := extracted.([]string); ok && len(channelIDs) > 0 {
				return nil
			}
		case types.List:
			if !val.IsNull() && !val.IsUnknown() && len(val.Elements()) > 0 {
				return nil
			}
		case types.Tuple:
			if !val.IsNull() && !val.IsUnknown() && len(val.Elements()) > 0 {
				return nil
			}
		}
	}
	return fmt.Errorf("channel_ids is required for Slack integration and must be a non-empty list")
}

// PrepareCreateRequest converts the Terraform config to Slack create request
func (s *SlackVendor) PrepareCreateRequest(config types.Map) (interface{}, error) {
	if err := s.ValidateConfig(config); err != nil {
		return nil, err
	}

	channelIDs := s.extractChannelIDs(config)
	return &slackIntegrationRequest{
		ChannelIDs: channelIDs,
	}, nil
}

// PrepareUpdateRequest converts the Terraform config to Slack update request
func (s *SlackVendor) PrepareUpdateRequest(config types.Map) (interface{}, error) {
	// Same as create for Slack
	return s.PrepareCreateRequest(config)
}

// ParseCreateResponse parses the Slack create response
func (s *SlackVendor) ParseCreateResponse(resp io.Reader) (*entities.ExternalKnowledgeModel, error) {
	var r slackIntegrationResponse
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		return nil, err
	}

	base := BaseExternalKnowledge{
		UUID:            r.UUID,
		Org:             r.Org,
		StartDate:       r.StartDate,
		IntegrationType: s.name,
	}

	configElements := map[string]attr.Value{
		"channel_ids": ConvertToTerraformValue(r.ChannelIDs),
	}

	return CreateExternalKnowledgeModel(base, s.name, configElements), nil
}

// ParseReadResponse parses the Slack read response
func (s *SlackVendor) ParseReadResponse(resp io.Reader) (*entities.ExternalKnowledgeModel, error) {
	var r slackIntegration
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		return nil, err
	}

	configElements := map[string]attr.Value{
		"channel_ids": ConvertToTerraformValue(r.ChannelIDs),
	}

	return CreateExternalKnowledgeModel(r.BaseExternalKnowledge, s.name, configElements), nil
}

// ParseUpdateResponse parses the Slack update response
func (s *SlackVendor) ParseUpdateResponse(resp io.Reader, currentModel *entities.ExternalKnowledgeModel) error {
	var r slackIntegrationResponse
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		return err
	}

	// Update the model with the response
	currentModel.Id = types.StringValue(r.UUID)
	currentModel.Org = types.StringValue(r.Org)
	currentModel.StartDate = types.StringValue(r.StartDate)

	// Update channel_ids in config
	configElements := map[string]attr.Value{
		"channel_ids": ConvertToTerraformValue(r.ChannelIDs),
	}

	// Create object type and value
	attrTypes := make(map[string]attr.Type)
	for key, val := range configElements {
		attrTypes[key] = val.Type(context.Background())
	}
	configObj, _ := types.ObjectValue(attrTypes, configElements)
	currentModel.Config = types.DynamicValue(configObj)

	return nil
}

// ParseListResponse parses the Slack list response
func (s *SlackVendor) ParseListResponse(resp io.Reader) ([]*entities.ExternalKnowledgeModel, error) {
	var slackList []*slackIntegration
	if err := json.NewDecoder(resp).Decode(&slackList); err != nil {
		return nil, err
	}

	result := make([]*entities.ExternalKnowledgeModel, 0, len(slackList))
	for _, item := range slackList {
		configElements := map[string]attr.Value{
			"channel_ids": ConvertToTerraformValue(item.ChannelIDs),
		}
		model := CreateExternalKnowledgeModel(item.BaseExternalKnowledge, s.name, configElements)
		result = append(result, model)
	}

	return result, nil
}

// extractChannelIDs extracts channel IDs from the config
func (s *SlackVendor) extractChannelIDs(config types.Map) []string {
	channelIDs := []string{}

	if v, ok := config.Elements()["channel_ids"]; ok {
		// Handle different value types
		switch val := v.(type) {
		case types.Dynamic:
			// Already a dynamic value, extract it
			if extracted := ExtractDynamicValue(val); extracted != nil {
				if ids, ok := extracted.([]string); ok {
					channelIDs = ids
				}
			}
		case types.List:
			// Direct list value
			for _, elem := range val.Elements() {
				if strVal, ok := elem.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
					channelIDs = append(channelIDs, strVal.ValueString())
				}
			}
		case types.Tuple:
			// Tuple value (common when defined in HCL)
			for _, elem := range val.Elements() {
				if strVal, ok := elem.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
					channelIDs = append(channelIDs, strVal.ValueString())
				}
			}
		default:
			// Try to handle other cases
			fmt.Printf("DEBUG: Unexpected type for channel_ids: %T\n", val)
		}
	}

	return channelIDs
}
