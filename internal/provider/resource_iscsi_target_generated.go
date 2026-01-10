package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type IscsiTargetResource struct {
	client *client.Client
}

type IscsiTargetResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Alias types.String `tfsdk:"alias"`
	Mode types.String `tfsdk:"mode"`
	Groups types.List `tfsdk:"groups"`
	AuthNetworks types.List `tfsdk:"auth_networks"`
	IscsiParameters types.String `tfsdk:"iscsi_parameters"`
}

func NewIscsiTargetResource() resource.Resource {
	return &IscsiTargetResource{}
}

func (r *IscsiTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_target"
}

func (r *IscsiTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Target.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name of the iSCSI target (maximum 120 characters).",
			},
			"alias": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional alias name for the iSCSI target.",
			},
			"mode": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Protocol mode for the target.  * `ISCSI`: iSCSI protocol only * `FC`: Fibre Channel protocol only * ",
			},
			"groups": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of portal-initiator group associations for this target.",
			},
			"auth_networks": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of network addresses allowed to access this target.",
			},
			"iscsi_parameters": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional iSCSI-specific parameters for this target.",
			},
		},
	}
}

func (r *IscsiTargetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiTargetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Alias.IsNull() {
		params["alias"] = data.Alias.ValueString()
	}
	if !data.Mode.IsNull() {
		params["mode"] = data.Mode.ValueString()
	}
	if !data.Groups.IsNull() {
		var groupsList []string
		data.Groups.ElementsAs(ctx, &groupsList, false)
		params["groups"] = groupsList
	}
	if !data.AuthNetworks.IsNull() {
		var auth_networksList []string
		data.AuthNetworks.ElementsAs(ctx, &auth_networksList, false)
		params["auth_networks"] = auth_networksList
	}
	if !data.IscsiParameters.IsNull() {
		params["iscsi_parameters"] = data.IscsiParameters.ValueString()
	}

	result, err := r.client.Call("iscsi.target.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_target: %s", err))
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

func (r *IscsiTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("iscsi.target.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_target: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["alias"]; ok && v != nil {
			data.Alias = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["mode"]; ok && v != nil {
			data.Mode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["groups"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Groups, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["auth_networks"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.AuthNetworks, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["iscsi_parameters"]; ok && v != nil {
			data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiTargetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiTargetResourceModel
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
	if !data.Alias.IsNull() {
		params["alias"] = data.Alias.ValueString()
	}
	if !data.Mode.IsNull() {
		params["mode"] = data.Mode.ValueString()
	}
	if !data.Groups.IsNull() {
		var groupsList []string
		data.Groups.ElementsAs(ctx, &groupsList, false)
		params["groups"] = groupsList
	}
	if !data.AuthNetworks.IsNull() {
		var auth_networksList []string
		data.AuthNetworks.ElementsAs(ctx, &auth_networksList, false)
		params["auth_networks"] = auth_networksList
	}
	if !data.IscsiParameters.IsNull() {
		params["iscsi_parameters"] = data.IscsiParameters.ValueString()
	}

	_, err = r.client.Call("iscsi.target.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_target: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("iscsi.target.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_target: %s", err))
		return
	}
}
