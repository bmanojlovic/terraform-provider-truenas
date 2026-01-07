package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type CloudBackupResource struct {
	client *client.Client
}

type CloudBackupResourceModel struct {
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
	Password types.String `tfsdk:"password"`
	KeepLast types.Int64 `tfsdk:"keep_last"`
	TransferSetting types.String `tfsdk:"transfer_setting"`
	AbsolutePaths types.Bool `tfsdk:"absolute_paths"`
	CachePath types.String `tfsdk:"cache_path"`
	RateLimit types.String `tfsdk:"rate_limit"`
}

func NewCloudBackupResource() resource.Resource {
	return &CloudBackupResource{}
}

func (r *CloudBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_backup"
}

func (r *CloudBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS cloud_backup resource",
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
			"password": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"keep_last": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"transfer_setting": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"absolute_paths": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"cache_path": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"rate_limit": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *CloudBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"description": data.Description.ValueString(),
		"path": data.Path.ValueString(),
		"credentials": data.Credentials.ValueInt64(),
		"pre_script": data.PreScript.ValueString(),
		"post_script": data.PostScript.ValueString(),
		"snapshot": data.Snapshot.ValueBool(),
		"args": data.Args.ValueString(),
		"enabled": data.Enabled.ValueBool(),
		"password": data.Password.ValueString(),
		"keep_last": data.KeepLast.ValueInt64(),
		"transfer_setting": data.TransferSetting.ValueString(),
		"absolute_paths": data.AbsolutePaths.ValueBool(),
		"cache_path": data.CachePath.ValueString(),
		"rate_limit": data.RateLimit.ValueString(),
	}

	result, err := r.client.Call("cloud_backup.create", params)
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

func (r *CloudBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("cloud_backup.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"description": data.Description.ValueString(),
		"path": data.Path.ValueString(),
		"credentials": data.Credentials.ValueInt64(),
		"pre_script": data.PreScript.ValueString(),
		"post_script": data.PostScript.ValueString(),
		"snapshot": data.Snapshot.ValueBool(),
		"args": data.Args.ValueString(),
		"enabled": data.Enabled.ValueBool(),
		"password": data.Password.ValueString(),
		"keep_last": data.KeepLast.ValueInt64(),
		"transfer_setting": data.TransferSetting.ValueString(),
		"absolute_paths": data.AbsolutePaths.ValueBool(),
		"cache_path": data.CachePath.ValueString(),
		"rate_limit": data.RateLimit.ValueString(),
	}

	_, err := r.client.Call("cloud_backup.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("cloud_backup.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
