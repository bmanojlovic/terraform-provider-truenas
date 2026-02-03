package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &ReplicationsDataSource{}

func NewReplicationsDataSource() datasource.DataSource {
	return &ReplicationsDataSource{}
}

type ReplicationsDataSource struct {
	client *client.Client
}

type ReplicationsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type ReplicationsItemModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Direction types.String `tfsdk:"direction"`
	Transport types.String `tfsdk:"transport"`
	SshCredentials types.String `tfsdk:"ssh_credentials"`
	NetcatActiveSide types.String `tfsdk:"netcat_active_side"`
	NetcatActiveSideListenAddress types.String `tfsdk:"netcat_active_side_listen_address"`
	NetcatActiveSidePortMin types.Int64 `tfsdk:"netcat_active_side_port_min"`
	NetcatActiveSidePortMax types.Int64 `tfsdk:"netcat_active_side_port_max"`
	NetcatPassiveSideConnectAddress types.String `tfsdk:"netcat_passive_side_connect_address"`
	Sudo types.Bool `tfsdk:"sudo"`
	TargetDataset types.String `tfsdk:"target_dataset"`
	Recursive types.Bool `tfsdk:"recursive"`
	Properties types.Bool `tfsdk:"properties"`
	PropertiesOverride types.String `tfsdk:"properties_override"`
	Replicate types.Bool `tfsdk:"replicate"`
	Encryption types.Bool `tfsdk:"encryption"`
	EncryptionInherit types.Bool `tfsdk:"encryption_inherit"`
	EncryptionKey types.String `tfsdk:"encryption_key"`
	EncryptionKeyFormat types.String `tfsdk:"encryption_key_format"`
	EncryptionKeyLocation types.String `tfsdk:"encryption_key_location"`
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
	Compression types.String `tfsdk:"compression"`
	SpeedLimit types.Int64 `tfsdk:"speed_limit"`
	LargeBlock types.Bool `tfsdk:"large_block"`
	Embed types.Bool `tfsdk:"embed"`
	Compressed types.Bool `tfsdk:"compressed"`
	Retries types.Int64 `tfsdk:"retries"`
	LoggingLevel types.String `tfsdk:"logging_level"`
	Enabled types.Bool `tfsdk:"enabled"`
	State types.String `tfsdk:"state"`
	Job types.String `tfsdk:"job"`
	HasEncryptedDatasetKeys types.Bool `tfsdk:"has_encrypted_dataset_keys"`
}

func (d *ReplicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replications"
}

