package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type ReplicationResource struct {
	client *client.Client
}

type ReplicationResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Direction types.String `tfsdk:"direction"`
	Transport types.String `tfsdk:"transport"`
	SshCredentials types.String `tfsdk:"ssh_credentials"`
	NetcatActiveSide types.String `tfsdk:"netcat_active_side"`
	NetcatActiveSideListenAddress types.String `tfsdk:"netcat_active_side_listen_address"`
	NetcatActiveSidePortMin types.String `tfsdk:"netcat_active_side_port_min"`
	NetcatActiveSidePortMax types.String `tfsdk:"netcat_active_side_port_max"`
	NetcatPassiveSideConnectAddress types.String `tfsdk:"netcat_passive_side_connect_address"`
	Sudo types.Bool `tfsdk:"sudo"`
	SourceDatasets types.List `tfsdk:"source_datasets"`
	TargetDataset types.String `tfsdk:"target_dataset"`
	Recursive types.Bool `tfsdk:"recursive"`
	Exclude types.List `tfsdk:"exclude"`
	Properties types.Bool `tfsdk:"properties"`
	PropertiesExclude types.List `tfsdk:"properties_exclude"`
	PropertiesOverride types.Object `tfsdk:"properties_override"`
	Replicate types.Bool `tfsdk:"replicate"`
	Encryption types.Bool `tfsdk:"encryption"`
	EncryptionInherit types.String `tfsdk:"encryption_inherit"`
	EncryptionKey types.String `tfsdk:"encryption_key"`
	EncryptionKeyFormat types.String `tfsdk:"encryption_key_format"`
	EncryptionKeyLocation types.String `tfsdk:"encryption_key_location"`
	PeriodicSnapshotTasks types.List `tfsdk:"periodic_snapshot_tasks"`
	NamingSchema types.List `tfsdk:"naming_schema"`
	AlsoIncludeNamingSchema types.List `tfsdk:"also_include_naming_schema"`
	NameRegex types.String `tfsdk:"name_regex"`
	Auto types.Bool `tfsdk:"auto"`
	Schedule types.String `tfsdk:"schedule"`
	RestrictSchedule types.String `tfsdk:"restrict_schedule"`
	OnlyMatchingSchedule types.Bool `tfsdk:"only_matching_schedule"`
	AllowFromScratch types.Bool `tfsdk:"allow_from_scratch"`
	Readonly types.String `tfsdk:"readonly"`
	HoldPendingSnapshots types.Bool `tfsdk:"hold_pending_snapshots"`
	RetentionPolicy types.String `tfsdk:"retention_policy"`
	LifetimeValue types.String `tfsdk:"lifetime_value"`
	LifetimeUnit types.String `tfsdk:"lifetime_unit"`
	Lifetimes types.List `tfsdk:"lifetimes"`
	Compression types.String `tfsdk:"compression"`
	SpeedLimit types.String `tfsdk:"speed_limit"`
	LargeBlock types.Bool `tfsdk:"large_block"`
	Embed types.Bool `tfsdk:"embed"`
	Compressed types.Bool `tfsdk:"compressed"`
	Retries types.Int64 `tfsdk:"retries"`
	LoggingLevel types.String `tfsdk:"logging_level"`
	Enabled types.Bool `tfsdk:"enabled"`
}

func NewReplicationResource() resource.Resource {
	return &ReplicationResource{}
}

