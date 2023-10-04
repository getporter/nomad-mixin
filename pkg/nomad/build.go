package nomad

import (
	"context"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the nomad mixin in porter.yaml
//mixins:
//	- nomad:
//		address: ""
//		region: ""
//		namespace: ""
//		httpAuth: ""
//
//		# tls
//		caCert: ""
//		caPath: ""
//		clientCert: ""
//		clientKey: ""
//		tlsServerName: ""
//		tlsSkipVerify:
//
//		# acl
//		token: ""

type MixinConfig struct {
}

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build(ctx context.Context) error {
	return nil
}
