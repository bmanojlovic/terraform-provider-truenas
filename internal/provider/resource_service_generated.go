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

type ServiceResource struct {
	client *client.Client
}

type ServiceResourceModel struct {
	ID types.String `tfsdk:"id"`
	StartOnCreate types.Bool `tfsdk:"start_on_create"`
	Enable types.Bool `tfsdk:"enable"`
}

func NewServiceResource() resource.Resource {
	return &ServiceResource{}
}

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS service resource (update-only)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"start_on_create": schema.BoolAttribute{
				Optional: true,
				Description: "Start the resource immediately after creation (default: true if not specified)",
			},
			"enable": schema.BoolAttribute{
				Required: true,
				Optional: false,
			},
		},
	}
}

func (r *ServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update-only resource: ID must be provided
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing ID", "ID is required for update-only resources. Use terraform import to manage existing resources.")
		return
	}

	// Convert string ID to integer for API call
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID must be a valid integer")
		return
	}

	params := map[string]interface{}{}
	params["enable"] = data.Enable.ValueBool()

	_, err = r.client.Call("service.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for API call
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID must be a valid integer")
		return
	}

	result, err := r.client.Call("service.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if val, ok := resultMap["enable"]; ok && val != nil {
			if boolVal, ok := val.(bool); ok {
				data.Enable = types.BoolValue(boolVal)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to integer for API call
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID must be a valid integer")
		return
	}

	params := map[string]interface{}{}
	params["enable"] = data.Enable.ValueBool()

	result, err := r.client.Call("service.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	// Update may return new ID
	if result != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", result))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Update-only resource: just remove from state, don't delete on server
}

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID (e.g., "ssh" for SSH service)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
