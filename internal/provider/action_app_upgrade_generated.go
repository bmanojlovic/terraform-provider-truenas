package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type AppUpgradeActionResource struct {
	client *client.Client
}

type AppUpgradeActionResourceModel struct {
	ID types.String `tfsdk:"id"`
	ResourceID types.String `tfsdk:"resource_id"`
	AppName types.String `tfsdk:"app_name"`
	Options types.String `tfsdk:"options"`
}

func NewAppUpgradeActionResource() resource.Resource {
	return &AppUpgradeActionResource{}
}

func (r *AppUpgradeActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_upgrade_action"
}

func (r *AppUpgradeActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Executes upgrade action on app resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the resource to perform action on",
			},
			"app_name": schema.StringAttribute{
				Optional: true,
			},
			"options": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *AppUpgradeActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppUpgradeActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AppUpgradeActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
		if !data.AppName.IsNull() {
			params["app_name"] = data.AppName.ValueString()
		}
		if !data.Options.IsNull() {
			params["options"] = data.Options.ValueString()
		}

	_, err := r.client.Call("app/upgrade", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute upgrade: %s", err.Error()))
		return
	}

	// Use timestamp as ID since actions are ephemeral
	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppUpgradeActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are ephemeral - nothing to read
	var data AppUpgradeActionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *AppUpgradeActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Actions are immutable - re-execute on update
	var data AppUpgradeActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
		if !data.AppName.IsNull() {
			params["app_name"] = data.AppName.ValueString()
		}
		if !data.Options.IsNull() {
			params["options"] = data.Options.ValueString()
		}

	_, err := r.client.Call("app/upgrade", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute upgrade: %s", err.Error()))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppUpgradeActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Actions cannot be undone - just remove from state
}
