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

var _ datasource.DataSource = &NvmetSubsyssDataSource{}

func NewNvmetSubsyssDataSource() datasource.DataSource {
	return &NvmetSubsyssDataSource{}
}

type NvmetSubsyssDataSource struct {
	client *client.Client
}

type NvmetSubsyssDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type NvmetSubsyssItemModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Subnqn       types.String `tfsdk:"subnqn"`
	Serial       types.String `tfsdk:"serial"`
	AllowAnyHost types.Bool   `tfsdk:"allow_any_host"`
	PiEnable     types.Bool   `tfsdk:"pi_enable"`
	QidMax       types.Int64  `tfsdk:"qid_max"`
	IeeeOui      types.String `tfsdk:"ieee_oui"`
	Ana          types.Bool   `tfsdk:"ana"`
}

func (d *NvmetSubsyssDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_subsyss"
}

func (d *NvmetSubsyssDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query nvmet_subsyss",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of nvmet_subsyss resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human readable name for the subsystem.  If `subnqn` is not provided on creation, then this name will",
						},
						"subnqn": schema.StringAttribute{
							Computed:    true,
							Description: "NVMe Qualified Name (NQN) for the subsystem.  If not provided during creation, will be auto-generate",
						},
						"serial": schema.StringAttribute{
							Computed:    true,
							Description: "Serial number assigned to the subsystem.",
						},
						"allow_any_host": schema.BoolAttribute{
							Computed:    true,
							Description: "Any host can access the storage associated with this subsystem (i.e. no access control).",
						},
						"pi_enable": schema.BoolAttribute{
							Computed:    true,
							Description: "Enable Protection Information (PI) for data integrity checking.",
						},
						"qid_max": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of queue IDs allowed for this subsystem.",
						},
						"ieee_oui": schema.StringAttribute{
							Computed:    true,
							Description: "IEEE Organizationally Unique Identifier for the subsystem.",
						},
						"ana": schema.BoolAttribute{
							Computed:    true,
							Description: "If set to either `True` or `False`, then *override* the global `ana` setting from `nvmet.global.conf",
						},
					},
				},
			},
		},
	}
}

func (d *NvmetSubsyssDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NvmetSubsyssDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NvmetSubsyssDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("nvmet.subsys.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query nvmet_subsyss: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]NvmetSubsyssItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := NvmetSubsyssItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["subnqn"]; ok && v != nil {
			itemModel.Subnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["serial"]; ok && v != nil {
			itemModel.Serial = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["allow_any_host"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.AllowAnyHost = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["pi_enable"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.PiEnable = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["qid_max"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.QidMax = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["ieee_oui"]; ok && v != nil {
			itemModel.IeeeOui = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ana"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Ana = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"allow_any_host": types.BoolType,
			"ana":            types.BoolType,
			"id":             types.StringType,
			"ieee_oui":       types.StringType,
			"name":           types.StringType,
			"pi_enable":      types.BoolType,
			"qid_max":        types.Int64Type,
			"serial":         types.StringType,
			"subnqn":         types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
