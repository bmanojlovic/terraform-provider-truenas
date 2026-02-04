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

var _ datasource.DataSource = &RsynctasksDataSource{}

func NewRsynctasksDataSource() datasource.DataSource {
	return &RsynctasksDataSource{}
}

type RsynctasksDataSource struct {
	client *client.Client
}

type RsynctasksDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type RsynctasksItemModel struct {
	ID             types.String `tfsdk:"id"`
	Path           types.String `tfsdk:"path"`
	User           types.String `tfsdk:"user"`
	Mode           types.String `tfsdk:"mode"`
	Remotehost     types.String `tfsdk:"remotehost"`
	Remoteport     types.Int64  `tfsdk:"remoteport"`
	Remotemodule   types.String `tfsdk:"remotemodule"`
	SshCredentials types.String `tfsdk:"ssh_credentials"`
	Remotepath     types.String `tfsdk:"remotepath"`
	Direction      types.String `tfsdk:"direction"`
	Desc           types.String `tfsdk:"desc"`
	Schedule       types.String `tfsdk:"schedule"`
	Recursive      types.Bool   `tfsdk:"recursive"`
	Times          types.Bool   `tfsdk:"times"`
	Compress       types.Bool   `tfsdk:"compress"`
	Archive        types.Bool   `tfsdk:"archive"`
	Delete         types.Bool   `tfsdk:"delete"`
	Quiet          types.Bool   `tfsdk:"quiet"`
	Preserveperm   types.Bool   `tfsdk:"preserveperm"`
	Preserveattr   types.Bool   `tfsdk:"preserveattr"`
	Delayupdates   types.Bool   `tfsdk:"delayupdates"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Locked         types.Bool   `tfsdk:"locked"`
	Job            types.String `tfsdk:"job"`
}

func (d *RsynctasksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rsynctasks"
}

func (d *RsynctasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query rsynctasks",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of rsynctasks resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"path": schema.StringAttribute{
							Computed:    true,
							Description: "Local filesystem path to synchronize.",
						},
						"user": schema.StringAttribute{
							Computed:    true,
							Description: "Username to run the rsync task as.",
						},
						"mode": schema.StringAttribute{
							Computed:    true,
							Description: "Operating mechanism for Rsync, i.e. Rsync Module mode or Rsync SSH mode.",
						},
						"remotehost": schema.StringAttribute{
							Computed:    true,
							Description: "IP address or hostname of the remote system. If username differs on the remote host, \"username@remot",
						},
						"remoteport": schema.Int64Attribute{
							Computed:    true,
							Description: "Port number for SSH connection. Only applies when `mode` is SSH.",
						},
						"remotemodule": schema.StringAttribute{
							Computed:    true,
							Description: "Name of remote module, this attribute should be specified when `mode` is set to MODULE.",
						},
						"ssh_credentials": schema.StringAttribute{
							Computed:    true,
							Description: "In SSH mode, if `ssh_credentials` (a keychain credential of `SSH_CREDENTIALS` type) is specified the",
						},
						"remotepath": schema.StringAttribute{
							Computed:    true,
							Description: "Path on the remote system to synchronize with.",
						},
						"direction": schema.StringAttribute{
							Computed:    true,
							Description: "Specify if data should be PULLED or PUSHED from the remote system.",
						},
						"desc": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the rsync task.",
						},
						"schedule": schema.StringAttribute{
							Computed:    true,
							Description: "Cron schedule for when the rsync task should run.",
						},
						"recursive": schema.BoolAttribute{
							Computed:    true,
							Description: "Recursively transfer subdirectories.",
						},
						"times": schema.BoolAttribute{
							Computed:    true,
							Description: "Preserve modification times of files.",
						},
						"compress": schema.BoolAttribute{
							Computed:    true,
							Description: "Reduce the size of the data to be transmitted.",
						},
						"archive": schema.BoolAttribute{
							Computed:    true,
							Description: "Make rsync run recursively, preserving symlinks, permissions, modification times, group, and special",
						},
						"delete": schema.BoolAttribute{
							Computed:    true,
							Description: "Delete files in the destination directory that do not exist in the source directory.",
						},
						"quiet": schema.BoolAttribute{
							Computed:    true,
							Description: "Suppress informational messages from rsync.",
						},
						"preserveperm": schema.BoolAttribute{
							Computed:    true,
							Description: "Preserve original file permissions.",
						},
						"preserveattr": schema.BoolAttribute{
							Computed:    true,
							Description: "Preserve extended attributes of files.",
						},
						"delayupdates": schema.BoolAttribute{
							Computed:    true,
							Description: "Delay updating destination files until all transfers are complete.",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this rsync task is enabled.",
						},
						"locked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this rsync task is currently locked (running).",
						},
						"job": schema.StringAttribute{
							Computed:    true,
							Description: "Information about the currently running job. `null` if no job is running.",
						},
					},
				},
			},
		},
	}
}

func (d *RsynctasksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RsynctasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RsynctasksDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("rsynctask.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query rsynctasks: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]RsynctasksItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := RsynctasksItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			itemModel.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user"]; ok && v != nil {
			itemModel.User = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mode"]; ok && v != nil {
			itemModel.Mode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["remotehost"]; ok && v != nil {
			itemModel.Remotehost = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["remoteport"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Remoteport = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["remotemodule"]; ok && v != nil {
			itemModel.Remotemodule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ssh_credentials"]; ok && v != nil {
			itemModel.SshCredentials = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["remotepath"]; ok && v != nil {
			itemModel.Remotepath = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["direction"]; ok && v != nil {
			itemModel.Direction = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["desc"]; ok && v != nil {
			itemModel.Desc = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			itemModel.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["recursive"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Recursive = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["times"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Times = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["compress"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Compress = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["archive"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Archive = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["delete"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Delete = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["quiet"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Quiet = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["preserveperm"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Preserveperm = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["preserveattr"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Preserveattr = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["delayupdates"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Delayupdates = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Locked = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["job"]; ok && v != nil {
			itemModel.Job = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"archive":         types.BoolType,
			"compress":        types.BoolType,
			"delayupdates":    types.BoolType,
			"delete":          types.BoolType,
			"desc":            types.StringType,
			"direction":       types.StringType,
			"enabled":         types.BoolType,
			"id":              types.StringType,
			"job":             types.StringType,
			"locked":          types.BoolType,
			"mode":            types.StringType,
			"path":            types.StringType,
			"preserveattr":    types.BoolType,
			"preserveperm":    types.BoolType,
			"quiet":           types.BoolType,
			"recursive":       types.BoolType,
			"remotehost":      types.StringType,
			"remotemodule":    types.StringType,
			"remotepath":      types.StringType,
			"remoteport":      types.Int64Type,
			"schedule":        types.StringType,
			"ssh_credentials": types.StringType,
			"times":           types.BoolType,
			"user":            types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
