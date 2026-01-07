package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type CloudBackupRestoreActionResource struct {
	client *client.Client
}

type CloudBackupRestoreActionResourceModel struct {
	ID types.String `tfsdk:"id"`
	ResourceID types.String `tfsdk:"resource_id"`
	SnapshotId types.String `tfsdk:"snapshot_id"`
	Subfolder types.String `tfsdk:"subfolder"`
	DestinationPath types.String `tfsdk:"destination_path"`
	Options types.String `tfsdk:"options"`
}

func NewCloudBackupRestoreActionResource() resource.Resource {
	return &CloudBackupRestoreActionResource{}
}

func (r *CloudBackupRestoreActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_backup_restore_action"
}

func (r *CloudBackupRestoreActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Executes restore action on cloud_backup resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the resource to perform action on",
			},
			"snapshot_id": schema.StringAttribute{
				Optional: true,
			},
			"subfolder": schema.StringAttribute{
				Optional: true,
			},
			"destination_path": schema.StringAttribute{
				Optional: true,
			},
			"options": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *CloudBackupRestoreActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CloudBackupRestoreActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudBackupRestoreActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
		if !data.SnapshotId.IsNull() {
			params["snapshot_id"] = data.SnapshotId.ValueString()
		}
		if !data.Subfolder.IsNull() {
			params["subfolder"] = data.Subfolder.ValueString()
		}
		if !data.DestinationPath.IsNull() {
			params["destination_path"] = data.DestinationPath.ValueString()
		}
		if !data.Options.IsNull() {
			params["options"] = data.Options.ValueString()
		}

	_, err := r.client.Call("cloud_backup/restore", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute restore: %s", err.Error()))
		return
	}

	// Use timestamp as ID since actions are ephemeral
	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupRestoreActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are ephemeral - nothing to read
	var data CloudBackupRestoreActionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *CloudBackupRestoreActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Actions are immutable - re-execute on update
	var data CloudBackupRestoreActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
		if !data.SnapshotId.IsNull() {
			params["snapshot_id"] = data.SnapshotId.ValueString()
		}
		if !data.Subfolder.IsNull() {
			params["subfolder"] = data.Subfolder.ValueString()
		}
		if !data.DestinationPath.IsNull() {
			params["destination_path"] = data.DestinationPath.ValueString()
		}
		if !data.Options.IsNull() {
			params["options"] = data.Options.ValueString()
		}

	_, err := r.client.Call("cloud_backup/restore", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute restore: %s", err.Error()))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudBackupRestoreActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Actions cannot be undone - just remove from state
}
