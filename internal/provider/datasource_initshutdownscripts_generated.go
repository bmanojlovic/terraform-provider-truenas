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

var _ datasource.DataSource = &InitshutdownscriptsDataSource{}

func NewInitshutdownscriptsDataSource() datasource.DataSource {
	return &InitshutdownscriptsDataSource{}
}

type InitshutdownscriptsDataSource struct {
	client *client.Client
}

type InitshutdownscriptsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type InitshutdownscriptsItemModel struct {
	ID      types.String `tfsdk:"id"`
	Type    types.String `tfsdk:"type"`
	Command types.String `tfsdk:"command"`
	Script  types.String `tfsdk:"script"`
	When    types.String `tfsdk:"when"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Timeout types.Int64  `tfsdk:"timeout"`
	Comment types.String `tfsdk:"comment"`
}

func (d *InitshutdownscriptsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_initshutdownscripts"
}

func (d *InitshutdownscriptsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query initshutdownscripts",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of initshutdownscripts resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of init/shutdown script to execute.  * `COMMAND`: Execute a single command * `SCRIPT`: Execute ",
						},
						"command": schema.StringAttribute{
							Computed:    true,
							Description: "Must be given if `type=\"COMMAND\"`.",
						},
						"script": schema.StringAttribute{
							Computed:    true,
							Description: "Must be given if `type=\"SCRIPT\"`.",
						},
						"when": schema.StringAttribute{
							Computed:    true,
							Description: "* \"PREINIT\": Early in the boot process before all services have started. * \"POSTINIT\": Late in the b",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the init/shutdown script is enabled to execute.",
						},
						"timeout": schema.Int64Attribute{
							Computed:    true,
							Description: "An integer time in seconds that the system should wait for the execution of the script/command.  A h",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Optional comment describing the purpose of this script.",
						},
					},
				},
			},
		},
	}
}

func (d *InitshutdownscriptsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InitshutdownscriptsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InitshutdownscriptsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("initshutdownscript.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query initshutdownscripts: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]InitshutdownscriptsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := InitshutdownscriptsItemModel{}
		if v, ok := resultMap["type"]; ok && v != nil {
			itemModel.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["command"]; ok && v != nil {
			itemModel.Command = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["script"]; ok && v != nil {
			itemModel.Script = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["when"]; ok && v != nil {
			itemModel.When = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["timeout"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.Timeout = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"command": types.StringType,
			"comment": types.StringType,
			"enabled": types.BoolType,
			"id":      types.StringType,
			"script":  types.StringType,
			"timeout": types.Int64Type,
			"type":    types.StringType,
			"when":    types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
