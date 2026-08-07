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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type PoolDatasetResource struct {
	client *client.Client
}

type PoolDatasetResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Comments              types.String `tfsdk:"comments"`
	Sync                  types.String `tfsdk:"sync"`
	Snapdev               types.String `tfsdk:"snapdev"`
	Compression           types.String `tfsdk:"compression"`
	Exec                  types.String `tfsdk:"exec"`
	Managedby             types.String `tfsdk:"managedby"`
	QuotaWarning          types.Int64  `tfsdk:"quota_warning"`
	QuotaCritical         types.Int64  `tfsdk:"quota_critical"`
	RefquotaWarning       types.Int64  `tfsdk:"refquota_warning"`
	RefquotaCritical      types.Int64  `tfsdk:"refquota_critical"`
	Reservation           types.Int64  `tfsdk:"reservation"`
	Refreservation        types.Int64  `tfsdk:"refreservation"`
	SpecialSmallBlockSize types.Int64  `tfsdk:"special_small_block_size"`
	Copies                types.Int64  `tfsdk:"copies"`
	Snapdir               types.String `tfsdk:"snapdir"`
	Deduplication         types.String `tfsdk:"deduplication"`
	Checksum              types.String `tfsdk:"checksum"`
	Readonly              types.String `tfsdk:"readonly"`
	ShareType             types.String `tfsdk:"share_type"`
	EncryptionOptions     types.String `tfsdk:"encryption_options"`
	Encryption            types.Bool   `tfsdk:"encryption"`
	InheritEncryption     types.Bool   `tfsdk:"inherit_encryption"`
	UserProperties        types.List   `tfsdk:"user_properties"`
	CreateAncestors       types.Bool   `tfsdk:"create_ancestors"`
	Type                  types.String `tfsdk:"type"`
	Aclmode               types.String `tfsdk:"aclmode"`
	Acltype               types.String `tfsdk:"acltype"`
	Atime                 types.String `tfsdk:"atime"`
	Casesensitivity       types.String `tfsdk:"casesensitivity"`
	Quota                 types.Int64  `tfsdk:"quota"`
	Refquota              types.Int64  `tfsdk:"refquota"`
	Recordsize            types.String `tfsdk:"recordsize"`
	ForceSize             types.Bool   `tfsdk:"force_size"`
	Sparse                types.Bool   `tfsdk:"sparse"`
	Volsize               types.Int64  `tfsdk:"volsize"`
	Volblocksize          types.String `tfsdk:"volblocksize"`
	UserPropertiesUpdate  types.List   `tfsdk:"user_properties_update"`
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
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "The name of the dataset to create.",
			},
			"comments": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Comments or description for the dataset.",
			},
			"sync": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Synchronous write behavior for the dataset.",
			},
			"snapdev": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Controls visibility of volume snapshots under /dev/zvol/.",
			},
			"compression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Compression algorithm to use for the dataset. Higher numbered variants provide better compression   ",
			},
			"exec": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether files in this dataset can be executed.",
			},
			"managedby": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifies which service or system manages this dataset.",
			},
			"quota_warning": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Percentage of dataset quota at which to issue a warning. 0-100 or 'INHERIT'.",
			},
			"quota_critical": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Percentage of dataset quota at which to issue a critical alert. 0-100 or 'INHERIT'.",
			},
			"refquota_warning": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Percentage of reference quota at which to issue a warning. 0-100 or 'INHERIT'.",
			},
			"refquota_critical": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Percentage of reference quota at which to issue a critical alert. 0-100 or 'INHERIT'.",
			},
			"reservation": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum disk space guaranteed to this dataset and its children in bytes.",
			},
			"refreservation": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum disk space guaranteed to this dataset itself in bytes.",
			},
			"special_small_block_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Size threshold below which blocks are stored on special vdevs.",
			},
			"copies": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of copies of data blocks to maintain for redundancy.",
			},
			"snapdir": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Controls visibility of the `.zfs/snapshot` directory. 'DISABLED' hides snapshots, 'VISIBLE' shows th",
			},
			"deduplication": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Deduplication setting. 'ON' enables dedup, 'VERIFY' enables with checksum verification, 'OFF' disabl",
			},
			"checksum": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Checksum algorithm to verify data integrity. Higher security algorithms like SHA256 provide better  ",
			},
			"readonly": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the dataset is read-only.",
			},
			"share_type": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Optimization type for the dataset based on its intended use.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"encryption_options": schema.StringAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Configuration for encryption of dataset for `name` pool.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"encryption": schema.BoolAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Create a ZFS encrypted root dataset for `name` pool. There is 1 case where ZFS encryption is not all",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"inherit_encryption": schema.BoolAttribute{
				Required:      false,
				Optional:      true,
				Description:   "Whether to inherit encryption settings from the parent dataset.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"user_properties": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Custom user-defined properties to set on the dataset.",
			},
			"create_ancestors": schema.BoolAttribute{
				Required:    false,
				Optional:    true,
				Description: "Whether to create any missing parent datasets.",
			},
			"type": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Type of dataset to create - volume (zvol).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"aclmode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "How Access Control Lists are handled when chmod is used.",
			},
			"acltype": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of Access Control List system to use.",
			},
			"atime": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether file access times are updated when files are accessed.",
			},
			"casesensitivity": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "File name case sensitivity setting.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"quota": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum disk space this dataset and its children can consume in bytes.",
			},
			"refquota": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum disk space this dataset itself can consume in bytes.",
			},
			"recordsize": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The suggested block size for files in this filesystem dataset.",
			},
			"force_size": schema.BoolAttribute{
				Required:    false,
				Optional:    true,
				Description: "Force creation even if the size is not optimal.",
			},
			"sparse": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Whether to use sparse (thin) provisioning for the volume.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown(), boolplanmodifier.RequiresReplace()},
			},
			"volsize": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The volume size in bytes; supposed to be a multiple of the block size.",
			},
			"volblocksize": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Defaults to `128K` if the parent pool is a DRAID pool or `16K` otherwise.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"user_properties_update": schema.ListAttribute{
				Required:    false,
				Optional:    true,
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
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		params["name"] = data.Name.ValueString()
	}
	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		params["comments"] = data.Comments.ValueString()
	}
	if !data.Sync.IsNull() && !data.Sync.IsUnknown() {
		params["sync"] = data.Sync.ValueString()
	}
	if !data.Snapdev.IsNull() && !data.Snapdev.IsUnknown() {
		params["snapdev"] = data.Snapdev.ValueString()
	}
	if !data.Compression.IsNull() && !data.Compression.IsUnknown() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.Exec.IsNull() && !data.Exec.IsUnknown() {
		params["exec"] = data.Exec.ValueString()
	}
	if !data.Managedby.IsNull() && !data.Managedby.IsUnknown() {
		params["managedby"] = data.Managedby.ValueString()
	}
	if !data.QuotaWarning.IsNull() && !data.QuotaWarning.IsUnknown() {
		params["quota_warning"] = data.QuotaWarning.ValueInt64()
	}
	if !data.QuotaCritical.IsNull() && !data.QuotaCritical.IsUnknown() {
		params["quota_critical"] = data.QuotaCritical.ValueInt64()
	}
	if !data.RefquotaWarning.IsNull() && !data.RefquotaWarning.IsUnknown() {
		params["refquota_warning"] = data.RefquotaWarning.ValueInt64()
	}
	if !data.RefquotaCritical.IsNull() && !data.RefquotaCritical.IsUnknown() {
		params["refquota_critical"] = data.RefquotaCritical.ValueInt64()
	}
	if !data.Reservation.IsNull() && !data.Reservation.IsUnknown() {
		params["reservation"] = data.Reservation.ValueInt64()
	}
	if !data.Refreservation.IsNull() && !data.Refreservation.IsUnknown() {
		params["refreservation"] = data.Refreservation.ValueInt64()
	}
	if !data.SpecialSmallBlockSize.IsNull() && !data.SpecialSmallBlockSize.IsUnknown() {
		params["special_small_block_size"] = data.SpecialSmallBlockSize.ValueInt64()
	}
	if !data.Copies.IsNull() && !data.Copies.IsUnknown() {
		params["copies"] = data.Copies.ValueInt64()
	}
	if !data.Snapdir.IsNull() && !data.Snapdir.IsUnknown() {
		params["snapdir"] = data.Snapdir.ValueString()
	}
	if !data.Deduplication.IsNull() && !data.Deduplication.IsUnknown() {
		params["deduplication"] = data.Deduplication.ValueString()
	}
	if !data.Checksum.IsNull() && !data.Checksum.IsUnknown() {
		params["checksum"] = data.Checksum.ValueString()
	}
	if !data.Readonly.IsNull() && !data.Readonly.IsUnknown() {
		params["readonly"] = data.Readonly.ValueString()
	}
	if !data.ShareType.IsNull() && !data.ShareType.IsUnknown() {
		params["share_type"] = data.ShareType.ValueString()
	}
	if !data.EncryptionOptions.IsNull() && !data.EncryptionOptions.IsUnknown() {
		var encryption_optionsObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.EncryptionOptions.ValueString()), &encryption_optionsObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse encryption_options: %s", err))
			return
		}
		params["encryption_options"] = encryption_optionsObj
	}
	if !data.Encryption.IsNull() && !data.Encryption.IsUnknown() {
		params["encryption"] = data.Encryption.ValueBool()
	}
	if !data.InheritEncryption.IsNull() && !data.InheritEncryption.IsUnknown() {
		params["inherit_encryption"] = data.InheritEncryption.ValueBool()
	}
	if !data.UserProperties.IsNull() && !data.UserProperties.IsUnknown() {
		var user_propertiesList []string
		data.UserProperties.ElementsAs(ctx, &user_propertiesList, false)
		var user_propertiesObjs []map[string]interface{}
		for _, jsonStr := range user_propertiesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse user_properties item: %s", err))
				return
			}
			user_propertiesObjs = append(user_propertiesObjs, obj)
		}
		params["user_properties"] = user_propertiesObjs
	}
	if !data.CreateAncestors.IsNull() && !data.CreateAncestors.IsUnknown() {
		params["create_ancestors"] = data.CreateAncestors.ValueBool()
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		params["type"] = data.Type.ValueString()
	}
	if !data.Aclmode.IsNull() && !data.Aclmode.IsUnknown() {
		params["aclmode"] = data.Aclmode.ValueString()
	}
	if !data.Acltype.IsNull() && !data.Acltype.IsUnknown() {
		params["acltype"] = data.Acltype.ValueString()
	}
	if !data.Atime.IsNull() && !data.Atime.IsUnknown() {
		params["atime"] = data.Atime.ValueString()
	}
	if !data.Casesensitivity.IsNull() && !data.Casesensitivity.IsUnknown() {
		params["casesensitivity"] = data.Casesensitivity.ValueString()
	}
	if !data.Quota.IsNull() && !data.Quota.IsUnknown() {
		params["quota"] = data.Quota.ValueInt64()
	}
	if !data.Refquota.IsNull() && !data.Refquota.IsUnknown() {
		params["refquota"] = data.Refquota.ValueInt64()
	}
	if !data.Recordsize.IsNull() && !data.Recordsize.IsUnknown() {
		params["recordsize"] = data.Recordsize.ValueString()
	}
	if !data.ForceSize.IsNull() && !data.ForceSize.IsUnknown() {
		params["force_size"] = data.ForceSize.ValueBool()
	}
	if !data.Sparse.IsNull() && !data.Sparse.IsUnknown() {
		params["sparse"] = data.Sparse.ValueBool()
	}
	if !data.Volsize.IsNull() && !data.Volsize.IsUnknown() {
		params["volsize"] = data.Volsize.ValueInt64()
	}
	if !data.Volblocksize.IsNull() && !data.Volblocksize.IsUnknown() {
		params["volblocksize"] = data.Volblocksize.ValueString()
	}
	if !data.UserPropertiesUpdate.IsNull() && !data.UserPropertiesUpdate.IsUnknown() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		var user_properties_updateObjs []map[string]interface{}
		for _, jsonStr := range user_properties_updateList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse user_properties_update item: %s", err))
				return
			}
			user_properties_updateObjs = append(user_properties_updateObjs, obj)
		}
		params["user_properties_update"] = user_properties_updateObjs
	}

	result, err := r.client.Call("pool.dataset.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create pool_dataset: %s", err))
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
	id := data.ID.ValueString()
	result, err = r.client.Call("pool.dataset.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back pool_dataset: %s", err))
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
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["casesensitivity"]; ok {
		switch val := v.(type) {
		case string:
			data.Casesensitivity = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["sparse"]; ok {
		if bv, ok := v.(bool); ok {
			data.Sparse = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["volblocksize"]; ok {
		switch val := v.(type) {
		case string:
			data.Volblocksize = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Volblocksize = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Volblocksize = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Comments.IsUnknown() {
		data.Comments = types.StringNull()
	}
	if data.Sync.IsUnknown() {
		data.Sync = types.StringNull()
	}
	if data.Snapdev.IsUnknown() {
		data.Snapdev = types.StringNull()
	}
	if data.Compression.IsUnknown() {
		data.Compression = types.StringNull()
	}
	if data.Exec.IsUnknown() {
		data.Exec = types.StringNull()
	}
	if data.Managedby.IsUnknown() {
		data.Managedby = types.StringNull()
	}
	if data.QuotaWarning.IsUnknown() {
		data.QuotaWarning = types.Int64Null()
	}
	if data.QuotaCritical.IsUnknown() {
		data.QuotaCritical = types.Int64Null()
	}
	if data.RefquotaWarning.IsUnknown() {
		data.RefquotaWarning = types.Int64Null()
	}
	if data.RefquotaCritical.IsUnknown() {
		data.RefquotaCritical = types.Int64Null()
	}
	if data.Reservation.IsUnknown() {
		data.Reservation = types.Int64Null()
	}
	if data.Refreservation.IsUnknown() {
		data.Refreservation = types.Int64Null()
	}
	if data.SpecialSmallBlockSize.IsUnknown() {
		data.SpecialSmallBlockSize = types.Int64Null()
	}
	if data.Copies.IsUnknown() {
		data.Copies = types.Int64Null()
	}
	if data.Snapdir.IsUnknown() {
		data.Snapdir = types.StringNull()
	}
	if data.Deduplication.IsUnknown() {
		data.Deduplication = types.StringNull()
	}
	if data.Checksum.IsUnknown() {
		data.Checksum = types.StringNull()
	}
	if data.Readonly.IsUnknown() {
		data.Readonly = types.StringNull()
	}
	if data.ShareType.IsUnknown() {
		data.ShareType = types.StringNull()
	}
	if data.EncryptionOptions.IsUnknown() {
		data.EncryptionOptions = types.StringNull()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.InheritEncryption.IsUnknown() {
		data.InheritEncryption = types.BoolNull()
	}
	if data.UserProperties.IsUnknown() {
		data.UserProperties, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.CreateAncestors.IsUnknown() {
		data.CreateAncestors = types.BoolNull()
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.Aclmode.IsUnknown() {
		data.Aclmode = types.StringNull()
	}
	if data.Acltype.IsUnknown() {
		data.Acltype = types.StringNull()
	}
	if data.Atime.IsUnknown() {
		data.Atime = types.StringNull()
	}
	if data.Casesensitivity.IsUnknown() {
		data.Casesensitivity = types.StringNull()
	}
	if data.Quota.IsUnknown() {
		data.Quota = types.Int64Null()
	}
	if data.Refquota.IsUnknown() {
		data.Refquota = types.Int64Null()
	}
	if data.Recordsize.IsUnknown() {
		data.Recordsize = types.StringNull()
	}
	if data.ForceSize.IsUnknown() {
		data.ForceSize = types.BoolNull()
	}
	if data.Sparse.IsUnknown() {
		data.Sparse = types.BoolNull()
	}
	if data.Volsize.IsUnknown() {
		data.Volsize = types.Int64Null()
	}
	if data.Volblocksize.IsUnknown() {
		data.Volblocksize = types.StringNull()
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
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read pool_dataset: %s", err))
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
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["casesensitivity"]; ok {
		switch val := v.(type) {
		case string:
			data.Casesensitivity = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["sparse"]; ok {
		if bv, ok := v.(bool); ok {
			data.Sparse = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["volblocksize"]; ok {
		switch val := v.(type) {
		case string:
			data.Volblocksize = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Volblocksize = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Volblocksize = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Comments.IsUnknown() {
		data.Comments = types.StringNull()
	}
	if data.Sync.IsUnknown() {
		data.Sync = types.StringNull()
	}
	if data.Snapdev.IsUnknown() {
		data.Snapdev = types.StringNull()
	}
	if data.Compression.IsUnknown() {
		data.Compression = types.StringNull()
	}
	if data.Exec.IsUnknown() {
		data.Exec = types.StringNull()
	}
	if data.Managedby.IsUnknown() {
		data.Managedby = types.StringNull()
	}
	if data.QuotaWarning.IsUnknown() {
		data.QuotaWarning = types.Int64Null()
	}
	if data.QuotaCritical.IsUnknown() {
		data.QuotaCritical = types.Int64Null()
	}
	if data.RefquotaWarning.IsUnknown() {
		data.RefquotaWarning = types.Int64Null()
	}
	if data.RefquotaCritical.IsUnknown() {
		data.RefquotaCritical = types.Int64Null()
	}
	if data.Reservation.IsUnknown() {
		data.Reservation = types.Int64Null()
	}
	if data.Refreservation.IsUnknown() {
		data.Refreservation = types.Int64Null()
	}
	if data.SpecialSmallBlockSize.IsUnknown() {
		data.SpecialSmallBlockSize = types.Int64Null()
	}
	if data.Copies.IsUnknown() {
		data.Copies = types.Int64Null()
	}
	if data.Snapdir.IsUnknown() {
		data.Snapdir = types.StringNull()
	}
	if data.Deduplication.IsUnknown() {
		data.Deduplication = types.StringNull()
	}
	if data.Checksum.IsUnknown() {
		data.Checksum = types.StringNull()
	}
	if data.Readonly.IsUnknown() {
		data.Readonly = types.StringNull()
	}
	if data.ShareType.IsUnknown() {
		data.ShareType = types.StringNull()
	}
	if data.EncryptionOptions.IsUnknown() {
		data.EncryptionOptions = types.StringNull()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.InheritEncryption.IsUnknown() {
		data.InheritEncryption = types.BoolNull()
	}
	if data.UserProperties.IsUnknown() {
		data.UserProperties, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.CreateAncestors.IsUnknown() {
		data.CreateAncestors = types.BoolNull()
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.Aclmode.IsUnknown() {
		data.Aclmode = types.StringNull()
	}
	if data.Acltype.IsUnknown() {
		data.Acltype = types.StringNull()
	}
	if data.Atime.IsUnknown() {
		data.Atime = types.StringNull()
	}
	if data.Casesensitivity.IsUnknown() {
		data.Casesensitivity = types.StringNull()
	}
	if data.Quota.IsUnknown() {
		data.Quota = types.Int64Null()
	}
	if data.Refquota.IsUnknown() {
		data.Refquota = types.Int64Null()
	}
	if data.Recordsize.IsUnknown() {
		data.Recordsize = types.StringNull()
	}
	if data.ForceSize.IsUnknown() {
		data.ForceSize = types.BoolNull()
	}
	if data.Sparse.IsUnknown() {
		data.Sparse = types.BoolNull()
	}
	if data.Volsize.IsUnknown() {
		data.Volsize = types.Int64Null()
	}
	if data.Volblocksize.IsUnknown() {
		data.Volblocksize = types.StringNull()
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
	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		params["comments"] = data.Comments.ValueString()
	}
	if !data.Sync.IsNull() && !data.Sync.IsUnknown() {
		params["sync"] = data.Sync.ValueString()
	}
	if !data.Snapdev.IsNull() && !data.Snapdev.IsUnknown() {
		params["snapdev"] = data.Snapdev.ValueString()
	}
	if !data.Compression.IsNull() && !data.Compression.IsUnknown() {
		params["compression"] = data.Compression.ValueString()
	}
	if !data.Exec.IsNull() && !data.Exec.IsUnknown() {
		params["exec"] = data.Exec.ValueString()
	}
	if !data.Managedby.IsNull() && !data.Managedby.IsUnknown() {
		params["managedby"] = data.Managedby.ValueString()
	}
	if !data.QuotaWarning.IsNull() && !data.QuotaWarning.IsUnknown() {
		params["quota_warning"] = data.QuotaWarning.ValueInt64()
	}
	if !data.QuotaCritical.IsNull() && !data.QuotaCritical.IsUnknown() {
		params["quota_critical"] = data.QuotaCritical.ValueInt64()
	}
	if !data.RefquotaWarning.IsNull() && !data.RefquotaWarning.IsUnknown() {
		params["refquota_warning"] = data.RefquotaWarning.ValueInt64()
	}
	if !data.RefquotaCritical.IsNull() && !data.RefquotaCritical.IsUnknown() {
		params["refquota_critical"] = data.RefquotaCritical.ValueInt64()
	}
	if !data.Reservation.IsNull() && !data.Reservation.IsUnknown() {
		params["reservation"] = data.Reservation.ValueInt64()
	}
	if !data.Refreservation.IsNull() && !data.Refreservation.IsUnknown() {
		params["refreservation"] = data.Refreservation.ValueInt64()
	}
	if !data.SpecialSmallBlockSize.IsNull() && !data.SpecialSmallBlockSize.IsUnknown() {
		params["special_small_block_size"] = data.SpecialSmallBlockSize.ValueInt64()
	}
	if !data.Copies.IsNull() && !data.Copies.IsUnknown() {
		params["copies"] = data.Copies.ValueInt64()
	}
	if !data.Snapdir.IsNull() && !data.Snapdir.IsUnknown() {
		params["snapdir"] = data.Snapdir.ValueString()
	}
	if !data.Deduplication.IsNull() && !data.Deduplication.IsUnknown() {
		params["deduplication"] = data.Deduplication.ValueString()
	}
	if !data.Checksum.IsNull() && !data.Checksum.IsUnknown() {
		params["checksum"] = data.Checksum.ValueString()
	}
	if !data.Readonly.IsNull() && !data.Readonly.IsUnknown() {
		params["readonly"] = data.Readonly.ValueString()
	}
	if !data.UserProperties.IsNull() && !data.UserProperties.IsUnknown() {
		var user_propertiesList []string
		data.UserProperties.ElementsAs(ctx, &user_propertiesList, false)
		var user_propertiesObjs []map[string]interface{}
		for _, jsonStr := range user_propertiesList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse user_properties item: %s", err))
				return
			}
			user_propertiesObjs = append(user_propertiesObjs, obj)
		}
		params["user_properties"] = user_propertiesObjs
	}
	if !data.CreateAncestors.IsNull() && !data.CreateAncestors.IsUnknown() {
		params["create_ancestors"] = data.CreateAncestors.ValueBool()
	}
	if !data.ForceSize.IsNull() && !data.ForceSize.IsUnknown() {
		params["force_size"] = data.ForceSize.ValueBool()
	}
	if !data.Volsize.IsNull() && !data.Volsize.IsUnknown() {
		params["volsize"] = data.Volsize.ValueInt64()
	}
	if !data.Aclmode.IsNull() && !data.Aclmode.IsUnknown() {
		params["aclmode"] = data.Aclmode.ValueString()
	}
	if !data.Acltype.IsNull() && !data.Acltype.IsUnknown() {
		params["acltype"] = data.Acltype.ValueString()
	}
	if !data.Atime.IsNull() && !data.Atime.IsUnknown() {
		params["atime"] = data.Atime.ValueString()
	}
	if !data.Quota.IsNull() && !data.Quota.IsUnknown() {
		params["quota"] = data.Quota.ValueInt64()
	}
	if !data.Refquota.IsNull() && !data.Refquota.IsUnknown() {
		params["refquota"] = data.Refquota.ValueInt64()
	}
	if !data.Recordsize.IsNull() && !data.Recordsize.IsUnknown() {
		params["recordsize"] = data.Recordsize.ValueString()
	}
	if !data.UserPropertiesUpdate.IsNull() && !data.UserPropertiesUpdate.IsUnknown() {
		var user_properties_updateList []string
		data.UserPropertiesUpdate.ElementsAs(ctx, &user_properties_updateList, false)
		var user_properties_updateObjs []map[string]interface{}
		for _, jsonStr := range user_properties_updateList {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
				resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse user_properties_update item: %s", err))
				return
			}
			user_properties_updateObjs = append(user_properties_updateObjs, obj)
		}
		params["user_properties_update"] = user_properties_updateObjs
	}

	_, err = r.client.Call("pool.dataset.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update pool_dataset: %s", err))
		return
	}

	data.ID = state.ID

	// Read back to populate computed fields
	result, readErr := r.client.Call("pool.dataset.get_instance", id)
	if readErr != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Updated but failed to read back pool_dataset: %s", readErr))
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
	if v, ok := resultMap["type"]; ok {
		switch val := v.(type) {
		case string:
			data.Type = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Type = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Type = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["casesensitivity"]; ok {
		switch val := v.(type) {
		case string:
			data.Casesensitivity = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Casesensitivity = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if v, ok := resultMap["sparse"]; ok {
		if bv, ok := v.(bool); ok {
			data.Sparse = types.BoolValue(bv)
		}
	}
	if v, ok := resultMap["volblocksize"]; ok {
		switch val := v.(type) {
		case string:
			data.Volblocksize = types.StringValue(val)
		case map[string]interface{}:
			if strVal, ok := val["value"]; ok && strVal != nil {
				data.Volblocksize = types.StringValue(fmt.Sprintf("%v", strVal))
			}
		default:
			data.Volblocksize = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	if data.Comments.IsUnknown() {
		data.Comments = types.StringNull()
	}
	if data.Sync.IsUnknown() {
		data.Sync = types.StringNull()
	}
	if data.Snapdev.IsUnknown() {
		data.Snapdev = types.StringNull()
	}
	if data.Compression.IsUnknown() {
		data.Compression = types.StringNull()
	}
	if data.Exec.IsUnknown() {
		data.Exec = types.StringNull()
	}
	if data.Managedby.IsUnknown() {
		data.Managedby = types.StringNull()
	}
	if data.QuotaWarning.IsUnknown() {
		data.QuotaWarning = types.Int64Null()
	}
	if data.QuotaCritical.IsUnknown() {
		data.QuotaCritical = types.Int64Null()
	}
	if data.RefquotaWarning.IsUnknown() {
		data.RefquotaWarning = types.Int64Null()
	}
	if data.RefquotaCritical.IsUnknown() {
		data.RefquotaCritical = types.Int64Null()
	}
	if data.Reservation.IsUnknown() {
		data.Reservation = types.Int64Null()
	}
	if data.Refreservation.IsUnknown() {
		data.Refreservation = types.Int64Null()
	}
	if data.SpecialSmallBlockSize.IsUnknown() {
		data.SpecialSmallBlockSize = types.Int64Null()
	}
	if data.Copies.IsUnknown() {
		data.Copies = types.Int64Null()
	}
	if data.Snapdir.IsUnknown() {
		data.Snapdir = types.StringNull()
	}
	if data.Deduplication.IsUnknown() {
		data.Deduplication = types.StringNull()
	}
	if data.Checksum.IsUnknown() {
		data.Checksum = types.StringNull()
	}
	if data.Readonly.IsUnknown() {
		data.Readonly = types.StringNull()
	}
	if data.ShareType.IsUnknown() {
		data.ShareType = types.StringNull()
	}
	if data.EncryptionOptions.IsUnknown() {
		data.EncryptionOptions = types.StringNull()
	}
	if data.Encryption.IsUnknown() {
		data.Encryption = types.BoolNull()
	}
	if data.InheritEncryption.IsUnknown() {
		data.InheritEncryption = types.BoolNull()
	}
	if data.UserProperties.IsUnknown() {
		data.UserProperties, _ = types.ListValue(types.StringType, []attr.Value{})
	}
	if data.CreateAncestors.IsUnknown() {
		data.CreateAncestors = types.BoolNull()
	}
	if data.Type.IsUnknown() {
		data.Type = types.StringNull()
	}
	if data.Aclmode.IsUnknown() {
		data.Aclmode = types.StringNull()
	}
	if data.Acltype.IsUnknown() {
		data.Acltype = types.StringNull()
	}
	if data.Atime.IsUnknown() {
		data.Atime = types.StringNull()
	}
	if data.Casesensitivity.IsUnknown() {
		data.Casesensitivity = types.StringNull()
	}
	if data.Quota.IsUnknown() {
		data.Quota = types.Int64Null()
	}
	if data.Refquota.IsUnknown() {
		data.Refquota = types.Int64Null()
	}
	if data.Recordsize.IsUnknown() {
		data.Recordsize = types.StringNull()
	}
	if data.ForceSize.IsUnknown() {
		data.ForceSize = types.BoolNull()
	}
	if data.Sparse.IsUnknown() {
		data.Sparse = types.BoolNull()
	}
	if data.Volsize.IsUnknown() {
		data.Volsize = types.Int64Null()
	}
	if data.Volblocksize.IsUnknown() {
		data.Volblocksize = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolDatasetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolDatasetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := []interface{}{data.ID.ValueString(), map[string]interface{}{}}

	_, err := r.client.Call("pool.dataset.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete pool_dataset: %s", err))
		return
	}
}
