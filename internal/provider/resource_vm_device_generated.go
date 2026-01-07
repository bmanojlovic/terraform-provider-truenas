package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type VmDeviceResource struct {
	client *client.Client
}

type VmDeviceResourceModel struct {
	ID types.String `tfsdk:"id"`
	Attributes types.String `tfsdk:"attributes"`
	Vm types.Int64 `tfsdk:"vm"`
	Order types.String `tfsdk:"order"`
}

func NewVmDeviceResource() resource.Resource {
	return &VmDeviceResource{}
}

func (r *VmDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_device"
}

func (r *VmDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS vm_device resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"attributes": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"vm": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"order": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *VmDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VmDeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"attributes": data.Attributes.ValueString(),
		"vm": data.Vm.ValueInt64(),
		"order": data.Order.ValueString(),
	}

	result, err := r.client.Call("vm/device.create", params)
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

func (r *VmDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("vm/device.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{
		"attributes": data.Attributes.ValueString(),
		"vm": data.Vm.ValueInt64(),
		"order": data.Order.ValueString(),
	}

	_, err := r.client.Call("vm/device.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("vm/device.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
