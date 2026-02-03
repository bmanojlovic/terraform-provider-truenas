package provider

import (
	"context"
	"fmt"
	"strings"
"strconv"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FilesystemAcltemplateResource struct {
	client *client.Client
}

type FilesystemAcltemplateResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Acltype types.String `tfsdk:"acltype"`
	Acl types.List `tfsdk:"acl"`
	Comment types.String `tfsdk:"comment"`
}

func NewFilesystemAcltemplateResource() resource.Resource {
	return &FilesystemAcltemplateResource{}
}

func (r *FilesystemAcltemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_acltemplate"
}

func (r *FilesystemAcltemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FilesystemAcltemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new filesystem ACL template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Human-readable name for the ACL template.",
			},
			"acltype": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "ACL type this template provides.",
			},
			"acl": schema.ListAttribute{
				Required: true,
				Optional: false,
				ElementType: types.StringType,
				Description: "Array of Access Control Entries defined by this template.",
			},
			"comment": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optional descriptive comment about the template's purpose.",
			},
		},
	}
}

func (r *FilesystemAcltemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FilesystemAcltemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FilesystemAcltemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Acltype.IsNull() {
		params["acltype"] = data.Acltype.ValueString()
	}
	if !data.Acl.IsNull() {
		var aclList []string
		data.Acl.ElementsAs(ctx, &aclList, false)
		params["acl"] = aclList
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	result, err := r.client.Call("filesystem.acltemplate.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create filesystem_acltemplate: %s", err))
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemAcltemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FilesystemAcltemplateResourceModel
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

	result, err := r.client.Call("filesystem.acltemplate.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read filesystem_acltemplate: %s", err))
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
		if v, ok := resultMap["acltype"]; ok && v != nil {
			switch val := v.(type) {
			case string:
				data.Acltype = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Acltype = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Acltype = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["acl"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.Acl, _ = types.ListValue(types.StringType, strVals)
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemAcltemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FilesystemAcltemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FilesystemAcltemplateResourceModel
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
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Acltype.IsNull() {
		params["acltype"] = data.Acltype.ValueString()
	}
	if !data.Acl.IsNull() {
		var aclList []string
		data.Acl.ElementsAs(ctx, &aclList, false)
		params["acl"] = aclList
	}
	if !data.Comment.IsNull() {
		params["comment"] = data.Comment.ValueString()
	}

	_, err = r.client.Call("filesystem.acltemplate.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update filesystem_acltemplate: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilesystemAcltemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FilesystemAcltemplateResourceModel
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

	_, err = r.client.Call("filesystem.acltemplate.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete filesystem_acltemplate: %s", err))
		return
	}
}
