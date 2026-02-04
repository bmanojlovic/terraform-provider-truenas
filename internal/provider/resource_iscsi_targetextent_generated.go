package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IscsiTargetextentResource struct {
	client *client.Client
}

type IscsiTargetextentResourceModel struct {
	ID types.String `tfsdk:"id"`
	Target types.Int64 `tfsdk:"target"`
	Lunid types.Int64 `tfsdk:"lunid"`
	Extent types.Int64 `tfsdk:"extent"`
}

func NewIscsiTargetextentResource() resource.Resource {
	return &IscsiTargetextentResource{}
}

func (r *IscsiTargetextentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_targetextent"
}

func (r *IscsiTargetextentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiTargetextentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an Associated Target.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"target": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the iSCSI target to associate with the extent.",
			},
			"lunid": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Logical Unit Number (LUN) ID for presenting the extent to the target.",
			},
			"extent": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the iSCSI extent to associate with the target.",
			},
		},
	}
}

func (r *IscsiTargetextentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
		return
	}
	r.client = client
}

func (r *IscsiTargetextentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Target.IsNull() && !data.Target.IsUnknown() {
		params["target"] = data.Target.ValueInt64()
	}
	if !data.Lunid.IsNull() && !data.Lunid.IsUnknown() {
		params["lunid"] = data.Lunid.ValueInt64()
	}
	if !data.Extent.IsNull() && !data.Extent.IsUnknown() {
		params["extent"] = data.Extent.ValueInt64()
	}

	result, err := r.client.Call("iscsi.targetextent.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_targetextent: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}


	// Read back to populate computed fields
	var id interface{}
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("iscsi.targetextent.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back iscsi_targetextent: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["id"]; ok && v != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["target"]; ok {
			switch val := v.(type) {
			case float64:
				data.Target = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Target = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["lunid"]; ok {
			if v == nil {
				data.Lunid = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.Lunid = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.Lunid = types.Int64Value(int64(fv)) }
					}
				}
			}
		}
		if v, ok := resultMap["extent"]; ok {
			switch val := v.(type) {
			case float64:
				data.Extent = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Extent = types.Int64Value(int64(fv)) }
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetextentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("iscsi.targetextent.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_targetextent: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["id"]; ok && v != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["target"]; ok {
			switch val := v.(type) {
			case float64:
				data.Target = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Target = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["lunid"]; ok {
			if v == nil {
				data.Lunid = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.Lunid = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.Lunid = types.Int64Value(int64(fv)) }
					}
				}
			}
		}
		if v, ok := resultMap["extent"]; ok {
			switch val := v.(type) {
			case float64:
				data.Extent = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.Extent = types.Int64Value(int64(fv)) }
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetextentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Target.IsNull() && !data.Target.IsUnknown() {
		params["target"] = data.Target.ValueInt64()
	}
	if !data.Lunid.IsNull() && !data.Lunid.IsUnknown() {
		params["lunid"] = data.Lunid.ValueInt64()
	}
	if !data.Extent.IsNull() && !data.Extent.IsUnknown() {
		params["extent"] = data.Extent.ValueInt64()
	}

	_, err = r.client.Call("iscsi.targetextent.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_targetextent: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetextentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiTargetextentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	id = []interface{}{id, map[string]interface{}{}}

	_, err = r.client.Call("iscsi.targetextent.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_targetextent: %s", err))
		return
	}
}
