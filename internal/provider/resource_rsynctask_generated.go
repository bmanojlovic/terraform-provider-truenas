package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type RsynctaskResource struct {
	client *client.Client
}

type RsynctaskResourceModel struct {
	ID types.String `tfsdk:"id"`
	Path types.String `tfsdk:"path"`
	User types.String `tfsdk:"user"`
	Mode types.String `tfsdk:"mode"`
	Remotehost types.String `tfsdk:"remotehost"`
	Remoteport types.String `tfsdk:"remoteport"`
	Remotemodule types.String `tfsdk:"remotemodule"`
	SshCredentials types.String `tfsdk:"ssh_credentials"`
	Remotepath types.String `tfsdk:"remotepath"`
	Direction types.String `tfsdk:"direction"`
	Desc types.String `tfsdk:"desc"`
	Schedule types.Object `tfsdk:"schedule"`
	Recursive types.Bool `tfsdk:"recursive"`
	Times types.Bool `tfsdk:"times"`
	Compress types.Bool `tfsdk:"compress"`
	Archive types.Bool `tfsdk:"archive"`
	Delete types.Bool `tfsdk:"delete"`
	Quiet types.Bool `tfsdk:"quiet"`
	Preserveperm types.Bool `tfsdk:"preserveperm"`
	Preserveattr types.Bool `tfsdk:"preserveattr"`
	Delayupdates types.Bool `tfsdk:"delayupdates"`
	Extra types.List `tfsdk:"extra"`
	Enabled types.Bool `tfsdk:"enabled"`
	ValidateRpath types.Bool `tfsdk:"validate_rpath"`
	SshKeyscan types.Bool `tfsdk:"ssh_keyscan"`
}

func NewRsynctaskResource() resource.Resource {
	return &RsynctaskResource{}
}

func (r *RsynctaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rsynctask"
}

func (r *RsynctaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS rsynctask resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"user": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"mode": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"remotehost": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"remoteport": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"remotemodule": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"ssh_credentials": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"remotepath": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"direction": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"desc": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"schedule": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"recursive": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"times": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"compress": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"archive": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"delete": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"quiet": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"preserveperm": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"preserveattr": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"delayupdates": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"extra": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"validate_rpath": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"ssh_keyscan": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *RsynctaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RsynctaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RsynctaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"path": data.Path.ValueString(),
		"user": data.User.ValueString(),
		"mode": data.Mode.ValueString(),
		"remotehost": data.Remotehost.ValueString(),
		"remoteport": data.Remoteport.ValueString(),
		"remotemodule": data.Remotemodule.ValueString(),
		"ssh_credentials": data.SshCredentials.ValueString(),
		"remotepath": data.Remotepath.ValueString(),
		"direction": data.Direction.ValueString(),
		"desc": data.Desc.ValueString(),
		"recursive": data.Recursive.ValueBool(),
		"times": data.Times.ValueBool(),
		"compress": data.Compress.ValueBool(),
		"archive": data.Archive.ValueBool(),
		"delete": data.Delete.ValueBool(),
		"quiet": data.Quiet.ValueBool(),
		"preserveperm": data.Preserveperm.ValueBool(),
		"preserveattr": data.Preserveattr.ValueBool(),
		"delayupdates": data.Delayupdates.ValueBool(),
		"enabled": data.Enabled.ValueBool(),
		"validate_rpath": data.ValidateRpath.ValueBool(),
		"ssh_keyscan": data.SshKeyscan.ValueBool(),
	}

	result, err := r.client.Call("rsynctask.create", params)
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

func (r *RsynctaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RsynctaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("rsynctask.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RsynctaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"path": data.Path.ValueString(),
		"user": data.User.ValueString(),
		"mode": data.Mode.ValueString(),
		"remotehost": data.Remotehost.ValueString(),
		"remoteport": data.Remoteport.ValueString(),
		"remotemodule": data.Remotemodule.ValueString(),
		"ssh_credentials": data.SshCredentials.ValueString(),
		"remotepath": data.Remotepath.ValueString(),
		"direction": data.Direction.ValueString(),
		"desc": data.Desc.ValueString(),
		"recursive": data.Recursive.ValueBool(),
		"times": data.Times.ValueBool(),
		"compress": data.Compress.ValueBool(),
		"archive": data.Archive.ValueBool(),
		"delete": data.Delete.ValueBool(),
		"quiet": data.Quiet.ValueBool(),
		"preserveperm": data.Preserveperm.ValueBool(),
		"preserveattr": data.Preserveattr.ValueBool(),
		"delayupdates": data.Delayupdates.ValueBool(),
		"enabled": data.Enabled.ValueBool(),
		"validate_rpath": data.ValidateRpath.ValueBool(),
		"ssh_keyscan": data.SshKeyscan.ValueBool(),
	}

	_, err := r.client.Call("rsynctask.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RsynctaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("rsynctask.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
