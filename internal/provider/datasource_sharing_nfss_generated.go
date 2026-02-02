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

var _ datasource.DataSource = &SharingNfssDataSource{}

func NewSharingNfssDataSource() datasource.DataSource {
	return &SharingNfssDataSource{}
}

type SharingNfssDataSource struct {
	client *client.Client
}

type SharingNfssDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type SharingNfssItemModel struct {
	ID              types.String `tfsdk:"id"`
	Path            types.String `tfsdk:"path"`
	Comment         types.String `tfsdk:"comment"`
	Ro              types.Bool   `tfsdk:"ro"`
	MaprootUser     types.String `tfsdk:"maproot_user"`
	MaprootGroup    types.String `tfsdk:"maproot_group"`
	MapallUser      types.String `tfsdk:"mapall_user"`
	MapallGroup     types.String `tfsdk:"mapall_group"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Locked          types.Bool   `tfsdk:"locked"`
	ExposeSnapshots types.Bool   `tfsdk:"expose_snapshots"`
}

func (d *SharingNfssDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sharing_nfss"
}

func (d *SharingNfssDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query sharing_nfss",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of sharing_nfss resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"path": schema.StringAttribute{
							Computed:    true,
							Description: "Local path to be exported. ",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "User comment associated with share. ",
						},
						"ro": schema.BoolAttribute{
							Computed:    true,
							Description: "Export the share as read only. ",
						},
						"maproot_user": schema.StringAttribute{
							Computed:    true,
							Description: "Map root user client to a specified user. ",
						},
						"maproot_group": schema.StringAttribute{
							Computed:    true,
							Description: "Map root group client to a specified group. ",
						},
						"mapall_user": schema.StringAttribute{
							Computed:    true,
							Description: "Map all client users to a specified user. ",
						},
						"mapall_group": schema.StringAttribute{
							Computed:    true,
							Description: "Map all client groups to a specified group. ",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Enable or disable the share. ",
						},
						"locked": schema.BoolAttribute{
							Computed:    true,
							Description: "Read-only value indicating whether the share is located on a locked dataset.  Returns:     - True: T",
						},
						"expose_snapshots": schema.BoolAttribute{
							Computed:    true,
							Description: "Enterprise feature to enable access to the ZFS snapshot directory for the export. Export path must b",
						},
					},
				},
			},
		},
	}
}

func (d *SharingNfssDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SharingNfssDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SharingNfssDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("sharing.nfs.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query sharing_nfss: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]SharingNfssItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := SharingNfssItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			itemModel.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ro"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Ro = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["maproot_user"]; ok && v != nil {
			itemModel.MaprootUser = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["maproot_group"]; ok && v != nil {
			itemModel.MaprootGroup = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mapall_user"]; ok && v != nil {
			itemModel.MapallUser = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mapall_group"]; ok && v != nil {
			itemModel.MapallGroup = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Enabled = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Locked = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["expose_snapshots"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.ExposeSnapshots = types.BoolValue(bv)
			}
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"comment":          types.StringType,
			"enabled":          types.BoolType,
			"expose_snapshots": types.BoolType,
			"id":               types.StringType,
			"locked":           types.BoolType,
			"mapall_group":     types.StringType,
			"mapall_user":      types.StringType,
			"maproot_group":    types.StringType,
			"maproot_user":     types.StringType,
			"path":             types.StringType,
			"ro":               types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
