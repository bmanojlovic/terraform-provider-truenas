package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
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

func (r *AppRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AppRegistryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an app registry entry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Human-readable name for the container registry.",
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional description of the container registry or `null`.",
			},
			"username": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Username for registry authentication (masked for security).",
			},
			"password": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Password or access token for registry authentication (masked for security).",
			},
			"uri": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Container registry URI endpoint (defaults to Docker Hub).",
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Username.IsNull() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}
	if !data.Uri.IsNull() {
		params["uri"] = data.Uri.ValueString()
	}

	result, err := r.client.Call("app.registry.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create app_registry: %s", err))
		return
	}

	// Extract ID from result
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

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("app.registry.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read app_registry: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["description"]; ok && v != nil {
			data.Description = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["uri"]; ok && v != nil {
			data.Uri = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		params["description"] = data.Description.ValueString()
	}
	if !data.Username.IsNull() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}
	if !data.Uri.IsNull() {
		params["uri"] = data.Uri.ValueString()
	}

	_, err = r.client.Call("app.registry.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update app_registry: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("app.registry.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete app_registry: %s", err))
		return
	}
}
