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

type NvmetHostResource struct {
	client *client.Client
}

type NvmetHostResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Hostnqn       types.String `tfsdk:"hostnqn"`
	DhchapKey     types.String `tfsdk:"dhchap_key"`
	DhchapCtrlKey types.String `tfsdk:"dhchap_ctrl_key"`
	DhchapDhgroup types.String `tfsdk:"dhchap_dhgroup"`
	DhchapHash    types.String `tfsdk:"dhchap_hash"`
}

func NewNvmetHostResource() resource.Resource {
	return &NvmetHostResource{}
}

func (r *NvmetHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_host"
}

func (r *NvmetHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NvmetHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an NVMe target `host`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"hostnqn": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "NQN of the host that will connect to this TrueNAS. ",
			},
			"dhchap_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If set, the secret that the host must present when connecting.  A suitable secret can be generated u",
			},
			"dhchap_ctrl_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If set, the secret that this TrueNAS will present to the host when the host is connecting (Bi-Direct",
			},
			"dhchap_dhgroup": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If selected, the DH (Diffie-Hellman) key exchange built on top of CHAP to be used for authentication",
			},
			"dhchap_hash": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "HMAC (Hashed Message Authentication Code) to be used in conjunction if a `dhchap_dhgroup` is selecte",
			},
		},
	}
}

func (r *NvmetHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NvmetHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Hostnqn.IsNull() && !data.Hostnqn.IsUnknown() {
		params["hostnqn"] = data.Hostnqn.ValueString()
	}
	if !data.DhchapKey.IsNull() && !data.DhchapKey.IsUnknown() {
		params["dhchap_key"] = data.DhchapKey.ValueString()
	}
	if !data.DhchapCtrlKey.IsNull() && !data.DhchapCtrlKey.IsUnknown() {
		params["dhchap_ctrl_key"] = data.DhchapCtrlKey.ValueString()
	}
	if !data.DhchapDhgroup.IsNull() && !data.DhchapDhgroup.IsUnknown() {
		params["dhchap_dhgroup"] = data.DhchapDhgroup.ValueString()
	}
	if !data.DhchapHash.IsNull() && !data.DhchapHash.IsUnknown() {
		params["dhchap_hash"] = data.DhchapHash.ValueString()
	}

	result, err := r.client.Call("nvmet.host.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create nvmet_host: %s", err))
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
	result, err = r.client.Call("nvmet.host.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back nvmet_host: %s", err))
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
	if v, ok := resultMap["hostnqn"]; ok {
		switch val := v.(type) {
		case string:
			data.Hostnqn = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Hostnqn = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Hostnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["dhchap_key"]; ok {
		if v == nil {
			data.DhchapKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["dhchap_ctrl_key"]; ok {
		if v == nil {
			data.DhchapCtrlKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapCtrlKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["dhchap_dhgroup"]; ok {
		if v == nil {
			data.DhchapDhgroup = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapDhgroup = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.DhchapKey.IsUnknown() {
		data.DhchapKey = types.StringNull()
	}
	if data.DhchapCtrlKey.IsUnknown() {
		data.DhchapCtrlKey = types.StringNull()
	}
	if data.DhchapDhgroup.IsUnknown() {
		data.DhchapDhgroup = types.StringNull()
	}
	if data.DhchapHash.IsUnknown() {
		data.DhchapHash = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetHostResourceModel
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

	result, err := r.client.Call("nvmet.host.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read nvmet_host: %s", err))
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
	if v, ok := resultMap["hostnqn"]; ok {
		switch val := v.(type) {
		case string:
			data.Hostnqn = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Hostnqn = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Hostnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["dhchap_key"]; ok {
		if v == nil {
			data.DhchapKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["dhchap_ctrl_key"]; ok {
		if v == nil {
			data.DhchapCtrlKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapCtrlKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["dhchap_dhgroup"]; ok {
		if v == nil {
			data.DhchapDhgroup = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapDhgroup = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.DhchapKey.IsUnknown() {
		data.DhchapKey = types.StringNull()
	}
	if data.DhchapCtrlKey.IsUnknown() {
		data.DhchapCtrlKey = types.StringNull()
	}
	if data.DhchapDhgroup.IsUnknown() {
		data.DhchapDhgroup = types.StringNull()
	}
	if data.DhchapHash.IsUnknown() {
		data.DhchapHash = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetHostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NvmetHostResourceModel
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
	if !data.Hostnqn.IsNull() && !data.Hostnqn.IsUnknown() {
		params["hostnqn"] = data.Hostnqn.ValueString()
	}
	if !data.DhchapKey.IsNull() && !data.DhchapKey.IsUnknown() {
		params["dhchap_key"] = data.DhchapKey.ValueString()
	}
	if !data.DhchapCtrlKey.IsNull() && !data.DhchapCtrlKey.IsUnknown() {
		params["dhchap_ctrl_key"] = data.DhchapCtrlKey.ValueString()
	}
	if !data.DhchapDhgroup.IsNull() && !data.DhchapDhgroup.IsUnknown() {
		params["dhchap_dhgroup"] = data.DhchapDhgroup.ValueString()
	}
	if !data.DhchapHash.IsNull() && !data.DhchapHash.IsUnknown() {
		params["dhchap_hash"] = data.DhchapHash.ValueString()
	}

	_, err = r.client.Call("nvmet.host.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update nvmet_host: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("nvmet.host.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back nvmet_host: %s", readErr))
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
	if v, ok := resultMap["hostnqn"]; ok {
		switch val := v.(type) {
		case string:
			data.Hostnqn = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Hostnqn = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Hostnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["dhchap_key"]; ok {
		if v == nil {
			data.DhchapKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["dhchap_ctrl_key"]; ok {
		if v == nil {
			data.DhchapCtrlKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapCtrlKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["dhchap_dhgroup"]; ok {
		if v == nil {
			data.DhchapDhgroup = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.DhchapDhgroup = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.DhchapKey.IsUnknown() {
		data.DhchapKey = types.StringNull()
	}
	if data.DhchapCtrlKey.IsUnknown() {
		data.DhchapCtrlKey = types.StringNull()
	}
	if data.DhchapDhgroup.IsUnknown() {
		data.DhchapDhgroup = types.StringNull()
	}
	if data.DhchapHash.IsUnknown() {
		data.DhchapHash = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetHostResourceModel
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

	_, err = r.client.Call("nvmet.host.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete nvmet_host: %s", err))
		return
	}
}
