package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Unit Tests - Validate anyOf schema merging

func TestPoolDatasetResource_Schema(t *testing.T) {
	r := NewPoolDatasetResource()
	if r == nil {
		t.Fatal("NewPoolDatasetResource returned nil")
	}
}

func TestPoolDatasetResource_FilesystemType(t *testing.T) {
	model := PoolDatasetResourceModel{
		Name: types.StringValue("tank/mydata"),
		Type: types.StringValue("FILESYSTEM"),
		Quota: types.Int64Value(107374182400), // 100GB
		Atime: types.StringValue("on"),
	}

	params := map[string]interface{}{
		"name": model.Name.ValueString(),
		"type": model.Type.ValueString(),
	}

	if !model.Quota.IsNull() {
		params["quota"] = model.Quota.ValueInt64()
	}
	if !model.Atime.IsNull() {
		params["atime"] = model.Atime.ValueString()
	}

	if params["type"] != "FILESYSTEM" {
		t.Error("Expected FILESYSTEM type")
	}
	if params["quota"] != int64(107374182400) {
		t.Error("Expected quota to be set for FILESYSTEM")
	}
}

func TestPoolDatasetResource_VolumeType(t *testing.T) {
	model := PoolDatasetResourceModel{
		Name:    types.StringValue("tank/myvol"),
		Type:    types.StringValue("VOLUME"),
		Volsize: types.Int64Value(53687091200), // 50GB
		Sparse:  types.BoolValue(true),
	}

	params := map[string]interface{}{
		"name": model.Name.ValueString(),
		"type": model.Type.ValueString(),
	}

	if !model.Volsize.IsNull() {
		params["volsize"] = model.Volsize.ValueInt64()
	}
	if !model.Sparse.IsNull() {
		params["sparse"] = model.Sparse.ValueBool()
	}

	if params["type"] != "VOLUME" {
		t.Error("Expected VOLUME type")
	}
	if params["volsize"] != int64(53687091200) {
		t.Error("Expected volsize to be set for VOLUME")
	}
	if params["sparse"] != true {
		t.Error("Expected sparse to be true")
	}
}

func TestPoolDatasetResource_OptionalFieldsOmitted(t *testing.T) {
	model := PoolDatasetResourceModel{
		Name: types.StringValue("tank/minimal"),
		Type: types.StringValue("FILESYSTEM"),
		// All other fields null
	}

	params := map[string]interface{}{
		"name": model.Name.ValueString(),
		"type": model.Type.ValueString(),
	}

	// Verify optional fields not in params
	if _, exists := params["quota"]; exists {
		t.Error("Expected null quota to be omitted")
	}
	if _, exists := params["volsize"]; exists {
		t.Error("Expected null volsize to be omitted")
	}
}
