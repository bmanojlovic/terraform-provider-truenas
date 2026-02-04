package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type FcportResource struct {
	client *client.Client
}

type FcportResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Port     types.String `tfsdk:"port"`
	TargetId types.Int64  `tfsdk:"target_id"`
}

func NewFcportResource() resource.Resource {
	return &FcportResource{}
}

func (r *FcportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fcport"
}

func (r *FcportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FcportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates mapping between a FC port and a target.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"port": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Alias name for the Fibre Channel port.",
			},
			"target_id": schema.Int64Attribute{
				Required:    true,
				Optional:    false,
				Description: "ID of the target to associate with this FC port.",
			},
		},
	}
}

func (r *FcportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FcportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FcportResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		params["port"] = data.Port.ValueString()
	}
	if !data.TargetId.IsNull() && !data.TargetId.IsUnknown() {
		params["target_id"] = data.TargetId.ValueInt64()
	}

	result, err := r.client.Call("fcport.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create fcport: %s", err))
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
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("fcport.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back fcport: %s", err))
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
	if v, ok := resultMap["port"]; ok {
		switch val := v.(type) {
		case string:
			data.Port = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Port = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Port = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["target_id"]; ok {
		switch val := v.(type) {
		case float64:
			data.TargetId = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.TargetId = types.Int64Value(int64(fv))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FcportResourceModel
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

	result, err := r.client.Call("fcport.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read fcport: %s", err))
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
	if v, ok := resultMap["port"]; ok {
		switch val := v.(type) {
		case string:
			data.Port = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Port = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Port = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["target_id"]; ok {
		switch val := v.(type) {
		case float64:
			data.TargetId = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.TargetId = types.Int64Value(int64(fv))
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FcportResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FcportResourceModel
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
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		params["port"] = data.Port.ValueString()
	}
	if !data.TargetId.IsNull() && !data.TargetId.IsUnknown() {
		params["target_id"] = data.TargetId.ValueInt64()
	}

	_, err = r.client.Call("fcport.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update fcport: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FcportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FcportResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("fcport.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete fcport: %s", err))
		return
	}
}
