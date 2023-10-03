package nomad

import (
	"context"
	"fmt"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
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
	Address   string `yaml:"address"`
	Region    string `yaml:"region"`
	Namespace string `yaml:"namespace"`
	HttpAuth  string `yaml:"httpAuth"`

	CaCert        string `yaml:"caCert"`
	CaPath        string `yaml:"caPath"`
	ClientCert    string `yaml:"clientCert"`
	ClientKey     string `yaml:"clientKey"`
	TlsServerName string `yaml:"tlsServerName"`
	TlsSkipVerify bool   `yaml:"tlsSkipVerify"`

	// acl
	Token string `yaml:"token"`
}

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build(ctx context.Context) error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(m.Out, "ENV NOMAD_ADDR=%s\n", input.Config.Address)
	fmt.Fprintf(m.Out, "ENV NOMAD_REGION=%s\n", input.Config.Region)
	fmt.Fprintf(m.Out, "ENV NOMAD_NAMESPACE=%s\n", input.Config.Namespace)
	fmt.Fprintf(m.Out, "ENV NOMAD_HTTP_AUTH=%s\n", input.Config.HttpAuth)

	fmt.Fprintf(m.Out, "ENV NOMAD_CACERT=%s\n", input.Config.CaCert)
	fmt.Fprintf(m.Out, "ENV NOMAD_CAPATH=%s\n", input.Config.CaPath)
	fmt.Fprintf(m.Out, "ENV NOMAD_CLIENT_CERT=%s\n", input.Config.ClientCert)
	fmt.Fprintf(m.Out, "ENV NOMAD_CLIENT_KEY=%s\n", input.Config.ClientKey)
	fmt.Fprintf(m.Out, "ENV NOMAD_TLS_SERVER_NAME=%s\n", input.Config.TlsServerName)
	fmt.Fprintf(m.Out, "ENV NOMAD_SKIP_VERIFY=%t\n", input.Config.TlsSkipVerify)
	fmt.Fprintf(m.Out, "ENV NOMAD_TOKEN=%s\n", input.Config.Token)

	return nil
}
