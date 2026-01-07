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

type RsynctaskRunActionResource struct {
	client *client.Client
}

type RsynctaskRunActionResourceModel struct {
	ID types.String `tfsdk:"id"`
	ResourceID types.String `tfsdk:"resource_id"`
}

func NewRsynctaskRunActionResource() resource.Resource {
	return &RsynctaskRunActionResource{}
}

func (r *RsynctaskRunActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rsynctask_run_action"
}

func (r *RsynctaskRunActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Executes run action on rsynctask resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the resource to perform action on",
			},
		},
	}
}

func (r *RsynctaskRunActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RsynctaskRunActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RsynctaskRunActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No additional parameters

	_, err := r.client.Call("rsynctask/id/{id_}/run", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute run: %s", err.Error()))
		return
	}

	// Use timestamp as ID since actions are ephemeral
	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskRunActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Actions are ephemeral - nothing to read
	var data RsynctaskRunActionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *RsynctaskRunActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Actions are immutable - re-execute on update
	var data RsynctaskRunActionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No additional parameters

	_, err := r.client.Call("rsynctask/id/{id_}/run", data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Action Failed", fmt.Sprintf("Failed to execute run: %s", err.Error()))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s-%d", data.ResourceID.ValueString(), time.Now().Unix()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskRunActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Actions cannot be undone - just remove from state
}
