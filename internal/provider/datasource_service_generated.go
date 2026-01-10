package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"strconv"
	
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ datasource.DataSource = &ServiceDataSource{}

func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

type ServiceDataSource struct {
	client *client.Client
}

type ServiceDataSourceModel struct {
	ID types.String `tfsdk:"id"`
	Service types.String `tfsdk:"service"`
	Enable types.Bool `tfsdk:"enable"`
	State types.String `tfsdk:"state"`
	Pids types.List `tfsdk:"pids"`
}

func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns instance matching `id`. If `id` is not found, Validation error is raised.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"service": schema.StringAttribute{
				Computed: true,
				Description: "Name of the system service.",
			},
			"enable": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the service is enabled to start on boot.",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Current state of the service (e.g., 'RUNNING', 'STOPPED').",
			},
			"pids": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
				Description: "Array of process IDs associated with this service.",
			},
		},
	}
}

func (d *ServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Call("service.get_instance", func() int { id, _ := strconv.Atoi(data.ID.ValueString()); return id }())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read service: %s", err.Error()))
		return
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

		if v, ok := resultMap["service"]; ok && v != nil {
			data.Service = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Enable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			data.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["pids"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Pids, _ = types.ListValue(types.StringType, strVals)
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
