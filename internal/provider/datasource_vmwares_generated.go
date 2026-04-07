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

var _ datasource.DataSource = &VmwaresDataSource{}

func NewVmwaresDataSource() datasource.DataSource {
	return &VmwaresDataSource{}
}

type VmwaresDataSource struct {
	client *client.Client
}

type VmwaresDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type VmwaresItemModel struct {
	ID types.String `tfsdk:"id"`
	Datastore types.String `tfsdk:"datastore"`
	Filesystem types.String `tfsdk:"filesystem"`
	Hostname types.String `tfsdk:"hostname"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	State types.String `tfsdk:"state"`
}

func (d *VmwaresDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vmwares"
}

func (d *VmwaresDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query vmwares",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of vmwares resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"datastore": schema.StringAttribute{
				Computed: true,
				Description: "Valid datastore name which exists on the VMWare host.",
			},
			"filesystem": schema.StringAttribute{
				Computed: true,
				Description: "ZFS filesystem or dataset to use for VMware storage.",
			},
			"hostname": schema.StringAttribute{
				Computed: true,
				Description: "Valid IP address / hostname of a VMWare host. When clustering, this is the vCenter server for the cl",
			},
			"username": schema.StringAttribute{
				Computed: true,
				Description: "Credentials used to authorize access to the VMWare host.",
			},
			"password": schema.StringAttribute{
				Computed: true,
				Description: "Password for VMware host authentication.",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "Current connection and synchronization state with the VMware host.",
			},
					},
				},
			},
		},
	}
}

func (d *VmwaresDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VmwaresDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VmwaresDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("vmware.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query vmwares: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]VmwaresItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := VmwaresItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["datastore"]; ok && v != nil {
			itemModel.Datastore = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["filesystem"]; ok && v != nil {
			itemModel.Filesystem = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hostname"]; ok && v != nil {
			itemModel.Hostname = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			itemModel.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			itemModel.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["state"]; ok && v != nil {
			itemModel.State = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"datastore": types.StringType,
			"filesystem": types.StringType,
			"hostname": types.StringType,
			"id": types.StringType,
			"password": types.StringType,
			"state": types.StringType,
			"username": types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
