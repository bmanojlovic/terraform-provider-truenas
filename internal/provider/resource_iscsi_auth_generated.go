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

type IscsiAuthResource struct {
	client *client.Client
}

type IscsiAuthResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Tag           types.Int64  `tfsdk:"tag"`
	User          types.String `tfsdk:"user"`
	Secret        types.String `tfsdk:"secret"`
	Peeruser      types.String `tfsdk:"peeruser"`
	Peersecret    types.String `tfsdk:"peersecret"`
	DiscoveryAuth types.String `tfsdk:"discovery_auth"`
}

func NewIscsiAuthResource() resource.Resource {
	return &IscsiAuthResource{}
}

func (r *IscsiAuthResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_auth"
}

func (r *IscsiAuthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiAuthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Authorized Access.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"tag": schema.Int64Attribute{
				Required:    true,
				Optional:    false,
				Description: "Numeric tag used to associate this credential with iSCSI targets.",
			},
			"user": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Username for iSCSI CHAP authentication.",
			},
			"secret": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Password/secret for iSCSI CHAP authentication.",
			},
			"peeruser": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Username for mutual CHAP authentication or empty string if not configured.",
			},
			"peersecret": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Password/secret for mutual CHAP authentication or empty string if not configured.",
			},
			"discovery_auth": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Authentication method for target discovery. If \"CHAP_MUTUAL\" is selected for target discovery, it is",
			},
		},
	}
}

func (r *IscsiAuthResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Tag.IsNull() && !data.Tag.IsUnknown() {
		params["tag"] = data.Tag.ValueInt64()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Secret.IsNull() && !data.Secret.IsUnknown() {
		params["secret"] = data.Secret.ValueString()
	}
	if !data.Peeruser.IsNull() && !data.Peeruser.IsUnknown() {
		params["peeruser"] = data.Peeruser.ValueString()
	}
	if !data.Peersecret.IsNull() && !data.Peersecret.IsUnknown() {
		params["peersecret"] = data.Peersecret.ValueString()
	}
	if !data.DiscoveryAuth.IsNull() && !data.DiscoveryAuth.IsUnknown() {
		params["discovery_auth"] = data.DiscoveryAuth.ValueString()
	}

	result, err := r.client.Call("iscsi.auth.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_auth: %s", err))
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
	result, err = r.client.Call("iscsi.auth.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back iscsi_auth: %s", err))
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
	if v, ok := resultMap["tag"]; ok {
		switch val := v.(type) {
		case float64:
			data.Tag = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Tag = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["user"]; ok {
		switch val := v.(type) {
		case string:
			data.User = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.User = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.User = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["secret"]; ok {
		switch val := v.(type) {
		case string:
			data.Secret = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Secret = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Secret = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Peeruser.IsUnknown() {
		data.Peeruser = types.StringNull()
	}
	if data.Peersecret.IsUnknown() {
		data.Peersecret = types.StringNull()
	}
	if data.DiscoveryAuth.IsUnknown() {
		data.DiscoveryAuth = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiAuthResourceModel
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

	result, err := r.client.Call("iscsi.auth.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_auth: %s", err))
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
	if v, ok := resultMap["tag"]; ok {
		switch val := v.(type) {
		case float64:
			data.Tag = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Tag = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["user"]; ok {
		switch val := v.(type) {
		case string:
			data.User = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.User = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.User = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["secret"]; ok {
		switch val := v.(type) {
		case string:
			data.Secret = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Secret = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Secret = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Peeruser.IsUnknown() {
		data.Peeruser = types.StringNull()
	}
	if data.Peersecret.IsUnknown() {
		data.Peersecret = types.StringNull()
	}
	if data.DiscoveryAuth.IsUnknown() {
		data.DiscoveryAuth = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiAuthResourceModel
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
	if !data.Tag.IsNull() && !data.Tag.IsUnknown() {
		params["tag"] = data.Tag.ValueInt64()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Secret.IsNull() && !data.Secret.IsUnknown() {
		params["secret"] = data.Secret.ValueString()
	}
	if !data.Peeruser.IsNull() && !data.Peeruser.IsUnknown() {
		params["peeruser"] = data.Peeruser.ValueString()
	}
	if !data.Peersecret.IsNull() && !data.Peersecret.IsUnknown() {
		params["peersecret"] = data.Peersecret.ValueString()
	}
	if !data.DiscoveryAuth.IsNull() && !data.DiscoveryAuth.IsUnknown() {
		params["discovery_auth"] = data.DiscoveryAuth.ValueString()
	}

	_, err = r.client.Call("iscsi.auth.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_auth: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("iscsi.auth.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back iscsi_auth: %s", readErr))
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
	if v, ok := resultMap["tag"]; ok {
		switch val := v.(type) {
		case float64:
			data.Tag = types.Int64Value(int64(val))
		case map[string]interface{}:
			if parsed, ok := val["parsed"]; ok && parsed != nil {
				if fv, ok := parsed.(float64); ok {
					data.Tag = types.Int64Value(int64(fv))
				}
			}
		}
	}
	if v, ok := resultMap["user"]; ok {
		switch val := v.(type) {
		case string:
			data.User = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.User = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.User = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["secret"]; ok {
		switch val := v.(type) {
		case string:
			data.Secret = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Secret = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Secret = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Peeruser.IsUnknown() {
		data.Peeruser = types.StringNull()
	}
	if data.Peersecret.IsUnknown() {
		data.Peersecret = types.StringNull()
	}
	if data.DiscoveryAuth.IsUnknown() {
		data.DiscoveryAuth = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("iscsi.auth.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_auth: %s", err))
		return
	}
}
