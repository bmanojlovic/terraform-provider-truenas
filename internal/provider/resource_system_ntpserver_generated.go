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

func (r *SystemNtpserverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SystemNtpserverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Add an NTP Server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"address": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Hostname or IP address of the NTP server.",
			},
			"burst": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Send a burst of packets when the server is reachable.",
			},
			"iburst": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Send a burst of packets when the server is unreachable.",
			},
			"prefer": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Mark this server as preferred for time synchronization.",
			},
			"minpoll": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Minimum polling interval (log2 seconds).",
			},
			"maxpoll": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum polling interval (log2 seconds).",
			},
			"force": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Force creation even if the server is unreachable.",
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
	if !data.Address.IsNull() {
		params["address"] = data.Address.ValueString()
	}
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

	result, err := r.client.Call("system.ntpserver.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create system_ntpserver: %s", err))
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

func (r *SystemNtpserverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("system.ntpserver.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read system_ntpserver: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["address"]; ok && v != nil {
			data.Address = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["burst"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Burst = types.BoolValue(bv) }
		}
		if v, ok := resultMap["iburst"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Iburst = types.BoolValue(bv) }
		}
		if v, ok := resultMap["prefer"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Prefer = types.BoolValue(bv) }
		}
		if v, ok := resultMap["minpoll"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Minpoll = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["maxpoll"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Maxpoll = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["force"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Force = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemNtpserverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SystemNtpserverResourceModel
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
	if !data.Address.IsNull() {
		params["address"] = data.Address.ValueString()
	}
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

	_, err = r.client.Call("system.ntpserver.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update system_ntpserver: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemNtpserverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("system.ntpserver.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete system_ntpserver: %s", err))
		return
	}
}
