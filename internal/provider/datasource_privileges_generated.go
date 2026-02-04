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

var _ datasource.DataSource = &PrivilegesDataSource{}

func NewPrivilegesDataSource() datasource.DataSource {
	return &PrivilegesDataSource{}
}

type PrivilegesDataSource struct {
	client *client.Client
}

type PrivilegesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type PrivilegesItemModel struct {
	ID          types.String `tfsdk:"id"`
	BuiltinName types.String `tfsdk:"builtin_name"`
	Name        types.String `tfsdk:"name"`
	WebShell    types.Bool   `tfsdk:"web_shell"`
}

func (d *PrivilegesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privileges"
}

func (d *PrivilegesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query privileges",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of privileges resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"builtin_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the built-in privilege if this is a system privilege. `null` for custom privileges.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Display name of the privilege.",
						},
						"web_shell": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this privilege grants access to the web shell.",
						},
					},
				},
			},
		},
	}
}

func (d *PrivilegesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PrivilegesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PrivilegesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("privilege.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query privileges: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]PrivilegesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := PrivilegesItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["builtin_name"]; ok && v != nil {
			itemModel.BuiltinName = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["web_shell"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.WebShell = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"builtin_name": types.StringType,
			"id":           types.StringType,
			"name":         types.StringType,
			"web_shell":    types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
