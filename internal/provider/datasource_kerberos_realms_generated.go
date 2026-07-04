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

var _ datasource.DataSource = &KerberosRealmsDataSource{}

func NewKerberosRealmsDataSource() datasource.DataSource {
	return &KerberosRealmsDataSource{}
}

type KerberosRealmsDataSource struct {
	client *client.Client
}

type KerberosRealmsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type KerberosRealmsItemModel struct {
	ID         types.String `tfsdk:"id"`
	Realm      types.String `tfsdk:"realm"`
	PrimaryKdc types.String `tfsdk:"primary_kdc"`
}

func (d *KerberosRealmsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_realms"
}

func (d *KerberosRealmsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query kerberos_realms",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of kerberos_realms resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"realm": schema.StringAttribute{
							Computed:    true,
							Description: "Kerberos realm name. This is external to TrueNAS and is case-sensitive.     The general convention f",
						},
						"primary_kdc": schema.StringAttribute{
							Computed:    true,
							Description: "The master Kerberos domain controller for this realm. TrueNAS uses this as a fallback if it cannot g",
						},
					},
				},
			},
		},
	}
}

func (d *KerberosRealmsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *KerberosRealmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KerberosRealmsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("kerberos.realm.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query kerberos_realms: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]KerberosRealmsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := KerberosRealmsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["realm"]; ok && v != nil {
			itemModel.Realm = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["primary_kdc"]; ok && v != nil {
			itemModel.PrimaryKdc = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"primary_kdc": types.StringType,
			"realm":       types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
