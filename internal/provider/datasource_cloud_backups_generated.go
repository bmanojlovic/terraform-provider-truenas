package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &CloudBackupsDataSource{}

func NewCloudBackupsDataSource() datasource.DataSource {
	return &CloudBackupsDataSource{}
}

type CloudBackupsDataSource struct {
	client *client.Client
}

type CloudBackupsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type CloudBackupsItemModel struct {
	ID              types.String `tfsdk:"id"`
	Description     types.String `tfsdk:"description"`
	Path            types.String `tfsdk:"path"`
	Credentials     types.String `tfsdk:"credentials"`
	Attributes      types.String `tfsdk:"attributes"`
	Schedule        types.String `tfsdk:"schedule"`
	PreScript       types.String `tfsdk:"pre_script"`
	PostScript      types.String `tfsdk:"post_script"`
	Snapshot        types.Bool   `tfsdk:"snapshot"`
	Args            types.String `tfsdk:"args"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Job             types.String `tfsdk:"job"`
	Locked          types.Bool   `tfsdk:"locked"`
	Password        types.String `tfsdk:"password"`
	KeepLast        types.Int64  `tfsdk:"keep_last"`
	TransferSetting types.String `tfsdk:"transfer_setting"`
	AbsolutePaths   types.Bool   `tfsdk:"absolute_paths"`
	CachePath       types.String `tfsdk:"cache_path"`
	RateLimit       types.Int64  `tfsdk:"rate_limit"`
}

func (d *CloudBackupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_backups"
}

func (d *CloudBackupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query cloud_backups",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of cloud_backups resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the task to display in the UI.",
						},
						"path": schema.StringAttribute{
							Computed:    true,
							Description: "The local path to back up beginning with `/mnt` or `/dev/zvol`.",
						},
						"credentials": schema.StringAttribute{
							Computed:    true,
							Description: "Cloud credentials to use for each backup.",
						},
						"attributes": schema.StringAttribute{
							Computed:    true,
							Description: "Additional information for each backup, e.g. bucket name.",
						},
						"schedule": schema.StringAttribute{
							Computed:    true,
							Description: "Cron schedule dictating when the task should run.",
						},
						"pre_script": schema.StringAttribute{
							Computed:    true,
							Description: "A Bash script to run immediately before every backup.",
						},
						"post_script": schema.StringAttribute{
							Computed:    true,
							Description: "A Bash script to run immediately after every backup if it succeeds.",
						},
						"snapshot": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to create a temporary snapshot of the dataset before every backup.",
						},
						"args": schema.StringAttribute{
							Computed:    true,
							Description: "(Slated for removal).",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Can enable/disable the task.",
						},
						"job": schema.StringAttribute{
							Computed:    true,
							Description: "Information regarding the task's job state, e.g. progress.",
						},
						"locked": schema.BoolAttribute{
							Computed:    true,
							Description: "A locked task cannot run.",
						},
						"password": schema.StringAttribute{
							Computed:    true,
							Description: "Password for the remote repository.",
						},
						"keep_last": schema.Int64Attribute{
							Computed:    true,
							Description: "How many of the most recent backup snapshots to keep after each backup.",
						},
						"transfer_setting": schema.StringAttribute{
							Computed:    true,
							Description: "* DEFAULT:     * pack size given by `$RESTIC_PACK_SIZE` (default 16 MiB)     * read concurrency give",
						},
						"absolute_paths": schema.BoolAttribute{
							Computed:    true,
							Description: "Preserve absolute paths in each backup (cannot be set when `snapshot=True`).",
						},
						"cache_path": schema.StringAttribute{
							Computed:    true,
							Description: "Cache path. If not set, performance may degrade.",
						},
						"rate_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum upload/download rate in KiB/s. Passed to `restic --limit-upload` on `cloud_backup.sync` and ",
						},
					},
				},
			},
		},
	}
}

func (d *CloudBackupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudBackupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudBackupsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("cloud_backup.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query cloud_backups: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]CloudBackupsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := CloudBackupsItemModel{}
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
			if bv, ok := v.(bool); ok {
				itemModel.Snapshot = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["args"]; ok && v != nil {
			itemModel.Args = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["job"]; ok && v != nil {
			itemModel.Job = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Locked = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			itemModel.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["keep_last"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.KeepLast = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["transfer_setting"]; ok && v != nil {
			itemModel.TransferSetting = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["absolute_paths"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.AbsolutePaths = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["cache_path"]; ok && v != nil {
			itemModel.CachePath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["rate_limit"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.RateLimit = types.Int64Value(int64(fv))
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"absolute_paths":   types.BoolType,
			"args":             types.StringType,
			"attributes":       types.StringType,
			"cache_path":       types.StringType,
			"credentials":      types.StringType,
			"description":      types.StringType,
			"enabled":          types.BoolType,
			"id":               types.StringType,
			"job":              types.StringType,
			"keep_last":        types.Int64Type,
			"locked":           types.BoolType,
			"password":         types.StringType,
			"path":             types.StringType,
			"post_script":      types.StringType,
			"pre_script":       types.StringType,
			"rate_limit":       types.Int64Type,
			"schedule":         types.StringType,
			"snapshot":         types.BoolType,
			"transfer_setting": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
