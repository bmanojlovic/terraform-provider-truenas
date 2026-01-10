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

var _ datasource.DataSource = &PoolDatasetDataSource{}

func NewPoolDatasetDataSource() datasource.DataSource {
	return &PoolDatasetDataSource{}
}

type PoolDatasetDataSource struct {
	client *client.Client
}

type PoolDatasetDataSourceModel struct {
	ID types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
	Pool types.String `tfsdk:"pool"`
	Encrypted types.Bool `tfsdk:"encrypted"`
	EncryptionRoot types.String `tfsdk:"encryption_root"`
	KeyLoaded types.String `tfsdk:"key_loaded"`
	Children types.List `tfsdk:"children"`
	UserProperties types.String `tfsdk:"user_properties"`
	Locked types.Bool `tfsdk:"locked"`
	Comments types.String `tfsdk:"comments"`
	QuotaWarning types.String `tfsdk:"quota_warning"`
	QuotaCritical types.String `tfsdk:"quota_critical"`
	RefquotaWarning types.String `tfsdk:"refquota_warning"`
	RefquotaCritical types.String `tfsdk:"refquota_critical"`
	Managedby types.String `tfsdk:"managedby"`
	Deduplication types.String `tfsdk:"deduplication"`
	Aclmode types.String `tfsdk:"aclmode"`
	Acltype types.String `tfsdk:"acltype"`
	Xattr types.String `tfsdk:"xattr"`
	Atime types.String `tfsdk:"atime"`
	Casesensitivity types.String `tfsdk:"casesensitivity"`
	Checksum types.String `tfsdk:"checksum"`
	Exec types.String `tfsdk:"exec"`
	Sync types.String `tfsdk:"sync"`
	Compression types.String `tfsdk:"compression"`
	Compressratio types.String `tfsdk:"compressratio"`
	Origin types.String `tfsdk:"origin"`
	Quota types.String `tfsdk:"quota"`
	Refquota types.String `tfsdk:"refquota"`
	Reservation types.String `tfsdk:"reservation"`
	Refreservation types.String `tfsdk:"refreservation"`
	Copies types.String `tfsdk:"copies"`
	Snapdir types.String `tfsdk:"snapdir"`
	Readonly types.String `tfsdk:"readonly"`
	Recordsize types.String `tfsdk:"recordsize"`
	Sparse types.String `tfsdk:"sparse"`
	Volsize types.String `tfsdk:"volsize"`
	Volblocksize types.String `tfsdk:"volblocksize"`
	KeyFormat types.String `tfsdk:"key_format"`
	EncryptionAlgorithm types.String `tfsdk:"encryption_algorithm"`
	Used types.String `tfsdk:"used"`
	Usedbychildren types.String `tfsdk:"usedbychildren"`
	Usedbydataset types.String `tfsdk:"usedbydataset"`
	Usedbyrefreservation types.String `tfsdk:"usedbyrefreservation"`
	Usedbysnapshots types.String `tfsdk:"usedbysnapshots"`
	Available types.String `tfsdk:"available"`
	SpecialSmallBlockSize types.String `tfsdk:"special_small_block_size"`
	Pbkdf2Iters types.String `tfsdk:"pbkdf2iters"`
	Creation types.String `tfsdk:"creation"`
	Snapdev types.String `tfsdk:"snapdev"`
	Mountpoint types.String `tfsdk:"mountpoint"`
}

func (d *PoolDatasetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_dataset"
}

