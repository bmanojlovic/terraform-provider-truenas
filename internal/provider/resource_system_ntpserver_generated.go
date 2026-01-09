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

type SystemNtpserverResource struct {
	client *client.Client
}

type SystemNtpserverResourceModel struct {
	ID types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
	Burst types.Bool `tfsdk:"burst"`
	Iburst types.Bool `tfsdk:"iburst"`
	Prefer types.Bool `tfsdk:"prefer"`
	Minpoll types.Int64 `tfsdk:"minpoll"`
	Maxpoll types.Int64 `tfsdk:"maxpoll"`
	Force types.Bool `tfsdk:"force"`
}

func NewSystemNtpserverResource() resource.Resource {
	return &SystemNtpserverResource{}
}

func (r *SystemNtpserverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_ntpserver"
}

func (r *SystemNtpserverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS system_ntpserver resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"address": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"burst": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"iburst": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"prefer": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"minpoll": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"maxpoll": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"force": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *SystemNtpserverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SystemNtpserverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["address"] = data.Address.ValueString()
	if !data.Burst.IsNull() {
		params["burst"] = data.Burst.ValueBool()
	}
	if !data.Iburst.IsNull() {
		params["iburst"] = data.Iburst.ValueBool()
	}
	if !data.Prefer.IsNull() {
		params["prefer"] = data.Prefer.ValueBool()
	}
	if !data.Minpoll.IsNull() {
		params["minpoll"] = data.Minpoll.ValueInt64()
	}
	if !data.Maxpoll.IsNull() {
		params["maxpoll"] = data.Maxpoll.ValueInt64()
	}
	if !data.Force.IsNull() {
		params["force"] = data.Force.ValueBool()
	}

	result, err := r.client.Call("system/ntpserver.create", params)
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

func (r *SystemNtpserverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("system/ntpserver.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemNtpserverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get ID from current state (not plan)
	var state SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["address"] = data.Address.ValueString()
	if !data.Burst.IsNull() {
		params["burst"] = data.Burst.ValueBool()
	}
	if !data.Iburst.IsNull() {
		params["iburst"] = data.Iburst.ValueBool()
	}
	if !data.Prefer.IsNull() {
		params["prefer"] = data.Prefer.ValueBool()
	}
	if !data.Minpoll.IsNull() {
		params["minpoll"] = data.Minpoll.ValueInt64()
	}
	if !data.Maxpoll.IsNull() {
		params["maxpoll"] = data.Maxpoll.ValueInt64()
	}
	if !data.Force.IsNull() {
		params["force"] = data.Force.ValueBool()
	}

	// Convert string ID to integer for TrueNAS API
	resourceID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}

	_, err = r.client.Call("system/ntpserver.update", []interface{}{resourceID, params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	// Preserve the ID in the new state
	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemNtpserverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("system/ntpserver.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
