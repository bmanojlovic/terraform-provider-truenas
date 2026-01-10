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

type VmwareResource struct {
	client *client.Client
}

type VmwareResourceModel struct {
	ID types.String `tfsdk:"id"`
	Datastore types.String `tfsdk:"datastore"`
	Filesystem types.String `tfsdk:"filesystem"`
	Hostname types.String `tfsdk:"hostname"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func NewVmwareResource() resource.Resource {
	return &VmwareResource{}
}

func (r *VmwareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vmware"
}

func (r *VmwareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *VmwareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create VMWare snapshot.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"datastore": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Valid datastore name which exists on the VMWare host.",
			},
			"filesystem": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "ZFS filesystem or dataset to use for VMware storage.",
			},
			"hostname": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Valid IP address / hostname of a VMWare host. When clustering, this is the vCenter server for the cl",
			},
			"username": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Credentials used to authorize access to the VMWare host.",
			},
			"password": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Password for VMware host authentication.",
			},
		},
	}
}

func (r *VmwareResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VmwareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VmwareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Datastore.IsNull() {
		params["datastore"] = data.Datastore.ValueString()
	}
	if !data.Filesystem.IsNull() {
		params["filesystem"] = data.Filesystem.ValueString()
	}
	if !data.Hostname.IsNull() {
		params["hostname"] = data.Hostname.ValueString()
	}
	if !data.Username.IsNull() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}

	result, err := r.client.Call("vmware.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create vmware: %s", err))
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

func (r *VmwareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VmwareResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("vmware.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read vmware: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["datastore"]; ok && v != nil {
			data.Datastore = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["filesystem"]; ok && v != nil {
			data.Filesystem = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["hostname"]; ok && v != nil {
			data.Hostname = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["username"]; ok && v != nil {
			data.Username = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["password"]; ok && v != nil {
			data.Password = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmwareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmwareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VmwareResourceModel
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
	if !data.Datastore.IsNull() {
		params["datastore"] = data.Datastore.ValueString()
	}
	if !data.Filesystem.IsNull() {
		params["filesystem"] = data.Filesystem.ValueString()
	}
	if !data.Hostname.IsNull() {
		params["hostname"] = data.Hostname.ValueString()
	}
	if !data.Username.IsNull() {
		params["username"] = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		params["password"] = data.Password.ValueString()
	}

	_, err = r.client.Call("vmware.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update vmware: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmwareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmwareResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("vmware.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete vmware: %s", err))
		return
	}
}
