package provider

import (
	"context"
	"fmt"


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
	Audit types.Object `tfsdk:"audit"`
	Options types.String `tfsdk:"options"`
}

func NewSharingSmbResource() resource.Resource {
	return &SharingSmbResource{}
}

func (r *SharingSmbResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sharing_smb"
}

func (r *SharingSmbResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS sharing_smb resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"purpose": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"readonly": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"browsable": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"access_based_share_enumeration": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"audit": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"options": schema.StringAttribute{
				Required: false,
				Optional: true,
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
	params["name"] = data.Name.ValueString()
	params["path"] = data.Path.ValueString()
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
	if !data.Options.IsNull() {
		params["options"] = data.Options.ValueString()
	}

	result, err := r.client.Call("sharing/smb.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

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

	_, err := r.client.Call("sharing/smb.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingSmbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SharingSmbResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Purpose.IsNull() {
		params["purpose"] = data.Purpose.ValueString()
	}
	params["name"] = data.Name.ValueString()
	params["path"] = data.Path.ValueString()
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
	if !data.Options.IsNull() {
		params["options"] = data.Options.ValueString()
	}

	_, err := r.client.Call("sharing/smb.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharingSmbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SharingSmbResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("sharing/smb.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
