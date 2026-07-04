package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type IscsiExtentResource struct {
	client *client.Client
}

type IscsiExtentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Disk           types.String `tfsdk:"disk"`
	Serial         types.String `tfsdk:"serial"`
	Path           types.String `tfsdk:"path"`
	Filesize       types.Int64  `tfsdk:"filesize"`
	Blocksize      types.Int64  `tfsdk:"blocksize"`
	Pblocksize     types.Bool   `tfsdk:"pblocksize"`
	AvailThreshold types.Int64  `tfsdk:"avail_threshold"`
	Comment        types.String `tfsdk:"comment"`
	InsecureTpc    types.Bool   `tfsdk:"insecure_tpc"`
	Xen            types.Bool   `tfsdk:"xen"`
	Rpm            types.String `tfsdk:"rpm"`
	Ro             types.Bool   `tfsdk:"ro"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	ProductId      types.String `tfsdk:"product_id"`
}

func NewIscsiExtentResource() resource.Resource {
	return &IscsiExtentResource{}
}

func (r *IscsiExtentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_extent"
}

func (r *IscsiExtentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiExtentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Extent.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Name of the iSCSI extent.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of the extent storage backend.",
			},
			"disk": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Disk device to use for the extent or `null` if using a file.",
			},
			"serial": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Serial number for the extent or `null` to auto-generate.",
			},
			"path": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "File path for file-based extents or `null` if using a disk.",
			},
			"filesize": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Size of the file-based extent in bytes.",
			},
			"blocksize": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Block size for the extent in bytes.",
			},
			"pblocksize": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to use physical block size reporting.",
			},
			"avail_threshold": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Available space threshold percentage or `null` to disable.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional comment describing the extent.",
			},
			"insecure_tpc": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable insecure Third Party Copy (TPC) operations.",
			},
			"xen": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable Xen compatibility mode.",
			},
			"rpm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Reported RPM type for the extent.",
			},
			"ro": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the extent is read-only.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the extent is enabled and available for use.",
			},
			"product_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Product ID string for the extent or `null` for default.",
			},
		},
	}
}

func (r *IscsiExtentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiExtentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Disk.IsNull() && !data.Disk.IsUnknown() {
		params["disk"] = data.Disk.ValueString()
	}
	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {
		params["serial"] = data.Serial.ValueString()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Filesize.IsNull() && !data.Filesize.IsUnknown() {
		params["filesize"] = data.Filesize.ValueInt64()
	}
	if !data.Blocksize.IsNull() && !data.Blocksize.IsUnknown() {
		params["blocksize"] = data.Blocksize.ValueInt64()
	}
	if !data.Pblocksize.IsNull() && !data.Pblocksize.IsUnknown() {
		params["pblocksize"] = data.Pblocksize.ValueBool()
	}
	if !data.AvailThreshold.IsNull() && !data.AvailThreshold.IsUnknown() {
		params["avail_threshold"] = data.AvailThreshold.ValueInt64()
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.InsecureTpc.IsNull() && !data.InsecureTpc.IsUnknown() {
		params["insecure_tpc"] = data.InsecureTpc.ValueBool()
	}
	if !data.Xen.IsNull() && !data.Xen.IsUnknown() {
		params["xen"] = data.Xen.ValueBool()
	}
	if !data.Rpm.IsNull() && !data.Rpm.IsUnknown() {
		params["rpm"] = data.Rpm.ValueString()
	}
	if !data.Ro.IsNull() && !data.Ro.IsUnknown() {
		params["ro"] = data.Ro.ValueBool()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ProductId.IsNull() && !data.ProductId.IsUnknown() {
		params["product_id"] = data.ProductId.ValueString()
	}

	result, err := r.client.Call("iscsi.extent.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_extent: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	// Read back to populate computed fields
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("iscsi.extent.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back iscsi_extent: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["disk"]; ok {
		if v == nil {
			data.Disk = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Disk = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Disk = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Disk = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["serial"]; ok {
		if v == nil {
			data.Serial = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Serial = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Serial = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Serial = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["path"]; ok {
		if v == nil {
			data.Path = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["avail_threshold"]; ok {
		if v == nil {
			data.AvailThreshold = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.AvailThreshold = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.AvailThreshold = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["product_id"]; ok {
		if v == nil {
			data.ProductId = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ProductId = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ProductId = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ProductId = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.Disk.IsUnknown() {
		data.Disk = types.StringNull()
	}
	if data.Serial.IsUnknown() {
		data.Serial = types.StringNull()
	}
	if data.Path.IsUnknown() {
		data.Path = types.StringNull()
	}
	if data.Filesize.IsUnknown() {
		data.Filesize = types.Int64Null()
	}
	if data.Blocksize.IsUnknown() {
		data.Blocksize = types.Int64Null()
	}
	if data.Pblocksize.IsUnknown() {
		data.Pblocksize = types.BoolNull()
	}
	if data.AvailThreshold.IsUnknown() {
		data.AvailThreshold = types.Int64Null()
	}
	if data.Comment.IsUnknown() {
		data.Comment = types.StringNull()
	}
	if data.InsecureTpc.IsUnknown() {
		data.InsecureTpc = types.BoolNull()
	}
	if data.Xen.IsUnknown() {
		data.Xen = types.BoolNull()
	}
	if data.Rpm.IsUnknown() {
		data.Rpm = types.StringNull()
	}
	if data.Ro.IsUnknown() {
		data.Ro = types.BoolNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}
	if data.ProductId.IsUnknown() {
		data.ProductId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("iscsi.extent.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_extent: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["disk"]; ok {
		if v == nil {
			data.Disk = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Disk = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Disk = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Disk = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["serial"]; ok {
		if v == nil {
			data.Serial = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Serial = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Serial = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Serial = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["path"]; ok {
		if v == nil {
			data.Path = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["avail_threshold"]; ok {
		if v == nil {
			data.AvailThreshold = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.AvailThreshold = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.AvailThreshold = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["product_id"]; ok {
		if v == nil {
			data.ProductId = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ProductId = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ProductId = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ProductId = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.Disk.IsUnknown() {
		data.Disk = types.StringNull()
	}
	if data.Serial.IsUnknown() {
		data.Serial = types.StringNull()
	}
	if data.Path.IsUnknown() {
		data.Path = types.StringNull()
	}
	if data.Filesize.IsUnknown() {
		data.Filesize = types.Int64Null()
	}
	if data.Blocksize.IsUnknown() {
		data.Blocksize = types.Int64Null()
	}
	if data.Pblocksize.IsUnknown() {
		data.Pblocksize = types.BoolNull()
	}
	if data.AvailThreshold.IsUnknown() {
		data.AvailThreshold = types.Int64Null()
	}
	if data.Comment.IsUnknown() {
		data.Comment = types.StringNull()
	}
	if data.InsecureTpc.IsUnknown() {
		data.InsecureTpc = types.BoolNull()
	}
	if data.Xen.IsUnknown() {
		data.Xen = types.BoolNull()
	}
	if data.Rpm.IsUnknown() {
		data.Rpm = types.StringNull()
	}
	if data.Ro.IsUnknown() {
		data.Ro = types.BoolNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}
	if data.ProductId.IsUnknown() {
		data.ProductId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Disk.IsNull() && !data.Disk.IsUnknown() {
		params["disk"] = data.Disk.ValueString()
	}
	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {
		params["serial"] = data.Serial.ValueString()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.Filesize.IsNull() && !data.Filesize.IsUnknown() {
		params["filesize"] = data.Filesize.ValueInt64()
	}
	if !data.Blocksize.IsNull() && !data.Blocksize.IsUnknown() {
		params["blocksize"] = data.Blocksize.ValueInt64()
	}
	if !data.Pblocksize.IsNull() && !data.Pblocksize.IsUnknown() {
		params["pblocksize"] = data.Pblocksize.ValueBool()
	}
	if !data.AvailThreshold.IsNull() && !data.AvailThreshold.IsUnknown() {
		params["avail_threshold"] = data.AvailThreshold.ValueInt64()
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}
	if !data.InsecureTpc.IsNull() && !data.InsecureTpc.IsUnknown() {
		params["insecure_tpc"] = data.InsecureTpc.ValueBool()
	}
	if !data.Xen.IsNull() && !data.Xen.IsUnknown() {
		params["xen"] = data.Xen.ValueBool()
	}
	if !data.Rpm.IsNull() && !data.Rpm.IsUnknown() {
		params["rpm"] = data.Rpm.ValueString()
	}
	if !data.Ro.IsNull() && !data.Ro.IsUnknown() {
		params["ro"] = data.Ro.ValueBool()
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ProductId.IsNull() && !data.ProductId.IsUnknown() {
		params["product_id"] = data.ProductId.ValueString()
	}

	_, err = r.client.Call("iscsi.extent.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_extent: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("iscsi.extent.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back iscsi_extent: %s", readErr))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["name"]; ok {
		switch val := v.(type) {
		case string:
			data.Name = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["disk"]; ok {
		if v == nil {
			data.Disk = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Disk = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Disk = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Disk = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["serial"]; ok {
		if v == nil {
			data.Serial = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Serial = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Serial = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Serial = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["path"]; ok {
		if v == nil {
			data.Path = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if v, ok := resultMap["avail_threshold"]; ok {
		if v == nil {
			data.AvailThreshold = types.Int64Null()
		} else {
			switch val := v.(type) {
			case float64:
				data.AvailThreshold = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok {
						data.AvailThreshold = types.Int64Value(int64(fv))
					}
				}
			}
		}
	}
	if v, ok := resultMap["product_id"]; ok {
		if v == nil {
			data.ProductId = types.StringNull()
		} else {
			switch val := v.(type) {
			case string:
				data.ProductId = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.ProductId = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.ProductId = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.Disk.IsUnknown() {
		data.Disk = types.StringNull()
	}
	if data.Serial.IsUnknown() {
		data.Serial = types.StringNull()
	}
	if data.Path.IsUnknown() {
		data.Path = types.StringNull()
	}
	if data.Filesize.IsUnknown() {
		data.Filesize = types.Int64Null()
	}
	if data.Blocksize.IsUnknown() {
		data.Blocksize = types.Int64Null()
	}
	if data.Pblocksize.IsUnknown() {
		data.Pblocksize = types.BoolNull()
	}
	if data.AvailThreshold.IsUnknown() {
		data.AvailThreshold = types.Int64Null()
	}
	if data.Comment.IsUnknown() {
		data.Comment = types.StringNull()
	}
	if data.InsecureTpc.IsUnknown() {
		data.InsecureTpc = types.BoolNull()
	}
	if data.Xen.IsUnknown() {
		data.Xen = types.BoolNull()
	}
	if data.Rpm.IsUnknown() {
		data.Rpm = types.StringNull()
	}
	if data.Ro.IsUnknown() {
		data.Ro = types.BoolNull()
	}
	if data.Enabled.IsUnknown() {
		data.Enabled = types.BoolNull()
	}
	if data.ProductId.IsUnknown() {
		data.ProductId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiExtentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiExtentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	id = []interface{}{id, map[string]interface{}{}}

	_, err = r.client.Call("iscsi.extent.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_extent: %s", err))
		return
	}
}
