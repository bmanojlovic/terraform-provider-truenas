package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bmanojlovic/terraform-provider-truenas/internal/client"
)

var _ provider.Provider = &TrueNASProvider{}

type TrueNASProvider struct {
	version string
}

type TrueNASProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *TrueNASProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "truenas"
	resp.Version = p.version
}

func (p *TrueNASProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "TrueNAS host address",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "TrueNAS API token",
				Optional:            true,
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

	host := os.Getenv("TRUENAS_HOST")
	token := os.Getenv("TRUENAS_TOKEN")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddError(
			"Missing TrueNAS Host",
			"The provider cannot create the TrueNAS client as there is a missing or empty value for the TrueNAS host. "+
				"Set the host value in the configuration or use the TRUENAS_HOST environment variable.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing TrueNAS Token",
			"The provider cannot create the TrueNAS client as there is a missing or empty value for the TrueNAS API token. "+
				"Set the token value in the configuration or use the TRUENAS_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := client.NewClient(host, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create TrueNAS Client",
			"An unexpected error occurred when creating the TrueNAS client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"TrueNAS Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *TrueNASProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
NewAcmeDnsAuthenticatorResource,
		NewAlertserviceResource,
		NewApiKeyResource,
		NewAppResource,
		NewAppRegistryResource,
		NewCertificateResource,
		NewCloudBackupResource,
		NewCloudsyncResource,
		NewCloudsyncCredentialsResource,
		NewCronjobResource,
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
		NewNvmetPortSubsysResource,
		NewNvmetSubsysResource,
		NewPoolResource,
		NewPoolDatasetResource,
		NewPoolScrubResource,
		NewPoolSnapshotResource,
		NewPoolSnapshottaskResource,
		NewPrivilegeResource,
		NewReplicationResource,
		NewReportingExportersResource,
		NewRsynctaskResource,
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
	}
}

func (p *TrueNASProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVmDataSource,
		NewPoolDataSource,
		NewPoolDatasetDataSource,
		NewDiskDataSource,
		NewUserDataSource,
		NewGroupDataSource,
		NewInterfaceDataSource,
		NewServiceDataSource,
		NewVmsDataSource,
		NewPoolsDataSource,
		NewPoolDatasetsDataSource,
		NewDisksDataSource,
		NewUsersDataSource,
		NewGroupsDataSource,
		NewInterfacesDataSource,
		NewServicesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TrueNASProvider{
			version: version,
		}
	}
}