func (d *ReplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query replications",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of replications resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Name for replication task.",
			},
			"direction": schema.StringAttribute{
				Computed: true,
				Description: "Whether task will `PUSH` or `PULL` snapshots.",
			},
			"transport": schema.StringAttribute{
				Computed: true,
				Description: "Method of snapshots transfer.  * `SSH` transfers snapshots via SSH connection. This method is suppor",
			},
			"ssh_credentials": schema.StringAttribute{
				Computed: true,
				Description: "Keychain Credential of type `SSH_CREDENTIALS`.",
			},
			"netcat_active_side": schema.StringAttribute{
				Computed: true,
				Description: "Which side actively establishes the netcat connection for `SSH+NETCAT` transport.  * `LOCAL`: Local ",
			},
			"netcat_active_side_listen_address": schema.StringAttribute{
				Computed: true,
				Description: "IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"netcat_active_side_port_min": schema.Int64Attribute{
				Computed: true,
				Description: "Minimum port number in the range for netcat connections. `null` if not applicable.",
			},
			"netcat_active_side_port_max": schema.Int64Attribute{
				Computed: true,
				Description: "Maximum port number in the range for netcat connections. `null` if not applicable.",
			},
			"netcat_passive_side_connect_address": schema.StringAttribute{
				Computed: true,
				Description: "IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"sudo": schema.BoolAttribute{
				Computed: true,
				Description: "`SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs",
			},
			"target_dataset": schema.StringAttribute{
				Computed: true,
				Description: "Dataset to put snapshots into.",
			},
			"recursive": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to recursively replicate child datasets.",
			},
			"properties": schema.BoolAttribute{
				Computed: true,
				Description: "Send dataset properties along with snapshots.",
			},
			"properties_override": schema.StringAttribute{
				Computed: true,
				Description: "Object mapping dataset property names to override values during replication.",
			},
			"replicate": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to use full ZFS replication.",
			},
			"encryption": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to enable encryption for the replicated datasets.",
			},
			"encryption_inherit": schema.BoolAttribute{
				Computed: true,
				Description: "Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled.",
			},
			"encryption_key": schema.StringAttribute{
				Computed: true,
				Description: "Encryption key for replicated datasets. `null` if not specified.",
			},
			"encryption_key_format": schema.StringAttribute{
				Computed: true,
				Description: "Format of the encryption key.  * `HEX`: Hexadecimal-encoded key * `PASSPHRASE`: Text passphrase * `n",
			},
			"encryption_key_location": schema.StringAttribute{
				Computed: true,
				Description: "Filesystem path where encryption key is stored. `null` if not using key file.",
			},
			"name_regex": schema.StringAttribute{
				Computed: true,
				Description: "Replicate all snapshots which names match specified regular expression.",
			},
			"auto": schema.BoolAttribute{
				Computed: true,
				Description: "Allow replication to run automatically on schedule or after bound periodic snapshot task.",
			},
			"schedule": schema.StringAttribute{
				Computed: true,
				Description: "Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot task",
			},
			"restrict_schedule": schema.StringAttribute{
				Computed: true,
				Description: "Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have p",
			},
			"only_matching_schedule": schema.BoolAttribute{
				Computed: true,
				Description: "Will only replicate snapshots that match `schedule` or `restrict_schedule`.",
			},
			"allow_from_scratch": schema.BoolAttribute{
				Computed: true,
				Description: "Will destroy all snapshots on target side and replicate everything from scratch if none of the snaps",
			},
			"readonly": schema.StringAttribute{
				Computed: true,
				Description: "Controls destination datasets readonly property.  * `SET`: Set all destination datasets to readonly=",
			},
			"hold_pending_snapshots": schema.BoolAttribute{
				Computed: true,
				Description: "Prevent source snapshots from being deleted by retention of replication fails for some reason.",
			},
			"retention_policy": schema.StringAttribute{
				Computed: true,
				Description: "How to delete old snapshots on target side:  * `SOURCE`: Delete snapshots that are absent on source ",
			},
			"lifetime_value": schema.Int64Attribute{
				Computed: true,
				Description: "Number of time units to retain snapshots for custom retention policy. Only applies when `retention_p",
			},
			"lifetime_unit": schema.StringAttribute{
				Computed: true,
				Description: "Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` i",
			},
			"compression": schema.StringAttribute{
				Computed: true,
				Description: "Compresses SSH stream. Available only for SSH transport.",
			},
			"speed_limit": schema.Int64Attribute{
				Computed: true,
				Description: "Limits speed of SSH stream. Available only for SSH transport.",
			},
			"large_block": schema.BoolAttribute{
				Computed: true,
				Description: "Enable large block support for ZFS send streams.",
			},
			"embed": schema.BoolAttribute{
				Computed: true,
				Description: "Enable embedded block support for ZFS send streams.",
			},
			"compressed": schema.BoolAttribute{
				Computed: true,
				Description: "Enable compressed ZFS send streams.",
			},
			"retries": schema.Int64Attribute{
				Computed: true,
				Description: "Number of retries before considering replication failed.",
			},
			"logging_level": schema.StringAttribute{
				Computed: true,
				Description: "Log level for replication task execution. Controls verbosity of replication logs.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this replication task is enabled.",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Current state information for the replication task.",
			},
			"job": schema.StringAttribute{
				Computed: true,
				Description: "Information about the currently running job. `null` if no job is running.",
			},
			"has_encrypted_dataset_keys": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this replication task has encrypted dataset keys available.",
			},
					},
				},
			},
		},
	}
}

