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

var _ datasource.DataSource = &TunablesDataSource{}

func NewTunablesDataSource() datasource.DataSource {
	return &TunablesDataSource{}
}

type TunablesDataSource struct {
	client *client.Client
}

type TunablesDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type TunablesItemModel struct {
	ID types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Var types.String `tfsdk:"var"`
	Value types.String `tfsdk:"value"`
	Comment types.String `tfsdk:"comment"`
	Enabled types.Bool `tfsdk:"enabled"`
	UpdateInitramfs types.Bool `tfsdk:"update_initramfs"`
	OrigValue types.String `tfsdk:"orig_value"`
}

func (d *TunablesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tunables"
}

func (d *TunablesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query tunables",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of tunables resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"type": schema.StringAttribute{
				Computed: true,
				Description: "* `SYSCTL`: `var` is a sysctl name (e.g. `kernel.watchdog`) and `value` is its corresponding value (",
			},
			"var": schema.StringAttribute{
				Computed: true,
				Description: "Name or identifier of the system parameter to tune.",
			},
			"value": schema.StringAttribute{
				Computed: true,
				Description: "Value to assign to the tunable parameter.",
			},
			"comment": schema.StringAttribute{
				Computed: true,
				Description: "Optional descriptive comment explaining the purpose of this tunable.",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this tunable is active and should be applied.",
			},
			"update_initramfs": schema.BoolAttribute{
				Computed: true,
				Description: "If `false`, then initramfs will not be updated after creating a ZFS tunable and you will need to run",
			},
			"orig_value": schema.StringAttribute{
				Computed: true,
				Description: "Original system value of the parameter before this tunable was applied.",
			},
					},
				},
			},
		},
	}
}

func (d *TunablesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TunablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TunablesDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("tunable.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query tunables: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]TunablesItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := TunablesItemModel{}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["var"]; ok && v != nil {
			itemModel.Var = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["value"]; ok && v != nil {
			itemModel.Value = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["update_initramfs"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.UpdateInitramfs = types.BoolValue(bv) }
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["orig_value"]; ok && v != nil {
			itemModel.OrigValue = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"comment": types.StringType,
			"enabled": types.BoolType,
			"id": types.StringType,
			"orig_value": types.StringType,
			"type": types.StringType,
			"update_initramfs": types.BoolType,
			"value": types.StringType,
			"var": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
