package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type IscsiPortalResource struct {
	client *client.Client
}

type IscsiPortalResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Listen  types.List   `tfsdk:"listen"`
	Comment types.String `tfsdk:"comment"`
}

func NewIscsiPortalResource() resource.Resource {
	return &IscsiPortalResource{}
}

func (r *IscsiPortalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_portal"
}

func (r *IscsiPortalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IscsiPortalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new iSCSI Portal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"listen": schema.ListAttribute{
				Required:    true,
				Optional:    false,
				ElementType: types.StringType,
				Description: "Array of IP addresses for the portal to listen on.",
			},
			"comment": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Optional comment describing the portal.",
			},
		},
	}
}

func (r *IscsiPortalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IscsiPortalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IscsiPortalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Listen.IsNull() && !data.Listen.IsUnknown() {
		var listenList []string
		data.Listen.ElementsAs(ctx, &listenList, false)
		var listenObjs []map[string]interface{}
		for _, jsonStr := range listenList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse listen item: %s", err))
				return
			}
			listenObjs = append(listenObjs, obj)
		}
		params["listen"] = listenObjs
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}

	result, err := r.client.Call("iscsi.portal.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create iscsi_portal: %s", err))
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
	result, err = r.client.Call("iscsi.portal.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back iscsi_portal: %s", err))
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
	if v, ok := resultMap["listen"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.Listen, _ = types.ListValue(types.StringType, strVals)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiPortalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IscsiPortalResourceModel
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

	result, err := r.client.Call("iscsi.portal.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read iscsi_portal: %s", err))
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
	if v, ok := resultMap["listen"]; ok {
		if arr, ok := v.([]interface{}); ok {
			strVals := make([]attr.Value, len(arr))
			for i, item := range arr {
				strVals[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
			data.Listen, _ = types.ListValue(types.StringType, strVals)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiPortalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IscsiPortalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IscsiPortalResourceModel
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
	if !data.Listen.IsNull() && !data.Listen.IsUnknown() {
		var listenList []string
		data.Listen.ElementsAs(ctx, &listenList, false)
		var listenObjs []map[string]interface{}
		for _, jsonStr := range listenList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse listen item: %s", err))
				return
			}
			listenObjs = append(listenObjs, obj)
		}
		params["listen"] = listenObjs
	}
	if !data.Comment.IsNull() && !data.Comment.IsUnknown() {
		params["comment"] = data.Comment.ValueString()
	}

	_, err = r.client.Call("iscsi.portal.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update iscsi_portal: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IscsiPortalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IscsiPortalResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}

	_, err = r.client.Call("iscsi.portal.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete iscsi_portal: %s", err))
		return
	}
}
