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

var _ datasource.DataSource = &AppRegistrysDataSource{}

func NewAppRegistrysDataSource() datasource.DataSource {
	return &AppRegistrysDataSource{}
}

type AppRegistrysDataSource struct {
	client *client.Client
}

type AppRegistrysDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type AppRegistrysItemModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Uri types.String `tfsdk:"uri"`
}

func (d *AppRegistrysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_registrys"
}

func (d *AppRegistrysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query app_registrys",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of app_registrys resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable name for the container registry.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Optional description of the container registry or `null`.",
			},
			"username": schema.StringAttribute{
				Computed: true,
				Description: "Username for registry authentication (masked for security).",
			},
			"password": schema.StringAttribute{
				Computed: true,
				Description: "Password or access token for registry authentication (masked for security).",
			},
			"uri": schema.StringAttribute{
				Computed: true,
				Description: "Container registry URI endpoint.",
			},
					},
				},
			},
		},
	}
}

func (d *AppRegistrysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppRegistrysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppRegistrysDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("app.registry.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query app_registrys: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]AppRegistrysItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := AppRegistrysItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			itemModel.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			itemModel.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["uri"]; ok && v != nil {
			itemModel.Uri = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"id": types.StringType,
			"name": types.StringType,
			"password": types.StringType,
			"uri": types.StringType,
			"username": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
