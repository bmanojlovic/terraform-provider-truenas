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

type RsynctaskResource struct {
	client *client.Client
}

type RsynctaskResourceModel struct {
	ID types.String `tfsdk:"id"`
	Path types.String `tfsdk:"path"`
	User types.String `tfsdk:"user"`
	Mode types.String `tfsdk:"mode"`
	Remotehost types.String `tfsdk:"remotehost"`
	Remoteport types.Int64 `tfsdk:"remoteport"`
	Remotemodule types.String `tfsdk:"remotemodule"`
	SshCredentials types.Int64 `tfsdk:"ssh_credentials"`
	Remotepath types.String `tfsdk:"remotepath"`
	Direction types.String `tfsdk:"direction"`
	Desc types.String `tfsdk:"desc"`
	Schedule types.String `tfsdk:"schedule"`
	Recursive types.Bool `tfsdk:"recursive"`
	Times types.Bool `tfsdk:"times"`
	Compress types.Bool `tfsdk:"compress"`
	Archive types.Bool `tfsdk:"archive"`
	Delete types.Bool `tfsdk:"delete"`
	Quiet types.Bool `tfsdk:"quiet"`
	Preserveperm types.Bool `tfsdk:"preserveperm"`
	Preserveattr types.Bool `tfsdk:"preserveattr"`
	Delayupdates types.Bool `tfsdk:"delayupdates"`
	Extra types.List `tfsdk:"extra"`
	Enabled types.Bool `tfsdk:"enabled"`
	ValidateRpath types.Bool `tfsdk:"validate_rpath"`
	SshKeyscan types.Bool `tfsdk:"ssh_keyscan"`
}

func NewRsynctaskResource() resource.Resource {
	return &RsynctaskResource{}
}

func (r *RsynctaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rsynctask"
}

func (r *RsynctaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RsynctaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Rsync Task.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Resource ID"},
			"path": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Local filesystem path to synchronize.",
			},
			"user": schema.StringAttribute{
				Required: true,
				Optional: false,
				Description: "Username to run the rsync task as.",
			},
			"mode": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Operating mechanism for Rsync, i.e. Rsync Module mode or Rsync SSH mode.",
			},
			"remotehost": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "IP address or hostname of the remote system. If username differs on the remote host, \"username@remot",
			},
			"remoteport": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Port number for SSH connection. Only applies when `mode` is SSH.",
			},
			"remotemodule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Name of remote module, this attribute should be specified when `mode` is set to MODULE.",
			},
			"ssh_credentials": schema.Int64Attribute{
				Required: false,
				Optional: true,
				Description: "Keychain credential ID for SSH authentication. `null` to use user's SSH keys.",
			},
			"remotepath": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Path on the remote system to synchronize with.",
			},
			"direction": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Specify if data should be PULLED or PUSHED from the remote system.",
			},
			"desc": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Description of the rsync task.",
			},
			"schedule": schema.StringAttribute{
				Required: false,
				Optional: true,
				Description: "Cron schedule for when the rsync task should run.",
			},
			"recursive": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Recursively transfer subdirectories.",
			},
			"times": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Preserve modification times of files.",
			},
			"compress": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Reduce the size of the data to be transmitted.",
			},
			"archive": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Make rsync run recursively, preserving symlinks, permissions, modification times, group, and special",
			},
			"delete": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Delete files in the destination directory that do not exist in the source directory.",
			},
			"quiet": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Suppress informational messages from rsync.",
			},
			"preserveperm": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Preserve original file permissions.",
			},
			"preserveattr": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Preserve extended attributes of files.",
			},
			"delayupdates": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Delay updating destination files until all transfers are complete.",
			},
			"extra": schema.ListAttribute{
				Required: false,
				Optional: true,
				ElementType: types.StringType,
				Description: "Array of additional rsync command-line options.",
			},
			"enabled": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Whether this rsync task is enabled.",
			},
			"validate_rpath": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Validate the existence of the remote path.",
			},
			"ssh_keyscan": schema.BoolAttribute{
				Required: false,
				Optional: true,
				Description: "Automatically add remote host key to user's known_hosts file.",
			},
		},
	}
}

