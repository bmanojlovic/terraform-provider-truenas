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

var _ datasource.DataSource = &NvmetHostsDataSource{}

func NewNvmetHostsDataSource() datasource.DataSource {
	return &NvmetHostsDataSource{}
}

type NvmetHostsDataSource struct {
	client *client.Client
}

type NvmetHostsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type NvmetHostsItemModel struct {
	ID            types.String `tfsdk:"id"`
	Hostnqn       types.String `tfsdk:"hostnqn"`
	DhchapKey     types.String `tfsdk:"dhchap_key"`
	DhchapCtrlKey types.String `tfsdk:"dhchap_ctrl_key"`
	DhchapDhgroup types.String `tfsdk:"dhchap_dhgroup"`
	DhchapHash    types.String `tfsdk:"dhchap_hash"`
}

func (d *NvmetHostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_hosts"
}

func (d *NvmetHostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query nvmet_hosts",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of nvmet_hosts resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"hostnqn": schema.StringAttribute{
							Computed:    true,
							Description: "NQN of the host that will connect to this TrueNAS. ",
						},
						"dhchap_key": schema.StringAttribute{
							Computed:    true,
							Description: "If set, the secret that the host must present when connecting.  A suitable secret can be generated u",
						},
						"dhchap_ctrl_key": schema.StringAttribute{
							Computed:    true,
							Description: "If set, the secret that this TrueNAS will present to the host when the host is connecting (Bi-Direct",
						},
						"dhchap_dhgroup": schema.StringAttribute{
							Computed:    true,
							Description: "If selected, the DH (Diffie-Hellman) key exchange built on top of CHAP to be used for authentication",
						},
						"dhchap_hash": schema.StringAttribute{
							Computed:    true,
							Description: "HMAC (Hashed Message Authentication Code) to be used in conjunction if a `dhchap_dhgroup` is selecte",
						},
					},
				},
			},
		},
	}
}

func (d *NvmetHostsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NvmetHostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NvmetHostsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("nvmet.host.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query nvmet_hosts: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]NvmetHostsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := NvmetHostsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hostnqn"]; ok && v != nil {
			itemModel.Hostnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_key"]; ok && v != nil {
			itemModel.DhchapKey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_ctrl_key"]; ok && v != nil {
			itemModel.DhchapCtrlKey = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_dhgroup"]; ok && v != nil {
			itemModel.DhchapDhgroup = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dhchap_hash"]; ok && v != nil {
			itemModel.DhchapHash = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"dhchap_ctrl_key": types.StringType,
			"dhchap_dhgroup":  types.StringType,
			"dhchap_hash":     types.StringType,
			"dhchap_key":      types.StringType,
			"hostnqn":         types.StringType,
			"id":              types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
