package provider

import (
	"context"
	"fmt"


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
	QidMax types.String `tfsdk:"qid_max"`
	IeeeOui types.String `tfsdk:"ieee_oui"`
	Ana types.String `tfsdk:"ana"`
}

func NewNvmetSubsysResource() resource.Resource {
	return &NvmetSubsysResource{}
}

func (r *NvmetSubsysResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_subsys"
}

func (r *NvmetSubsysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS nvmet_subsys resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"subnqn": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"allow_any_host": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
			"pi_enable": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"qid_max": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"ieee_oui": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"ana": schema.StringAttribute{
				Required: false,
				Optional: true,
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
	params["name"] = data.Name.ValueString()
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
		params["qid_max"] = data.QidMax.ValueString()
	}
	if !data.IeeeOui.IsNull() {
		params["ieee_oui"] = data.IeeeOui.ValueString()
	}
	if !data.Ana.IsNull() {
		params["ana"] = data.Ana.ValueString()
	}

	result, err := r.client.Call("nvmet/subsys.create", params)
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

func (r *NvmetSubsysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("nvmet/subsys.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetSubsysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
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
		params["qid_max"] = data.QidMax.ValueString()
	}
	if !data.IeeeOui.IsNull() {
		params["ieee_oui"] = data.IeeeOui.ValueString()
	}
	if !data.Ana.IsNull() {
		params["ana"] = data.Ana.ValueString()
	}

	_, err := r.client.Call("nvmet/subsys.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetSubsysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetSubsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("nvmet/subsys.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
