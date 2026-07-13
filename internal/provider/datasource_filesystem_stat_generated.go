package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &FilesystemStatDataSource{}

func NewFilesystemStatDataSource() datasource.DataSource {
	return &FilesystemStatDataSource{}
}

type FilesystemStatDataSource struct {
	client *client.Client
}

type FilesystemStatDataSourceModel struct {
	Path           types.String  `tfsdk:"path"`
	Exists         types.Bool    `tfsdk:"exists"`
	Realpath       types.String  `tfsdk:"realpath"`
	Type           types.String  `tfsdk:"type"`
	Size           types.Int64   `tfsdk:"size"`
	AllocationSize types.Int64   `tfsdk:"allocation_size"`
	Mode           types.Int64   `tfsdk:"mode"`
	MountId        types.Int64   `tfsdk:"mount_id"`
	Uid            types.Int64   `tfsdk:"uid"`
	Gid            types.Int64   `tfsdk:"gid"`
	Atime          types.Float64 `tfsdk:"atime"`
	Mtime          types.Float64 `tfsdk:"mtime"`
	Ctime          types.Float64 `tfsdk:"ctime"`
	Btime          types.Float64 `tfsdk:"btime"`
	Dev            types.Int64   `tfsdk:"dev"`
	Inode          types.Int64   `tfsdk:"inode"`
	Nlink          types.Int64   `tfsdk:"nlink"`
	Acl            types.Bool    `tfsdk:"acl"`
	IsMountpoint   types.Bool    `tfsdk:"is_mountpoint"`
	IsCtldir       types.Bool    `tfsdk:"is_ctldir"`
	User           types.String  `tfsdk:"user"`
	Group          types.String  `tfsdk:"group"`
}

func (d *FilesystemStatDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_stat"
}

func (d *FilesystemStatDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Return filesystem information for a given path.",
		Attributes: map[string]schema.Attribute{
			"path":            schema.StringAttribute{Required: true, Description: "Absolute filesystem path to get statistics for."},
			"exists":          schema.BoolAttribute{Computed: true, Description: "Whether the path exists."},
			"realpath":        schema.StringAttribute{Computed: true, Description: "Canonical path of the entry, eliminating any symbolic links. "},
			"type":            schema.StringAttribute{Computed: true, Description: "Type of filesystem entry."},
			"size":            schema.Int64Attribute{Computed: true, Description: "Size in bytes of a plain file. This corresonds with stx_size. "},
			"allocation_size": schema.Int64Attribute{Computed: true, Description: "Allocated size of file. Calculated by multiplying stx_blocks by 512. "},
			"mode":            schema.Int64Attribute{Computed: true, Description: "Entry's mode including file type information and file permission bits. This corresponds with stx_mod"},
			"mount_id":        schema.Int64Attribute{Computed: true, Description: "The mount ID of the mount containing the entry. This corresponds to the number in first     field of"},
			"uid":             schema.Int64Attribute{Computed: true, Description: "User ID of the entry's owner. This corresponds with stx_uid. "},
			"gid":             schema.Int64Attribute{Computed: true, Description: "Group ID of the entry's owner. This corresponds with stx_gid. "},
			"atime":           schema.Float64Attribute{Computed: true, Description: "Time of last access. Corresponds with stx_atime. This is mutable from userspace. "},
			"mtime":           schema.Float64Attribute{Computed: true, Description: "Time of last modification. Corresponds with stx_mtime. This is mutable from userspace. "},
			"ctime":           schema.Float64Attribute{Computed: true, Description: "Time of last status change. Corresponds with stx_ctime. "},
			"btime":           schema.Float64Attribute{Computed: true, Description: "Time of creation. Corresponds with stx_btime. "},
			"dev":             schema.Int64Attribute{Computed: true, Description: "The ID of the device containing the filesystem where the file resides. This is not sufficient to uni"},
			"inode":           schema.Int64Attribute{Computed: true, Description: "The inode number of the file. This corresponds with stx_ino. "},
			"nlink":           schema.Int64Attribute{Computed: true, Description: "Number of hard links. Corresponds with stx_nlinks. "},
			"acl":             schema.BoolAttribute{Computed: true, Description: "Specifies whether ACL is present on the entry. If this is the case then file permission     bits as "},
			"is_mountpoint":   schema.BoolAttribute{Computed: true, Description: "Specifies whether the entry is also the mountpoint of a filesystem. "},
			"is_ctldir":       schema.BoolAttribute{Computed: true, Description: "Specifies whether the entry is located within the ZFS ctldir (for example a snapshot). "},
			"user":            schema.StringAttribute{Computed: true, Description: "Username associated with `uid`. Will be None if the User ID does not map to existing user. "},
			"group":           schema.StringAttribute{Computed: true, Description: "Groupname associated with `gid`. Will be None if the Group ID does not map to existing group. "},
		},
	}
}

func (d *FilesystemStatDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FilesystemStatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FilesystemStatDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("filesystem.stat", data.Path.ValueString())
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
	if v, ok := resultMap["realpath"]; ok {
		data.Realpath = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["type"]; ok {
		data.Type = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["size"]; ok {
		if fv, ok := v.(float64); ok {
			data.Size = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["allocation_size"]; ok {
		if fv, ok := v.(float64); ok {
			data.AllocationSize = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["mode"]; ok {
		if fv, ok := v.(float64); ok {
			data.Mode = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["mount_id"]; ok {
		if fv, ok := v.(float64); ok {
			data.MountId = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["uid"]; ok {
		if fv, ok := v.(float64); ok {
			data.Uid = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["gid"]; ok {
		if fv, ok := v.(float64); ok {
			data.Gid = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["atime"]; ok {
		if fv, ok := v.(float64); ok {
			data.Atime = types.Float64Value(fv)
		}
	}
	if v, ok := resultMap["mtime"]; ok {
		if fv, ok := v.(float64); ok {
			data.Mtime = types.Float64Value(fv)
		}
	}
	if v, ok := resultMap["ctime"]; ok {
		if fv, ok := v.(float64); ok {
			data.Ctime = types.Float64Value(fv)
		}
	}
	if v, ok := resultMap["btime"]; ok {
		if fv, ok := v.(float64); ok {
			data.Btime = types.Float64Value(fv)
		}
	}
	if v, ok := resultMap["dev"]; ok {
		if fv, ok := v.(float64); ok {
			data.Dev = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["inode"]; ok {
		if fv, ok := v.(float64); ok {
			data.Inode = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["nlink"]; ok {
		if fv, ok := v.(float64); ok {
			data.Nlink = types.Int64Value(int64(fv))
		}
	}
	if v, ok := resultMap["acl"]; ok {
		if bv, ok := v.(bool); ok {
			data.Acl = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["is_mountpoint"]; ok {
		if bv, ok := v.(bool); ok {
			data.IsMountpoint = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["is_ctldir"]; ok {
		if bv, ok := v.(bool); ok {
			data.IsCtldir = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["user"]; ok {
		data.User = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["group"]; ok {
		data.Group = types.StringValue(fmt.Sprintf("%v", v))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
