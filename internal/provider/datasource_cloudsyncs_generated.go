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

var _ datasource.DataSource = &CloudsyncsDataSource{}

func NewCloudsyncsDataSource() datasource.DataSource {
	return &CloudsyncsDataSource{}
}

type CloudsyncsDataSource struct {
	client *client.Client
}

type CloudsyncsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type CloudsyncsItemModel struct {
	ID types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Path types.String `tfsdk:"path"`
	Credentials types.String `tfsdk:"credentials"`
	Attributes types.String `tfsdk:"attributes"`
	Schedule types.String `tfsdk:"schedule"`
	PreScript types.String `tfsdk:"pre_script"`
	PostScript types.String `tfsdk:"post_script"`
	Snapshot types.Bool `tfsdk:"snapshot"`
	Args types.String `tfsdk:"args"`
	Enabled types.Bool `tfsdk:"enabled"`
	Job types.String `tfsdk:"job"`
	Locked types.Bool `tfsdk:"locked"`
	Transfers types.Int64 `tfsdk:"transfers"`
	Direction types.String `tfsdk:"direction"`
	TransferMode types.String `tfsdk:"transfer_mode"`
	Encryption types.Bool `tfsdk:"encryption"`
	FilenameEncryption types.Bool `tfsdk:"filename_encryption"`
	EncryptionPassword types.String `tfsdk:"encryption_password"`
	EncryptionSalt types.String `tfsdk:"encryption_salt"`
	CreateEmptySrcDirs types.Bool `tfsdk:"create_empty_src_dirs"`
	FollowSymlinks types.Bool `tfsdk:"follow_symlinks"`
}

func (d *CloudsyncsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudsyncs"
}

func (d *CloudsyncsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query cloudsyncs",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of cloudsyncs resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "The name of the task to display in the UI.",
			},
			"path": schema.StringAttribute{
				Computed: true,
				Description: "The local path to back up beginning with `/mnt` or `/dev/zvol`.",
			},
			"credentials": schema.StringAttribute{
				Computed: true,
				Description: "Cloud credentials to use for each backup.",
			},
			"attributes": schema.StringAttribute{
				Computed: true,
				Description: "Additional information for each backup, e.g. bucket name.",
			},
			"schedule": schema.StringAttribute{
				Computed: true,
				Description: "Cron schedule dictating when the task should run.",
			},
			"pre_script": schema.StringAttribute{
				Computed: true,
				Description: "A Bash script to run immediately before every backup.",
			},
			"post_script": schema.StringAttribute{
				Computed: true,
				Description: "A Bash script to run immediately after every backup if it succeeds.",
			},
			"snapshot": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to create a temporary snapshot of the dataset before every backup.",
			},
			"args": schema.StringAttribute{
				Computed: true,
				Description: "(Slated for removal).",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Can enable/disable the task.",
			},
			"job": schema.StringAttribute{
				Computed: true,
				Description: "Information regarding the task's job state, e.g. progress.",
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				Description: "A locked task cannot run.",
			},
			"transfers": schema.Int64Attribute{
				Computed: true,
				Description: "Maximum number of parallel file transfers. `null` for default.",
			},
			"direction": schema.StringAttribute{
				Computed: true,
				Description: "Direction of the cloud sync operation.  * `PUSH`: Upload local files to cloud storage * `PULL`: Down",
			},
			"transfer_mode": schema.StringAttribute{
				Computed: true,
				Description: "How files are transferred between local and cloud storage.  * `SYNC`: Synchronize directories (add n",
			},
			"encryption": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to encrypt files before uploading to cloud storage.",
			},
			"filename_encryption": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to encrypt filenames in addition to file contents.",
			},
			"encryption_password": schema.StringAttribute{
				Computed: true,
				Description: "Password for client-side encryption. Empty string if encryption is disabled.",
			},
			"encryption_salt": schema.StringAttribute{
				Computed: true,
				Description: "Salt value for encryption key derivation. Empty string if encryption is disabled.",
			},
			"create_empty_src_dirs": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to create empty directories in the destination that exist in the source.",
			},
			"follow_symlinks": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to follow symbolic links and sync the files they point to.",
			},
					},
				},
			},
		},
	}
}

func (d *CloudsyncsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudsyncsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudsyncsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("cloudsync.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query cloudsyncs: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]CloudsyncsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := CloudsyncsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			itemModel.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["credentials"]; ok && v != nil {
			itemModel.Credentials = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			itemModel.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			itemModel.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pre_script"]; ok && v != nil {
			itemModel.PreScript = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["post_script"]; ok && v != nil {
			itemModel.PostScript = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapshot"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Snapshot = types.BoolValue(bv) }
		}
		if v, ok := resultMap["args"]; ok && v != nil {
			itemModel.Args = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["job"]; ok && v != nil {
			itemModel.Job = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Locked = types.BoolValue(bv) }
		}
		if v, ok := resultMap["transfers"]; ok && v != nil {
			if fv, ok := v.(float64); ok { itemModel.Transfers = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["direction"]; ok && v != nil {
			itemModel.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["transfer_mode"]; ok && v != nil {
			itemModel.TransferMode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Encryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["filename_encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.FilenameEncryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["encryption_password"]; ok && v != nil {
			itemModel.EncryptionPassword = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_salt"]; ok && v != nil {
			itemModel.EncryptionSalt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["create_empty_src_dirs"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.CreateEmptySrcDirs = types.BoolValue(bv) }
		}
		if v, ok := resultMap["follow_symlinks"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.FollowSymlinks = types.BoolValue(bv) }
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"args": types.StringType,
			"attributes": types.StringType,
			"create_empty_src_dirs": types.BoolType,
			"credentials": types.StringType,
			"description": types.StringType,
			"direction": types.StringType,
			"enabled": types.BoolType,
			"encryption": types.BoolType,
			"encryption_password": types.StringType,
			"encryption_salt": types.StringType,
			"filename_encryption": types.BoolType,
			"follow_symlinks": types.BoolType,
			"id": types.StringType,
			"job": types.StringType,
			"locked": types.BoolType,
			"path": types.StringType,
			"post_script": types.StringType,
			"pre_script": types.StringType,
			"schedule": types.StringType,
			"snapshot": types.BoolType,
			"transfer_mode": types.StringType,
			"transfers": types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
