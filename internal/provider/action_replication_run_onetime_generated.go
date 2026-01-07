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

type ReplicationRunOnetimeActionResource struct {
	client *client.Client
}

type ReplicationRunOnetimeActionResourceModel struct {
	ID types.String `tfsdk:"id"`
	ResourceID types.String `tfsdk:"resource_id"`
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
	RestrictSchedule types.String `tfsdk:"restrict_schedule"`
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
	ExcludeMountpointProperty types.Bool `tfsdk:"exclude_mountpoint_property"`
	OnlyFromScratch types.Bool `tfsdk:"only_from_scratch"`
	Mount types.Bool `tfsdk:"mount"`
}

func NewReplicationRunOnetimeActionResource() resource.Resource {
	return &ReplicationRunOnetimeActionResource{}
}

func (r *ReplicationRunOnetimeActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_run_onetime_action"
}

func (r *ReplicationRunOnetimeActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Executes run_onetime action on replication resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the resource to perform action on",
			},
			"direction": schema.StringAttribute{
				Optional: true,
			},
			"transport": schema.StringAttribute{
				Optional: true,
			},
			"ssh_credentials": schema.StringAttribute{
				Optional: true,
			},
			"netcat_active_side": schema.StringAttribute{
				Optional: true,
			},
			"netcat_active_side_listen_address": schema.StringAttribute{
				Optional: true,
			},
			"netcat_active_side_port_min": schema.StringAttribute{
				Optional: true,
			},
			"netcat_active_side_port_max": schema.StringAttribute{
				Optional: true,
			},
			"netcat_passive_side_connect_address": schema.StringAttribute{
				Optional: true,
			},
			"sudo": schema.BoolAttribute{
				Optional: true,
			},
			"source_datasets": schema.ListAttribute{
				Optional: true,
			},
			"target_dataset": schema.StringAttribute{
				Optional: true,
			},
			"recursive": schema.BoolAttribute{
				Optional: true,
			},
			"exclude": schema.ListAttribute{
				Optional: true,
			},
			"properties": schema.BoolAttribute{
				Optional: true,
			},
			"properties_exclude": schema.ListAttribute{
				Optional: true,
			},
			"properties_override": schema.ObjectAttribute{
				Optional: true,
			},
			"replicate": schema.BoolAttribute{
				Optional: true,
			},
			"encryption": schema.BoolAttribute{
				Optional: true,
			},
			"encryption_inherit": schema.StringAttribute{
				Optional: true,
			},
			"encryption_key": schema.StringAttribute{
				Optional: true,
			},
			"encryption_key_format": schema.StringAttribute{
				Optional: true,
			},
			"encryption_key_location": schema.StringAttribute{
				Optional: true,
			},
			"periodic_snapshot_tasks": schema.ListAttribute{
				Optional: true,
			},
			"naming_schema": schema.ListAttribute{
				Optional: true,
			},
			"also_include_naming_schema": schema.ListAttribute{
				Optional: true,
			},
			"name_regex": schema.StringAttribute{
				Optional: true,
			},
			"restrict_schedule": schema.StringAttribute{
				Optional: true,
			},
			"allow_from_scratch": schema.BoolAttribute{
				Optional: true,
			},
			"readonly": schema.StringAttribute{
				Optional: true,
			},
			"hold_pending_snapshots": schema.BoolAttribute{
				Optional: true,
			},
			"retention_policy": schema.StringAttribute{
				Optional: true,
			},
			"lifetime_value": schema.StringAttribute{
				Optional: true,
			},
			"lifetime_unit": schema.StringAttribute{
				Optional: true,
			},
			"lifetimes": schema.ListAttribute{
				Optional: true,
			},
			"compression": schema.StringAttribute{
				Optional: true,
			},
			"speed_limit": schema.StringAttribute{
				Optional: true,
			},
			"large_block": schema.BoolAttribute{
				Optional: true,
			},
			"embed": schema.BoolAttribute{
				Optional: true,
			},
			"compressed": schema.BoolAttribute{
				Optional: true,
			},
			"retries": schema.Int64Attribute{
				Optional: true,
			},
			"logging_level": schema.StringAttribute{
				Optional: true,
			},
			"exclude_mountpoint_property": schema.BoolAttribute{
				Optional: true,
			},
			"only_from_scratch": schema.BoolAttribute{
				Optional: true,
			},
			"mount": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *ReplicationRunOnetimeActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ReplicationRunOnetimeActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ReplicationRunOnetimeActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
		if !data.Direction.IsNull() {
			params["direction"] = data.Direction.ValueString()
		}
		if !data.Transport.IsNull() {
			params["transport"] = data.Transport.ValueString()
		}
		if !data.SshCredentials.IsNull() {
			params["ssh_credentials"] = data.SshCredentials.ValueString()
		}
		if !data.NetcatActiveSide.IsNull() {
			params["netcat_active_side"] = data.NetcatActiveSide.ValueString()
		}
		if !data.NetcatActiveSideListenAddress.IsNull() {
			params["netcat_active_side_listen_address"] = data.NetcatActiveSideListenAddress.ValueString()
		}
		if !data.NetcatActiveSidePortMin.IsNull() {
			params["netcat_active_side_port_min"] = data.NetcatActiveSidePortMin.ValueString()
		}
		if !data.NetcatActiveSidePortMax.IsNull() {
			params["netcat_active_side_port_max"] = data.NetcatActiveSidePortMax.ValueString()
		}
		if !data.NetcatPassiveSideConnectAddress.IsNull() {
			params["netcat_passive_side_connect_address"] = data.NetcatPassiveSideConnectAddress.ValueString()
		}
		if !data.Sudo.IsNull() {
			params["sudo"] = data.Sudo.ValueBool()
		}
		if !data.TargetDataset.IsNull() {
			params["target_dataset"] = data.TargetDataset.ValueString()
		}
		if !data.Recursive.IsNull() {
			params["recursive"] = data.Recursive.ValueBool()
		}
		if !data.Properties.IsNull() {
			params["properties"] = data.Properties.ValueBool()
		}
		if !data.Replicate.IsNull() {
			params["replicate"] = data.Replicate.ValueBool()
		}
		if !data.Encryption.IsNull() {
			params["encryption"] = data.Encryption.ValueBool()
		}
		if !data.EncryptionInherit.IsNull() {
			params["encryption_inherit"] = data.EncryptionInherit.ValueString()
		}
		if !data.EncryptionKey.IsNull() {
			params["encryption_key"] = data.EncryptionKey.ValueString()
		}
		if !data.EncryptionKeyFormat.IsNull() {
			params["encryption_key_format"] = data.EncryptionKeyFormat.ValueString()
		}
		if !data.EncryptionKeyLocation.IsNull() {
			params["encryption_key_location"] = data.EncryptionKeyLocation.ValueString()
		}
		if !data.NameRegex.IsNull() {
			params["name_regex"] = data.NameRegex.ValueString()
		}
		if !data.RestrictSchedule.IsNull() {
			params["restrict_schedule"] = data.RestrictSchedule.ValueString()
		}
		if !data.AllowFromScratch.IsNull() {
			params["allow_from_scratch"] = data.AllowFromScratch.ValueBool()
		}
		if !data.Readonly.IsNull() {
			params["readonly"] = data.Readonly.ValueString()
		}
		if !data.HoldPendingSnapshots.IsNull() {
			params["hold_pending_snapshots"] = data.HoldPendingSnapshots.ValueBool()
		}
		if !data.RetentionPolicy.IsNull() {
			params["retention_policy"] = data.RetentionPolicy.ValueString()
		}
		if !data.LifetimeValue.IsNull() {
			params["lifetime_value"] = data.LifetimeValue.ValueString()
		}
		if !data.LifetimeUnit.IsNull() {
			params["lifetime_unit"] = data.LifetimeUnit.ValueString()
		}
		if !data.Compression.IsNull() {
			params["compression"] = data.Compression.ValueString()
		}
		if !data.SpeedLimit.IsNull() {
			params["speed_limit"] = data.SpeedLimit.ValueString()
		}
		if !data.LargeBlock.IsNull() {
			params["large_block"] = data.LargeBlock.ValueBool()
		}
		if !data.Embed.IsNull() {
			params["embed"] = data.Embed.ValueBool()
		}
		if !data.Compressed.IsNull() {
			params["compressed"] = data.Compressed.ValueBool()
		}
		if !data.Retries.IsNull() {
			params["retries"] = data.Retries.ValueInt64()
		}
		if !data.LoggingLevel.IsNull() {
			params["logging_level"] = data.LoggingLevel.ValueString()
		}
		if !data.ExcludeMountpointProperty.IsNull() {
			params["exclude_mountpoint_property"] = data.ExcludeMountpointProperty.ValueBool()
		}
		if !data.OnlyFromScratch.IsNull() {
			params["only_from_scratch"] = data.OnlyFromScratch.ValueBool()
		}
		if !data.Mount.IsNull() {
			params["mount"] = data.Mount.ValueBool()
		}

	_, err := r.client.Call("replication/run_onetime", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute run_onetime: %s", err.Error()))
		return
	}

	// Use timestamp as ID since actions are ephemeral
	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationRunOnetimeActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are ephemeral - nothing to read
	var data ReplicationRunOnetimeActionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *ReplicationRunOnetimeActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Actions are immutable - re-execute on update
	var data ReplicationRunOnetimeActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
		if !data.Direction.IsNull() {
			params["direction"] = data.Direction.ValueString()
		}
		if !data.Transport.IsNull() {
			params["transport"] = data.Transport.ValueString()
		}
		if !data.SshCredentials.IsNull() {
			params["ssh_credentials"] = data.SshCredentials.ValueString()
		}
		if !data.NetcatActiveSide.IsNull() {
			params["netcat_active_side"] = data.NetcatActiveSide.ValueString()
		}
		if !data.NetcatActiveSideListenAddress.IsNull() {
			params["netcat_active_side_listen_address"] = data.NetcatActiveSideListenAddress.ValueString()
		}
		if !data.NetcatActiveSidePortMin.IsNull() {
			params["netcat_active_side_port_min"] = data.NetcatActiveSidePortMin.ValueString()
		}
		if !data.NetcatActiveSidePortMax.IsNull() {
			params["netcat_active_side_port_max"] = data.NetcatActiveSidePortMax.ValueString()
		}
		if !data.NetcatPassiveSideConnectAddress.IsNull() {
			params["netcat_passive_side_connect_address"] = data.NetcatPassiveSideConnectAddress.ValueString()
		}
		if !data.Sudo.IsNull() {
			params["sudo"] = data.Sudo.ValueBool()
		}
		if !data.TargetDataset.IsNull() {
			params["target_dataset"] = data.TargetDataset.ValueString()
		}
		if !data.Recursive.IsNull() {
			params["recursive"] = data.Recursive.ValueBool()
		}
		if !data.Properties.IsNull() {
			params["properties"] = data.Properties.ValueBool()
		}
		if !data.Replicate.IsNull() {
			params["replicate"] = data.Replicate.ValueBool()
		}
		if !data.Encryption.IsNull() {
			params["encryption"] = data.Encryption.ValueBool()
		}
		if !data.EncryptionInherit.IsNull() {
			params["encryption_inherit"] = data.EncryptionInherit.ValueString()
		}
		if !data.EncryptionKey.IsNull() {
			params["encryption_key"] = data.EncryptionKey.ValueString()
		}
		if !data.EncryptionKeyFormat.IsNull() {
			params["encryption_key_format"] = data.EncryptionKeyFormat.ValueString()
		}
		if !data.EncryptionKeyLocation.IsNull() {
			params["encryption_key_location"] = data.EncryptionKeyLocation.ValueString()
		}
		if !data.NameRegex.IsNull() {
			params["name_regex"] = data.NameRegex.ValueString()
		}
		if !data.RestrictSchedule.IsNull() {
			params["restrict_schedule"] = data.RestrictSchedule.ValueString()
		}
		if !data.AllowFromScratch.IsNull() {
			params["allow_from_scratch"] = data.AllowFromScratch.ValueBool()
		}
		if !data.Readonly.IsNull() {
			params["readonly"] = data.Readonly.ValueString()
		}
		if !data.HoldPendingSnapshots.IsNull() {
			params["hold_pending_snapshots"] = data.HoldPendingSnapshots.ValueBool()
		}
		if !data.RetentionPolicy.IsNull() {
			params["retention_policy"] = data.RetentionPolicy.ValueString()
		}
		if !data.LifetimeValue.IsNull() {
			params["lifetime_value"] = data.LifetimeValue.ValueString()
		}
		if !data.LifetimeUnit.IsNull() {
			params["lifetime_unit"] = data.LifetimeUnit.ValueString()
		}
		if !data.Compression.IsNull() {
			params["compression"] = data.Compression.ValueString()
		}
		if !data.SpeedLimit.IsNull() {
			params["speed_limit"] = data.SpeedLimit.ValueString()
		}
		if !data.LargeBlock.IsNull() {
			params["large_block"] = data.LargeBlock.ValueBool()
		}
		if !data.Embed.IsNull() {
			params["embed"] = data.Embed.ValueBool()
		}
		if !data.Compressed.IsNull() {
			params["compressed"] = data.Compressed.ValueBool()
		}
		if !data.Retries.IsNull() {
			params["retries"] = data.Retries.ValueInt64()
		}
		if !data.LoggingLevel.IsNull() {
			params["logging_level"] = data.LoggingLevel.ValueString()
		}
		if !data.ExcludeMountpointProperty.IsNull() {
			params["exclude_mountpoint_property"] = data.ExcludeMountpointProperty.ValueBool()
		}
		if !data.OnlyFromScratch.IsNull() {
			params["only_from_scratch"] = data.OnlyFromScratch.ValueBool()
		}
		if !data.Mount.IsNull() {
			params["mount"] = data.Mount.ValueBool()
		}

	_, err := r.client.Call("replication/run_onetime", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute run_onetime: %s", err.Error()))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationRunOnetimeActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Actions cannot be undone - just remove from state
}
