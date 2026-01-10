package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	SshCredentials types.Int64 `tfsdk:"ssh_credentials"`
	NetcatActiveSide types.String `tfsdk:"netcat_active_side"`
	NetcatActiveSideListenAddress types.String `tfsdk:"netcat_active_side_listen_address"`
	NetcatActiveSidePortMin types.Int64 `tfsdk:"netcat_active_side_port_min"`
	NetcatActiveSidePortMax types.Int64 `tfsdk:"netcat_active_side_port_max"`
	NetcatPassiveSideConnectAddress types.String `tfsdk:"netcat_passive_side_connect_address"`
	Sudo types.Bool `tfsdk:"sudo"`
	SourceDatasets types.List `tfsdk:"source_datasets"`
	TargetDataset types.String `tfsdk:"target_dataset"`
	Recursive types.Bool `tfsdk:"recursive"`
	Exclude types.List `tfsdk:"exclude"`
	Properties types.Bool `tfsdk:"properties"`
	PropertiesExclude types.List `tfsdk:"properties_exclude"`
	PropertiesOverride types.String `tfsdk:"properties_override"`
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
	LifetimeValue types.Int64 `tfsdk:"lifetime_value"`
	LifetimeUnit types.String `tfsdk:"lifetime_unit"`
	Lifetimes types.List `tfsdk:"lifetimes"`
	Compression types.String `tfsdk:"compression"`
	SpeedLimit types.Int64 `tfsdk:"speed_limit"`
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

