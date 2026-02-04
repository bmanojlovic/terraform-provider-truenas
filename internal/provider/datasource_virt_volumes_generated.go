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

var _ datasource.DataSource = &VirtVolumesDataSource{}

func NewVirtVolumesDataSource() datasource.DataSource {
	return &VirtVolumesDataSource{}
}

type VirtVolumesDataSource struct {
	client *client.Client
}

type VirtVolumesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type VirtVolumesItemModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	StoragePool types.String `tfsdk:"storage_pool"`
	ContentType types.String `tfsdk:"content_type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Type        types.String `tfsdk:"type"`
	Config      types.String `tfsdk:"config"`
}

func (d *VirtVolumesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virt_volumes"
}

func (d *VirtVolumesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query virt_volumes",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of virt_volumes resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable name of the virtualization volume.",
						},
						"storage_pool": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the storage pool containing this volume.",
						},
						"content_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of content stored in this volume (e.g., 'BLOCK', 'ISO').",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when this volume was created.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Volume type indicating its storage backend and characteristics.",
						},
						"config": schema.StringAttribute{
							Computed:    true,
							Description: "Object containing detailed configuration parameters for this volume.",
						},
					},
				},
			},
		},
	}
}

func (d *VirtVolumesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VirtVolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VirtVolumesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("virt.volume.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query virt_volumes: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]VirtVolumesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := VirtVolumesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["storage_pool"]; ok && v != nil {
			itemModel.StoragePool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["content_type"]; ok && v != nil {
			itemModel.ContentType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["created_at"]; ok && v != nil {
			itemModel.CreatedAt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["config"]; ok && v != nil {
			itemModel.Config = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"config":       types.StringType,
			"content_type": types.StringType,
			"created_at":   types.StringType,
			"id":           types.StringType,
			"name":         types.StringType,
			"storage_pool": types.StringType,
			"type":         types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
