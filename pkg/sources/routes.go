package sources

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/crypto/ssh"
)

// GetRoutes returns a map of all routes that are available for the given metadata.
func (m Metadata) GetRoutes() Routes {
	routes := make(Routes)
	if m.InstanceID != "" {
		routes.registerVersionedOpenStackMetadataRoute("/instance-id", render.String{Format: m.InstanceID})
	}
	if m.InstanceType != "" {
		routes.registerVersionedOpenStackMetadataRoute("/instance-type", render.String{Format: m.InstanceType})
	}

	if m.LocalHostname != "" {
		routes.registerVersionedOpenStackMetadataRoute("/local-hostname", render.String{Format: m.LocalHostname})
	}
	if m.PublicHostname != "" {
		routes.registerVersionedOpenStackMetadataRoute("/hostname", render.String{Format: m.PublicHostname})
		routes.registerVersionedOpenStackMetadataRoute("/public-hostname", render.String{Format: m.PublicHostname})
	}

	if len(m.PublicKeys) > 0 {
		for id, publicKey := range m.PublicKeys {
			routes.registerVersionedOpenStackMetadataRoute(fmt.Sprintf("/public-keys/%s", id), render.String{Format: string(ssh.MarshalAuthorizedKey(publicKey))})
		}
	}

	if m.UserData != nil {
		routes.registerOpenStackRoute("/user_data", render.Data{Data: m.UserData})
		routes.registerVersionedOpenStackMetadataRoute("/user-data", render.Data{Data: m.UserData})
	}

	if m.Password != nil && *m.Password != "" {
		routes.registerOpenStackRoute("/password", render.String{Format: *m.Password})
	}

	if len(m.Interfaces) > 0 {
		routes.registerOpenStackRoute("/network_data.json", m.OpenStackNetworkData())
	}

	if m.VendorData != nil {
		routes.registerOpenStackRoute("/vendor_data.json", m.OpenStackVendorData(m.VendorData))
	}

	if m.VendorData2 != nil {
		routes.registerOpenStackRoute("/vendor_data2.json", m.OpenStackVendorData(m.VendorData))
	}

	routes.registerOpenStackRoute("/meta_data.json", m.OpenStackMetaData())

	return routes
}
