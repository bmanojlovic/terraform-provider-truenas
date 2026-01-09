package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type AppRegistryResource struct {
	client *client.Client
}

type AppRegistryResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Uri types.String `tfsdk:"uri"`
}

func NewAppRegistryResource() resource.Resource {
	return &AppRegistryResource{}
}

func (r *AppRegistryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_registry"
}

func (r *AppRegistryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS app_registry resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"username": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"password": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"uri": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *AppRegistryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
		return
	}
	r.client = client
}

func (r *AppRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	params["username"] = data.Username.ValueString()
	params["password"] = data.Password.ValueString()
	if !data.Uri.IsNull() {
		params["uri"] = data.Uri.ValueString()
	}

	result, err := r.client.Call("app.registry.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("app.registry.get_instance", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	params["username"] = data.Username.ValueString()
	params["password"] = data.Password.ValueString()
	if !data.Uri.IsNull() {
		params["uri"] = data.Uri.ValueString()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("app.registry.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("app.registry.delete", resourceID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
