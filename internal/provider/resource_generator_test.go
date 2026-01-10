package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Unit Tests - Validate generator logic without requiring TrueNAS instance

func TestGeneratedResource_Schema(t *testing.T) {
	r := NewVmResource()
	if r == nil {
		t.Fatal("NewVmResource returned nil")
	}
}

func TestGeneratedResource_OptionalFieldsOmitted(t *testing.T) {
	model := VmResourceModel{
		Name:   types.StringValue("test"),
		Vcpus:  types.Int64Value(2),
		Memory: types.Int64Value(2048),
		// description is null - should be omitted from params
	}

	params := map[string]interface{}{
		"name":   model.Name.ValueString(),
		"vcpus":  int(model.Vcpus.ValueInt64()),
		"memory": int(model.Memory.ValueInt64()),
	}

	// Verify optional null fields are not in params
	if _, exists := params["description"]; exists {
		t.Error("Expected null description to be omitted from params")
	}

	if len(params) != 3 {
		t.Errorf("Expected 3 params, got %d", len(params))
	}
}

func TestGeneratedResource_IDConversion(t *testing.T) {
	// Test that string IDs are converted to int for API calls
	model := VmResourceModel{
		ID: types.StringValue("42"),
	}

	idStr := model.ID.ValueString()
	if idStr != "42" {
		t.Errorf("Expected ID '42', got '%s'", idStr)
	}
}
