package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &FilesystemGetZfsAttributesDataSource{}

func NewFilesystemGetZfsAttributesDataSource() datasource.DataSource {
	return &FilesystemGetZfsAttributesDataSource{}
}

type FilesystemGetZfsAttributesDataSource struct {
	client *client.Client
}

type FilesystemGetZfsAttributesDataSourceModel struct {
	Path       types.String `tfsdk:"path"`
	Exists     types.Bool   `tfsdk:"exists"`
	Readonly   types.Bool   `tfsdk:"readonly"`
	Hidden     types.Bool   `tfsdk:"hidden"`
	System     types.Bool   `tfsdk:"system"`
	Archive    types.Bool   `tfsdk:"archive"`
	Immutable  types.Bool   `tfsdk:"immutable"`
	Nounlink   types.Bool   `tfsdk:"nounlink"`
	Appendonly types.Bool   `tfsdk:"appendonly"`
	Offline    types.Bool   `tfsdk:"offline"`
	Sparse     types.Bool   `tfsdk:"sparse"`
}

func (d *FilesystemGetZfsAttributesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_get_zfs_attributes"
}

func (d *FilesystemGetZfsAttributesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get the current ZFS attributes for the file at the given path",
		Attributes: map[string]schema.Attribute{
			"path":       schema.StringAttribute{Required: true, Description: "Path to the file to get ZFS attributes for."},
			"exists":     schema.BoolAttribute{Computed: true, Description: "Whether the path exists."},
			"readonly":   schema.BoolAttribute{Computed: true, Description: "READONLY MS-DOS attribute. When set, file may not be written to (toggling     does not impact existi"},
			"hidden":     schema.BoolAttribute{Computed: true, Description: "HIDDEN MS-DOS attribute. When set, the SMB HIDDEN flag is set and file     is \"hidden\" from the pers"},
			"system":     schema.BoolAttribute{Computed: true, Description: "SYSTEM MS-DOS attribute. Is presented to SMB clients, but has no impact on local filesystem. "},
			"archive":    schema.BoolAttribute{Computed: true, Description: "ARCHIVE MS-DOS attribute. Value is reset to True whenever file is modified. "},
			"immutable":  schema.BoolAttribute{Computed: true, Description: "File may not be altered or deleted. Also appears as IMMUTABLE in attributes in     `filesystem.stat`"},
			"nounlink":   schema.BoolAttribute{Computed: true, Description: "File may be altered but not deleted. "},
			"appendonly": schema.BoolAttribute{Computed: true, Description: "File may only be opened with O_APPEND flag. Also appears as APPEND in     attributes in `filesystem."},
			"offline":    schema.BoolAttribute{Computed: true, Description: "OFFLINE MS-DOS attribute. Is presented to SMB clients, but has no impact on local filesystem. "},
			"sparse":     schema.BoolAttribute{Computed: true, Description: "SPARSE MS-DOS attribute. Is presented to SMB clients, but has no impact on local filesystem. "},
		},
	}
}

func (d *FilesystemGetZfsAttributesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FilesystemGetZfsAttributesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FilesystemGetZfsAttributesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("filesystem.get_zfs_attributes", data.Path.ValueString())
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
	if v, ok := resultMap["readonly"]; ok {
		if bv, ok := v.(bool); ok {
			data.Readonly = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["hidden"]; ok {
		if bv, ok := v.(bool); ok {
			data.Hidden = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["system"]; ok {
		if bv, ok := v.(bool); ok {
			data.System = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["archive"]; ok {
		if bv, ok := v.(bool); ok {
			data.Archive = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["immutable"]; ok {
		if bv, ok := v.(bool); ok {
			data.Immutable = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["nounlink"]; ok {
		if bv, ok := v.(bool); ok {
			data.Nounlink = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["appendonly"]; ok {
		if bv, ok := v.(bool); ok {
			data.Appendonly = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["offline"]; ok {
		if bv, ok := v.(bool); ok {
			data.Offline = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["sparse"]; ok {
		if bv, ok := v.(bool); ok {
			data.Sparse = types.BoolValue(bv)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
