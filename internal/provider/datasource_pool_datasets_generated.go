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

var _ datasource.DataSource = &PoolDatasetsDataSource{}

func NewPoolDatasetsDataSource() datasource.DataSource {
	return &PoolDatasetsDataSource{}
}

type PoolDatasetsDataSource struct {
	client *client.Client
}

type PoolDatasetsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type PoolDatasetsItemModel struct {
	ID types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
	Pool types.String `tfsdk:"pool"`
	Encrypted types.Bool `tfsdk:"encrypted"`
	EncryptionRoot types.String `tfsdk:"encryption_root"`
	KeyLoaded types.String `tfsdk:"key_loaded"`
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

func (d *PoolDatasetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_datasets"
}

func (d *PoolDatasetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query Pool Datasets with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of pool_datasets resources",
				NestedObject: schema.NestedAttributeObject{
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
				},
			},
		},
	}
}

func (d *PoolDatasetsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolDatasetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PoolDatasetsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("pool.dataset.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query pool_datasets: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]PoolDatasetsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := PoolDatasetsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pool"]; ok && v != nil {
			itemModel.Pool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encrypted"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Encrypted = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_root"]; ok && v != nil {
			itemModel.EncryptionRoot = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["key_loaded"]; ok && v != nil {
			itemModel.KeyLoaded = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user_properties"]; ok && v != nil {
			itemModel.UserProperties = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Locked = types.BoolValue(bv) }
		}
		if v, ok := resultMap["comments"]; ok && v != nil {
			itemModel.Comments = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota_warning"]; ok && v != nil {
			itemModel.QuotaWarning = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota_critical"]; ok && v != nil {
			itemModel.QuotaCritical = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refquota_warning"]; ok && v != nil {
			itemModel.RefquotaWarning = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refquota_critical"]; ok && v != nil {
			itemModel.RefquotaCritical = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["managedby"]; ok && v != nil {
			itemModel.Managedby = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["deduplication"]; ok && v != nil {
			itemModel.Deduplication = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["aclmode"]; ok && v != nil {
			itemModel.Aclmode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["acltype"]; ok && v != nil {
			itemModel.Acltype = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["xattr"]; ok && v != nil {
			itemModel.Xattr = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["atime"]; ok && v != nil {
			itemModel.Atime = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["casesensitivity"]; ok && v != nil {
			itemModel.Casesensitivity = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["checksum"]; ok && v != nil {
			itemModel.Checksum = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["exec"]; ok && v != nil {
			itemModel.Exec = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sync"]; ok && v != nil {
			itemModel.Sync = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["compression"]; ok && v != nil {
			itemModel.Compression = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["compressratio"]; ok && v != nil {
			itemModel.Compressratio = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["origin"]; ok && v != nil {
			itemModel.Origin = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota"]; ok && v != nil {
			itemModel.Quota = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refquota"]; ok && v != nil {
			itemModel.Refquota = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["reservation"]; ok && v != nil {
			itemModel.Reservation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["refreservation"]; ok && v != nil {
			itemModel.Refreservation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["copies"]; ok && v != nil {
			itemModel.Copies = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapdir"]; ok && v != nil {
			itemModel.Snapdir = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			itemModel.Readonly = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recordsize"]; ok && v != nil {
			itemModel.Recordsize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sparse"]; ok && v != nil {
			itemModel.Sparse = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["volsize"]; ok && v != nil {
			itemModel.Volsize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["volblocksize"]; ok && v != nil {
			itemModel.Volblocksize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["key_format"]; ok && v != nil {
			itemModel.KeyFormat = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_algorithm"]; ok && v != nil {
			itemModel.EncryptionAlgorithm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["used"]; ok && v != nil {
			itemModel.Used = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbychildren"]; ok && v != nil {
			itemModel.Usedbychildren = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbydataset"]; ok && v != nil {
			itemModel.Usedbydataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbyrefreservation"]; ok && v != nil {
			itemModel.Usedbyrefreservation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["usedbysnapshots"]; ok && v != nil {
			itemModel.Usedbysnapshots = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["available"]; ok && v != nil {
			itemModel.Available = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["special_small_block_size"]; ok && v != nil {
			itemModel.SpecialSmallBlockSize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pbkdf2iters"]; ok && v != nil {
			itemModel.Pbkdf2Iters = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["creation"]; ok && v != nil {
			itemModel.Creation = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapdev"]; ok && v != nil {
			itemModel.Snapdev = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mountpoint"]; ok && v != nil {
			itemModel.Mountpoint = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"aclmode": types.StringType,
			"acltype": types.StringType,
			"atime": types.StringType,
			"available": types.StringType,
			"casesensitivity": types.StringType,
			"checksum": types.StringType,
			"comments": types.StringType,
			"compression": types.StringType,
			"compressratio": types.StringType,
			"copies": types.StringType,
			"creation": types.StringType,
			"deduplication": types.StringType,
			"encrypted": types.BoolType,
			"encryption_algorithm": types.StringType,
			"encryption_root": types.StringType,
			"exec": types.StringType,
			"id": types.StringType,
			"key_format": types.StringType,
			"key_loaded": types.StringType,
			"locked": types.BoolType,
			"managedby": types.StringType,
			"mountpoint": types.StringType,
			"name": types.StringType,
			"origin": types.StringType,
			"pbkdf2iters": types.StringType,
			"pool": types.StringType,
			"quota": types.StringType,
			"quota_critical": types.StringType,
			"quota_warning": types.StringType,
			"readonly": types.StringType,
			"recordsize": types.StringType,
			"refquota": types.StringType,
			"refquota_critical": types.StringType,
			"refquota_warning": types.StringType,
			"refreservation": types.StringType,
			"reservation": types.StringType,
			"snapdev": types.StringType,
			"snapdir": types.StringType,
			"sparse": types.StringType,
			"special_small_block_size": types.StringType,
			"sync": types.StringType,
			"type": types.StringType,
			"used": types.StringType,
			"usedbychildren": types.StringType,
			"usedbydataset": types.StringType,
			"usedbyrefreservation": types.StringType,
			"usedbysnapshots": types.StringType,
			"user_properties": types.StringType,
			"volblocksize": types.StringType,
			"volsize": types.StringType,
			"xattr": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
