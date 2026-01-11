package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

type PoolDatasetResource struct {
	client *client.Client
}

type PoolDatasetResourceModel struct {
	ID types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Comments types.String `tfsdk:"comments"`
	Sync types.String `tfsdk:"sync"`
	Snapdev types.String `tfsdk:"snapdev"`
	Compression types.String `tfsdk:"compression"`
	Exec types.String `tfsdk:"exec"`
	Managedby types.String `tfsdk:"managedby"`
	QuotaWarning types.Int64 `tfsdk:"quota_warning"`
	QuotaCritical types.Int64 `tfsdk:"quota_critical"`
	RefquotaWarning types.Int64 `tfsdk:"refquota_warning"`
	RefquotaCritical types.Int64 `tfsdk:"refquota_critical"`
	Reservation types.Int64 `tfsdk:"reservation"`
	Refreservation types.Int64 `tfsdk:"refreservation"`
	SpecialSmallBlockSize types.Int64 `tfsdk:"special_small_block_size"`
	Copies types.Int64 `tfsdk:"copies"`
	Snapdir types.String `tfsdk:"snapdir"`
	Deduplication types.String `tfsdk:"deduplication"`
	Checksum types.String `tfsdk:"checksum"`
	Readonly types.String `tfsdk:"readonly"`
	ShareType types.String `tfsdk:"share_type"`
	EncryptionOptions types.String `tfsdk:"encryption_options"`
	Encryption types.Bool `tfsdk:"encryption"`
	InheritEncryption types.Bool `tfsdk:"inherit_encryption"`
	UserProperties types.List `tfsdk:"user_properties"`
	CreateAncestors types.Bool `tfsdk:"create_ancestors"`
	Type types.String `tfsdk:"type"`
	Aclmode types.String `tfsdk:"aclmode"`
	Acltype types.String `tfsdk:"acltype"`
	Atime types.String `tfsdk:"atime"`
	Casesensitivity types.String `tfsdk:"casesensitivity"`
	Quota types.Int64 `tfsdk:"quota"`
	Refquota types.Int64 `tfsdk:"refquota"`
	Recordsize types.String `tfsdk:"recordsize"`
	ForceSize types.Bool `tfsdk:"force_size"`
	Sparse types.Bool `tfsdk:"sparse"`
	Volsize types.Int64 `tfsdk:"volsize"`
	Volblocksize types.String `tfsdk:"volblocksize"`
	UserPropertiesUpdate types.List `tfsdk:"user_properties_update"`
}

func NewPoolDatasetResource() resource.Resource {
	return &PoolDatasetResource{}
}

func (r *PoolDatasetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_dataset"
}

func (r *PoolDatasetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PoolDatasetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a dataset/zvol.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "The name of the dataset to create.",
			},
			"comments": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Comments or description for the dataset.",
			},
			"sync": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Synchronous write behavior for the dataset.",
			},
			"snapdev": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Controls visibility of volume snapshots under /dev/zvol/.",
			},
			"compression": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Compression algorithm to use for the dataset. Higher numbered variants provide better compression   ",
			},
			"exec": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Whether files in this dataset can be executed.",
			},
			"managedby": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Identifies which service or system manages this dataset.",
			},
			"quota_warning": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Percentage of dataset quota at which to issue a warning. 0-100 or 'INHERIT'.",
			},
			"quota_critical": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Percentage of dataset quota at which to issue a critical alert. 0-100 or 'INHERIT'.",
			},
			"refquota_warning": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Percentage of reference quota at which to issue a warning. 0-100 or 'INHERIT'.",
			},
			"refquota_critical": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Percentage of reference quota at which to issue a critical alert. 0-100 or 'INHERIT'.",
			},
			"reservation": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Minimum disk space guaranteed to this dataset and its children in bytes.",
			},
			"refreservation": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Minimum disk space guaranteed to this dataset itself in bytes.",
			},
			"special_small_block_size": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Size threshold below which blocks are stored on special vdevs.",
			},
			"copies": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Number of copies of data blocks to maintain for redundancy.",
			},
			"snapdir": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Controls visibility of the `.zfs/snapshot` directory. 'DISABLED' hides snapshots, 'VISIBLE' shows th",
			},
			"deduplication": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Deduplication setting. 'ON' enables dedup, 'VERIFY' enables with checksum verification, 'OFF' disabl",
			},
			"checksum": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Checksum algorithm to verify data integrity. Higher security algorithms like SHA256 provide better  ",
			},
			"readonly": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Whether the dataset is read-only.",
			},
			"share_type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Optimization type for the dataset based on its intended use.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"encryption_options": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Configuration for encryption of dataset for `name` pool.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Create a ZFS encrypted root dataset for `name` pool. There is 1 case where ZFS encryption is not all",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"inherit_encryption": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to inherit encryption settings from the parent dataset.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"user_properties": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Custom user-defined properties to set on the dataset.",
			},
			"create_ancestors": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to create any missing parent datasets.",
			},
			"type": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Type of dataset to create - volume (zvol).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"aclmode": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "How Access Control Lists are handled when chmod is used.",
			},
			"acltype": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The type of Access Control List system to use.",
			},
			"atime": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Whether file access times are updated when files are accessed.",
			},
			"casesensitivity": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "File name case sensitivity setting.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"quota": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum disk space this dataset and its children can consume in bytes.",
			},
			"refquota": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Maximum disk space this dataset itself can consume in bytes.",
			},
			"recordsize": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "The suggested block size for files in this filesystem dataset.",
			},
			"force_size": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Force creation even if the size is not optimal.",
			},
			"sparse": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether to use sparse (thin) provisioning for the volume.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"volsize": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "The volume size in bytes; supposed to be a multiple of the block size.",
			},
			"volblocksize": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Defaults to `128K` if the parent pool is a DRAID pool or `16K` otherwise.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_properties_update": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of user property updates to apply to the dataset.",
			},
		},
	}
}

