package provider

import (
	"context"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type PoolResource struct {
	client *client.Client
}

type PoolResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Encryption types.Bool `tfsdk:"encryption"`
	DedupTableQuota types.String `tfsdk:"dedup_table_quota"`
	DedupTableQuotaValue types.Int64 `tfsdk:"dedup_table_quota_value"`
	Deduplication types.String `tfsdk:"deduplication"`
	Checksum types.String `tfsdk:"checksum"`
	EncryptionOptions types.String `tfsdk:"encryption_options"`
	Topology types.String `tfsdk:"topology"`
	AllowDuplicateSerials types.Bool `tfsdk:"allow_duplicate_serials"`
	Autotrim types.String `tfsdk:"autotrim"`
}

func NewPoolResource() resource.Resource {
	return &PoolResource{}
}

func (r *PoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool"
}

func (r *PoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new ZFS Pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Name for the new storage pool.",
			},
			"encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "If set, create a ZFS encrypted root dataset for this pool.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"dedup_table_quota": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "How to manage the deduplication table quota allocation.",
			},
			"dedup_table_quota_value": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Custom quota value in bytes when `dedup_table_quota` is set to CUSTOM.",
			},
			"deduplication": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Make sure no block of data is duplicated in the pool. If set to `VERIFY` and two blocks have similar",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"checksum": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Checksum algorithm to use for data integrity verification.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"encryption_options": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Specify configuration for encryption of root dataset.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"topology": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Updated topology configuration for adding new vdevs to the pool.",
			},
			"allow_duplicate_serials": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to allow disks with duplicate serial numbers in the pool.",
			},
			"autotrim": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to enable automatic TRIM operations on the pool.",
			},
		},
	}
}

func (r *PoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Encryption.IsNull() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.DedupTableQuota.IsNull() {
		params["dedup_table_quota"] = data.DedupTableQuota.ValueString()
	}
	if !data.DedupTableQuotaValue.IsNull() {
		params["dedup_table_quota_value"] = data.DedupTableQuotaValue.ValueInt64()
	}
	if !data.Deduplication.IsNull() {
		params["deduplication"] = data.Deduplication.ValueString()
	}
	if !data.Checksum.IsNull() {
		params["checksum"] = data.Checksum.ValueString()
	}
	if !data.EncryptionOptions.IsNull() {
		var encryption_optionsObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.EncryptionOptions.ValueString()), &encryption_optionsObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse encryption_options: %s", err))
			return
		}
		params["encryption_options"] = encryption_optionsObj
	}
	if !data.Topology.IsNull() {
		var topologyObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Topology.ValueString()), &topologyObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse topology: %s", err))
			return
		}
		params["topology"] = topologyObj
	}
	if !data.AllowDuplicateSerials.IsNull() {
		params["allow_duplicate_serials"] = data.AllowDuplicateSerials.ValueBool()
	}
	if !data.Autotrim.IsNull() {
		params["autotrim"] = data.Autotrim.ValueString()
	}

	result, err := r.client.CallWithJob("pool.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create pool: %s", err))
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

func (r *PoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	result, err := r.client.Call("pool.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Encryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["dedup_table_quota"]; ok && v != nil {
			data.DedupTableQuota = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["dedup_table_quota_value"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.DedupTableQuotaValue = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["deduplication"]; ok && v != nil {
			data.Deduplication = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["checksum"]; ok && v != nil {
			data.Checksum = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_options"]; ok && v != nil {
			data.EncryptionOptions = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["topology"]; ok && v != nil {
			data.Topology = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["allow_duplicate_serials"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.AllowDuplicateSerials = types.BoolValue(bv) }
		}
		if v, ok := resultMap["autotrim"]; ok && v != nil {
			data.Autotrim = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(state.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	params := map[string]interface{}{}
	if !data.DedupTableQuota.IsNull() {
		params["dedup_table_quota"] = data.DedupTableQuota.ValueString()
	}
	if !data.DedupTableQuotaValue.IsNull() {
		params["dedup_table_quota_value"] = data.DedupTableQuotaValue.ValueInt64()
	}
	if !data.Topology.IsNull() {
		var topologyObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Topology.ValueString()), &topologyObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse topology: %s", err))
			return
		}
		params["topology"] = topologyObj
	}
	if !data.AllowDuplicateSerials.IsNull() {
		params["allow_duplicate_serials"] = data.AllowDuplicateSerials.ValueBool()
	}
	if !data.Autotrim.IsNull() {
		params["autotrim"] = data.Autotrim.ValueString()
	}

	_, err = r.client.CallWithJob("pool.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update pool: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {{
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}}

	_, err = r.client.Call("pool.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool: %s", err))
		return
	}
}
