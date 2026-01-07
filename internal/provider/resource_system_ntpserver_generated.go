package provider

import (
	"context"
	"fmt"

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

	params := map[string]interface{}{
		"address": data.Address.ValueString(),
		"burst": data.Burst.ValueBool(),
		"iburst": data.Iburst.ValueBool(),
		"prefer": data.Prefer.ValueBool(),
		"minpoll": data.Minpoll.ValueInt64(),
		"maxpoll": data.Maxpoll.ValueInt64(),
		"force": data.Force.ValueBool(),
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

	params := map[string]interface{}{
		"address": data.Address.ValueString(),
		"burst": data.Burst.ValueBool(),
		"iburst": data.Iburst.ValueBool(),
		"prefer": data.Prefer.ValueBool(),
		"minpoll": data.Minpoll.ValueInt64(),
		"maxpoll": data.Maxpoll.ValueInt64(),
		"force": data.Force.ValueBool(),
	}

	_, err := r.client.Call("system/ntpserver.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
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
