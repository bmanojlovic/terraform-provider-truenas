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

var _ datasource.DataSource = &SharingSmbsDataSource{}

func NewSharingSmbsDataSource() datasource.DataSource {
	return &SharingSmbsDataSource{}
}

type SharingSmbsDataSource struct {
	client *client.Client
}

type SharingSmbsDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type SharingSmbsItemModel struct {
	ID types.String `tfsdk:"id"`
	Purpose types.String `tfsdk:"purpose"`
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
	Enabled types.Bool `tfsdk:"enabled"`
	Comment types.String `tfsdk:"comment"`
	Readonly types.Bool `tfsdk:"readonly"`
	Browsable types.Bool `tfsdk:"browsable"`
	AccessBasedShareEnumeration types.Bool `tfsdk:"access_based_share_enumeration"`
	Locked types.Bool `tfsdk:"locked"`
	Audit types.String `tfsdk:"audit"`
	Options types.String `tfsdk:"options"`
}

func (d *SharingSmbsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sharing_smbs"
}

func (d *SharingSmbsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query sharing_smbs",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of sharing_smbs resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
			"purpose": schema.StringAttribute{
				Computed: true,
				Description: "This parameter sets the purpose of the SMB share. It controls how the SMB share behaves and what fea",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Description: "SMB share name. SMB share names are case-insensitive and must be unique, and are subject     to the ",
			},
			"path": schema.StringAttribute{
				Computed: true,
				Description: "Local server path to share by using the SMB protocol. The path must start with `/mnt/` and must be i",
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
				Description: "If unset, the SMB share is not available over the SMB protocol. ",
			},
			"comment": schema.StringAttribute{
				Computed: true,
				Description: "Text field that is seen next to a share when an SMB client requests a list of SMB shares on the True",
			},
			"readonly": schema.BoolAttribute{
				Computed: true,
				Description: "If set, SMB clients cannot create or change files and directories in the SMB share.  NOTE: If set, t",
			},
			"browsable": schema.BoolAttribute{
				Computed: true,
				Description: "If set, the share is included when an SMB client requests a list of SMB shares on the TrueNAS server",
			},
			"access_based_share_enumeration": schema.BoolAttribute{
				Computed: true,
				Description: "If set, the share is only included when an SMB client requests a list of shares on the SMB server if",
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				Description: "Read-only value indicating whether the share is located on a locked dataset.  Returns:     - True: T",
			},
			"audit": schema.StringAttribute{
				Computed: true,
				Description: "Audit configuration for monitoring SMB share access and operations.",
			},
			"options": schema.StringAttribute{
				Computed: true,
				Description: "Additional configuration related to the configured SMB share purpose. If null, then the default     ",
			},
					},
				},
			},
		},
	}
}

func (d *SharingSmbsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SharingSmbsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SharingSmbsDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("sharing.smb.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query sharing_smbs: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]SharingSmbsItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := SharingSmbsItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["purpose"]; ok && v != nil {
			itemModel.Purpose = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["path"]; ok && v != nil {
			itemModel.Path = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["enabled"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Enabled = types.BoolValue(bv) }
		}
		if v, ok := resultMap["comment"]; ok && v != nil {
			itemModel.Comment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Readonly = types.BoolValue(bv) }
		}
		if v, ok := resultMap["browsable"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Browsable = types.BoolValue(bv) }
		}
		if v, ok := resultMap["access_based_share_enumeration"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.AccessBasedShareEnumeration = types.BoolValue(bv) }
		}
		if v, ok := resultMap["locked"]; ok && v != nil {
			if bv, ok := v.(bool); ok { itemModel.Locked = types.BoolValue(bv) }
		}
		if v, ok := resultMap["audit"]; ok && v != nil {
			itemModel.Audit = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["options"]; ok && v != nil {
			itemModel.Options = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"access_based_share_enumeration": types.BoolType,
			"audit": types.StringType,
			"browsable": types.BoolType,
			"comment": types.StringType,
			"enabled": types.BoolType,
			"id": types.StringType,
			"locked": types.BoolType,
			"name": types.StringType,
			"options": types.StringType,
			"path": types.StringType,
			"purpose": types.StringType,
			"readonly": types.BoolType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
