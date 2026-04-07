package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NvmetPortSubsysResource struct {
	client *client.Client
}

type NvmetPortSubsysResourceModel struct {
	ID types.String `tfsdk:"id"`
	PortId types.Int64 `tfsdk:"port_id"`
	SubsysId types.Int64 `tfsdk:"subsys_id"`
}

func NewNvmetPortSubsysResource() resource.Resource {
	return &NvmetPortSubsysResource{}
}

func (r *NvmetPortSubsysResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_port_subsys"
}

func (r *NvmetPortSubsysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NvmetPortSubsysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an association between a `port` and a subsystem (`subsys`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"port_id": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the NVMe-oF port to associate.",
			},
			"subsys_id": schema.Int64Attribute{
				Required: true,
				Optional: false,
				Description: "ID of the NVMe-oF subsystem to make accessible.",
			},
		},
	}
}

func (r *NvmetPortSubsysResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NvmetPortSubsysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NvmetPortSubsysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.PortId.IsNull() && !data.PortId.IsUnknown() {
		params["port_id"] = data.PortId.ValueInt64()
	}
	if !data.SubsysId.IsNull() && !data.SubsysId.IsUnknown() {
		params["subsys_id"] = data.SubsysId.ValueInt64()
	}

	result, err := r.client.Call("nvmet.port_subsys.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create nvmet_port_subsys: %s", err))
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
	result, err = r.client.Call("nvmet.port_subsys.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back nvmet_port_subsys: %s", err))
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
		if v, ok := resultMap["port_id"]; ok {
			switch val := v.(type) {
			case float64:
				data.PortId = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.PortId = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["subsys_id"]; ok {
			switch val := v.(type) {
			case float64:
				data.SubsysId = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.SubsysId = types.Int64Value(int64(fv)) }
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetPortSubsysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NvmetPortSubsysResourceModel
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

	result, err := r.client.Call("nvmet.port_subsys.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read nvmet_port_subsys: %s", err))
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
		if v, ok := resultMap["port_id"]; ok {
			switch val := v.(type) {
			case float64:
				data.PortId = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.PortId = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["subsys_id"]; ok {
			switch val := v.(type) {
			case float64:
				data.SubsysId = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.SubsysId = types.Int64Value(int64(fv)) }
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetPortSubsysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NvmetPortSubsysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NvmetPortSubsysResourceModel
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
	if !data.PortId.IsNull() && !data.PortId.IsUnknown() {
		params["port_id"] = data.PortId.ValueInt64()
	}
	if !data.SubsysId.IsNull() && !data.SubsysId.IsUnknown() {
		params["subsys_id"] = data.SubsysId.ValueInt64()
	}

	_, err = r.client.Call("nvmet.port_subsys.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update nvmet_port_subsys: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("nvmet.port_subsys.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back nvmet_port_subsys: %s", readErr))
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
		if v, ok := resultMap["port_id"]; ok {
			switch val := v.(type) {
			case float64:
				data.PortId = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.PortId = types.Int64Value(int64(fv)) }
				}
			}
		}
		if v, ok := resultMap["subsys_id"]; ok {
			switch val := v.(type) {
			case float64:
				data.SubsysId = types.Int64Value(int64(val))
			case map[string]interface{}:
				if parsed, ok := val["parsed"]; ok && parsed != nil {
					if fv, ok := parsed.(float64); ok { data.SubsysId = types.Int64Value(int64(fv)) }
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NvmetPortSubsysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NvmetPortSubsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("nvmet.port_subsys.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete nvmet_port_subsys: %s", err))
		return
	}
}
