package provider

import (
	"context"
	"fmt"
	"strings"
"encoding/json"
	"time"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Autostart types.Bool `tfsdk:"autostart"`
	Cpu types.String `tfsdk:"cpu"`
	Devices types.List `tfsdk:"devices"`
	Memory types.Int64 `tfsdk:"memory"`
	PrivilegedMode types.Bool `tfsdk:"privileged_mode"`
	VncPort types.Int64 `tfsdk:"vnc_port"`
	EnableVnc types.Bool `tfsdk:"enable_vnc"`
	VncPassword types.String `tfsdk:"vnc_password"`
	SecureBoot types.Bool `tfsdk:"secure_boot"`
	ImageOs types.String `tfsdk:"image_os"`
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"storage_pool": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Storage pool under which to allocate root filesystem. Must be one of the pools     listed in virt.gl",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"image": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Image identifier to use for creating the instance.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"root_disk_size": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Size of the root disk in GB (minimum 5GB) or `null` to keep current size.",
			},
			"root_disk_io_bus": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "I/O bus type for the root disk or `null` to keep current setting.",
			},
			"remote": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Remote image source to use.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"instance_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Type of instance to create.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"environment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Environment variables to set inside the instance.",
			},
			"autostart": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the instance should automatically start when the host boots.",
			},
			"cpu": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "CPU allocation specification or `null` for automatic allocation.",
			},
			"devices": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
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
			"vnc_port": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "TCP port number for VNC access (5900-65535) or `null` to disable VNC.",
			},
			"enable_vnc": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable VNC remote access for the instance.",
			},
			"vnc_password": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Setting vnc_password to null will unset VNC password.",
			},
			"secure_boot": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable UEFI Secure Boot (VMs only).",
			},
			"image_os": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Operating system type for the instance or `null` for auto-detection.",
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
		var environmentObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Environment.ValueString()), &environmentObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse environment: %s", err))
			return
		}
		params["environment"] = environmentObj
	}
	if !data.Autostart.IsNull() {
		params["autostart"] = data.Autostart.ValueBool()
	}
	if !data.Cpu.IsNull() {
		params["cpu"] = data.Cpu.ValueString()
	}
	if !data.Devices.IsNull() {
		var devicesList []string
		data.Devices.ElementsAs(ctx, &devicesList, false)
		params["devices"] = devicesList
	}
	if !data.Memory.IsNull() {
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.PrivilegedMode.IsNull() {
		params["privileged_mode"] = data.PrivilegedMode.ValueBool()
	}
	if !data.VncPort.IsNull() {
		params["vnc_port"] = data.VncPort.ValueInt64()
	}
	if !data.EnableVnc.IsNull() {
		params["enable_vnc"] = data.EnableVnc.ValueBool()
	}
	if !data.VncPassword.IsNull() {
		params["vnc_password"] = data.VncPassword.ValueString()
	}
	if !data.SecureBoot.IsNull() {
		params["secure_boot"] = data.SecureBoot.ValueBool()
	}
	if !data.ImageOs.IsNull() {
		params["image_os"] = data.ImageOs.ValueString()
	}

	result, err := r.client.CallWithJob("virt.instance.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create virt_instance: %s", err))
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

	startOnCreate := true
	if !data.StartOnCreate.IsNull() { startOnCreate = data.StartOnCreate.ValueBool() }
	if startOnCreate {
		_, err = r.client.Call("virt.instance.start", data.ID.ValueString())
		if err != nil { resp.Diagnostics.AddWarning("Start Failed", fmt.Sprintf("Resource created but failed to start: %s", err.Error())) }
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VirtInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VirtInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	result, err := r.client.Call("virt.instance.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read virt_instance: %s", err))
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
		if v, ok := resultMap["name"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Name = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Name = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Name = types.StringValue(fmt.Sprintf("%v", v))
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

	var id interface{}
	var err error
	id = state.ID.ValueString()

	params := map[string]interface{}{}
	if !data.Environment.IsNull() {
		var environmentObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Environment.ValueString()), &environmentObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse environment: %s", err))
			return
		}
		params["environment"] = environmentObj
	}
	if !data.Autostart.IsNull() {
		params["autostart"] = data.Autostart.ValueBool()
	}
	if !data.Cpu.IsNull() {
		params["cpu"] = data.Cpu.ValueString()
	}
	if !data.Memory.IsNull() {
		params["memory"] = data.Memory.ValueInt64()
	}
	if !data.VncPort.IsNull() {
		params["vnc_port"] = data.VncPort.ValueInt64()
	}
	if !data.EnableVnc.IsNull() {
		params["enable_vnc"] = data.EnableVnc.ValueBool()
	}
	if !data.VncPassword.IsNull() {
		params["vnc_password"] = data.VncPassword.ValueString()
	}
	if !data.SecureBoot.IsNull() {
		params["secure_boot"] = data.SecureBoot.ValueBool()
	}
	if !data.RootDiskSize.IsNull() {
		params["root_disk_size"] = data.RootDiskSize.ValueInt64()
	}
	if !data.RootDiskIoBus.IsNull() {
		params["root_disk_io_bus"] = data.RootDiskIoBus.ValueString()
	}
	if !data.ImageOs.IsNull() {
		params["image_os"] = data.ImageOs.ValueString()
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

	var id interface{}
	var err error
	id = data.ID.ValueString()

	_, _ = r.client.Call("virt.instance.stop", data.ID.ValueString())
	time.Sleep(2 * time.Second)

	_, err = r.client.CallWithJob("virt.instance.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete virt_instance: %s", err))
		return
	}
}
