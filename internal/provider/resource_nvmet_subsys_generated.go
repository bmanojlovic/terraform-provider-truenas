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

type NvmetSubsysResource struct {
	client *client.Client
}

type NvmetSubsysResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Subnqn types.String `tfsdk:"subnqn"`
	AllowAnyHost types.Bool `tfsdk:"allow_any_host"`
	PiEnable types.String `tfsdk:"pi_enable"`
	QidMax types.Int64 `tfsdk:"qid_max"`
	IeeeOui types.String `tfsdk:"ieee_oui"`
	Ana types.String `tfsdk:"ana"`
}

func NewNvmetSubsysResource() resource.Resource {
	return &NvmetSubsysResource{}
}

func (r *NvmetSubsysResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_subsys"
}

func (r *NvmetSubsysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NvmetSubsysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a NVMe target subsystem (`subsys`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Human readable name for the subsystem.  If `subnqn` is not provided on creation, then this name will",
			},
			"subnqn": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "NVMe Qualified Name (NQN) for the subsystem.  Must be a valid NQN format if provided.",
			},
			"allow_any_host": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Any host can access the storage associated with this subsystem (i.e. no access control).",
			},
			"pi_enable": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Enable Protection Information (PI) for data integrity checking.",
			},
			"qid_max": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum number of queue IDs allowed for this subsystem.",
			},
			"ieee_oui": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "IEEE Organizationally Unique Identifier for the subsystem.",
			},
			"ana": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "If set to either `True` or `False`, then *override* the global `ana` setting from `nvmet.global.conf",
			},
		},
	}
}

func (r *NvmetSubsysResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NvmetSubsysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Subnqn.IsNull() {
		params["subnqn"] = data.Subnqn.ValueString()
	}
	if !data.AllowAnyHost.IsNull() {
		params["allow_any_host"] = data.AllowAnyHost.ValueBool()
	}
	if !data.PiEnable.IsNull() {
		params["pi_enable"] = data.PiEnable.ValueString()
	}
	if !data.QidMax.IsNull() {
		params["qid_max"] = data.QidMax.ValueInt64()
	}
	if !data.IeeeOui.IsNull() {
		params["ieee_oui"] = data.IeeeOui.ValueString()
	}
	if !data.Ana.IsNull() {
		params["ana"] = data.Ana.ValueString()
	}

	result, err := r.client.Call("nvmet.subsys.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create nvmet_subsys: %s", err))
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

func (r *NvmetSubsysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("nvmet.subsys.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read nvmet_subsys: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["subnqn"]; ok && v != nil {
			data.Subnqn = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["allow_any_host"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AllowAnyHost = types.BoolValue(bv) }
		}
		if v, ok := resultMap["pi_enable"]; ok && v != nil {
			data.PiEnable = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["qid_max"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.QidMax = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["ieee_oui"]; ok && v != nil {
			data.IeeeOui = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["ana"]; ok && v != nil {
			data.Ana = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetSubsysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NvmetSubsysResourceModel
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
	if !data.Subnqn.IsNull() {
		params["subnqn"] = data.Subnqn.ValueString()
	}
	if !data.AllowAnyHost.IsNull() {
		params["allow_any_host"] = data.AllowAnyHost.ValueBool()
	}
	if !data.PiEnable.IsNull() {
		params["pi_enable"] = data.PiEnable.ValueString()
	}
	if !data.QidMax.IsNull() {
		params["qid_max"] = data.QidMax.ValueInt64()
	}
	if !data.IeeeOui.IsNull() {
		params["ieee_oui"] = data.IeeeOui.ValueString()
	}
	if !data.Ana.IsNull() {
		params["ana"] = data.Ana.ValueString()
	}

	_, err = r.client.Call("nvmet.subsys.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update nvmet_subsys: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetSubsysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("nvmet.subsys.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete nvmet_subsys: %s", err))
		return
	}
}
