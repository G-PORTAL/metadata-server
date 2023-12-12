# Metadata Server

This metadata server emulates the OpenStack Nova [metadata API](https://docs.openstack.org/nova/latest/user/metadata.html)
for being used outside a normal OpenStack environment. As example this can be used if running OpenStack Ironic as a standalone service.

This metadata server can handle [cloud-init](https://github.com/canonical/cloud-init) and [Cloudbase-Init](https://github.com/cloudbase/cloudbase-init)
at the moment.

## Add new sources

Please feel free to add your own API services as datasource. Just raise a pull request and implement the interface.

```go
package sources

type Source interface {
	Type() string
	Initialize(cfg SourceConfig) error
	GetMetadata(ip net.IP) (*Metadata, error)
	ReportLog(log ReportMessage) error
}
```

## Run the metadata server

To run the metadata server use Docker or build the binary and run it on your infrastructure.

### Docker Image

Our pre-build image is published automatically on [Docker Hub](https://hub.docker.com/r/gportal/metadata-server). 
It is using Alpine Linux as base image and is also compatible to run on OpenShift or Kubernetes.

Currently, we do not provide a helm chart for this image but might be coming.

```bash
docker run --rm \
    --pull always \
    --platform linux/amd64 \
    --name metadata-server \
    --publish 0.0.0.0:80:8080 \
    --volume `pwd`/default.config.yaml:/data/config.yaml \
    gportal/metadata-server:latest
```

## Provisioning

#### cloud-init

For getting the metadata server to work with Cloud-Init on bare metal you need at least version `23.1.1`
of cloud-init installed on the machine.

Example `/etc/cloud.cfg.d/01-metadata.cfg` configuration to use the metadata server with cloud-init. It
is important that only OpenStack is configured as datasource to bypass the platform checks of Cloud-Init.

```yaml
datasource_list: [ OpenStack ]

reporting:
  metadata:
    type: webhook
    endpoint: http://169.254.169.254/reporting/cloud-init
    timeout: 10
    retries: 1
```

### Cloudbase-Init

Example `C:\ProgramData\Cloudbase Solutions\Cloudbase-Init\conf\cloudbase-init.conf` configuration to use the metadata server with Cloudbase-Init.

```ini
[DEFAULT]
username=Administrator
groups=Administrators

first_logon_behaviour=no
inject_user_password=true

mtu_use_dhcp_config=true
ntp_use_dhcp_config=true

metadata_services=cloudbaseinit.metadata.services.httpservice.HttpService

bsdtar_path=C:\Program Files\Cloudbase Solutions\Cloudbase-Init\bin\bsdtar.exe
mtools_path=C:\Program Files\Cloudbase Solutions\Cloudbase-Init\bin\

log-dir=C:\Program Files\Cloudbase Solutions\Cloudbase-Init\log\
log-file=cloudbase-init.log

local_scripts_path=C:\Program Files\Cloudbase Solutions\Cloudbase-Init\LocalScripts\

default_log_levels=comtypes=INFO,suds=INFO,iso8601=WARN,requests=INFO
verbose=true
debug=true
```