func (d *ReplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *ReplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ReplicationsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("replication.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query replications: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]ReplicationsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := ReplicationsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["direction"]; ok && v != nil {
			itemModel.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["transport"]; ok && v != nil {
			itemModel.Transport = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ssh_credentials"]; ok && v != nil {
			itemModel.SshCredentials = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["netcat_active_side"]; ok && v != nil {
			itemModel.NetcatActiveSide = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["netcat_active_side_listen_address"]; ok && v != nil {
			itemModel.NetcatActiveSideListenAddress = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["netcat_active_side_port_min"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.NetcatActiveSidePortMin = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["netcat_active_side_port_max"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.NetcatActiveSidePortMax = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["netcat_passive_side_connect_address"]; ok && v != nil {
			itemModel.NetcatPassiveSideConnectAddress = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sudo"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Sudo = types.BoolValue(bv) }
		}
		if v, ok := resultMap["target_dataset"]; ok && v != nil {
			itemModel.TargetDataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recursive"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Recursive = types.BoolValue(bv) }
		}
		if v, ok := resultMap["properties"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Properties = types.BoolValue(bv) }
		}
		if v, ok := resultMap["properties_override"]; ok && v != nil {
			itemModel.PropertiesOverride = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["replicate"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Replicate = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Encryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_inherit"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.EncryptionInherit = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_key"]; ok && v != nil {
			itemModel.EncryptionKey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_key_format"]; ok && v != nil {
			itemModel.EncryptionKeyFormat = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_key_location"]; ok && v != nil {
			itemModel.EncryptionKeyLocation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name_regex"]; ok && v != nil {
			itemModel.NameRegex = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["auto"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Auto = types.BoolValue(bv) }
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			itemModel.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["restrict_schedule"]; ok && v != nil {
			itemModel.RestrictSchedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["only_matching_schedule"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.OnlyMatchingSchedule = types.BoolValue(bv) }
		}
		if v, ok := resultMap["allow_from_scratch"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.AllowFromScratch = types.BoolValue(bv) }
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			itemModel.Readonly = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hold_pending_snapshots"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.HoldPendingSnapshots = types.BoolValue(bv) }
		}
		if v, ok := resultMap["retention_policy"]; ok && v != nil {
			itemModel.RetentionPolicy = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["lifetime_value"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.LifetimeValue = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["lifetime_unit"]; ok && v != nil {
			itemModel.LifetimeUnit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["compression"]; ok && v != nil {
			itemModel.Compression = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["speed_limit"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.SpeedLimit = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["large_block"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.LargeBlock = types.BoolValue(bv) }
		}
		if v, ok := resultMap["embed"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Embed = types.BoolValue(bv) }
		}
		if v, ok := resultMap["compressed"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Compressed = types.BoolValue(bv) }
		}
		if v, ok := resultMap["retries"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Retries = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["logging_level"]; ok && v != nil {
			itemModel.LoggingLevel = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["job"]; ok && v != nil {
			itemModel.Job = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["has_encrypted_dataset_keys"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.HasEncryptedDatasetKeys = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"allow_from_scratch": types.BoolType,
			"auto": types.BoolType,
			"compressed": types.BoolType,
			"compression": types.StringType,
			"direction": types.StringType,
			"embed": types.BoolType,
			"enabled": types.BoolType,
			"encryption": types.BoolType,
			"encryption_inherit": types.BoolType,
			"encryption_key": types.StringType,
			"encryption_key_format": types.StringType,
			"encryption_key_location": types.StringType,
			"has_encrypted_dataset_keys": types.BoolType,
			"hold_pending_snapshots": types.BoolType,
			"id": types.StringType,
			"job": types.StringType,
			"large_block": types.BoolType,
			"lifetime_unit": types.StringType,
			"lifetime_value": types.Int64Type,
			"logging_level": types.StringType,
			"name": types.StringType,
			"name_regex": types.StringType,
			"netcat_active_side": types.StringType,
			"netcat_active_side_listen_address": types.StringType,
			"netcat_active_side_port_max": types.Int64Type,
			"netcat_active_side_port_min": types.Int64Type,
			"netcat_passive_side_connect_address": types.StringType,
			"only_matching_schedule": types.BoolType,
			"properties": types.BoolType,
			"properties_override": types.StringType,
			"readonly": types.StringType,
			"recursive": types.BoolType,
			"replicate": types.BoolType,
			"restrict_schedule": types.StringType,
			"retention_policy": types.StringType,
			"retries": types.Int64Type,
			"schedule": types.StringType,
			"speed_limit": types.Int64Type,
			"ssh_credentials": types.StringType,
			"state": types.StringType,
			"sudo": types.BoolType,
			"target_dataset": types.StringType,
			"transport": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