func (r *ReplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

func (r *ReplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS replication resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"direction": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"transport": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"ssh_credentials": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"netcat_active_side": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"netcat_active_side_listen_address": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"netcat_active_side_port_min": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"netcat_active_side_port_max": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"netcat_passive_side_connect_address": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"sudo": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"source_datasets": schema.ListAttribute{
				Required: true,
				Optional: false,
			},
			"target_dataset": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"recursive": schema.BoolAttribute{
				Required: true,
				Optional: false,
			},
			"exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"properties": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"properties_exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"properties_override": schema.ObjectAttribute{
				Required: false,
				Optional: true,
			},
			"replicate": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"encryption_inherit": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"encryption_key": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"encryption_key_format": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"encryption_key_location": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"periodic_snapshot_tasks": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"naming_schema": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"also_include_naming_schema": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"name_regex": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"auto": schema.BoolAttribute{
				Required: true,
				Optional: false,
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"restrict_schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"only_matching_schedule": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"allow_from_scratch": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"readonly": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"hold_pending_snapshots": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"retention_policy": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"lifetime_value": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"lifetime_unit": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"lifetimes": schema.ListAttribute{
				Required: false,
				Optional: true,
			},
			"compression": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"speed_limit": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"large_block": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"embed": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"compressed": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"retries": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"logging_level": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *ReplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ReplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"name": data.Name.ValueString(),
		"direction": data.Direction.ValueString(),
		"transport": data.Transport.ValueString(),
		"ssh_credentials": data.SshCredentials.ValueString(),
		"netcat_active_side": data.NetcatActiveSide.ValueString(),
		"netcat_active_side_listen_address": data.NetcatActiveSideListenAddress.ValueString(),
		"netcat_active_side_port_min": data.NetcatActiveSidePortMin.ValueString(),
		"netcat_active_side_port_max": data.NetcatActiveSidePortMax.ValueString(),
		"netcat_passive_side_connect_address": data.NetcatPassiveSideConnectAddress.ValueString(),
		"sudo": data.Sudo.ValueBool(),
		"target_dataset": data.TargetDataset.ValueString(),
		"recursive": data.Recursive.ValueBool(),
		"properties": data.Properties.ValueBool(),
		"replicate": data.Replicate.ValueBool(),
		"encryption": data.Encryption.ValueBool(),
		"encryption_inherit": data.EncryptionInherit.ValueString(),
		"encryption_key": data.EncryptionKey.ValueString(),
		"encryption_key_format": data.EncryptionKeyFormat.ValueString(),
		"encryption_key_location": data.EncryptionKeyLocation.ValueString(),
		"name_regex": data.NameRegex.ValueString(),
		"auto": data.Auto.ValueBool(),
		"schedule": data.Schedule.ValueString(),
		"restrict_schedule": data.RestrictSchedule.ValueString(),
		"only_matching_schedule": data.OnlyMatchingSchedule.ValueBool(),
		"allow_from_scratch": data.AllowFromScratch.ValueBool(),
		"readonly": data.Readonly.ValueString(),
		"hold_pending_snapshots": data.HoldPendingSnapshots.ValueBool(),
		"retention_policy": data.RetentionPolicy.ValueString(),
		"lifetime_value": data.LifetimeValue.ValueString(),
		"lifetime_unit": data.LifetimeUnit.ValueString(),
		"compression": data.Compression.ValueString(),
		"speed_limit": data.SpeedLimit.ValueString(),
		"large_block": data.LargeBlock.ValueBool(),
		"embed": data.Embed.ValueBool(),
		"compressed": data.Compressed.ValueBool(),
		"retries": data.Retries.ValueInt64(),
		"logging_level": data.LoggingLevel.ValueString(),
		"enabled": data.Enabled.ValueBool(),
	}

	result, err := r.client.Call("replication.create", params)
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

func (r *ReplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("replication.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ReplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"name": data.Name.ValueString(),
		"direction": data.Direction.ValueString(),
		"transport": data.Transport.ValueString(),
		"ssh_credentials": data.SshCredentials.ValueString(),
		"netcat_active_side": data.NetcatActiveSide.ValueString(),
		"netcat_active_side_listen_address": data.NetcatActiveSideListenAddress.ValueString(),
		"netcat_active_side_port_min": data.NetcatActiveSidePortMin.ValueString(),
		"netcat_active_side_port_max": data.NetcatActiveSidePortMax.ValueString(),
		"netcat_passive_side_connect_address": data.NetcatPassiveSideConnectAddress.ValueString(),
		"sudo": data.Sudo.ValueBool(),
		"target_dataset": data.TargetDataset.ValueString(),
		"recursive": data.Recursive.ValueBool(),
		"properties": data.Properties.ValueBool(),
		"replicate": data.Replicate.ValueBool(),
		"encryption": data.Encryption.ValueBool(),
		"encryption_inherit": data.EncryptionInherit.ValueString(),
		"encryption_key": data.EncryptionKey.ValueString(),
		"encryption_key_format": data.EncryptionKeyFormat.ValueString(),
		"encryption_key_location": data.EncryptionKeyLocation.ValueString(),
		"name_regex": data.NameRegex.ValueString(),
		"auto": data.Auto.ValueBool(),
		"schedule": data.Schedule.ValueString(),
		"restrict_schedule": data.RestrictSchedule.ValueString(),
		"only_matching_schedule": data.OnlyMatchingSchedule.ValueBool(),
		"allow_from_scratch": data.AllowFromScratch.ValueBool(),
		"readonly": data.Readonly.ValueString(),
		"hold_pending_snapshots": data.HoldPendingSnapshots.ValueBool(),
		"retention_policy": data.RetentionPolicy.ValueString(),
		"lifetime_value": data.LifetimeValue.ValueString(),
		"lifetime_unit": data.LifetimeUnit.ValueString(),
		"compression": data.Compression.ValueString(),
		"speed_limit": data.SpeedLimit.ValueString(),
		"large_block": data.LargeBlock.ValueBool(),
		"embed": data.Embed.ValueBool(),
		"compressed": data.Compressed.ValueBool(),
		"retries": data.Retries.ValueInt64(),
		"logging_level": data.LoggingLevel.ValueString(),
		"enabled": data.Enabled.ValueBool(),
	}

	_, err := r.client.Call("replication.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("replication.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