func (r *RsynctaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RsynctaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RsynctaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]interface{}{}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
		params["mode"] = data.Mode.ValueString()
	}
	if !data.Remotehost.IsNull() && !data.Remotehost.IsUnknown() {
		params["remotehost"] = data.Remotehost.ValueString()
	}
	if !data.Remoteport.IsNull() && !data.Remoteport.IsUnknown() {
		params["remoteport"] = data.Remoteport.ValueInt64()
	}
	if !data.Remotemodule.IsNull() && !data.Remotemodule.IsUnknown() {
		params["remotemodule"] = data.Remotemodule.ValueString()
	}
	if !data.SshCredentials.IsNull() && !data.SshCredentials.IsUnknown() {
		params["ssh_credentials"] = data.SshCredentials.ValueInt64()
	}
	if !data.Remotepath.IsNull() && !data.Remotepath.IsUnknown() {
		params["remotepath"] = data.Remotepath.ValueString()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.Desc.IsNull() && !data.Desc.IsUnknown() {
		params["desc"] = data.Desc.ValueString()
	}
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.Recursive.IsNull() && !data.Recursive.IsUnknown() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Times.IsNull() && !data.Times.IsUnknown() {
		params["times"] = data.Times.ValueBool()
	}
	if !data.Compress.IsNull() && !data.Compress.IsUnknown() {
		params["compress"] = data.Compress.ValueBool()
	}
	if !data.Archive.IsNull() && !data.Archive.IsUnknown() {
		params["archive"] = data.Archive.ValueBool()
	}
	if !data.Delete.IsNull() && !data.Delete.IsUnknown() {
		params["delete"] = data.Delete.ValueBool()
	}
	if !data.Quiet.IsNull() && !data.Quiet.IsUnknown() {
		params["quiet"] = data.Quiet.ValueBool()
	}
	if !data.Preserveperm.IsNull() && !data.Preserveperm.IsUnknown() {
		params["preserveperm"] = data.Preserveperm.ValueBool()
	}
	if !data.Preserveattr.IsNull() && !data.Preserveattr.IsUnknown() {
		params["preserveattr"] = data.Preserveattr.ValueBool()
	}
	if !data.Delayupdates.IsNull() && !data.Delayupdates.IsUnknown() {
		params["delayupdates"] = data.Delayupdates.ValueBool()
	}
	if !data.Extra.IsNull() && !data.Extra.IsUnknown() {
		var extraList []string
		data.Extra.ElementsAs(ctx, &extraList, false)
		params["extra"] = extraList
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ValidateRpath.IsNull() && !data.ValidateRpath.IsUnknown() {
		params["validate_rpath"] = data.ValidateRpath.ValueBool()
	}
	if !data.SshKeyscan.IsNull() && !data.SshKeyscan.IsUnknown() {
		params["ssh_keyscan"] = data.SshKeyscan.ValueBool()
	}

	result, err := r.client.Call("rsynctask.create", params)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Unable to create rsynctask: %s", err))
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
	var id interface{}
	id, err = strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Cannot parse ID: %s", err))
		return
	}
	result, err = r.client.Call("rsynctask.get_instance", id)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Created but failed to read back rsynctask: %s", err))
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
		if v, ok := resultMap["path"]; ok {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["user"]; ok {
			switch val := v.(type) {
			case string:
				data.User = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.User = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.User = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["remotehost"]; ok {
			if v == nil {
				data.Remotehost = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Remotehost = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Remotehost = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Remotehost = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["remoteport"]; ok {
			if v == nil {
				data.Remoteport = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.Remoteport = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.Remoteport = types.Int64Value(int64(fv)) }
					}
				}
			}
		}
		if v, ok := resultMap["remotemodule"]; ok {
			if v == nil {
				data.Remotemodule = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Remotemodule = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Remotemodule = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Remotemodule = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["ssh_credentials"]; ok {
			if v == nil {
				data.SshCredentials = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.SshCredentials = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.SshCredentials = types.Int64Value(int64(fv)) }
					}
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RsynctaskResourceModel
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

	result, err := r.client.Call("rsynctask.get_instance", id)
	if err != nil {
		// Check if resource was deleted outside Terraform (ENOENT = entity not found)
		if strings.Contains(err.Error(), "[ENOENT]") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read rsynctask: %s", err))
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
		if v, ok := resultMap["path"]; ok {
			switch val := v.(type) {
			case string:
				data.Path = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.Path = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.Path = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["user"]; ok {
			switch val := v.(type) {
			case string:
				data.User = types.StringValue(val)
			case map[string]interface{}:
				if strVal, ok := val["value"]; ok && strVal != nil {
					data.User = types.StringValue(fmt.Sprintf("%v", strVal))
				}
			default:
				data.User = types.StringValue(fmt.Sprintf("%v", v))
			}
		}
		if v, ok := resultMap["remotehost"]; ok {
			if v == nil {
				data.Remotehost = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Remotehost = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Remotehost = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Remotehost = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["remoteport"]; ok {
			if v == nil {
				data.Remoteport = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.Remoteport = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.Remoteport = types.Int64Value(int64(fv)) }
					}
				}
			}
		}
		if v, ok := resultMap["remotemodule"]; ok {
			if v == nil {
				data.Remotemodule = types.StringNull()
			} else {
				switch val := v.(type) {
				case string:
					data.Remotemodule = types.StringValue(val)
				case map[string]interface{}:
					if strVal, ok := val["value"]; ok && strVal != nil {
						data.Remotemodule = types.StringValue(fmt.Sprintf("%v", strVal))
					}
				default:
					data.Remotemodule = types.StringValue(fmt.Sprintf("%v", v))
				}
			}
		}
		if v, ok := resultMap["ssh_credentials"]; ok {
			if v == nil {
				data.SshCredentials = types.Int64Null()
			} else {
				switch val := v.(type) {
				case float64:
					data.SshCredentials = types.Int64Value(int64(val))
				case map[string]interface{}:
					if parsed, ok := val["parsed"]; ok && parsed != nil {
						if fv, ok := parsed.(float64); ok { data.SshCredentials = types.Int64Value(int64(fv)) }
					}
				}
			}
		}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RsynctaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state RsynctaskResourceModel
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
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		params["path"] = data.Path.ValueString()
	}
	if !data.User.IsNull() && !data.User.IsUnknown() {
		params["user"] = data.User.ValueString()
	}
	if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
		params["mode"] = data.Mode.ValueString()
	}
	if !data.Remotehost.IsNull() && !data.Remotehost.IsUnknown() {
		params["remotehost"] = data.Remotehost.ValueString()
	}
	if !data.Remoteport.IsNull() && !data.Remoteport.IsUnknown() {
		params["remoteport"] = data.Remoteport.ValueInt64()
	}
	if !data.Remotemodule.IsNull() && !data.Remotemodule.IsUnknown() {
		params["remotemodule"] = data.Remotemodule.ValueString()
	}
	if !data.SshCredentials.IsNull() && !data.SshCredentials.IsUnknown() {
		params["ssh_credentials"] = data.SshCredentials.ValueInt64()
	}
	if !data.Remotepath.IsNull() && !data.Remotepath.IsUnknown() {
		params["remotepath"] = data.Remotepath.ValueString()
	}
	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		params["direction"] = data.Direction.ValueString()
	}
	if !data.Desc.IsNull() && !data.Desc.IsUnknown() {
		params["desc"] = data.Desc.ValueString()
	}
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		var scheduleObj map[string]interface{}
		if err := json.Unmarshal([]byte(data.Schedule.ValueString()), &scheduleObj); err != nil {
			resp.Diagnostics.AddError("JSON Parse Error", fmt.Sprintf("Failed to parse schedule: %s", err))
			return
		}
		params["schedule"] = scheduleObj
	}
	if !data.Recursive.IsNull() && !data.Recursive.IsUnknown() {
		params["recursive"] = data.Recursive.ValueBool()
	}
	if !data.Times.IsNull() && !data.Times.IsUnknown() {
		params["times"] = data.Times.ValueBool()
	}
	if !data.Compress.IsNull() && !data.Compress.IsUnknown() {
		params["compress"] = data.Compress.ValueBool()
	}
	if !data.Archive.IsNull() && !data.Archive.IsUnknown() {
		params["archive"] = data.Archive.ValueBool()
	}
	if !data.Delete.IsNull() && !data.Delete.IsUnknown() {
		params["delete"] = data.Delete.ValueBool()
	}
	if !data.Quiet.IsNull() && !data.Quiet.IsUnknown() {
		params["quiet"] = data.Quiet.ValueBool()
	}
	if !data.Preserveperm.IsNull() && !data.Preserveperm.IsUnknown() {
		params["preserveperm"] = data.Preserveperm.ValueBool()
	}
	if !data.Preserveattr.IsNull() && !data.Preserveattr.IsUnknown() {
		params["preserveattr"] = data.Preserveattr.ValueBool()
	}
	if !data.Delayupdates.IsNull() && !data.Delayupdates.IsUnknown() {
		params["delayupdates"] = data.Delayupdates.ValueBool()
	}
	if !data.Extra.IsNull() && !data.Extra.IsUnknown() {
		var extraList []string
		data.Extra.ElementsAs(ctx, &extraList, false)
		params["extra"] = extraList
	}
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		params["enabled"] = data.Enabled.ValueBool()
	}
	if !data.ValidateRpath.IsNull() && !data.ValidateRpath.IsUnknown() {
		params["validate_rpath"] = data.ValidateRpath.ValueBool()
	}
	if !data.SshKeyscan.IsNull() && !data.SshKeyscan.IsUnknown() {
		params["ssh_keyscan"] = data.SshKeyscan.ValueBool()
	}

	_, err = r.client.Call("rsynctask.update", []interface{}{id, params})
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Unable to update rsynctask: %s", err))
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RsynctaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RsynctaskResourceModel
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

	_, err = r.client.Call("rsynctask.delete", id)
	if err != nil {
		// Ignore ENOENT - resource already deleted
		if strings.Contains(err.Error(), "[ENOENT]") {
			return
		}
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete rsynctask: %s", err))
		return
	}
}
