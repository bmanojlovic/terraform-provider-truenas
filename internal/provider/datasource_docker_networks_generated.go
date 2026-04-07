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

var _ datasource.DataSource = &DockerNetworksDataSource{}

func NewDockerNetworksDataSource() datasource.DataSource {
	return &DockerNetworksDataSource{}
}

type DockerNetworksDataSource struct {
	client *client.Client
}

type DockerNetworksDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type DockerNetworksItemModel struct {
	ID types.String `tfsdk:"id"`
	Ipam types.String `tfsdk:"ipam"`
	Labels types.String `tfsdk:"labels"`
	Created types.String `tfsdk:"created"`
	Driver types.String `tfsdk:"driver"`
	Name types.String `tfsdk:"name"`
	Scope types.String `tfsdk:"scope"`
	ShortId types.String `tfsdk:"short_id"`
}

func (d *DockerNetworksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_docker_networks"
}

func (d *DockerNetworksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query all docker networks",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of docker_networks resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"ipam": schema.StringAttribute{
				Computed: true,
				Description: "IP Address Management configuration for the network or `null`.",
			},
			"labels": schema.StringAttribute{
				Computed: true,
				Description: "Metadata labels attached to the network or `null`.",
			},
			"created": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when the network was created or `null`.",
			},
			"driver": schema.StringAttribute{
				Computed: true,
				Description: "Network driver type (bridge, host, overlay, etc.) or `null`.",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "Human-readable name of the network or `null`.",
			},
			"scope": schema.StringAttribute{
				Computed: true,
				Description: "Network scope (local, global, swarm) or `null`.",
			},
			"short_id": schema.StringAttribute{
				Computed: true,
				Description: "Shortened network identifier or `null`.",
			},
					},
				},
			},
		},
	}
}

func (d *DockerNetworksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DockerNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DockerNetworksDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("docker.network.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query docker_networks: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]DockerNetworksItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := DockerNetworksItemModel{}
		if v, ok := resultMap["ipam"]; ok && v != nil {
			itemModel.Ipam = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["labels"]; ok && v != nil {
			itemModel.Labels = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["created"]; ok && v != nil {
			itemModel.Created = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["driver"]; ok && v != nil {
			itemModel.Driver = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["scope"]; ok && v != nil {
			itemModel.Scope = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["short_id"]; ok && v != nil {
			itemModel.ShortId = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"created": types.StringType,
			"driver": types.StringType,
			"id": types.StringType,
			"ipam": types.StringType,
			"labels": types.StringType,
			"name": types.StringType,
			"scope": types.StringType,
			"short_id": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