func (r *PoolDatasetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolDatasetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolDatasetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Name.IsNull() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Comments.IsNull() {
		params["comments"] = data.Comments.ValueString()
	}
	if !data.Sync.IsNull() {
		params["sync"] = data.Sync.ValueString()
	}
	if !data.Snapdev.IsNull() {
		params["snapdev"] = data.Snapdev.ValueString()
	}
	if !data.Compression.IsNull() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.Exec.IsNull() {
		params["exec"] = data.Exec.ValueString()
	}
	if !data.Managedby.IsNull() {
		params["managedby"] = data.Managedby.ValueString()
	}
	if !data.QuotaWarning.IsNull() {
		params["quota_warning"] = data.QuotaWarning.ValueInt64()
	}
	if !data.QuotaCritical.IsNull() {
		params["quota_critical"] = data.QuotaCritical.ValueInt64()
	}
	if !data.RefquotaWarning.IsNull() {
		params["refquota_warning"] = data.RefquotaWarning.ValueInt64()
	}
	if !data.RefquotaCritical.IsNull() {
		params["refquota_critical"] = data.RefquotaCritical.ValueInt64()
	}
	if !data.Reservation.IsNull() {
		params["reservation"] = data.Reservation.ValueInt64()
	}
	if !data.Refreservation.IsNull() {
		params["refreservation"] = data.Refreservation.ValueInt64()
	}
	if !data.SpecialSmallBlockSize.IsNull() {
		params["special_small_block_size"] = data.SpecialSmallBlockSize.ValueInt64()
	}
	if !data.Copies.IsNull() {
		params["copies"] = data.Copies.ValueInt64()
	}
	if !data.Snapdir.IsNull() {
		params["snapdir"] = data.Snapdir.ValueString()
	}
	if !data.Deduplication.IsNull() {
		params["deduplication"] = data.Deduplication.ValueString()
	}
	if !data.Checksum.IsNull() {
		params["checksum"] = data.Checksum.ValueString()
	}
	if !data.Readonly.IsNull() {
		params["readonly"] = data.Readonly.ValueString()
	}
	if !data.ShareType.IsNull() {
		params["share_type"] = data.ShareType.ValueString()
	}
	if !data.EncryptionOptions.IsNull() {
		var encryption_optionsObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.EncryptionOptions.ValueString()), &encryption_optionsObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse encryption_options: %s", err))
			return
		}
		params["encryption_options"] = encryption_optionsObj
	}
	if !data.Encryption.IsNull() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.InheritEncryption.IsNull() {
		params["inherit_encryption"] = data.InheritEncryption.ValueBool()
	}
	if !data.UserProperties.IsNull() {
		var user_propertiesList []string
		data.UserProperties.ElementsAs(ctx, &user_propertiesList, false)
		params["user_properties"] = user_propertiesList
	}
	if !data.CreateAncestors.IsNull() {
		params["create_ancestors"] = data.CreateAncestors.ValueBool()
	}
	if !data.Type.IsNull() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Aclmode.IsNull() {
		params["aclmode"] = data.Aclmode.ValueString()
	}
	if !data.Acltype.IsNull() {
		params["acltype"] = data.Acltype.ValueString()
	}
	if !data.Atime.IsNull() {
		params["atime"] = data.Atime.ValueString()
	}
	if !data.Casesensitivity.IsNull() {
		params["casesensitivity"] = data.Casesensitivity.ValueString()
	}
	if !data.Quota.IsNull() {
		params["quota"] = data.Quota.ValueInt64()
	}
	if !data.Refquota.IsNull() {
		params["refquota"] = data.Refquota.ValueInt64()
	}
	if !data.Recordsize.IsNull() {
		params["recordsize"] = data.Recordsize.ValueString()
	}
	if !data.ForceSize.IsNull() {
		params["force_size"] = data.ForceSize.ValueBool()
	}
	if !data.Sparse.IsNull() {
		params["sparse"] = data.Sparse.ValueBool()
	}
	if !data.Volsize.IsNull() {
		params["volsize"] = data.Volsize.ValueInt64()
	}
	if !data.Volblocksize.IsNull() {
		params["volblocksize"] = data.Volblocksize.ValueString()
	}
	if !data.UserPropertiesUpdate.IsNull() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		params["user_properties_update"] = user_properties_updateList
	}

	result, err := r.client.Call("pool.dataset.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create pool_dataset: %s", err))
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

