package provider

import (
	"context"
	"fmt"
	"strconv"

	"time"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	Memory types.Int64 `tfsdk:"memory"`
	PrivilegedMode types.Bool `tfsdk:"privileged_mode"`
}

func NewVirtInstanceResource() resource.Resource {
	return &VirtInstanceResource{}
}

func (r *VirtInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virt_instance"
}

func (r *VirtInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *VirtInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new virtualized instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"start_on_create": schema.BoolAttribute{Optional: true, Description: "Start the resource immediately after creation (default: true)"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name for the new virtual instance.",
			},
			"source_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Source type for instance creation.",
			},
			"storage_pool": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Storage pool under which to allocate root filesystem. Must be one of the pools     listed in virt.gl",
			},
			"image": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Image identifier to use for creating the instance.",
			},
			"root_disk_size": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "This can be specified when creating VMs so the root device's size can be configured. Root device for",
			},
			"root_disk_io_bus": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "I/O bus type for the root disk (defaults to NVME for best performance).",
			},
			"remote": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Remote image source to use.",
			},
			"instance_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Type of instance to create.",
			},
			"environment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Environment variables to set inside the instance.",
			},
			"autostart": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the instance should automatically start when the host boots.",
			},
			"cpu": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "CPU allocation specification or `null` for automatic allocation.",
			},
			"devices": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Array of devices to attach to the instance.",
			},
			"memory": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Memory allocation in bytes or `null` for automatic allocation.",
			},
			"privileged_mode": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "This is only valid for containers and should only be set when container instance which is to be depl",
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.SourceType.IsNull() {
		params["source_type"] = data.SourceType.ValueString()
	}
	if !data.StoragePool.IsNull() {
		params["storage_pool"] = data.StoragePool.ValueString()
	}
	if !data.Image.IsNull() {
		params["image"] = data.Image.ValueString()
	}
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
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.PrivilegedMode.IsNull() {
		params["privileged_mode"] = data.PrivilegedMode.ValueBool()
	}

	result, err := r.client.CallWithJob("virt.instance.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create virt_instance: %s", err))
		return
	}

	// Extract ID from result
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
		vmID, err := strconv.Atoi(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
			return
		}
		_, err = r.client.Call("virt.instance.start", vmID)
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

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("virt.instance.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read virt_instance: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["source_type"]; ok && v != nil {
			data.SourceType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["storage_pool"]; ok && v != nil {
			data.StoragePool = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["image"]; ok && v != nil {
			data.Image = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["root_disk_size"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.RootDiskSize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["root_disk_io_bus"]; ok && v != nil {
			data.RootDiskIoBus = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["remote"]; ok && v != nil {
			data.Remote = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["instance_type"]; ok && v != nil {
			data.InstanceType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["environment"]; ok && v != nil {
			data.Environment = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["autostart"]; ok && v != nil {
			data.Autostart = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["cpu"]; ok && v != nil {
			data.Cpu = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["devices"]; ok && v != nil {
			data.Devices = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["memory"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Memory = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["privileged_mode"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.PrivilegedMode = types.BoolValue(bv) }
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VirtInstanceResourceModel
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
	if !data.SourceType.IsNull() {
		params["source_type"] = data.SourceType.ValueString()
	}
	if !data.StoragePool.IsNull() {
		params["storage_pool"] = data.StoragePool.ValueString()
	}
	if !data.Image.IsNull() {
		params["image"] = data.Image.ValueString()
	}
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
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.PrivilegedMode.IsNull() {
		params["privileged_mode"] = data.PrivilegedMode.ValueBool()
	}

	_, err = r.client.CallWithJob("virt.instance.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update virt_instance: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	// Stop VM before deletion if running
	vmID, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ID Conversion Error", fmt.Sprintf("Failed to convert ID to integer: %s", err.Error()))
		return
	}
	_, _ = r.client.Call("virt.instance.stop", vmID)  // Ignore errors - VM might already be stopped
	time.Sleep(2 * time.Second)  // Wait for VM to stop

	_, err = r.client.CallWithJob("virt.instance.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete virt_instance: %s", err))
		return
	}
}
