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

var _ datasource.DataSource = &AppsDataSource{}

func NewAppsDataSource() datasource.DataSource {
	return &AppsDataSource{}
}

type AppsDataSource struct {
	client *client.Client
}

type AppsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type AppsItemModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	State types.String `tfsdk:"state"`
	UpgradeAvailable types.Bool `tfsdk:"upgrade_available"`
	LatestVersion types.String `tfsdk:"latest_version"`
	ImageUpdatesAvailable types.Bool `tfsdk:"image_updates_available"`
	CustomApp types.Bool `tfsdk:"custom_app"`
	Migrated types.Bool `tfsdk:"migrated"`
	HumanVersion types.String `tfsdk:"human_version"`
	Version types.String `tfsdk:"version"`
	Metadata types.String `tfsdk:"metadata"`
	ActiveWorkloads types.String `tfsdk:"active_workloads"`
	Notes types.String `tfsdk:"notes"`
	Portals types.String `tfsdk:"portals"`
	VersionDetails types.String `tfsdk:"version_details"`
	Config types.String `tfsdk:"config"`
}

func (d *AppsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

func (d *AppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query all apps with `query-filters` and `query-options`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of apps resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "The display name of the application.",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Current operational state of the application.",
			},
			"upgrade_available": schema.BoolAttribute{
				Computed: true,
				Description: "Whether a newer version of the application is available for upgrade.",
			},
			"latest_version": schema.StringAttribute{
				Computed: true,
				Description: "The latest available version string, or `null` if no updates are available.",
			},
			"image_updates_available": schema.BoolAttribute{
				Computed: true,
				Description: "Whether newer Docker images are available for the containers in this application.",
			},
			"custom_app": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this is a custom application (`true`) or from a catalog (`false`).",
			},
			"migrated": schema.BoolAttribute{
				Computed: true,
				Description: "Whether this application has been migrated from kubernetes.",
			},
			"human_version": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable version string for display purposes.",
			},
			"version": schema.StringAttribute{
				Computed: true,
				Description: "Technical version identifier of the currently installed application.",
			},
			"metadata": schema.StringAttribute{
				Computed: true,
				Description: "Application metadata including description, category, and other catalog information.",
			},
			"active_workloads": schema.StringAttribute{
				Computed: true,
				Description: "Information about the running containers, ports, and resources used by this application.",
			},
			"notes": schema.StringAttribute{
				Computed: true,
				Description: "User-provided notes or comments about this application instance.",
			},
			"portals": schema.StringAttribute{
				Computed: true,
				Description: "Web portals and access points provided by the application (URLs, ports, etc.).",
			},
			"version_details": schema.StringAttribute{
				Computed: true,
				Description: "Detailed version information including changelog and upgrade notes. `null` if not available.",
			},
			"config": schema.StringAttribute{
				Computed: true,
				Description: "Current configuration values for the application. `null` if configuration is not requested.",
			},
					},
				},
			},
		},
	}
}

func (d *AppsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("app.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query apps: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]AppsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := AppsItemModel{}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["upgrade_available"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.UpgradeAvailable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["latest_version"]; ok && v != nil {
			itemModel.LatestVersion = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["image_updates_available"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.ImageUpdatesAvailable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["custom_app"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.CustomApp = types.BoolValue(bv) }
		}
		if v, ok := resultMap["migrated"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Migrated = types.BoolValue(bv) }
		}
		if v, ok := resultMap["human_version"]; ok && v != nil {
			itemModel.HumanVersion = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["version"]; ok && v != nil {
			itemModel.Version = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["metadata"]; ok && v != nil {
			itemModel.Metadata = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["active_workloads"]; ok && v != nil {
			itemModel.ActiveWorkloads = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["notes"]; ok && v != nil {
			itemModel.Notes = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["portals"]; ok && v != nil {
			itemModel.Portals = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["version_details"]; ok && v != nil {
			itemModel.VersionDetails = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["config"]; ok && v != nil {
			itemModel.Config = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"active_workloads": types.StringType,
			"config": types.StringType,
			"custom_app": types.BoolType,
			"human_version": types.StringType,
			"id": types.StringType,
			"image_updates_available": types.BoolType,
			"latest_version": types.StringType,
			"metadata": types.StringType,
			"migrated": types.BoolType,
			"name": types.StringType,
			"notes": types.StringType,
			"portals": types.StringType,
			"state": types.StringType,
			"upgrade_available": types.BoolType,
			"version": types.StringType,
			"version_details": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