func (r *PoolDatasetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolDatasetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	result, err := r.client.Call("pool.dataset.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool_dataset: %s", err))
		return
	}

	// Map result back to state
	if resultMap, ok := result.(map[string]interface{}); ok {
		if v, ok := resultMap["name"]; ok && v != nil {
			data.Name = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["comments"]; ok && v != nil {
			data.Comments = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["sync"]; ok && v != nil {
			data.Sync = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["snapdev"]; ok && v != nil {
			data.Snapdev = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["compression"]; ok && v != nil {
			data.Compression = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["exec"]; ok && v != nil {
			data.Exec = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["managedby"]; ok && v != nil {
			data.Managedby = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota_warning"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.QuotaWarning = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["quota_critical"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.QuotaCritical = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["refquota_warning"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.RefquotaWarning = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["refquota_critical"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.RefquotaCritical = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["reservation"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Reservation = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["refreservation"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Refreservation = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["special_small_block_size"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.SpecialSmallBlockSize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["copies"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Copies = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["snapdir"]; ok && v != nil {
			data.Snapdir = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["deduplication"]; ok && v != nil {
			data.Deduplication = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["checksum"]; ok && v != nil {
			data.Checksum = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["readonly"]; ok && v != nil {
			data.Readonly = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["share_type"]; ok && v != nil {
			data.ShareType = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption_options"]; ok && v != nil {
			data.EncryptionOptions = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Encryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["inherit_encryption"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.InheritEncryption = types.BoolValue(bv) }
		}
		if v, ok := resultMap["user_properties"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.UserProperties, _ = types.ListValue(types.StringType, strVals)
			}
		}
		if v, ok := resultMap["create_ancestors"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.CreateAncestors = types.BoolValue(bv) }
		}
		if v, ok := resultMap["type"]; ok && v != nil {
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["aclmode"]; ok && v != nil {
			data.Aclmode = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["acltype"]; ok && v != nil {
			data.Acltype = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["atime"]; ok && v != nil {
			data.Atime = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["casesensitivity"]; ok && v != nil {
			data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["quota"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Quota = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["refquota"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Refquota = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["recordsize"]; ok && v != nil {
			data.Recordsize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["force_size"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.ForceSize = types.BoolValue(bv) }
		}
		if v, ok := resultMap["sparse"]; ok && v != nil {
			if bv, ok := v.(bool); ok { data.Sparse = types.BoolValue(bv) }
		}
		if v, ok := resultMap["volsize"]; ok && v != nil {
			if fv, ok := v.(float64); ok { data.Volsize = types.Int64Value(int64(fv)) }
		}
		if v, ok := resultMap["volblocksize"]; ok && v != nil {
			data.Volblocksize = types.StringValue(fmt.Sprintf("%v", v))
		}
		if v, ok := resultMap["user_properties_update"]; ok && v != nil {
			if arr, ok := v.([]interface{}); ok {
				strVals := make([]attr.Value, len(arr))
				for i, item := range arr { strVals[i] = types.StringValue(fmt.Sprintf("%v", item)) }
				data.UserPropertiesUpdate, _ = types.ListValue(types.StringType, strVals)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolDatasetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PoolDatasetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PoolDatasetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = state.ID.ValueString()

	params := map[string]interface{}{}
	if !data.Comments.IsNull() {
		params["comments"] = data.Comments.ValueString()
	}
	if !data.Sync.IsNull() {
		params["sync"] = data.Sync.ValueString()
	}
	if !data.Snapdev.IsNull() {
		params["snapdev"] = data.Snapdev.ValueString()
	}
	if !data.Compression.IsNull() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.Exec.IsNull() {
		params["exec"] = data.Exec.ValueString()
	}
	if !data.Managedby.IsNull() {
		params["managedby"] = data.Managedby.ValueString()
	}
	if !data.QuotaWarning.IsNull() {
		params["quota_warning"] = data.QuotaWarning.ValueInt64()
	}
	if !data.QuotaCritical.IsNull() {
		params["quota_critical"] = data.QuotaCritical.ValueInt64()
	}
	if !data.RefquotaWarning.IsNull() {
		params["refquota_warning"] = data.RefquotaWarning.ValueInt64()
	}
	if !data.RefquotaCritical.IsNull() {
		params["refquota_critical"] = data.RefquotaCritical.ValueInt64()
	}
	if !data.Reservation.IsNull() {
		params["reservation"] = data.Reservation.ValueInt64()
	}
	if !data.Refreservation.IsNull() {
		params["refreservation"] = data.Refreservation.ValueInt64()
	}
	if !data.SpecialSmallBlockSize.IsNull() {
		params["special_small_block_size"] = data.SpecialSmallBlockSize.ValueInt64()
	}
	if !data.Copies.IsNull() {
		params["copies"] = data.Copies.ValueInt64()
	}
	if !data.Snapdir.IsNull() {
		params["snapdir"] = data.Snapdir.ValueString()
	}
	if !data.Deduplication.IsNull() {
		params["deduplication"] = data.Deduplication.ValueString()
	}
	if !data.Checksum.IsNull() {
		params["checksum"] = data.Checksum.ValueString()
	}
	if !data.Readonly.IsNull() {
		params["readonly"] = data.Readonly.ValueString()
	}
	if !data.UserProperties.IsNull() {
		var user_propertiesList []string
		data.UserProperties.ElementsAs(ctx, &user_propertiesList, false)
		params["user_properties"] = user_propertiesList
	}
	if !data.CreateAncestors.IsNull() {
		params["create_ancestors"] = data.CreateAncestors.ValueBool()
	}
	if !data.ForceSize.IsNull() {
		params["force_size"] = data.ForceSize.ValueBool()
	}
	if !data.Volsize.IsNull() {
		params["volsize"] = data.Volsize.ValueInt64()
	}
	if !data.Aclmode.IsNull() {
		params["aclmode"] = data.Aclmode.ValueString()
	}
	if !data.Acltype.IsNull() {
		params["acltype"] = data.Acltype.ValueString()
	}
	if !data.Atime.IsNull() {
		params["atime"] = data.Atime.ValueString()
	}
	if !data.Quota.IsNull() {
		params["quota"] = data.Quota.ValueInt64()
	}
	if !data.Refquota.IsNull() {
		params["refquota"] = data.Refquota.ValueInt64()
	}
	if !data.Recordsize.IsNull() {
		params["recordsize"] = data.Recordsize.ValueString()
	}
	if !data.UserPropertiesUpdate.IsNull() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		params["user_properties_update"] = user_properties_updateList
	}

	_, err = r.client.Call("pool.dataset.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update pool_dataset: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolDatasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolDatasetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id interface{}
	var err error
	id = data.ID.ValueString()

	_, err = r.client.Call("pool.dataset.delete", id)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool_dataset: %s", err))
		return
	}
}
