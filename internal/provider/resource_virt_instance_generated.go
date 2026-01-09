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

type VirtInstanceResource struct {
	client *client.Client
}

type VirtInstanceResourceModel struct {
	ID types.String `tfsdk:"id"`
	StartOnCreate types.Bool `tfsdk:"start_on_create"`
	Name types.String `tfsdk:"name"`
	SourceType types.String `tfsdk:"source_type"`
	StoragePool types.String `tfsdk:"storage_pool"`
	Image types.String `tfsdk:"image"`
	RootDiskSize types.Int64 `tfsdk:"root_disk_size"`
	RootDiskIoBus types.String `tfsdk:"root_disk_io_bus"`
	Remote types.String `tfsdk:"remote"`
	InstanceType types.String `tfsdk:"instance_type"`
	Environment types.String `tfsdk:"environment"`
	Autostart types.String `tfsdk:"autostart"`
	Cpu types.String `tfsdk:"cpu"`
	Devices types.String `tfsdk:"devices"`
	Memory types.String `tfsdk:"memory"`
	PrivilegedMode types.Bool `tfsdk:"privileged_mode"`
}

func NewVirtInstanceResource() resource.Resource {
	return &VirtInstanceResource{}
}

func (r *VirtInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virt_instance"
}

func (r *VirtInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TrueNAS virt_instance resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"start_on_create": schema.BoolAttribute{
				Optional: true,
				Description: "Start the resource immediately after creation (default: true if not specified)",
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"source_type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"storage_pool": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"image": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"root_disk_size": schema.Int64Attribute{
				Required: false,
				Optional: true,
			},
			"root_disk_io_bus": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"remote": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"instance_type": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"environment": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"autostart": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"cpu": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"devices": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"memory": schema.StringAttribute{
				Required: false,
				Optional: true,
			},
			"privileged_mode": schema.BoolAttribute{
				Required: false,
				Optional: true,
			},
		},
	}
}

func (r *VirtInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VirtInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
	if !data.SourceType.IsNull() {
		params["source_type"] = data.SourceType.ValueString()
	}
	if !data.StoragePool.IsNull() {
		params["storage_pool"] = data.StoragePool.ValueString()
	}
	params["image"] = data.Image.ValueString()
	if !data.RootDiskSize.IsNull() {
		params["root_disk_size"] = data.RootDiskSize.ValueInt64()
	}
	if !data.RootDiskIoBus.IsNull() {
		params["root_disk_io_bus"] = data.RootDiskIoBus.ValueString()
	}
	if !data.Remote.IsNull() {
		params["remote"] = data.Remote.ValueString()
	}
	if !data.InstanceType.IsNull() {
		params["instance_type"] = data.InstanceType.ValueString()
	}
	if !data.Environment.IsNull() {
		params["environment"] = data.Environment.ValueString()
	}
	if !data.Autostart.IsNull() {
		params["autostart"] = data.Autostart.ValueString()
	}
	if !data.Cpu.IsNull() {
		params["cpu"] = data.Cpu.ValueString()
	}
	if !data.Devices.IsNull() {
		params["devices"] = data.Devices.ValueString()
	}
	if !data.Memory.IsNull() {
		params["memory"] = data.Memory.ValueString()
	}
	if !data.PrivilegedMode.IsNull() {
		params["privileged_mode"] = data.PrivilegedMode.ValueBool()
	}

	result, err := r.client.Call("virt/instance.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if id, exists := resultMap["id"]; exists {
			data.ID = types.StringValue(fmt.Sprintf("%v", id))
		}
	}

	// Handle lifecycle action - start on create if requested
	startOnCreate := true  // default when not specified
	if !data.StartOnCreate.IsNull() {
		startOnCreate = data.StartOnCreate.ValueBool()
	}
	if startOnCreate {
		// Convert string ID to integer for TrueNAS API
		vmID, err := strconv.Atoi(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
			return
		}
		_, err = r.client.Call("virt/instance.start", vmID)
		if err != nil {
			resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error()))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("virt/instance.get_instance", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	params["name"] = data.Name.ValueString()
	if !data.SourceType.IsNull() {
		params["source_type"] = data.SourceType.ValueString()
	}
	if !data.StoragePool.IsNull() {
		params["storage_pool"] = data.StoragePool.ValueString()
	}
	params["image"] = data.Image.ValueString()
	if !data.RootDiskSize.IsNull() {
		params["root_disk_size"] = data.RootDiskSize.ValueInt64()
	}
	if !data.RootDiskIoBus.IsNull() {
		params["root_disk_io_bus"] = data.RootDiskIoBus.ValueString()
	}
	if !data.Remote.IsNull() {
		params["remote"] = data.Remote.ValueString()
	}
	if !data.InstanceType.IsNull() {
		params["instance_type"] = data.InstanceType.ValueString()
	}
	if !data.Environment.IsNull() {
		params["environment"] = data.Environment.ValueString()
	}
	if !data.Autostart.IsNull() {
		params["autostart"] = data.Autostart.ValueString()
	}
	if !data.Cpu.IsNull() {
		params["cpu"] = data.Cpu.ValueString()
	}
	if !data.Devices.IsNull() {
		params["devices"] = data.Devices.ValueString()
	}
	if !data.Memory.IsNull() {
		params["memory"] = data.Memory.ValueString()
	}
	if !data.PrivilegedMode.IsNull() {
		params["privileged_mode"] = data.PrivilegedMode.ValueBool()
	}

	_, err := r.client.Call("virt/instance.update", []interface{}{data.ID.ValueString(), params})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Call("virt/instance.delete", data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}