func (d *PoolDatasetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "The dataset type.",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "The dataset name without the pool prefix.",
			},
			"pool": schema.StringAttribute{
				Computed: true,
				Description: "The name of the ZFS pool containing this dataset.",
			},
			"encrypted": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the dataset is encrypted.",
			},
			"encryption_root": schema.StringAttribute{
				Computed: true,
				Description: "The root dataset where encryption is enabled. `null` if the dataset is not encrypted.",
			},
			"key_loaded": schema.StringAttribute{
				Computed: true,
				Description: "Whether the encryption key is currently loaded for encrypted datasets. `null` for unencrypted datase",
			},
			"children": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "Array of child dataset objects nested under this dataset.",
			},
			"user_properties": schema.StringAttribute{
				Computed: true,
				Description: "Custom user-defined ZFS properties set on this dataset as key-value pairs.",
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				Description: "Whether an encrypted dataset is currently locked (key not loaded).",
			},
			"comments": schema.StringAttribute{
				Computed: true,
				Description: "ZFS comments property for storing descriptive text about the dataset.",
			},
			"quota_warning": schema.StringAttribute{
				Computed: true,
				Description: "ZFS quota warning threshold property as a percentage.",
			},
			"quota_critical": schema.StringAttribute{
				Computed: true,
				Description: "ZFS quota critical threshold property as a percentage.",
			},
			"refquota_warning": schema.StringAttribute{
				Computed: true,
				Description: "ZFS reference quota warning threshold property as a percentage.",
			},
			"refquota_critical": schema.StringAttribute{
				Computed: true,
				Description: "ZFS reference quota critical threshold property as a percentage.",
			},
			"managedby": schema.StringAttribute{
				Computed: true,
				Description: "Identifies which service or system manages this dataset.",
			},
			"deduplication": schema.StringAttribute{
				Computed: true,
				Description: "ZFS deduplication setting - whether identical data blocks are stored only once.",
			},
			"aclmode": schema.StringAttribute{
				Computed: true,
				Description: "How Access Control Lists (ACLs) are handled when chmod is used.",
			},
			"acltype": schema.StringAttribute{
				Computed: true,
				Description: "The type of Access Control List system used (NFSV4, POSIX, or OFF).",
			},
			"xattr": schema.StringAttribute{
				Computed: true,
				Description: "Extended attributes storage method (on/off).",
			},
			"atime": schema.StringAttribute{
				Computed: true,
				Description: "Whether file access times are updated when files are accessed.",
			},
			"casesensitivity": schema.StringAttribute{
				Computed: true,
				Description: "File name case sensitivity setting (sensitive/insensitive).",
			},
			"checksum": schema.StringAttribute{
				Computed: true,
				Description: "Data integrity checksum algorithm used for this dataset.",
			},
			"exec": schema.StringAttribute{
				Computed: true,
				Description: "Whether files in this dataset can be executed.",
			},
			"sync": schema.StringAttribute{
				Computed: true,
				Description: "Synchronous write behavior (standard/always/disabled).",
			},
			"compression": schema.StringAttribute{
				Computed: true,
				Description: "Compression algorithm and level applied to data in this dataset.",
			},
			"compressratio": schema.StringAttribute{
				Computed: true,
				Description: "The achieved compression ratio as a decimal (e.g., '2.50x').",
			},
			"origin": schema.StringAttribute{
				Computed: true,
				Description: "The snapshot from which this clone was created. Empty for non-clone datasets.",
			},
			"quota": schema.StringAttribute{
				Computed: true,
				Description: "Maximum amount of disk space this dataset and its children can consume.",
			},
			"refquota": schema.StringAttribute{
				Computed: true,
				Description: "Maximum amount of disk space this dataset itself can consume (excluding children).",
			},
			"reservation": schema.StringAttribute{
				Computed: true,
				Description: "Minimum amount of disk space guaranteed to be available to this dataset and its children.",
			},
			"refreservation": schema.StringAttribute{
				Computed: true,
				Description: "Minimum amount of disk space guaranteed to be available to this dataset itself.",
			},
			"copies": schema.StringAttribute{
				Computed: true,
				Description: "Number of copies of data blocks to maintain for redundancy (1-3).",
			},
			"snapdir": schema.StringAttribute{
				Computed: true,
				Description: "Visibility of the .zfs/snapshot directory (visible/hidden).",
			},
			"readonly": schema.StringAttribute{
				Computed: true,
				Description: "Whether the dataset is read-only.",
			},
			"recordsize": schema.StringAttribute{
				Computed: true,
				Description: "The suggested block size for files in this filesystem dataset.",
			},
			"sparse": schema.StringAttribute{
				Computed: true,
				Description: "For volumes, whether to use sparse (thin) provisioning.",
			},
			"volsize": schema.StringAttribute{
				Computed: true,
				Description: "For volumes, the logical size of the volume.",
			},
			"volblocksize": schema.StringAttribute{
				Computed: true,
				Description: "For volumes, the block size used by the volume.",
			},
			"key_format": schema.StringAttribute{
				Computed: true,
				Description: "Format of the encryption key (hex/raw/passphrase). Only relevant for encrypted datasets.",
			},
			"encryption_algorithm": schema.StringAttribute{
				Computed: true,
				Description: "Encryption algorithm used (e.g., AES-256-GCM). Only relevant for encrypted datasets.",
			},
			"used": schema.StringAttribute{
				Computed: true,
				Description: "Total amount of disk space consumed by this dataset and all its children.",
			},
			"usedbychildren": schema.StringAttribute{
				Computed: true,
				Description: "Amount of disk space consumed by child datasets.",
			},
			"usedbydataset": schema.StringAttribute{
				Computed: true,
				Description: "Amount of disk space consumed by this dataset itself, excluding children and snapshots.",
			},
			"usedbyrefreservation": schema.StringAttribute{
				Computed: true,
				Description: "Amount of disk space consumed by the refreservation of this dataset.",
			},
			"usedbysnapshots": schema.StringAttribute{
				Computed: true,
				Description: "Amount of disk space consumed by snapshots of this dataset.",
			},
			"available": schema.StringAttribute{
				Computed: true,
				Description: "Amount of disk space available to this dataset and its children.",
			},
			"special_small_block_size": schema.StringAttribute{
				Computed: true,
				Description: "Size threshold below which blocks are stored on special vdevs if configured.",
			},
			"pbkdf2iters": schema.StringAttribute{
				Computed: true,
				Description: "Number of PBKDF2 iterations used for passphrase-based encryption keys.",
			},
			"creation": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when this dataset was created.",
			},
			"snapdev": schema.StringAttribute{
				Computed: true,
				Description: "Controls visibility of volume snapshots under /dev/zvol/<pool>/.",
			},
			"mountpoint": schema.StringAttribute{
				Computed: true,
				Description: "Filesystem path where this dataset is mounted. Null for unmounted datasets or volumes.",
			},
		},
	}
}

