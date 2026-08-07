package provider

import (
	"context"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type SystemNtpserverResource struct {
	client *client.Client
}

type SystemNtpserverResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
	Burst   types.Bool   `tfsdk:"burst"`
	Iburst  types.Bool   `tfsdk:"iburst"`
	Prefer  types.Bool   `tfsdk:"prefer"`
	Minpoll types.Int64  `tfsdk:"minpoll"`
	Maxpoll types.Int64  `tfsdk:"maxpoll"`
	Force   types.Bool   `tfsdk:"force"`
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
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"address": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Hostname or IP address of the NTP server.",
			},
			"burst": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Send a burst of packets when the server is reachable.",
			},
			"iburst": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Send a burst of packets when the server is unreachable.",
			},
			"prefer": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Mark this server as preferred for time synchronization.",
			},
			"minpoll": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum polling interval (log2 seconds).",
			},
			"maxpoll": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum polling interval (log2 seconds).",
			},
			"force": schema.BoolAttribute{
				Required:    false,
				Optional:    true,
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
	if !data.Address.IsNull() && !data.Address.IsUnknown() {
		params["address"] = data.Address.ValueString()
	}
	if !data.Burst.IsNull() && !data.Burst.IsUnknown() {
		params["burst"] = data.Burst.ValueBool()
	}
	if !data.Iburst.IsNull() && !data.Iburst.IsUnknown() {
		params["iburst"] = data.Iburst.ValueBool()
	}
	if !data.Prefer.IsNull() && !data.Prefer.IsUnknown() {
		params["prefer"] = data.Prefer.ValueBool()
	}
	if !data.Minpoll.IsNull() && !data.Minpoll.IsUnknown() {
		params["minpoll"] = data.Minpoll.ValueInt64()
	}
	if !data.Maxpoll.IsNull() && !data.Maxpoll.IsUnknown() {
		params["maxpoll"] = data.Maxpoll.ValueInt64()
	}
	if !data.Force.IsNull() && !data.Force.IsUnknown() {
		params["force"] = data.Force.ValueBool()
	}

	result, err := r.client.Call("system.ntpserver.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create system_ntpserver: %s", err))
		return
	}

	// Extract ID from result
	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists && id != nil {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Validate ID was set
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Create Error", "API did not return a valid ID")
		return
	}

	// Read back to populate computed fields
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("system.ntpserver.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back system_ntpserver: %s", err))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["address"]; ok {
		switch val := v.(type) {
		case string:
			data.Address = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Address = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Address = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Burst.IsUnknown() {
		data.Burst = types.BoolNull()
	}
	if data.Iburst.IsUnknown() {
		data.Iburst = types.BoolNull()
	}
	if data.Prefer.IsUnknown() {
		data.Prefer = types.BoolNull()
	}
	if data.Minpoll.IsUnknown() {
		data.Minpoll = types.Int64Null()
	}
	if data.Maxpoll.IsUnknown() {
		data.Maxpoll = types.Int64Null()
	}
	if data.Force.IsUnknown() {
		data.Force = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemNtpserverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SystemNtpserverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("system.ntpserver.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read system_ntpserver: %s", err))
		return
	}

	// Map result back to state
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["address"]; ok {
		switch val := v.(type) {
		case string:
			data.Address = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Address = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Address = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Burst.IsUnknown() {
		data.Burst = types.BoolNull()
	}
	if data.Iburst.IsUnknown() {
		data.Iburst = types.BoolNull()
	}
	if data.Prefer.IsUnknown() {
		data.Prefer = types.BoolNull()
	}
	if data.Minpoll.IsUnknown() {
		data.Minpoll = types.Int64Null()
	}
	if data.Maxpoll.IsUnknown() {
		data.Maxpoll = types.Int64Null()
	}
	if data.Force.IsUnknown() {
		data.Force = types.BoolNull()
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Address.IsNull() && !data.Address.IsUnknown() {
		params["address"] = data.Address.ValueString()
	}
	if !data.Burst.IsNull() && !data.Burst.IsUnknown() {
		params["burst"] = data.Burst.ValueBool()
	}
	if !data.Iburst.IsNull() && !data.Iburst.IsUnknown() {
		params["iburst"] = data.Iburst.ValueBool()
	}
	if !data.Prefer.IsNull() && !data.Prefer.IsUnknown() {
		params["prefer"] = data.Prefer.ValueBool()
	}
	if !data.Minpoll.IsNull() && !data.Minpoll.IsUnknown() {
		params["minpoll"] = data.Minpoll.ValueInt64()
	}
	if !data.Maxpoll.IsNull() && !data.Maxpoll.IsUnknown() {
		params["maxpoll"] = data.Maxpoll.ValueInt64()
	}
	if !data.Force.IsNull() && !data.Force.IsUnknown() {
		params["force"] = data.Force.ValueBool()
	}

	_, err = r.client.Call("system.ntpserver.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update system_ntpserver: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("system.ntpserver.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back system_ntpserver: %s", readErr))
		return
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Parse Error", "Failed to parse API response")
		return
	}

	if v, ok := resultMap["id"]; ok && v != nil {
		data.ID = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := resultMap["address"]; ok {
		switch val := v.(type) {
		case string:
			data.Address = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Address = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Address = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Burst.IsUnknown() {
		data.Burst = types.BoolNull()
	}
	if data.Iburst.IsUnknown() {
		data.Iburst = types.BoolNull()
	}
	if data.Prefer.IsUnknown() {
		data.Prefer = types.BoolNull()
	}
	if data.Minpoll.IsUnknown() {
		data.Minpoll = types.Int64Null()
	}
	if data.Maxpoll.IsUnknown() {
		data.Maxpoll = types.Int64Null()
	}
	if data.Force.IsUnknown() {
		data.Force = types.BoolNull()
	}

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
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete system_ntpserver: %s", err))
		return
	}
}
