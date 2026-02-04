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

var _ datasource.DataSource = &BootEnvironmentsDataSource{}

func NewBootEnvironmentsDataSource() datasource.DataSource {
	return &BootEnvironmentsDataSource{}
}

type BootEnvironmentsDataSource struct {
	client *client.Client
}

type BootEnvironmentsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type BootEnvironmentsItemModel struct {
	ID          types.String `tfsdk:"id"`
	Dataset     types.String `tfsdk:"dataset"`
	Active      types.Bool   `tfsdk:"active"`
	Activated   types.Bool   `tfsdk:"activated"`
	Created     types.String `tfsdk:"created"`
	UsedBytes   types.Int64  `tfsdk:"used_bytes"`
	Used        types.String `tfsdk:"used"`
	Keep        types.Bool   `tfsdk:"keep"`
	CanActivate types.Bool   `tfsdk:"can_activate"`
}

func (d *BootEnvironmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_boot_environments"
}

func (d *BootEnvironmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query boot_environments",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of boot_environments resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"dataset": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the zfs dataset that represents the boot environment.",
						},
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "This is the currently running boot environment.",
						},
						"activated": schema.BoolAttribute{
							Computed:    true,
							Description: "Use this boot environment on next boot.",
						},
						"created": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the boot environment was created.",
						},
						"used_bytes": schema.Int64Attribute{
							Computed:    true,
							Description: "The total amount of bytes used by the boot environment.",
						},
						"used": schema.StringAttribute{
							Computed:    true,
							Description: "The boot environment's used space in human readable format.",
						},
						"keep": schema.BoolAttribute{
							Computed:    true,
							Description: "When set to false, this makes the boot environment subject to     automatic deletion if the TrueNAS ",
						},
						"can_activate": schema.BoolAttribute{
							Computed:    true,
							Description: "The given boot environment may be activated.",
						},
					},
				},
			},
		},
	}
}

func (d *BootEnvironmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BootEnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BootEnvironmentsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("boot.environment.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query boot_environments: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]BootEnvironmentsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := BootEnvironmentsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dataset"]; ok && v != nil {
			itemModel.Dataset = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["active"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Active = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["activated"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Activated = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["created"]; ok && v != nil {
			itemModel.Created = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["used_bytes"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.UsedBytes = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["used"]; ok && v != nil {
			itemModel.Used = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["keep"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Keep = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["can_activate"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.CanActivate = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"activated":    types.BoolType,
			"active":       types.BoolType,
			"can_activate": types.BoolType,
			"created":      types.StringType,
			"dataset":      types.StringType,
			"id":           types.StringType,
			"keep":         types.BoolType,
			"used":         types.StringType,
			"used_bytes":   types.Int64Type,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