func (d *PoolDatasetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolDatasetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PoolDatasetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("pool.dataset.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read pool_dataset: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["type"]; ok && v != nil {
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pool"]; ok && v != nil {
			data.Pool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encrypted"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Encrypted = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_root"]; ok && v != nil {
			data.EncryptionRoot = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["key_loaded"]; ok && v != nil {
			data.KeyLoaded = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["children"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Children, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["user_properties"]; ok && v != nil {
			data.UserProperties = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Locked = types.BoolValue(bv) }
		}
		if v, ok := resultMap["comments"]; ok && v != nil {
			data.Comments = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota_warning"]; ok && v != nil {
			data.QuotaWarning = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota_critical"]; ok && v != nil {
			data.QuotaCritical = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refquota_warning"]; ok && v != nil {
			data.RefquotaWarning = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refquota_critical"]; ok && v != nil {
			data.RefquotaCritical = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["managedby"]; ok && v != nil {
			data.Managedby = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["deduplication"]; ok && v != nil {
			data.Deduplication = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["aclmode"]; ok && v != nil {
			data.Aclmode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["acltype"]; ok && v != nil {
			data.Acltype = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["xattr"]; ok && v != nil {
			data.Xattr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["atime"]; ok && v != nil {
			data.Atime = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["casesensitivity"]; ok && v != nil {
			data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["checksum"]; ok && v != nil {
			data.Checksum = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["exec"]; ok && v != nil {
			data.Exec = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sync"]; ok && v != nil {
			data.Sync = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["compression"]; ok && v != nil {
			data.Compression = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["compressratio"]; ok && v != nil {
			data.Compressratio = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["origin"]; ok && v != nil {
			data.Origin = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota"]; ok && v != nil {
			data.Quota = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refquota"]; ok && v != nil {
			data.Refquota = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["reservation"]; ok && v != nil {
			data.Reservation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refreservation"]; ok && v != nil {
			data.Refreservation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["copies"]; ok && v != nil {
			data.Copies = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapdir"]; ok && v != nil {
			data.Snapdir = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			data.Readonly = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recordsize"]; ok && v != nil {
			data.Recordsize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sparse"]; ok && v != nil {
			data.Sparse = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["volsize"]; ok && v != nil {
			data.Volsize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["volblocksize"]; ok && v != nil {
			data.Volblocksize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["key_format"]; ok && v != nil {
			data.KeyFormat = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_algorithm"]; ok && v != nil {
			data.EncryptionAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["used"]; ok && v != nil {
			data.Used = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbychildren"]; ok && v != nil {
			data.Usedbychildren = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbydataset"]; ok && v != nil {
			data.Usedbydataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbyrefreservation"]; ok && v != nil {
			data.Usedbyrefreservation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbysnapshots"]; ok && v != nil {
			data.Usedbysnapshots = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["available"]; ok && v != nil {
			data.Available = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["special_small_block_size"]; ok && v != nil {
			data.SpecialSmallBlockSize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pbkdf2iters"]; ok && v != nil {
			data.Pbkdf2Iters = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["creation"]; ok && v != nil {
			data.Creation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapdev"]; ok && v != nil {
			data.Snapdev = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mountpoint"]; ok && v != nil {
			data.Mountpoint = types.StringValue(fmt.Sprintf("%v", v))
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
