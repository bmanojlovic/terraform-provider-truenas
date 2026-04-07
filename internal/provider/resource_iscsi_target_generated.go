package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"encoding/json"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Alias.IsNull() && !data.Alias.IsUnknown() {
		params["alias"] = data.Alias.ValueString()
	}
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
		params["mode"] = data.Mode.ValueString()
	}
	if !data.Groups.IsNull() && !data.Groups.IsUnknown() {
		var groupsList []string
		data.Groups.ElementsAs(ctx, &groupsList, false)
		var groupsObjs []map[string]interface{}
		for _, jsonStr := range groupsList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse groups item: %s", err))
				return
			}
			groupsObjs = append(groupsObjs, obj)
		}
		params["groups"] = groupsObjs
	}
	if !data.AuthNetworks.IsNull() && !data.AuthNetworks.IsUnknown() {
		var auth_networksList []string
		data.AuthNetworks.ElementsAs(ctx, &auth_networksList, false)
		params["auth_networks"] = auth_networksList
	}
	if !data.IscsiParameters.IsNull() && !data.IscsiParameters.IsUnknown() {
		var iscsi_parametersObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.IscsiParameters.ValueString()), &iscsi_parametersObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse iscsi_parameters: %s", err))
			return
		}
		params["iscsi_parameters"] = iscsi_parametersObj
	}

	result, err := r.client.Call("iscsi.target.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_target: %s", err))
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
	result, err = r.client.Call("iscsi.target.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back iscsi_target: %s", err))
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
		if v, ok := resultMap["name"]; ok {
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
		if v, ok := resultMap["alias"]; ok {
			if v == nil {
				data.Alias = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Alias = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Alias = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Alias = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["iscsi_parameters"]; ok {
			if v == nil {
				data.IscsiParameters = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.IscsiParameters = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", v))
				}
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	result, err := r.client.Call("iscsi.target.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_target: %s", err))
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
		if v, ok := resultMap["name"]; ok {
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
		if v, ok := resultMap["alias"]; ok {
			if v == nil {
				data.Alias = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Alias = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Alias = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Alias = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["iscsi_parameters"]; ok {
			if v == nil {
				data.IscsiParameters = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.IscsiParameters = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", v))
				}
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

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Alias.IsNull() && !data.Alias.IsUnknown() {
		params["alias"] = data.Alias.ValueString()
	}
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
		params["mode"] = data.Mode.ValueString()
	}
	if !data.Groups.IsNull() && !data.Groups.IsUnknown() {
		var groupsList []string
		data.Groups.ElementsAs(ctx, &groupsList, false)
		var groupsObjs []map[string]interface{}
		for _, jsonStr := range groupsList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse groups item: %s", err))
				return
			}
			groupsObjs = append(groupsObjs, obj)
		}
		params["groups"] = groupsObjs
	}
	if !data.AuthNetworks.IsNull() && !data.AuthNetworks.IsUnknown() {
		var auth_networksList []string
		data.AuthNetworks.ElementsAs(ctx, &auth_networksList, false)
		params["auth_networks"] = auth_networksList
	}
	if !data.IscsiParameters.IsNull() && !data.IscsiParameters.IsUnknown() {
		var iscsi_parametersObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.IscsiParameters.ValueString()), &iscsi_parametersObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse iscsi_parameters: %s", err))
			return
		}
		params["iscsi_parameters"] = iscsi_parametersObj
	}

	_, err = r.client.Call("iscsi.target.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_target: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("iscsi.target.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back iscsi_target: %s", readErr))
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
		if v, ok := resultMap["name"]; ok {
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
		if v, ok := resultMap["alias"]; ok {
			if v == nil {
				data.Alias = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Alias = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Alias = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Alias = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["iscsi_parameters"]; ok {
			if v == nil {
				data.IscsiParameters = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.IscsiParameters = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.IscsiParameters = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiTargetResourceModel
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
	id = []interface{}{id, map[string]interface{}{}}

	_, err = r.client.Call("iscsi.target.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_target: %s", err))
		return
	}
}
