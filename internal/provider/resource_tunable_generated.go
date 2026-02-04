package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TunableResource struct {
	client *client.Client
}

type TunableResourceModel struct {
	ID types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Var types.String `tfsdk:"var"`
	Value types.String `tfsdk:"value"`
	Comment types.String `tfsdk:"comment"`
	Enabled types.Bool `tfsdk:"enabled"`
	UpdateInitramfs types.Bool `tfsdk:"update_initramfs"`
}

func NewTunableResource() resource.Resource {
	return &TunableResource{}
}

func (r *TunableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunable"
}

func (r *TunableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TunableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a tunable.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "* `SYSCTL`: `var` is a sysctl name (e.g. `kernel.watchdog`) and `value` is its corresponding value (",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"var": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name or identifier of the system parameter to tune.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"value": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Value to assign to the tunable parameter.",
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional descriptive comment explaining the purpose of this tunable.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether this tunable is active and should be applied.",
			},
			"update_initramfs": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If `false`, then initramfs will not be updated after creating a ZFS tunable and you will need to run",
			},
		},
	}
}

func (r *TunableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TunableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TunableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Var.IsNull() && !data.Var.IsUnknown() {
		params["var"] = data.Var.ValueString()
	}
	if !data.Value.IsNull() && !data.Value.IsUnknown() {
		params["value"] = data.Value.ValueString()
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.UpdateInitramfs.IsNull() && !data.UpdateInitramfs.IsUnknown() {
		params["update_initramfs"] = data.UpdateInitramfs.ValueBool()
	}

	result, err := r.client.CallWithJob("tunable.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create tunable: %s", err))
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
	result, err = r.client.Call("tunable.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back tunable: %s", err))
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
		if v, ok := resultMap["type"]; ok {
			switch val := v.(type) {
			case string:
				data.Type = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Type = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["value"]; ok {
			switch val := v.(type) {
			case string:
				data.Value = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Value = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Value = types.StringValue(fmt.Sprintf("%v", v))
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TunableResourceModel
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

	result, err := r.client.Call("tunable.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read tunable: %s", err))
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
		if v, ok := resultMap["type"]; ok {
			switch val := v.(type) {
			case string:
				data.Type = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Type = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["value"]; ok {
			switch val := v.(type) {
			case string:
				data.Value = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Value = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Value = types.StringValue(fmt.Sprintf("%v", v))
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TunableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TunableResourceModel
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
	if !data.Value.IsNull() && !data.Value.IsUnknown() {
		params["value"] = data.Value.ValueString()
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.UpdateInitramfs.IsNull() && !data.UpdateInitramfs.IsUnknown() {
		params["update_initramfs"] = data.UpdateInitramfs.ValueBool()
	}

	_, err = r.client.CallWithJob("tunable.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update tunable: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TunableResourceModel
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

	_, err = r.client.CallWithJob("tunable.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete tunable: %s", err))
		return
	}
}
