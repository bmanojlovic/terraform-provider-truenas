package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type ReplicationResource struct {
	client *client.Client
}

type ReplicationResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Direction                       types.String `tfsdk:"direction"`
	Transport                       types.String `tfsdk:"transport"`
	SshCredentials                  types.Int64  `tfsdk:"ssh_credentials"`
	NetcatActiveSide                types.String `tfsdk:"netcat_active_side"`
	NetcatActiveSideListenAddress   types.String `tfsdk:"netcat_active_side_listen_address"`
	NetcatActiveSidePortMin         types.Int64  `tfsdk:"netcat_active_side_port_min"`
	NetcatActiveSidePortMax         types.Int64  `tfsdk:"netcat_active_side_port_max"`
	NetcatPassiveSideConnectAddress types.String `tfsdk:"netcat_passive_side_connect_address"`
	Sudo                            types.Bool   `tfsdk:"sudo"`
	SourceDatasets                  types.List   `tfsdk:"source_datasets"`
	TargetDataset                   types.String `tfsdk:"target_dataset"`
	Recursive                       types.Bool   `tfsdk:"recursive"`
	Exclude                         types.List   `tfsdk:"exclude"`
	Properties                      types.Bool   `tfsdk:"properties"`
	PropertiesExclude               types.List   `tfsdk:"properties_exclude"`
	PropertiesOverride              types.String `tfsdk:"properties_override"`
	Replicate                       types.Bool   `tfsdk:"replicate"`
	Encryption                      types.Bool   `tfsdk:"encryption"`
	EncryptionInherit               types.Bool   `tfsdk:"encryption_inherit"`
	EncryptionKey                   types.String `tfsdk:"encryption_key"`
	EncryptionKeyFormat             types.String `tfsdk:"encryption_key_format"`
	EncryptionKeyLocation           types.String `tfsdk:"encryption_key_location"`
	PeriodicSnapshotTasks           types.List   `tfsdk:"periodic_snapshot_tasks"`
	NamingSchema                    types.List   `tfsdk:"naming_schema"`
	AlsoIncludeNamingSchema         types.List   `tfsdk:"also_include_naming_schema"`
	NameRegex                       types.String `tfsdk:"name_regex"`
	Auto                            types.Bool   `tfsdk:"auto"`
	Schedule                        types.String `tfsdk:"schedule"`
	RestrictSchedule                types.String `tfsdk:"restrict_schedule"`
	OnlyMatchingSchedule            types.Bool   `tfsdk:"only_matching_schedule"`
	AllowFromScratch                types.Bool   `tfsdk:"allow_from_scratch"`
	Readonly                        types.String `tfsdk:"readonly"`
	HoldPendingSnapshots            types.Bool   `tfsdk:"hold_pending_snapshots"`
	RetentionPolicy                 types.String `tfsdk:"retention_policy"`
	LifetimeValue                   types.Int64  `tfsdk:"lifetime_value"`
	LifetimeUnit                    types.String `tfsdk:"lifetime_unit"`
	Lifetimes                       types.List   `tfsdk:"lifetimes"`
	Compression                     types.String `tfsdk:"compression"`
	SpeedLimit                      types.Int64  `tfsdk:"speed_limit"`
	LargeBlock                      types.Bool   `tfsdk:"large_block"`
	Embed                           types.Bool   `tfsdk:"embed"`
	Compressed                      types.Bool   `tfsdk:"compressed"`
	Retries                         types.Int64  `tfsdk:"retries"`
	LoggingLevel                    types.String `tfsdk:"logging_level"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

func NewReplicationResource() resource.Resource {
	return &ReplicationResource{}
}

func (r *ReplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

func (r *ReplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ReplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Replication Task that will push or pull ZFS snapshots to or from remote host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Name for replication task.",
			},
			"direction": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Whether task will `PUSH` or `PULL` snapshots.",
			},
			"transport": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Method of snapshots transfer.  * `SSH` transfers snapshots via SSH connection. This method is suppor",
			},
			"ssh_credentials": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Keychain Credential ID of type `SSH_CREDENTIALS`.",
			},
			"netcat_active_side": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Which side actively establishes the netcat connection for `SSH+NETCAT` transport.  * `LOCAL`: Local ",
			},
			"netcat_active_side_listen_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"netcat_active_side_port_min": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum port number in the range for netcat connections. `null` if not applicable.",
			},
			"netcat_active_side_port_max": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum port number in the range for netcat connections. `null` if not applicable.",
			},
			"netcat_passive_side_connect_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"sudo": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "`SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs",
			},
			"source_datasets": schema.ListAttribute{
				Required:    true,
				Optional:    false,
				ElementType: types.StringType,
				Description: "List of datasets to replicate snapshots from.",
			},
			"target_dataset": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Dataset to put snapshots into.",
			},
			"recursive": schema.BoolAttribute{
				Required:    true,
				Optional:    false,
				Description: "Whether to recursively replicate child datasets.",
			},
			"exclude": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of dataset patterns to exclude from replication.",
			},
			"properties": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Send dataset properties along with snapshots.",
			},
			"properties_exclude": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of dataset property names to exclude from replication.",
			},
			"properties_override": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Object mapping dataset property names to override values during replication.",
			},
			"replicate": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to use full ZFS replication.",
			},
			"encryption": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable encryption for the replicated datasets.",
			},
			"encryption_inherit": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled.",
			},
			"encryption_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Encryption key for replicated datasets. `null` if not specified.",
			},
			"encryption_key_format": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Format of the encryption key.  * `HEX`: Hexadecimal-encoded key * `PASSPHRASE`: Text passphrase * `n",
			},
			"encryption_key_location": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Filesystem path where encryption key is stored. `null` if not using key file.",
			},
			"periodic_snapshot_tasks": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of periodic snapshot task IDs that are sources of snapshots for this replication task. Only pus",
			},
			"naming_schema": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of naming schemas for pull replication.",
			},
			"also_include_naming_schema": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of naming schemas for push replication.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Replicate all snapshots which names match specified regular expression.",
			},
			"auto": schema.BoolAttribute{
				Required:    true,
				Optional:    false,
				Description: "Allow replication to run automatically on schedule or after bound periodic snapshot task.",
			},
			"schedule": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot task",
			},
			"restrict_schedule": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have p",
			},
			"only_matching_schedule": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Will only replicate snapshots that match `schedule` or `restrict_schedule`.",
			},
			"allow_from_scratch": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Will destroy all snapshots on target side and replicate everything from scratch if none of the snaps",
			},
			"readonly": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Controls destination datasets readonly property.  * `SET`: Set all destination datasets to readonly=",
			},
			"hold_pending_snapshots": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Prevent source snapshots from being deleted by retention of replication fails for some reason.",
			},
			"retention_policy": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "How to delete old snapshots on target side:  * `SOURCE`: Delete snapshots that are absent on source ",
			},
			"lifetime_value": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of time units to retain snapshots for custom retention policy. Only applies when `retention_p",
			},
			"lifetime_unit": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` i",
			},
			"lifetimes": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of different retention schedules with their own cron schedules and lifetime settings.",
			},
			"compression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Compresses SSH stream. Available only for SSH transport.",
			},
			"speed_limit": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Limits speed of SSH stream. Available only for SSH transport.",
			},
			"large_block": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable large block support for ZFS send streams.",
			},
			"embed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable embedded block support for ZFS send streams.",
			},
			"compressed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable compressed ZFS send streams.",
			},
			"retries": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of retries before considering replication failed.",
			},
			"logging_level": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Log level for replication task execution. Controls verbosity of replication logs.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this replication task is enabled.",
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

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.Transport.IsNull() && !data.Transport.IsUnknown() {
		params["transport"] = data.Transport.ValueString()
	}
	if !data.SshCredentials.IsNull() && !data.SshCredentials.IsUnknown() {
		params["ssh_credentials"] = data.SshCredentials.ValueInt64()
	}
	if !data.NetcatActiveSide.IsNull() && !data.NetcatActiveSide.IsUnknown() {
		params["netcat_active_side"] = data.NetcatActiveSide.ValueString()
	}
	if !data.NetcatActiveSideListenAddress.IsNull() && !data.NetcatActiveSideListenAddress.IsUnknown() {
		params["netcat_active_side_listen_address"] = data.NetcatActiveSideListenAddress.ValueString()
	}
	if !data.NetcatActiveSidePortMin.IsNull() && !data.NetcatActiveSidePortMin.IsUnknown() {
		params["netcat_active_side_port_min"] = data.NetcatActiveSidePortMin.ValueInt64()
	}
	if !data.NetcatActiveSidePortMax.IsNull() && !data.NetcatActiveSidePortMax.IsUnknown() {
		params["netcat_active_side_port_max"] = data.NetcatActiveSidePortMax.ValueInt64()
	}
	if !data.NetcatPassiveSideConnectAddress.IsNull() && !data.NetcatPassiveSideConnectAddress.IsUnknown() {
		params["netcat_passive_side_connect_address"] = data.NetcatPassiveSideConnectAddress.ValueString()
	}
	if !data.Sudo.IsNull() && !data.Sudo.IsUnknown() {
		params["sudo"] = data.Sudo.ValueBool()
	}
	if !data.SourceDatasets.IsNull() && !data.SourceDatasets.IsUnknown() {
		var source_datasetsList []string
		data.SourceDatasets.ElementsAs(ctx, &source_datasetsList, false)
		params["source_datasets"] = source_datasetsList
	}
	if !data.TargetDataset.IsNull() && !data.TargetDataset.IsUnknown() {
		params["target_dataset"] = data.TargetDataset.ValueString()
	}
	if !data.Recursive.IsNull() && !data.Recursive.IsUnknown() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Exclude.IsNull() && !data.Exclude.IsUnknown() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Properties.IsNull() && !data.Properties.IsUnknown() {
		params["properties"] = data.Properties.ValueBool()
	}
	if !data.PropertiesExclude.IsNull() && !data.PropertiesExclude.IsUnknown() {
		var properties_excludeList []string
		data.PropertiesExclude.ElementsAs(ctx, &properties_excludeList, false)
		params["properties_exclude"] = properties_excludeList
	}
	if !data.PropertiesOverride.IsNull() && !data.PropertiesOverride.IsUnknown() {
		var properties_overrideObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.PropertiesOverride.ValueString()), &properties_overrideObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse properties_override: %s", err))
			return
		}
		params["properties_override"] = properties_overrideObj
	}
	if !data.Replicate.IsNull() && !data.Replicate.IsUnknown() {
		params["replicate"] = data.Replicate.ValueBool()
	}
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.EncryptionInherit.IsNull() && !data.EncryptionInherit.IsUnknown() {
		params["encryption_inherit"] = data.EncryptionInherit.ValueBool()
	}
	if !data.EncryptionKey.IsNull() && !data.EncryptionKey.IsUnknown() {
		params["encryption_key"] = data.EncryptionKey.ValueString()
	}
	if !data.EncryptionKeyFormat.IsNull() && !data.EncryptionKeyFormat.IsUnknown() {
		params["encryption_key_format"] = data.EncryptionKeyFormat.ValueString()
	}
	if !data.EncryptionKeyLocation.IsNull() && !data.EncryptionKeyLocation.IsUnknown() {
		params["encryption_key_location"] = data.EncryptionKeyLocation.ValueString()
	}
	if !data.PeriodicSnapshotTasks.IsNull() && !data.PeriodicSnapshotTasks.IsUnknown() {
		var periodic_snapshot_tasksList []string
		data.PeriodicSnapshotTasks.ElementsAs(ctx, &periodic_snapshot_tasksList, false)
		params["periodic_snapshot_tasks"] = periodic_snapshot_tasksList
	}
	if !data.NamingSchema.IsNull() && !data.NamingSchema.IsUnknown() {
		var naming_schemaList []string
		data.NamingSchema.ElementsAs(ctx, &naming_schemaList, false)
		params["naming_schema"] = naming_schemaList
	}
	if !data.AlsoIncludeNamingSchema.IsNull() && !data.AlsoIncludeNamingSchema.IsUnknown() {
		var also_include_naming_schemaList []string
		data.AlsoIncludeNamingSchema.ElementsAs(ctx, &also_include_naming_schemaList, false)
		params["also_include_naming_schema"] = also_include_naming_schemaList
	}
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		params["name_regex"] = data.NameRegex.ValueString()
	}
	if !data.Auto.IsNull() && !data.Auto.IsUnknown() {
		params["auto"] = data.Auto.ValueBool()
	}
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.RestrictSchedule.IsNull() && !data.RestrictSchedule.IsUnknown() {
		var restrict_scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.RestrictSchedule.ValueString()), &restrict_scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse restrict_schedule: %s", err))
			return
		}
		params["restrict_schedule"] = restrict_scheduleObj
	}
	if !data.OnlyMatchingSchedule.IsNull() && !data.OnlyMatchingSchedule.IsUnknown() {
		params["only_matching_schedule"] = data.OnlyMatchingSchedule.ValueBool()
	}
	if !data.AllowFromScratch.IsNull() && !data.AllowFromScratch.IsUnknown() {
		params["allow_from_scratch"] = data.AllowFromScratch.ValueBool()
	}
	if !data.Readonly.IsNull() && !data.Readonly.IsUnknown() {
		params["readonly"] = data.Readonly.ValueString()
	}
	if !data.HoldPendingSnapshots.IsNull() && !data.HoldPendingSnapshots.IsUnknown() {
		params["hold_pending_snapshots"] = data.HoldPendingSnapshots.ValueBool()
	}
	if !data.RetentionPolicy.IsNull() && !data.RetentionPolicy.IsUnknown() {
		params["retention_policy"] = data.RetentionPolicy.ValueString()
	}
	if !data.LifetimeValue.IsNull() && !data.LifetimeValue.IsUnknown() {
		params["lifetime_value"] = data.LifetimeValue.ValueInt64()
	}
	if !data.LifetimeUnit.IsNull() && !data.LifetimeUnit.IsUnknown() {
		params["lifetime_unit"] = data.LifetimeUnit.ValueString()
	}
	if !data.Lifetimes.IsNull() && !data.Lifetimes.IsUnknown() {
		var lifetimesList []string
		data.Lifetimes.ElementsAs(ctx, &lifetimesList, false)
		var lifetimesObjs []map[string]interface{}
		for _, jsonStr := range lifetimesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse lifetimes item: %s", err))
				return
			}
			lifetimesObjs = append(lifetimesObjs, obj)
		}
		params["lifetimes"] = lifetimesObjs
	}
	if !data.Compression.IsNull() && !data.Compression.IsUnknown() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.SpeedLimit.IsNull() && !data.SpeedLimit.IsUnknown() {
		params["speed_limit"] = data.SpeedLimit.ValueInt64()
	}
	if !data.LargeBlock.IsNull() && !data.LargeBlock.IsUnknown() {
		params["large_block"] = data.LargeBlock.ValueBool()
	}
	if !data.Embed.IsNull() && !data.Embed.IsUnknown() {
		params["embed"] = data.Embed.ValueBool()
	}
	if !data.Compressed.IsNull() && !data.Compressed.IsUnknown() {
		params["compressed"] = data.Compressed.ValueBool()
	}
	if !data.Retries.IsNull() && !data.Retries.IsUnknown() {
		params["retries"] = data.Retries.ValueInt64()
	}
	if !data.LoggingLevel.IsNull() && !data.LoggingLevel.IsUnknown() {
		params["logging_level"] = data.LoggingLevel.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	result, err := r.client.Call("replication.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create replication: %s", err))
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
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("replication.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back replication: %s", err))
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
	if v, ok := resultMap["direction"]; ok {
		switch val := v.(type) {
		case string:
			data.Direction = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Direction = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transport"]; ok {
		switch val := v.(type) {
		case string:
			data.Transport = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Transport = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Transport = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["ssh_credentials"]; ok {
		if v == nil {
			data.SshCredentials = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.SshCredentials = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.SshCredentials = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_active_side"]; ok {
		if v == nil {
			data.NetcatActiveSide = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatActiveSide = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_listen_address"]; ok {
		if v == nil {
			data.NetcatActiveSideListenAddress = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatActiveSideListenAddress = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_port_min"]; ok {
		if v == nil {
			data.NetcatActiveSidePortMin = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.NetcatActiveSidePortMin = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.NetcatActiveSidePortMin = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_port_max"]; ok {
		if v == nil {
			data.NetcatActiveSidePortMax = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.NetcatActiveSidePortMax = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.NetcatActiveSidePortMax = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_passive_side_connect_address"]; ok {
		if v == nil {
			data.NetcatPassiveSideConnectAddress = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatPassiveSideConnectAddress = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["source_datasets"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.SourceDatasets, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["target_dataset"]; ok {
		switch val := v.(type) {
		case string:
			data.TargetDataset = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.TargetDataset = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.TargetDataset = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["recursive"]; ok {
		if bv, ok := v.(bool); ok {
			data.Recursive = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["encryption_inherit"]; ok {
		if v == nil {
			data.EncryptionInherit = types.BoolNull()
		} else {
			if bv, ok := v.(bool); ok {
				data.EncryptionInherit = types.BoolValue(bv)
			}
		}
	}
	if v, ok := resultMap["encryption_key"]; ok {
		if v == nil {
			data.EncryptionKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["encryption_key_format"]; ok {
		if v == nil {
			data.EncryptionKeyFormat = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKeyFormat = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["encryption_key_location"]; ok {
		if v == nil {
			data.EncryptionKeyLocation = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKeyLocation = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["name_regex"]; ok {
		if v == nil {
			data.NameRegex = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NameRegex = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NameRegex = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NameRegex = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["auto"]; ok {
		if bv, ok := v.(bool); ok {
			data.Auto = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["schedule"]; ok {
		if v == nil {
			data.Schedule = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Schedule = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Schedule = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["restrict_schedule"]; ok {
		if v == nil {
			data.RestrictSchedule = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.RestrictSchedule = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["retention_policy"]; ok {
		switch val := v.(type) {
		case string:
			data.RetentionPolicy = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["lifetime_value"]; ok {
		if v == nil {
			data.LifetimeValue = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.LifetimeValue = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.LifetimeValue = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["lifetime_unit"]; ok {
		if v == nil {
			data.LifetimeUnit = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LifetimeUnit = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["compression"]; ok {
		if v == nil {
			data.Compression = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Compression = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Compression = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Compression = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["speed_limit"]; ok {
		if v == nil {
			data.SpeedLimit = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.SpeedLimit = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.SpeedLimit = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["logging_level"]; ok {
		if v == nil {
			data.LoggingLevel = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LoggingLevel = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.SshCredentials.IsUnknown() {
		data.SshCredentials = types.Int64Null()
	}
	if data.NetcatActiveSide.IsUnknown() {
		data.NetcatActiveSide = types.StringNull()
	}
	if data.NetcatActiveSideListenAddress.IsUnknown() {
		data.NetcatActiveSideListenAddress = types.StringNull()
	}
	if data.NetcatActiveSidePortMin.IsUnknown() {
		data.NetcatActiveSidePortMin = types.Int64Null()
	}
	if data.NetcatActiveSidePortMax.IsUnknown() {
		data.NetcatActiveSidePortMax = types.Int64Null()
	}
	if data.NetcatPassiveSideConnectAddress.IsUnknown() {
		data.NetcatPassiveSideConnectAddress = types.StringNull()
	}
	if data.Sudo.IsUnknown() {
		data.Sudo = types.BoolNull()
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Properties.IsUnknown() {
		data.Properties = types.BoolNull()
	}
	if data.PropertiesExclude.IsUnknown() {
		data.PropertiesExclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.PropertiesOverride.IsUnknown() {
		data.PropertiesOverride = types.StringNull()
	}
	if data.Replicate.IsUnknown() {
		data.Replicate = types.BoolNull()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.EncryptionInherit.IsUnknown() {
		data.EncryptionInherit = types.BoolNull()
	}
	if data.EncryptionKey.IsUnknown() {
		data.EncryptionKey = types.StringNull()
	}
	if data.EncryptionKeyFormat.IsUnknown() {
		data.EncryptionKeyFormat = types.StringNull()
	}
	if data.EncryptionKeyLocation.IsUnknown() {
		data.EncryptionKeyLocation = types.StringNull()
	}
	if data.PeriodicSnapshotTasks.IsUnknown() {
		data.PeriodicSnapshotTasks, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.NamingSchema.IsUnknown() {
		data.NamingSchema, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.AlsoIncludeNamingSchema.IsUnknown() {
		data.AlsoIncludeNamingSchema, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.NameRegex.IsUnknown() {
		data.NameRegex = types.StringNull()
	}
	if data.Schedule.IsUnknown() {
		data.Schedule = types.StringNull()
	}
	if data.RestrictSchedule.IsUnknown() {
		data.RestrictSchedule = types.StringNull()
	}
	if data.OnlyMatchingSchedule.IsUnknown() {
		data.OnlyMatchingSchedule = types.BoolNull()
	}
	if data.AllowFromScratch.IsUnknown() {
		data.AllowFromScratch = types.BoolNull()
	}
	if data.Readonly.IsUnknown() {
		data.Readonly = types.StringNull()
	}
	if data.HoldPendingSnapshots.IsUnknown() {
		data.HoldPendingSnapshots = types.BoolNull()
	}
	if data.LifetimeValue.IsUnknown() {
		data.LifetimeValue = types.Int64Null()
	}
	if data.LifetimeUnit.IsUnknown() {
		data.LifetimeUnit = types.StringNull()
	}
	if data.Lifetimes.IsUnknown() {
		data.Lifetimes, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Compression.IsUnknown() {
		data.Compression = types.StringNull()
	}
	if data.SpeedLimit.IsUnknown() {
		data.SpeedLimit = types.Int64Null()
	}
	if data.LargeBlock.IsUnknown() {
		data.LargeBlock = types.BoolNull()
	}
	if data.Embed.IsUnknown() {
		data.Embed = types.BoolNull()
	}
	if data.Compressed.IsUnknown() {
		data.Compressed = types.BoolNull()
	}
	if data.Retries.IsUnknown() {
		data.Retries = types.Int64Null()
	}
	if data.LoggingLevel.IsUnknown() {
		data.LoggingLevel = types.StringNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ReplicationResourceModel
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

	result, err := r.client.Call("replication.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read replication: %s", err))
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
	if v, ok := resultMap["direction"]; ok {
		switch val := v.(type) {
		case string:
			data.Direction = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Direction = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transport"]; ok {
		switch val := v.(type) {
		case string:
			data.Transport = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Transport = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Transport = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["ssh_credentials"]; ok {
		if v == nil {
			data.SshCredentials = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.SshCredentials = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.SshCredentials = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_active_side"]; ok {
		if v == nil {
			data.NetcatActiveSide = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatActiveSide = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_listen_address"]; ok {
		if v == nil {
			data.NetcatActiveSideListenAddress = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatActiveSideListenAddress = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_port_min"]; ok {
		if v == nil {
			data.NetcatActiveSidePortMin = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.NetcatActiveSidePortMin = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.NetcatActiveSidePortMin = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_port_max"]; ok {
		if v == nil {
			data.NetcatActiveSidePortMax = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.NetcatActiveSidePortMax = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.NetcatActiveSidePortMax = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_passive_side_connect_address"]; ok {
		if v == nil {
			data.NetcatPassiveSideConnectAddress = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatPassiveSideConnectAddress = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["source_datasets"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.SourceDatasets, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["target_dataset"]; ok {
		switch val := v.(type) {
		case string:
			data.TargetDataset = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.TargetDataset = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.TargetDataset = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["recursive"]; ok {
		if bv, ok := v.(bool); ok {
			data.Recursive = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["encryption_inherit"]; ok {
		if v == nil {
			data.EncryptionInherit = types.BoolNull()
		} else {
			if bv, ok := v.(bool); ok {
				data.EncryptionInherit = types.BoolValue(bv)
			}
		}
	}
	if v, ok := resultMap["encryption_key"]; ok {
		if v == nil {
			data.EncryptionKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["encryption_key_format"]; ok {
		if v == nil {
			data.EncryptionKeyFormat = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKeyFormat = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["encryption_key_location"]; ok {
		if v == nil {
			data.EncryptionKeyLocation = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKeyLocation = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["name_regex"]; ok {
		if v == nil {
			data.NameRegex = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NameRegex = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NameRegex = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NameRegex = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["auto"]; ok {
		if bv, ok := v.(bool); ok {
			data.Auto = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["schedule"]; ok {
		if v == nil {
			data.Schedule = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Schedule = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Schedule = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["restrict_schedule"]; ok {
		if v == nil {
			data.RestrictSchedule = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.RestrictSchedule = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["retention_policy"]; ok {
		switch val := v.(type) {
		case string:
			data.RetentionPolicy = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["lifetime_value"]; ok {
		if v == nil {
			data.LifetimeValue = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.LifetimeValue = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.LifetimeValue = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["lifetime_unit"]; ok {
		if v == nil {
			data.LifetimeUnit = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LifetimeUnit = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["compression"]; ok {
		if v == nil {
			data.Compression = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Compression = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Compression = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Compression = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["speed_limit"]; ok {
		if v == nil {
			data.SpeedLimit = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.SpeedLimit = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.SpeedLimit = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["logging_level"]; ok {
		if v == nil {
			data.LoggingLevel = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LoggingLevel = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.SshCredentials.IsUnknown() {
		data.SshCredentials = types.Int64Null()
	}
	if data.NetcatActiveSide.IsUnknown() {
		data.NetcatActiveSide = types.StringNull()
	}
	if data.NetcatActiveSideListenAddress.IsUnknown() {
		data.NetcatActiveSideListenAddress = types.StringNull()
	}
	if data.NetcatActiveSidePortMin.IsUnknown() {
		data.NetcatActiveSidePortMin = types.Int64Null()
	}
	if data.NetcatActiveSidePortMax.IsUnknown() {
		data.NetcatActiveSidePortMax = types.Int64Null()
	}
	if data.NetcatPassiveSideConnectAddress.IsUnknown() {
		data.NetcatPassiveSideConnectAddress = types.StringNull()
	}
	if data.Sudo.IsUnknown() {
		data.Sudo = types.BoolNull()
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Properties.IsUnknown() {
		data.Properties = types.BoolNull()
	}
	if data.PropertiesExclude.IsUnknown() {
		data.PropertiesExclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.PropertiesOverride.IsUnknown() {
		data.PropertiesOverride = types.StringNull()
	}
	if data.Replicate.IsUnknown() {
		data.Replicate = types.BoolNull()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.EncryptionInherit.IsUnknown() {
		data.EncryptionInherit = types.BoolNull()
	}
	if data.EncryptionKey.IsUnknown() {
		data.EncryptionKey = types.StringNull()
	}
	if data.EncryptionKeyFormat.IsUnknown() {
		data.EncryptionKeyFormat = types.StringNull()
	}
	if data.EncryptionKeyLocation.IsUnknown() {
		data.EncryptionKeyLocation = types.StringNull()
	}
	if data.PeriodicSnapshotTasks.IsUnknown() {
		data.PeriodicSnapshotTasks, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.NamingSchema.IsUnknown() {
		data.NamingSchema, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.AlsoIncludeNamingSchema.IsUnknown() {
		data.AlsoIncludeNamingSchema, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.NameRegex.IsUnknown() {
		data.NameRegex = types.StringNull()
	}
	if data.Schedule.IsUnknown() {
		data.Schedule = types.StringNull()
	}
	if data.RestrictSchedule.IsUnknown() {
		data.RestrictSchedule = types.StringNull()
	}
	if data.OnlyMatchingSchedule.IsUnknown() {
		data.OnlyMatchingSchedule = types.BoolNull()
	}
	if data.AllowFromScratch.IsUnknown() {
		data.AllowFromScratch = types.BoolNull()
	}
	if data.Readonly.IsUnknown() {
		data.Readonly = types.StringNull()
	}
	if data.HoldPendingSnapshots.IsUnknown() {
		data.HoldPendingSnapshots = types.BoolNull()
	}
	if data.LifetimeValue.IsUnknown() {
		data.LifetimeValue = types.Int64Null()
	}
	if data.LifetimeUnit.IsUnknown() {
		data.LifetimeUnit = types.StringNull()
	}
	if data.Lifetimes.IsUnknown() {
		data.Lifetimes, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Compression.IsUnknown() {
		data.Compression = types.StringNull()
	}
	if data.SpeedLimit.IsUnknown() {
		data.SpeedLimit = types.Int64Null()
	}
	if data.LargeBlock.IsUnknown() {
		data.LargeBlock = types.BoolNull()
	}
	if data.Embed.IsUnknown() {
		data.Embed = types.BoolNull()
	}
	if data.Compressed.IsUnknown() {
		data.Compressed = types.BoolNull()
	}
	if data.Retries.IsUnknown() {
		data.Retries = types.Int64Null()
	}
	if data.LoggingLevel.IsUnknown() {
		data.LoggingLevel = types.StringNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ReplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ReplicationResourceModel
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
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.Transport.IsNull() && !data.Transport.IsUnknown() {
		params["transport"] = data.Transport.ValueString()
	}
	if !data.SshCredentials.IsNull() && !data.SshCredentials.IsUnknown() {
		params["ssh_credentials"] = data.SshCredentials.ValueInt64()
	}
	if !data.NetcatActiveSide.IsNull() && !data.NetcatActiveSide.IsUnknown() {
		params["netcat_active_side"] = data.NetcatActiveSide.ValueString()
	}
	if !data.NetcatActiveSideListenAddress.IsNull() && !data.NetcatActiveSideListenAddress.IsUnknown() {
		params["netcat_active_side_listen_address"] = data.NetcatActiveSideListenAddress.ValueString()
	}
	if !data.NetcatActiveSidePortMin.IsNull() && !data.NetcatActiveSidePortMin.IsUnknown() {
		params["netcat_active_side_port_min"] = data.NetcatActiveSidePortMin.ValueInt64()
	}
	if !data.NetcatActiveSidePortMax.IsNull() && !data.NetcatActiveSidePortMax.IsUnknown() {
		params["netcat_active_side_port_max"] = data.NetcatActiveSidePortMax.ValueInt64()
	}
	if !data.NetcatPassiveSideConnectAddress.IsNull() && !data.NetcatPassiveSideConnectAddress.IsUnknown() {
		params["netcat_passive_side_connect_address"] = data.NetcatPassiveSideConnectAddress.ValueString()
	}
	if !data.Sudo.IsNull() && !data.Sudo.IsUnknown() {
		params["sudo"] = data.Sudo.ValueBool()
	}
	if !data.SourceDatasets.IsNull() && !data.SourceDatasets.IsUnknown() {
		var source_datasetsList []string
		data.SourceDatasets.ElementsAs(ctx, &source_datasetsList, false)
		params["source_datasets"] = source_datasetsList
	}
	if !data.TargetDataset.IsNull() && !data.TargetDataset.IsUnknown() {
		params["target_dataset"] = data.TargetDataset.ValueString()
	}
	if !data.Recursive.IsNull() && !data.Recursive.IsUnknown() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Exclude.IsNull() && !data.Exclude.IsUnknown() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Properties.IsNull() && !data.Properties.IsUnknown() {
		params["properties"] = data.Properties.ValueBool()
	}
	if !data.PropertiesExclude.IsNull() && !data.PropertiesExclude.IsUnknown() {
		var properties_excludeList []string
		data.PropertiesExclude.ElementsAs(ctx, &properties_excludeList, false)
		params["properties_exclude"] = properties_excludeList
	}
	if !data.PropertiesOverride.IsNull() && !data.PropertiesOverride.IsUnknown() {
		var properties_overrideObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.PropertiesOverride.ValueString()), &properties_overrideObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse properties_override: %s", err))
			return
		}
		params["properties_override"] = properties_overrideObj
	}
	if !data.Replicate.IsNull() && !data.Replicate.IsUnknown() {
		params["replicate"] = data.Replicate.ValueBool()
	}
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.EncryptionInherit.IsNull() && !data.EncryptionInherit.IsUnknown() {
		params["encryption_inherit"] = data.EncryptionInherit.ValueBool()
	}
	if !data.EncryptionKey.IsNull() && !data.EncryptionKey.IsUnknown() {
		params["encryption_key"] = data.EncryptionKey.ValueString()
	}
	if !data.EncryptionKeyFormat.IsNull() && !data.EncryptionKeyFormat.IsUnknown() {
		params["encryption_key_format"] = data.EncryptionKeyFormat.ValueString()
	}
	if !data.EncryptionKeyLocation.IsNull() && !data.EncryptionKeyLocation.IsUnknown() {
		params["encryption_key_location"] = data.EncryptionKeyLocation.ValueString()
	}
	if !data.PeriodicSnapshotTasks.IsNull() && !data.PeriodicSnapshotTasks.IsUnknown() {
		var periodic_snapshot_tasksList []string
		data.PeriodicSnapshotTasks.ElementsAs(ctx, &periodic_snapshot_tasksList, false)
		params["periodic_snapshot_tasks"] = periodic_snapshot_tasksList
	}
	if !data.NamingSchema.IsNull() && !data.NamingSchema.IsUnknown() {
		var naming_schemaList []string
		data.NamingSchema.ElementsAs(ctx, &naming_schemaList, false)
		params["naming_schema"] = naming_schemaList
	}
	if !data.AlsoIncludeNamingSchema.IsNull() && !data.AlsoIncludeNamingSchema.IsUnknown() {
		var also_include_naming_schemaList []string
		data.AlsoIncludeNamingSchema.ElementsAs(ctx, &also_include_naming_schemaList, false)
		params["also_include_naming_schema"] = also_include_naming_schemaList
	}
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		params["name_regex"] = data.NameRegex.ValueString()
	}
	if !data.Auto.IsNull() && !data.Auto.IsUnknown() {
		params["auto"] = data.Auto.ValueBool()
	}
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.RestrictSchedule.IsNull() && !data.RestrictSchedule.IsUnknown() {
		var restrict_scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.RestrictSchedule.ValueString()), &restrict_scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse restrict_schedule: %s", err))
			return
		}
		params["restrict_schedule"] = restrict_scheduleObj
	}
	if !data.OnlyMatchingSchedule.IsNull() && !data.OnlyMatchingSchedule.IsUnknown() {
		params["only_matching_schedule"] = data.OnlyMatchingSchedule.ValueBool()
	}
	if !data.AllowFromScratch.IsNull() && !data.AllowFromScratch.IsUnknown() {
		params["allow_from_scratch"] = data.AllowFromScratch.ValueBool()
	}
	if !data.Readonly.IsNull() && !data.Readonly.IsUnknown() {
		params["readonly"] = data.Readonly.ValueString()
	}
	if !data.HoldPendingSnapshots.IsNull() && !data.HoldPendingSnapshots.IsUnknown() {
		params["hold_pending_snapshots"] = data.HoldPendingSnapshots.ValueBool()
	}
	if !data.RetentionPolicy.IsNull() && !data.RetentionPolicy.IsUnknown() {
		params["retention_policy"] = data.RetentionPolicy.ValueString()
	}
	if !data.LifetimeValue.IsNull() && !data.LifetimeValue.IsUnknown() {
		params["lifetime_value"] = data.LifetimeValue.ValueInt64()
	}
	if !data.LifetimeUnit.IsNull() && !data.LifetimeUnit.IsUnknown() {
		params["lifetime_unit"] = data.LifetimeUnit.ValueString()
	}
	if !data.Lifetimes.IsNull() && !data.Lifetimes.IsUnknown() {
		var lifetimesList []string
		data.Lifetimes.ElementsAs(ctx, &lifetimesList, false)
		var lifetimesObjs []map[string]interface{}
		for _, jsonStr := range lifetimesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse lifetimes item: %s", err))
				return
			}
			lifetimesObjs = append(lifetimesObjs, obj)
		}
		params["lifetimes"] = lifetimesObjs
	}
	if !data.Compression.IsNull() && !data.Compression.IsUnknown() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.SpeedLimit.IsNull() && !data.SpeedLimit.IsUnknown() {
		params["speed_limit"] = data.SpeedLimit.ValueInt64()
	}
	if !data.LargeBlock.IsNull() && !data.LargeBlock.IsUnknown() {
		params["large_block"] = data.LargeBlock.ValueBool()
	}
	if !data.Embed.IsNull() && !data.Embed.IsUnknown() {
		params["embed"] = data.Embed.ValueBool()
	}
	if !data.Compressed.IsNull() && !data.Compressed.IsUnknown() {
		params["compressed"] = data.Compressed.ValueBool()
	}
	if !data.Retries.IsNull() && !data.Retries.IsUnknown() {
		params["retries"] = data.Retries.ValueInt64()
	}
	if !data.LoggingLevel.IsNull() && !data.LoggingLevel.IsUnknown() {
		params["logging_level"] = data.LoggingLevel.ValueString()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	_, err = r.client.Call("replication.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update replication: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("replication.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back replication: %s", readErr))
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
	if v, ok := resultMap["direction"]; ok {
		switch val := v.(type) {
		case string:
			data.Direction = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Direction = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["transport"]; ok {
		switch val := v.(type) {
		case string:
			data.Transport = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Transport = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Transport = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["ssh_credentials"]; ok {
		if v == nil {
			data.SshCredentials = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.SshCredentials = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.SshCredentials = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_active_side"]; ok {
		if v == nil {
			data.NetcatActiveSide = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatActiveSide = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_listen_address"]; ok {
		if v == nil {
			data.NetcatActiveSideListenAddress = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatActiveSideListenAddress = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_port_min"]; ok {
		if v == nil {
			data.NetcatActiveSidePortMin = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.NetcatActiveSidePortMin = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.NetcatActiveSidePortMin = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_active_side_port_max"]; ok {
		if v == nil {
			data.NetcatActiveSidePortMax = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.NetcatActiveSidePortMax = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.NetcatActiveSidePortMax = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["netcat_passive_side_connect_address"]; ok {
		if v == nil {
			data.NetcatPassiveSideConnectAddress = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NetcatPassiveSideConnectAddress = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["source_datasets"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.SourceDatasets, _ = types.ListValue(types.StringType, strVals)
		}
	}
	if v, ok := resultMap["target_dataset"]; ok {
		switch val := v.(type) {
		case string:
			data.TargetDataset = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.TargetDataset = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.TargetDataset = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["recursive"]; ok {
		if bv, ok := v.(bool); ok {
			data.Recursive = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["encryption_inherit"]; ok {
		if v == nil {
			data.EncryptionInherit = types.BoolNull()
		} else {
			if bv, ok := v.(bool); ok {
				data.EncryptionInherit = types.BoolValue(bv)
			}
		}
	}
	if v, ok := resultMap["encryption_key"]; ok {
		if v == nil {
			data.EncryptionKey = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKey = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["encryption_key_format"]; ok {
		if v == nil {
			data.EncryptionKeyFormat = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKeyFormat = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["encryption_key_location"]; ok {
		if v == nil {
			data.EncryptionKeyLocation = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.EncryptionKeyLocation = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["name_regex"]; ok {
		if v == nil {
			data.NameRegex = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.NameRegex = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.NameRegex = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.NameRegex = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["auto"]; ok {
		if bv, ok := v.(bool); ok {
			data.Auto = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["schedule"]; ok {
		if v == nil {
			data.Schedule = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Schedule = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Schedule = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["restrict_schedule"]; ok {
		if v == nil {
			data.RestrictSchedule = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.RestrictSchedule = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["retention_policy"]; ok {
		switch val := v.(type) {
		case string:
			data.RetentionPolicy = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["lifetime_value"]; ok {
		if v == nil {
			data.LifetimeValue = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.LifetimeValue = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.LifetimeValue = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["lifetime_unit"]; ok {
		if v == nil {
			data.LifetimeUnit = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LifetimeUnit = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["compression"]; ok {
		if v == nil {
			data.Compression = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Compression = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Compression = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Compression = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["speed_limit"]; ok {
		if v == nil {
			data.SpeedLimit = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.SpeedLimit = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.SpeedLimit = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["logging_level"]; ok {
		if v == nil {
			data.LoggingLevel = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.LoggingLevel = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.SshCredentials.IsUnknown() {
		data.SshCredentials = types.Int64Null()
	}
	if data.NetcatActiveSide.IsUnknown() {
		data.NetcatActiveSide = types.StringNull()
	}
	if data.NetcatActiveSideListenAddress.IsUnknown() {
		data.NetcatActiveSideListenAddress = types.StringNull()
	}
	if data.NetcatActiveSidePortMin.IsUnknown() {
		data.NetcatActiveSidePortMin = types.Int64Null()
	}
	if data.NetcatActiveSidePortMax.IsUnknown() {
		data.NetcatActiveSidePortMax = types.Int64Null()
	}
	if data.NetcatPassiveSideConnectAddress.IsUnknown() {
		data.NetcatPassiveSideConnectAddress = types.StringNull()
	}
	if data.Sudo.IsUnknown() {
		data.Sudo = types.BoolNull()
	}
	if data.Exclude.IsUnknown() {
		data.Exclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Properties.IsUnknown() {
		data.Properties = types.BoolNull()
	}
	if data.PropertiesExclude.IsUnknown() {
		data.PropertiesExclude, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.PropertiesOverride.IsUnknown() {
		data.PropertiesOverride = types.StringNull()
	}
	if data.Replicate.IsUnknown() {
		data.Replicate = types.BoolNull()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.EncryptionInherit.IsUnknown() {
		data.EncryptionInherit = types.BoolNull()
	}
	if data.EncryptionKey.IsUnknown() {
		data.EncryptionKey = types.StringNull()
	}
	if data.EncryptionKeyFormat.IsUnknown() {
		data.EncryptionKeyFormat = types.StringNull()
	}
	if data.EncryptionKeyLocation.IsUnknown() {
		data.EncryptionKeyLocation = types.StringNull()
	}
	if data.PeriodicSnapshotTasks.IsUnknown() {
		data.PeriodicSnapshotTasks, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.NamingSchema.IsUnknown() {
		data.NamingSchema, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.AlsoIncludeNamingSchema.IsUnknown() {
		data.AlsoIncludeNamingSchema, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.NameRegex.IsUnknown() {
		data.NameRegex = types.StringNull()
	}
	if data.Schedule.IsUnknown() {
		data.Schedule = types.StringNull()
	}
	if data.RestrictSchedule.IsUnknown() {
		data.RestrictSchedule = types.StringNull()
	}
	if data.OnlyMatchingSchedule.IsUnknown() {
		data.OnlyMatchingSchedule = types.BoolNull()
	}
	if data.AllowFromScratch.IsUnknown() {
		data.AllowFromScratch = types.BoolNull()
	}
	if data.Readonly.IsUnknown() {
		data.Readonly = types.StringNull()
	}
	if data.HoldPendingSnapshots.IsUnknown() {
		data.HoldPendingSnapshots = types.BoolNull()
	}
	if data.LifetimeValue.IsUnknown() {
		data.LifetimeValue = types.Int64Null()
	}
	if data.LifetimeUnit.IsUnknown() {
		data.LifetimeUnit = types.StringNull()
	}
	if data.Lifetimes.IsUnknown() {
		data.Lifetimes, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.Compression.IsUnknown() {
		data.Compression = types.StringNull()
	}
	if data.SpeedLimit.IsUnknown() {
		data.SpeedLimit = types.Int64Null()
	}
	if data.LargeBlock.IsUnknown() {
		data.LargeBlock = types.BoolNull()
	}
	if data.Embed.IsUnknown() {
		data.Embed = types.BoolNull()
	}
	if data.Compressed.IsUnknown() {
		data.Compressed = types.BoolNull()
	}
	if data.Retries.IsUnknown() {
		data.Retries = types.Int64Null()
	}
	if data.LoggingLevel.IsUnknown() {
		data.LoggingLevel = types.StringNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ReplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ReplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("replication.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete replication: %s", err))
		return
	}
}
