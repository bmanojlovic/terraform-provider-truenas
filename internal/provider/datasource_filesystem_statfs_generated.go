package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &FilesystemStatfsDataSource{}

func NewFilesystemStatfsDataSource() datasource.DataSource {
	return &FilesystemStatfsDataSource{}
}

type FilesystemStatfsDataSource struct {
	client *client.Client
}

type FilesystemStatfsDataSourceModel struct {
	Path           types.String `tfsdk:"path"`
	Exists         types.Bool   `tfsdk:"exists"`
	Fsid           types.String `tfsdk:"fsid"`
	Fstype         types.String `tfsdk:"fstype"`
	Source         types.String `tfsdk:"source"`
	Dest           types.String `tfsdk:"dest"`
	Blocksize      types.Int64  `tfsdk:"blocksize"`
	TotalBlocks    types.Int64  `tfsdk:"total_blocks"`
	FreeBlocks     types.Int64  `tfsdk:"free_blocks"`
	AvailBlocks    types.Int64  `tfsdk:"avail_blocks"`
	TotalBlocksStr types.String `tfsdk:"total_blocks_str"`
	FreeBlocksStr  types.String `tfsdk:"free_blocks_str"`
	AvailBlocksStr types.String `tfsdk:"avail_blocks_str"`
	Files          types.Int64  `tfsdk:"files"`
	FreeFiles      types.Int64  `tfsdk:"free_files"`
	NameMax        types.Int64  `tfsdk:"name_max"`
	TotalBytes     types.Int64  `tfsdk:"total_bytes"`
	FreeBytes      types.Int64  `tfsdk:"free_bytes"`
	AvailBytes     types.Int64  `tfsdk:"avail_bytes"`
	TotalBytesStr  types.String `tfsdk:"total_bytes_str"`
	FreeBytesStr   types.String `tfsdk:"free_bytes_str"`
	AvailBytesStr  types.String `tfsdk:"avail_bytes_str"`
}

func (d *FilesystemStatfsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_statfs"
}

func (d *FilesystemStatfsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Return stats from the filesystem of a given path.",
		Attributes: map[string]schema.Attribute{
			"path":             schema.StringAttribute{Required: true, Description: "Path on the filesystem to get statistics for."},
			"exists":           schema.BoolAttribute{Computed: true, Description: "Whether the path exists."},
			"fsid":             schema.StringAttribute{Computed: true, Description: "Unique filesystem ID as returned by statvfs. "},
			"fstype":           schema.StringAttribute{Computed: true, Description: "String representation of filesystem type from mountinfo. "},
			"source":           schema.StringAttribute{Computed: true, Description: "Source for the mounted filesystem. For ZFS this will be dataset name. "},
			"dest":             schema.StringAttribute{Computed: true, Description: "Local path on which filesystem is mounted. "},
			"blocksize":        schema.Int64Attribute{Computed: true, Description: "Filesystem block size as reported by statvfs. "},
			"total_blocks":     schema.Int64Attribute{Computed: true, Description: "Filesystem size as reported in blocksize blocks as reported by statvfs. "},
			"free_blocks":      schema.Int64Attribute{Computed: true, Description: "Number of free blocks as reported by statvfs. "},
			"avail_blocks":     schema.Int64Attribute{Computed: true, Description: "Number of available blocks as reported by statvfs. "},
			"total_blocks_str": schema.StringAttribute{Computed: true, Description: "Total filesystem size in blocks as a human-readable string."},
			"free_blocks_str":  schema.StringAttribute{Computed: true, Description: "Free blocks available as a human-readable string."},
			"avail_blocks_str": schema.StringAttribute{Computed: true, Description: "Available blocks for unprivileged users as a human-readable string."},
			"files":            schema.Int64Attribute{Computed: true, Description: "Number of inodes in use as reported by statvfs. "},
			"free_files":       schema.Int64Attribute{Computed: true, Description: "Number of free inodes as reported by statvfs. "},
			"name_max":         schema.Int64Attribute{Computed: true, Description: "Maximum filename length as reported by statvfs. "},
			"total_bytes":      schema.Int64Attribute{Computed: true, Description: "Total filesystem size in bytes."},
			"free_bytes":       schema.Int64Attribute{Computed: true, Description: "Free space available in bytes."},
			"avail_bytes":      schema.Int64Attribute{Computed: true, Description: "Available space for unprivileged users in bytes."},
			"total_bytes_str":  schema.StringAttribute{Computed: true, Description: "Total filesystem size in bytes as a human-readable string."},
			"free_bytes_str":   schema.StringAttribute{Computed: true, Description: "Free space available in bytes as a human-readable string."},
			"avail_bytes_str":  schema.StringAttribute{Computed: true, Description: "Available space for unprivileged users in bytes as a human-readable string."},
		},
	}
}

func (d *FilesystemStatfsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FilesystemStatfsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FilesystemStatfsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("filesystem.statfs", data.Path.ValueString())
	if err != nil {
		data.Exists = types.BoolValue(false)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	data.Exists = types.BoolValue(true)
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}
	if v, ok := resultMap["fsid"]; ok {
		data.Fsid = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["fstype"]; ok {
		data.Fstype = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["source"]; ok {
		data.Source = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["dest"]; ok {
		data.Dest = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["blocksize"]; ok {
		if fv, ok := v.(float64); ok {
			data.Blocksize = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["total_blocks"]; ok {
		if fv, ok := v.(float64); ok {
			data.TotalBlocks = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["free_blocks"]; ok {
		if fv, ok := v.(float64); ok {
			data.FreeBlocks = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["avail_blocks"]; ok {
		if fv, ok := v.(float64); ok {
			data.AvailBlocks = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["total_blocks_str"]; ok {
		data.TotalBlocksStr = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["free_blocks_str"]; ok {
		data.FreeBlocksStr = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["avail_blocks_str"]; ok {
		data.AvailBlocksStr = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["files"]; ok {
		if fv, ok := v.(float64); ok {
			data.Files = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["free_files"]; ok {
		if fv, ok := v.(float64); ok {
			data.FreeFiles = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["name_max"]; ok {
		if fv, ok := v.(float64); ok {
			data.NameMax = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["total_bytes"]; ok {
		if fv, ok := v.(float64); ok {
			data.TotalBytes = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["free_bytes"]; ok {
		if fv, ok := v.(float64); ok {
			data.FreeBytes = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["avail_bytes"]; ok {
		if fv, ok := v.(float64); ok {
			data.AvailBytes = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["total_bytes_str"]; ok {
		data.TotalBytesStr = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["free_bytes_str"]; ok {
		data.FreeBytesStr = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["avail_bytes_str"]; ok {
		data.AvailBytesStr = types.StringValue(fmt.Sprintf("%v", v))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
