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

var _ datasource.DataSource = &ApiKeysDataSource{}

func NewApiKeysDataSource() datasource.DataSource {
	return &ApiKeysDataSource{}
}

type ApiKeysDataSource struct {
	client *client.Client
}

type ApiKeysDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

type ApiKeysItemModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Username       types.String `tfsdk:"username"`
	UserIdentifier types.Int64  `tfsdk:"user_identifier"`
	Keyhash        types.String `tfsdk:"keyhash"`
	CreatedAt      types.String `tfsdk:"created_at"`
	ExpiresAt      types.String `tfsdk:"expires_at"`
	Local          types.Bool   `tfsdk:"local"`
	Revoked        types.Bool   `tfsdk:"revoked"`
	RevokedReason  types.String `tfsdk:"revoked_reason"`
}

func (d *ApiKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_keys"
}

func (d *ApiKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query api_keys",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of api_keys resources",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Required: true, Description: "Resource ID"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Human-readable name for the API key.",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "Username associated with the API key or `null` for system keys.",
						},
						"user_identifier": schema.Int64Attribute{
							Computed:    true,
							Description: "User ID (numeric) or SID (string) that owns this API key.",
						},
						"keyhash": schema.StringAttribute{
							Computed:    true,
							Description: "Hashed representation of the API key (masked for security).",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the API key was created.",
						},
						"expires_at": schema.StringAttribute{
							Computed:    true,
							Description: "Expiration timestamp for the API key or `null` for no expiration.",
						},
						"local": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this API key is for local system use only.",
						},
						"revoked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the API key has been revoked and is no longer valid.",
						},
						"revoked_reason": schema.StringAttribute{
							Computed:    true,
							Description: "Reason for API key revocation or `null` if not revoked.",
						},
					},
				},
			},
		},
	}
}

func (d *ApiKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApiKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApiKeysDataSourceModel

	// Call query method with empty filters to get all items
	result, err := d.client.Call("api_key.query", []interface{}{})
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to query api_keys: %s", err.Error()))
		return
	}

	resultList, ok := result.([]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response as list")
		return
	}

	// Convert results to items
	items := make([]ApiKeysItemModel, 0, len(resultList))
	for _, item := range resultList {
		resultMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemModel := ApiKeysItemModel{}
		if v, ok := resultMap["id"]; ok && v != nil {
			itemModel.ID = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["name"]; ok && v != nil {
			itemModel.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			itemModel.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user_identifier"]; ok && v != nil {
			if fv, ok := v.(float64); ok {
				itemModel.UserIdentifier = types.Int64Value(int64(fv))
			}
		}
		if v, ok := resultMap["keyhash"]; ok && v != nil {
			itemModel.Keyhash = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["created_at"]; ok && v != nil {
			itemModel.CreatedAt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["expires_at"]; ok && v != nil {
			itemModel.ExpiresAt = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["local"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Local = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["revoked"]; ok && v != nil {
			if bv, ok := v.(bool); ok {
				itemModel.Revoked = types.BoolValue(bv)
			}
		}
		if v, ok := resultMap["revoked_reason"]; ok && v != nil {
			itemModel.RevokedReason = types.StringValue(fmt.Sprintf("%v", v))
		}
		items = append(items, itemModel)
	}

	// Convert to types.List
	itemsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"created_at":      types.StringType,
			"expires_at":      types.StringType,
			"id":              types.StringType,
			"keyhash":         types.StringType,
			"local":           types.BoolType,
			"name":            types.StringType,
			"revoked":         types.BoolType,
			"revoked_reason":  types.StringType,
			"user_identifier": types.Int64Type,
			"username":        types.StringType,
		},
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Items = itemsValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
