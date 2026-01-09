package provider

import (
	"context"
	"fmt"


	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type CloudsyncResource struct {
	client *client.Client
}

type CloudsyncResourceModel struct {
	ID types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Path types.String `tfsdk:"path"`
	Credentials types.Int64 `tfsdk:"credentials"`
	Attributes types.Object `tfsdk:"attributes"`
	Schedule types.Object `tfsdk:"schedule"`
	PreScript types.String `tfsdk:"pre_script"`
	PostScript types.String `tfsdk:"post_script"`
	Snapshot types.Bool `tfsdk:"snapshot"`
	Include types.List `tfsdk:"include"`
	Exclude types.List `tfsdk:"exclude"`
	Args types.String `tfsdk:"args"`
	Enabled types.Bool `tfsdk:"enabled"`
	Bwlimit types.List `tfsdk:"bwlimit"`
	Transfers types.String `tfsdk:"transfers"`
	Direction types.String `tfsdk:"direction"`
	TransferMode types.String `tfsdk:"transfer_mode"`
	Encryption types.Bool `tfsdk:"encryption"`
	FilenameEncryption types.Bool `tfsdk:"filename_encryption"`
	EncryptionPassword types.String `tfsdk:"encryption_password"`
	EncryptionSalt types.String `tfsdk:"encryption_salt"`
	CreateEmptySrcDirs types.Bool `tfsdk:"create_empty_src_dirs"`
	FollowSymlinks types.Bool `tfsdk:"follow_symlinks"`
}

func NewCloudsyncResource() resource.Resource {
	return &CloudsyncResource{}
}

func (r *CloudsyncResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudsync"
}

func (r *CloudsyncResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS cloudsync resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"credentials": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"attributes": schema.ObjectAttribute{
				Required: true,
				Optional: false,
			},
			"schedule": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"pre_script": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"post_script": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"snapshot": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"include": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"args": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"bwlimit": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"transfers": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"direction": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"transfer_mode": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"filename_encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"encryption_password": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"encryption_salt": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"create_empty_src_dirs": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"follow_symlinks": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *CloudsyncResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudsyncResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	params["path"] = data.Path.ValueString()
	params["credentials"] = data.Credentials.ValueInt64()
	if !data.PreScript.IsNull() {
		params["pre_script"] = data.PreScript.ValueString()
	}
	if !data.PostScript.IsNull() {
		params["post_script"] = data.PostScript.ValueString()
	}
	if !data.Snapshot.IsNull() {
		params["snapshot"] = data.Snapshot.ValueBool()
	}
	if !data.Args.IsNull() {
		params["args"] = data.Args.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Transfers.IsNull() {
		params["transfers"] = data.Transfers.ValueString()
	}
	params["direction"] = data.Direction.ValueString()
	params["transfer_mode"] = data.TransferMode.ValueString()
	if !data.Encryption.IsNull() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.FilenameEncryption.IsNull() {
		params["filename_encryption"] = data.FilenameEncryption.ValueBool()
	}
	if !data.EncryptionPassword.IsNull() {
		params["encryption_password"] = data.EncryptionPassword.ValueString()
	}
	if !data.EncryptionSalt.IsNull() {
		params["encryption_salt"] = data.EncryptionSalt.ValueString()
	}
	if !data.CreateEmptySrcDirs.IsNull() {
		params["create_empty_src_dirs"] = data.CreateEmptySrcDirs.ValueBool()
	}
	if !data.FollowSymlinks.IsNull() {
		params["follow_symlinks"] = data.FollowSymlinks.ValueBool()
	}

	result, err := r.client.Call("cloudsync.create", params)
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

func (r *CloudsyncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("cloudsync.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudsyncResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	params["path"] = data.Path.ValueString()
	params["credentials"] = data.Credentials.ValueInt64()
	if !data.PreScript.IsNull() {
		params["pre_script"] = data.PreScript.ValueString()
	}
	if !data.PostScript.IsNull() {
		params["post_script"] = data.PostScript.ValueString()
	}
	if !data.Snapshot.IsNull() {
		params["snapshot"] = data.Snapshot.ValueBool()
	}
	if !data.Args.IsNull() {
		params["args"] = data.Args.ValueString()
	}
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.Transfers.IsNull() {
		params["transfers"] = data.Transfers.ValueString()
	}
	params["direction"] = data.Direction.ValueString()
	params["transfer_mode"] = data.TransferMode.ValueString()
	if !data.Encryption.IsNull() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.FilenameEncryption.IsNull() {
		params["filename_encryption"] = data.FilenameEncryption.ValueBool()
	}
	if !data.EncryptionPassword.IsNull() {
		params["encryption_password"] = data.EncryptionPassword.ValueString()
	}
	if !data.EncryptionSalt.IsNull() {
		params["encryption_salt"] = data.EncryptionSalt.ValueString()
	}
	if !data.CreateEmptySrcDirs.IsNull() {
		params["create_empty_src_dirs"] = data.CreateEmptySrcDirs.ValueBool()
	}
	if !data.FollowSymlinks.IsNull() {
		params["follow_symlinks"] = data.FollowSymlinks.ValueBool()
	}

	_, err := r.client.Call("cloudsync.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudsyncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudsyncResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("cloudsync.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
