package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"encoding/json"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SharingSmbResource struct {
	client *client.Client
}

type SharingSmbResourceModel struct {
	ID types.String `tfsdk:"id"`
	Purpose types.String `tfsdk:"purpose"`
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
	Enabled types.Bool `tfsdk:"enabled"`
	Comment types.String `tfsdk:"comment"`
	Readonly types.Bool `tfsdk:"readonly"`
	Browsable types.Bool `tfsdk:"browsable"`
	AccessBasedShareEnumeration types.Bool `tfsdk:"access_based_share_enumeration"`
	Audit types.String `tfsdk:"audit"`
	Options types.String `tfsdk:"options"`
}

func NewSharingSmbResource() resource.Resource {
	return &SharingSmbResource{}
}

func (r *SharingSmbResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sharing_smb"
}

func (r *SharingSmbResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SharingSmbResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS sharing_smb resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"purpose": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "This parameter sets the purpose of the SMB share. It controls how the SMB share behaves and what fea",
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "SMB share name. SMB share names are case-insensitive and must be unique, and are subject     to the ",
			},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Local server path to share by using the SMB protocol. The path must start with `/mnt/` and must be i",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If unset, the SMB share is not available over the SMB protocol. ",
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Text field that is seen next to a share when an SMB client requests a list of SMB shares on the True",
			},
			"readonly": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set, SMB clients cannot create or change files and directories in the SMB share.  NOTE: If set, t",
			},
			"browsable": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set, the share is included when an SMB client requests a list of SMB shares on the TrueNAS server",
			},
			"access_based_share_enumeration": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set, the share is only included when an SMB client requests a list of shares on the SMB server if",
			},
			"audit": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Audit configuration for monitoring SMB share access and operations.",
			},
			"options": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Additional configuration related to the configured SMB share purpose. If null, then the default     ",
			},
		},
	}
}

func (r *SharingSmbResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SharingSmbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SharingSmbResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Purpose.IsNull() && !data.Purpose.IsUnknown() {
		params["purpose"] = data.Purpose.ValueString()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Readonly.IsNull() && !data.Readonly.IsUnknown() {
		params["readonly"] = data.Readonly.ValueBool()
	}
	if !data.Browsable.IsNull() && !data.Browsable.IsUnknown() {
		params["browsable"] = data.Browsable.ValueBool()
	}
	if !data.AccessBasedShareEnumeration.IsNull() && !data.AccessBasedShareEnumeration.IsUnknown() {
		params["access_based_share_enumeration"] = data.AccessBasedShareEnumeration.ValueBool()
	}
	if !data.Audit.IsNull() && !data.Audit.IsUnknown() {
		var auditObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Audit.ValueString()), &auditObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse audit: %s", err))
			return
		}
		params["audit"] = auditObj
	}
	if !data.Options.IsNull() && !data.Options.IsUnknown() {
		params["options"] = data.Options.ValueString()
	}

	result, err := r.client.Call("sharing.smb.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create sharing_smb: %s", err))
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
	result, err = r.client.Call("sharing.smb.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back sharing_smb: %s", err))
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
		if v, ok := resultMap["name"]; ok {
			switch val := v.(type) {
			case string:
				data.Name = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Name = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["path"]; ok {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["options"]; ok {
			if v == nil {
				data.Options = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Options = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Options = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Options = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingSmbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SharingSmbResourceModel
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

	result, err := r.client.Call("sharing.smb.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read sharing_smb: %s", err))
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
		if v, ok := resultMap["name"]; ok {
			switch val := v.(type) {
			case string:
				data.Name = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Name = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["path"]; ok {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["options"]; ok {
			if v == nil {
				data.Options = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Options = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Options = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Options = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingSmbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SharingSmbResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SharingSmbResourceModel
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
	if !data.Purpose.IsNull() && !data.Purpose.IsUnknown() {
		params["purpose"] = data.Purpose.ValueString()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Readonly.IsNull() && !data.Readonly.IsUnknown() {
		params["readonly"] = data.Readonly.ValueBool()
	}
	if !data.Browsable.IsNull() && !data.Browsable.IsUnknown() {
		params["browsable"] = data.Browsable.ValueBool()
	}
	if !data.AccessBasedShareEnumeration.IsNull() && !data.AccessBasedShareEnumeration.IsUnknown() {
		params["access_based_share_enumeration"] = data.AccessBasedShareEnumeration.ValueBool()
	}
	if !data.Audit.IsNull() && !data.Audit.IsUnknown() {
		var auditObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Audit.ValueString()), &auditObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse audit: %s", err))
			return
		}
		params["audit"] = auditObj
	}
	if !data.Options.IsNull() && !data.Options.IsUnknown() {
		params["options"] = data.Options.ValueString()
	}

	_, err = r.client.Call("sharing.smb.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update sharing_smb: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingSmbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SharingSmbResourceModel
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

	_, err = r.client.Call("sharing.smb.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete sharing_smb: %s", err))
		return
	}
}
