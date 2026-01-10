package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
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
		MarkdownDescription: "None",
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
	if !data.Purpose.IsNull() {
		params["purpose"] = data.Purpose.ValueString()
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Readonly.IsNull() {
		params["readonly"] = data.Readonly.ValueBool()
	}
	if !data.Browsable.IsNull() {
		params["browsable"] = data.Browsable.ValueBool()
	}
	if !data.AccessBasedShareEnumeration.IsNull() {
		params["access_based_share_enumeration"] = data.AccessBasedShareEnumeration.ValueBool()
	}
	if !data.Audit.IsNull() {
		params["audit"] = data.Audit.ValueString()
	}
	if !data.Options.IsNull() {
		params["options"] = data.Options.ValueString()
	}

	result, err := r.client.Call("sharing.smb.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create sharing_smb: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
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

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("sharing.smb.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read sharing_smb: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["purpose"]; ok && v != nil {
			data.Purpose = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			data.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			data.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Readonly = types.BoolValue(bv) }
		}
		if v, ok := resultMap["browsable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Browsable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["access_based_share_enumeration"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AccessBasedShareEnumeration = types.BoolValue(bv) }
		}
		if v, ok := resultMap["audit"]; ok && v != nil {
			data.Audit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["options"]; ok && v != nil {
			data.Options = types.StringValue(fmt.Sprintf("%v", v))
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

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Purpose.IsNull() {
		params["purpose"] = data.Purpose.ValueString()
	}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Path.IsNull() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.Readonly.IsNull() {
		params["readonly"] = data.Readonly.ValueBool()
	}
	if !data.Browsable.IsNull() {
		params["browsable"] = data.Browsable.ValueBool()
	}
	if !data.AccessBasedShareEnumeration.IsNull() {
		params["access_based_share_enumeration"] = data.AccessBasedShareEnumeration.ValueBool()
	}
	if !data.Audit.IsNull() {
		params["audit"] = data.Audit.ValueString()
	}
	if !data.Options.IsNull() {
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

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("sharing.smb.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete sharing_smb: %s", err))
		return
	}
}
