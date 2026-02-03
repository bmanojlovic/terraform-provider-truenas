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

var _ datasource.DataSource = &CronjobsDataSource{}

func NewCronjobsDataSource() datasource.DataSource {
	return &CronjobsDataSource{}
}

type CronjobsDataSource struct {
	client *client.Client
}

type CronjobsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type CronjobsItemModel struct {
	ID types.String `tfsdk:"id"`
	Enabled types.Bool `tfsdk:"enabled"`
	Stderr types.Bool `tfsdk:"stderr"`
	Stdout types.Bool `tfsdk:"stdout"`
	Schedule types.String `tfsdk:"schedule"`
	Command types.String `tfsdk:"command"`
	Description types.String `tfsdk:"description"`
	User types.String `tfsdk:"user"`
}

func (d *CronjobsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cronjobs"
}

func (d *CronjobsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query cronjobs",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of cronjobs resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the cron job is active and will be executed.",
			},
			"stderr": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to IGNORE standard error (if `false`, it will be added to email).",
			},
			"stdout": schema.BoolAttribute{
				Computed: true,
				Description: "Whether to IGNORE standard output (if `false`, it will be added to email).",
			},
			"schedule": schema.StringAttribute{
				Computed: true,
				Description: "Cron schedule configuration for when the job runs.",
			},
			"command": schema.StringAttribute{
				Computed: true,
				Description: "Shell command or script to execute.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable description of what this cron job does.",
			},
			"user": schema.StringAttribute{
				Computed: true,
				Description: "System user account to run the command as.",
			},
					},
				},
			},
		},
	}
}

func (d *CronjobsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CronjobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CronjobsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("cronjob.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query cronjobs: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]CronjobsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := CronjobsItemModel{}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["stderr"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Stderr = types.BoolValue(bv) }
		}
		if v, ok := resultMap["stdout"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Stdout = types.BoolValue(bv) }
		}
		if v, ok := resultMap["schedule"]; ok && v != nil {
			itemModel.Schedule = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["command"]; ok && v != nil {
			itemModel.Command = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			itemModel.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user"]; ok && v != nil {
			itemModel.User = types.StringValue(fmt.Sprintf("%v", v))
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
			"description": types.StringType,
			"enabled": types.BoolType,
			"id": types.StringType,
			"schedule": types.StringType,
			"stderr": types.BoolType,
			"stdout": types.BoolType,
			"user": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