func (r *ReplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ReplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Replication Task that will push or pull ZFS snapshots to or from remote host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name for replication task.",
			},
			"direction": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Whether task will `PUSH` or `PULL` snapshots.",
			},
			"transport": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Method of snapshots transfer.  * `SSH` transfers snapshots via SSH connection. This method is suppor",
			},
			"ssh_credentials": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Keychain Credential ID of type `SSH_CREDENTIALS`.",
			},
			"netcat_active_side": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Which side actively establishes the netcat connection for `SSH+NETCAT` transport.  * `LOCAL`: Local ",
			},
			"netcat_active_side_listen_address": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"netcat_active_side_port_min": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Minimum port number in the range for netcat connections. `null` if not applicable.",
			},
			"netcat_active_side_port_max": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum port number in the range for netcat connections. `null` if not applicable.",
			},
			"netcat_passive_side_connect_address": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"sudo": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "`SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs",
			},
			"source_datasets": schema.ListAttribute{
				Required: true,
				Optional: false,
				ElementType: types.StringType,
				Description: "List of datasets to replicate snapshots from.",
			},
			"target_dataset": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Dataset to put snapshots into.",
			},
			"recursive": schema.BoolAttribute{
				Required: true,
				Optional: false,
				Description: "Whether to recursively replicate child datasets.",
			},
			"exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of dataset patterns to exclude from replication.",
			},
			"properties": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Send dataset properties along with snapshots.",
			},
			"properties_exclude": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of dataset property names to exclude from replication.",
			},
			"properties_override": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Object mapping dataset property names to override values during replication.",
			},
			"replicate": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to use full ZFS replication.",
			},
			"encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable encryption for the replicated datasets.",
			},
			"encryption_inherit": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled.",
			},
			"encryption_key": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Encryption key for replicated datasets. `null` if not specified.",
			},
			"encryption_key_format": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Format of the encryption key.  * `HEX`: Hexadecimal-encoded key * `PASSPHRASE`: Text passphrase * `n",
			},
			"encryption_key_location": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Filesystem path where encryption key is stored. `null` if not using key file.",
			},
			"periodic_snapshot_tasks": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of periodic snapshot task IDs that are sources of snapshots for this replication task. Only pus",
			},
			"naming_schema": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of naming schemas for pull replication.",
			},
			"also_include_naming_schema": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "List of naming schemas for push replication.",
			},
			"name_regex": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Replicate all snapshots which names match specified regular expression.",
			},
			"auto": schema.BoolAttribute{
				Required: true,
				Optional: false,
				Description: "Allow replication to run automatically on schedule or after bound periodic snapshot task.",
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot task",
			},
			"restrict_schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have p",
			},
			"only_matching_schedule": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Will only replicate snapshots that match `schedule` or `restrict_schedule`.",
			},
			"allow_from_scratch": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Will destroy all snapshots on target side and replicate everything from scratch if none of the snaps",
			},
			"readonly": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Controls destination datasets readonly property.  * `SET`: Set all destination datasets to readonly=",
			},
			"hold_pending_snapshots": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Prevent source snapshots from being deleted by retention of replication fails for some reason.",
			},
			"retention_policy": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "How to delete old snapshots on target side:  * `SOURCE`: Delete snapshots that are absent on source ",
			},
			"lifetime_value": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of time units to retain snapshots for custom retention policy. Only applies when `retention_p",
			},
			"lifetime_unit": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` i",
			},
			"lifetimes": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of different retention schedules with their own cron schedules and lifetime settings.",
			},
			"compression": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Compresses SSH stream. Available only for SSH transport.",
			},
			"speed_limit": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Limits speed of SSH stream. Available only for SSH transport.",
			},
			"large_block": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable large block support for ZFS send streams.",
			},
			"embed": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable embedded block support for ZFS send streams.",
			},
			"compressed": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Enable compressed ZFS send streams.",
			},
			"retries": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of retries before considering replication failed.",
			},
			"logging_level": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Log level for replication task execution. Controls verbosity of replication logs.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Direction.IsNull() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.Transport.IsNull() {
		params["transport"] = data.Transport.ValueString()
	}
	if !data.SshCredentials.IsNull() {
		params["ssh_credentials"] = data.SshCredentials.ValueInt64()
	}
	if !data.NetcatActiveSide.IsNull() {
		params["netcat_active_side"] = data.NetcatActiveSide.ValueString()
	}
	if !data.NetcatActiveSideListenAddress.IsNull() {
		params["netcat_active_side_listen_address"] = data.NetcatActiveSideListenAddress.ValueString()
	}
	if !data.NetcatActiveSidePortMin.IsNull() {
		params["netcat_active_side_port_min"] = data.NetcatActiveSidePortMin.ValueInt64()
	}
	if !data.NetcatActiveSidePortMax.IsNull() {
		params["netcat_active_side_port_max"] = data.NetcatActiveSidePortMax.ValueInt64()
	}
	if !data.NetcatPassiveSideConnectAddress.IsNull() {
		params["netcat_passive_side_connect_address"] = data.NetcatPassiveSideConnectAddress.ValueString()
	}
	if !data.Sudo.IsNull() {
		params["sudo"] = data.Sudo.ValueBool()
	}
	if !data.SourceDatasets.IsNull() {
		var source_datasetsList []string
		data.SourceDatasets.ElementsAs(ctx, &source_datasetsList, false)
		params["source_datasets"] = source_datasetsList
	}
	if !data.TargetDataset.IsNull() {
		params["target_dataset"] = data.TargetDataset.ValueString()
	}
	if !data.Recursive.IsNull() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Properties.IsNull() {
		params["properties"] = data.Properties.ValueBool()
	}
	if !data.PropertiesExclude.IsNull() {
		var properties_excludeList []string
		data.PropertiesExclude.ElementsAs(ctx, &properties_excludeList, false)
		params["properties_exclude"] = properties_excludeList
	}
	if !data.PropertiesOverride.IsNull() {
		params["properties_override"] = data.PropertiesOverride.ValueString()
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
	if !data.PeriodicSnapshotTasks.IsNull() {
		var periodic_snapshot_tasksList []string
		data.PeriodicSnapshotTasks.ElementsAs(ctx, &periodic_snapshot_tasksList, false)
		params["periodic_snapshot_tasks"] = periodic_snapshot_tasksList
	}
	if !data.NamingSchema.IsNull() {
		var naming_schemaList []string
		data.NamingSchema.ElementsAs(ctx, &naming_schemaList, false)
		params["naming_schema"] = naming_schemaList
	}
	if !data.AlsoIncludeNamingSchema.IsNull() {
		var also_include_naming_schemaList []string
		data.AlsoIncludeNamingSchema.ElementsAs(ctx, &also_include_naming_schemaList, false)
		params["also_include_naming_schema"] = also_include_naming_schemaList
	}
	if !data.NameRegex.IsNull() {
		params["name_regex"] = data.NameRegex.ValueString()
	}
	if !data.Auto.IsNull() {
		params["auto"] = data.Auto.ValueBool()
	}
	if !data.Schedule.IsNull() {
		params["schedule"] = data.Schedule.ValueString()
	}
	if !data.RestrictSchedule.IsNull() {
		params["restrict_schedule"] = data.RestrictSchedule.ValueString()
	}
	if !data.OnlyMatchingSchedule.IsNull() {
		params["only_matching_schedule"] = data.OnlyMatchingSchedule.ValueBool()
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
		params["lifetime_value"] = data.LifetimeValue.ValueInt64()
	}
	if !data.LifetimeUnit.IsNull() {
		params["lifetime_unit"] = data.LifetimeUnit.ValueString()
	}
	if !data.Lifetimes.IsNull() {
		var lifetimesList []string
		data.Lifetimes.ElementsAs(ctx, &lifetimesList, false)
		params["lifetimes"] = lifetimesList
	}
	if !data.Compression.IsNull() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.SpeedLimit.IsNull() {
		params["speed_limit"] = data.SpeedLimit.ValueInt64()
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
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	result, err := r.client.Call("replication.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create replication: %s", err))
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

func (r *ReplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	result, err := r.client.Call("replication.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read replication: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["direction"]; ok && v != nil {
			data.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["transport"]; ok && v != nil {
			data.Transport = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ssh_credentials"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.SshCredentials = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["netcat_active_side"]; ok && v != nil {
			data.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["netcat_active_side_listen_address"]; ok && v != nil {
			data.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["netcat_active_side_port_min"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.NetcatActiveSidePortMin = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["netcat_active_side_port_max"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.NetcatActiveSidePortMax = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["netcat_passive_side_connect_address"]; ok && v != nil {
			data.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sudo"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Sudo = types.BoolValue(bv) }
		}
		if v, ok := resultMap["source_datasets"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.SourceDatasets, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["target_dataset"]; ok && v != nil {
			data.TargetDataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recursive"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Recursive = types.BoolValue(bv) }
		}
		if v, ok := resultMap["exclude"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Exclude, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["properties"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Properties = types.BoolValue(bv) }
		}
		if v, ok := resultMap["properties_exclude"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.PropertiesExclude, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["properties_override"]; ok && v != nil {
			data.PropertiesOverride = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["replicate"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Replicate = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Encryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_inherit"]; ok && v != nil {
			data.EncryptionInherit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_key"]; ok && v != nil {
			data.EncryptionKey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_key_format"]; ok && v != nil {
			data.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_key_location"]; ok && v != nil {
			data.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["periodic_snapshot_tasks"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.PeriodicSnapshotTasks, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["naming_schema"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.NamingSchema, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["also_include_naming_schema"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.AlsoIncludeNamingSchema, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["name_regex"]; ok && v != nil {
			data.NameRegex = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["auto"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Auto = types.BoolValue(bv) }
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			data.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["restrict_schedule"]; ok && v != nil {
			data.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["only_matching_schedule"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.OnlyMatchingSchedule = types.BoolValue(bv) }
		}
		if v, ok := resultMap["allow_from_scratch"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AllowFromScratch = types.BoolValue(bv) }
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			data.Readonly = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hold_pending_snapshots"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.HoldPendingSnapshots = types.BoolValue(bv) }
		}
		if v, ok := resultMap["retention_policy"]; ok && v != nil {
			data.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lifetime_value"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.LifetimeValue = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["lifetime_unit"]; ok && v != nil {
			data.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lifetimes"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Lifetimes, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["compression"]; ok && v != nil {
			data.Compression = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["speed_limit"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.SpeedLimit = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["large_block"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.LargeBlock = types.BoolValue(bv) }
		}
		if v, ok := resultMap["embed"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Embed = types.BoolValue(bv) }
		}
		if v, ok := resultMap["compressed"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Compressed = types.BoolValue(bv) }
		}
		if v, ok := resultMap["retries"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Retries = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["logging_level"]; ok && v != nil {
			data.LoggingLevel = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enabled = types.BoolValue(bv) }
		}
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

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Direction.IsNull() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.Transport.IsNull() {
		params["transport"] = data.Transport.ValueString()
	}
	if !data.SshCredentials.IsNull() {
		params["ssh_credentials"] = data.SshCredentials.ValueInt64()
	}
	if !data.NetcatActiveSide.IsNull() {
		params["netcat_active_side"] = data.NetcatActiveSide.ValueString()
	}
	if !data.NetcatActiveSideListenAddress.IsNull() {
		params["netcat_active_side_listen_address"] = data.NetcatActiveSideListenAddress.ValueString()
	}
	if !data.NetcatActiveSidePortMin.IsNull() {
		params["netcat_active_side_port_min"] = data.NetcatActiveSidePortMin.ValueInt64()
	}
	if !data.NetcatActiveSidePortMax.IsNull() {
		params["netcat_active_side_port_max"] = data.NetcatActiveSidePortMax.ValueInt64()
	}
	if !data.NetcatPassiveSideConnectAddress.IsNull() {
		params["netcat_passive_side_connect_address"] = data.NetcatPassiveSideConnectAddress.ValueString()
	}
	if !data.Sudo.IsNull() {
		params["sudo"] = data.Sudo.ValueBool()
	}
	if !data.SourceDatasets.IsNull() {
		var source_datasetsList []string
		data.SourceDatasets.ElementsAs(ctx, &source_datasetsList, false)
		params["source_datasets"] = source_datasetsList
	}
	if !data.TargetDataset.IsNull() {
		params["target_dataset"] = data.TargetDataset.ValueString()
	}
	if !data.Recursive.IsNull() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Exclude.IsNull() {
		var excludeList []string
		data.Exclude.ElementsAs(ctx, &excludeList, false)
		params["exclude"] = excludeList
	}
	if !data.Properties.IsNull() {
		params["properties"] = data.Properties.ValueBool()
	}
	if !data.PropertiesExclude.IsNull() {
		var properties_excludeList []string
		data.PropertiesExclude.ElementsAs(ctx, &properties_excludeList, false)
		params["properties_exclude"] = properties_excludeList
	}
	if !data.PropertiesOverride.IsNull() {
		params["properties_override"] = data.PropertiesOverride.ValueString()
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
	if !data.PeriodicSnapshotTasks.IsNull() {
		var periodic_snapshot_tasksList []string
		data.PeriodicSnapshotTasks.ElementsAs(ctx, &periodic_snapshot_tasksList, false)
		params["periodic_snapshot_tasks"] = periodic_snapshot_tasksList
	}
	if !data.NamingSchema.IsNull() {
		var naming_schemaList []string
		data.NamingSchema.ElementsAs(ctx, &naming_schemaList, false)
		params["naming_schema"] = naming_schemaList
	}
	if !data.AlsoIncludeNamingSchema.IsNull() {
		var also_include_naming_schemaList []string
		data.AlsoIncludeNamingSchema.ElementsAs(ctx, &also_include_naming_schemaList, false)
		params["also_include_naming_schema"] = also_include_naming_schemaList
	}
	if !data.NameRegex.IsNull() {
		params["name_regex"] = data.NameRegex.ValueString()
	}
	if !data.Auto.IsNull() {
		params["auto"] = data.Auto.ValueBool()
	}
	if !data.Schedule.IsNull() {
		params["schedule"] = data.Schedule.ValueString()
	}
	if !data.RestrictSchedule.IsNull() {
		params["restrict_schedule"] = data.RestrictSchedule.ValueString()
	}
	if !data.OnlyMatchingSchedule.IsNull() {
		params["only_matching_schedule"] = data.OnlyMatchingSchedule.ValueBool()
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
		params["lifetime_value"] = data.LifetimeValue.ValueInt64()
	}
	if !data.LifetimeUnit.IsNull() {
		params["lifetime_unit"] = data.LifetimeUnit.ValueString()
	}
	if !data.Lifetimes.IsNull() {
		var lifetimesList []string
		data.Lifetimes.ElementsAs(ctx, &lifetimesList, false)
		params["lifetimes"] = lifetimesList
	}
	if !data.Compression.IsNull() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.SpeedLimit.IsNull() {
		params["speed_limit"] = data.SpeedLimit.ValueInt64()
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
	if !data.Enabled.IsNull() {
		params["enabled"] = data.Enabled.ValueBool()
	}

	_, err = r.client.Call("replication.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update replication: %s", err))
		return
	}

	data.ID = state.ID
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
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete replication: %s", err))
		return
	}
}
