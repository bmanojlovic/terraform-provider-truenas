package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

type TrueNASProvider struct {
	version string
}

type TrueNASProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TrueNASProvider{
			version: version,
		}
	}
}

func (p *TrueNASProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "truenas"
	resp.Version = p.version
}

func (p *TrueNASProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "TrueNAS host (e.g., 192.168.1.100)",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "TrueNAS API token",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *TrueNASProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TrueNASProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := client.NewClient(data.Host.ValueString(), data.Token.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
	
	if err := client.Connect(); err != nil {
		resp.Diagnostics.AddError("Connection Error", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TrueNASProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAcmeDnsAuthenticatorResource,
		NewAlertserviceResource,
		NewApiKeyResource,
		NewAppResource,
		NewAppRedeployActionResource,
		NewAppRollbackActionResource,
		NewAppRollbackVersionsActionResource,
		NewAppUpgradeActionResource,
		NewAppUpgradeSummaryActionResource,
		NewAppRegistryResource,
		NewCertificateResource,
		NewCloudBackupResource,
		NewCloudBackupRestoreActionResource,
		NewCloudBackupSyncActionResource,
		NewCloudsyncResource,
		NewCloudsyncRestoreActionResource,
		NewCloudsyncSyncActionResource,
		NewCloudsyncSyncOnetimeActionResource,
		NewCloudsyncCredentialsResource,
		NewCronjobResource,
		NewCronjobRunActionResource,
		NewFcFcHostResource,
		NewFcportResource,
		NewFilesystemAcltemplateResource,
		NewGroupResource,
		NewInitshutdownscriptResource,
		NewInterfaceResource,
		NewIscsiAuthResource,
		NewIscsiExtentResource,
		NewIscsiInitiatorResource,
		NewIscsiPortalResource,
		NewIscsiTargetResource,
		NewIscsiTargetextentResource,
		NewJbofResource,
		NewKerberosKeytabResource,
		NewKerberosRealmResource,
		NewKeychaincredentialResource,
		NewNvmetHostResource,
		NewNvmetHostSubsysResource,
		NewNvmetNamespaceResource,
		NewNvmetPortResource,
		NewNvmetPortSubsysResource,
		NewNvmetSubsysResource,
		NewPoolDatasetResource,
		NewPoolScrubResource,
		NewPoolScrubRunActionResource,
		NewPoolScrubScrubActionResource,
		NewPoolSnapshotResource,
		NewPoolSnapshotRollbackActionResource,
		NewPoolSnapshottaskResource,
		NewPoolSnapshottaskRunActionResource,
		NewPrivilegeResource,
		NewReplicationResource,
		NewReplicationRestoreActionResource,
		NewReplicationRunActionResource,
		NewReplicationRunOnetimeActionResource,
		NewReportingExportersResource,
		NewRsynctaskResource,
		NewRsynctaskRunActionResource,
		NewSharingNfsResource,
		NewSharingSmbResource,
		NewStaticrouteResource,
		NewSystemNtpserverResource,
		NewTunableResource,
		NewUserResource,
		NewVirtInstanceResource,
		NewVirtVolumeResource,
		NewVmResource,
		NewVmDeviceResource,
		NewVmwareResource,
		NewDiskResource,
		NewServiceResource,
	}
}

func (p *TrueNASProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
