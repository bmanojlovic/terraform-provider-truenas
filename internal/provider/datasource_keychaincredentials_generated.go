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

var _ datasource.DataSource = &KeychaincredentialsDataSource{}

func NewKeychaincredentialsDataSource() datasource.DataSource {
	return &KeychaincredentialsDataSource{}
}

type KeychaincredentialsDataSource struct {
	client *client.Client
}

type KeychaincredentialsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type KeychaincredentialsItemModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Attributes types.String `tfsdk:"attributes"`
}

func (d *KeychaincredentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keychaincredentials"
}

func (d *KeychaincredentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query keychaincredentials",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of keychaincredentials resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Distinguishes this Keychain Credential from others.",
			},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "Type of credential stored in the keychain.  * `SSH_KEY_PAIR`: SSH public/private key pair * `SSH_CRE",
			},
			"attributes": schema.StringAttribute{
				Computed: true,
				Description: "Credential-specific configuration and authentication data.",
			},
					},
				},
			},
		},
	}
}

func (d *KeychaincredentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KeychaincredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KeychaincredentialsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("keychaincredential.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query keychaincredentials: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]KeychaincredentialsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := KeychaincredentialsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["attributes"]; ok && v != nil {
			itemModel.Attributes = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attributes": types.StringType,
			"id": types.StringType,
			"name": types.StringType,
			"type": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
