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

var _ datasource.DataSource = &AppImagesDataSource{}

func NewAppImagesDataSource() datasource.DataSource {
	return &AppImagesDataSource{}
}

type AppImagesDataSource struct {
	client *client.Client
}

type AppImagesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type AppImagesItemModel struct {
	ID              types.String `tfsdk:"id"`
	Size            types.Int64  `tfsdk:"size"`
	Dangling        types.Bool   `tfsdk:"dangling"`
	UpdateAvailable types.Bool   `tfsdk:"update_available"`
	Created         types.String `tfsdk:"created"`
	Author          types.String `tfsdk:"author"`
	Comment         types.String `tfsdk:"comment"`
}

func (d *AppImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_images"
}

func (d *AppImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query all docker images with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of app_images resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the container image in bytes.",
						},
						"dangling": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this is a dangling image (no tags or references).",
						},
						"update_available": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether a newer version of this image is available for download.",
						},
						"created": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the container image was created (ISO format).",
						},
						"author": schema.StringAttribute{
							Computed:    true,
							Description: "Author or maintainer of the container image.",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Comment or description provided by the image author.",
						},
					},
				},
			},
		},
	}
}

func (d *AppImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppImagesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("app.image.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query app_images: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]AppImagesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := AppImagesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["size"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Size = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["dangling"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Dangling = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["update_available"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.UpdateAvailable = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["created"]; ok && v != nil {
			itemModel.Created = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["author"]; ok && v != nil {
			itemModel.Author = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"author":           types.StringType,
			"comment":          types.StringType,
			"created":          types.StringType,
			"dangling":         types.BoolType,
			"id":               types.StringType,
			"size":             types.Int64Type,
			"update_available": types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
