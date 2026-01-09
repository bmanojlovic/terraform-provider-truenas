package provider

import (
	"context"
	"fmt"
	"strconv"

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
	RateLimit types.Int64 `tfsdk:"rate_limit"`
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
				ElementType: types.StringType,
				Required: false,
				Optional: true,
			},
			"exclude": schema.ListAttribute{
				ElementType: types.StringType,
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
			"rate_limit": schema.Int64Attribute{
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
	params["password"] = data.Password.ValueString()
	params["keep_last"] = data.KeepLast.ValueInt64()
	if !data.TransferSetting.IsNull() {
		params["transfer_setting"] = data.TransferSetting.ValueString()
	}
	if !data.AbsolutePaths.IsNull() {
		params["absolute_paths"] = data.AbsolutePaths.ValueBool()
	}
	if !data.CachePath.IsNull() {
		params["cache_path"] = data.CachePath.ValueString()
	}
	if !data.RateLimit.IsNull() {
		params["rate_limit"] = data.RateLimit.ValueInt64()
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

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("cloud_backup.get_instance", resourceID)
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

	// Get ID from current state (not plan)
	var state CloudBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
	params["password"] = data.Password.ValueString()
	params["keep_last"] = data.KeepLast.ValueInt64()
	if !data.TransferSetting.IsNull() {
		params["transfer_setting"] = data.TransferSetting.ValueString()
	}
	if !data.AbsolutePaths.IsNull() {
		params["absolute_paths"] = data.AbsolutePaths.ValueBool()
	}
	if !data.CachePath.IsNull() {
		params["cache_path"] = data.CachePath.ValueString()
	}
	if !data.RateLimit.IsNull() {
		params["rate_limit"] = data.RateLimit.ValueInt64()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("cloud_backup.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("cloud_backup.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
