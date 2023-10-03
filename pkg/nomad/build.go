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
//		NOMAD_ADDR: ""
//		NOMAD_REGION: ""
//		NOMAD_NAMESPACE: ""
//		NOMAD_HTTP_AUTH: ""
//
//		# tls
//		NOMAD_CACERT: ""
//		NOMAD_CAPATH: ""
//		NOMAD_CLIENT_CERT: ""
//		NOMAD_CLIENT_KEY: ""
//		NOMAD_TLS_SERVER_NAME: ""
//		NOMAD_SKIP_VERIFY: ""
//
//		# acl
//		NOMAD_TOKEN: ""

type MixinConfig struct {
	ServerAddress string `yaml:"NOMAD_ADDR,omitempty"`
	Region        string `yaml:"NOMAD_REGION,omitempty"`
	Namespace     string `yaml:"NOMAD_NAMESPACE,omitempty"`
	HttpAuth      string `yaml:"NOMAD_HTTP_AUTH,omitempty"`

	CaCert        string `yaml:"NOMAD_CACERT,omitempty"`
	CaPath        string `yaml:"NOMAD_CAPATH,omitempty"`
	ClientCert    string `yaml:"NOMAD_CLIENT_CERT,omitempty"`
	ClientKey     string `yaml:"NOMAD_CLIENT_KEY,omitempty"`
	TLSServerName string `yaml:"NOMAD_TLS_SERVER_NAME,omitempty"`
	SkipVerify    string `yaml:"NOMAD_SKIP_VERIFY,omitempty"`

	Token string `yaml:"NOMAD_TOKEN,omitempty"`
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
	fmt.Fprintf(m.Out, "ENV NOMAD_ADDR=%s\n", input.Config.ServerAddress)
	fmt.Fprintf(m.Out, "ENV NOMAD_REGION=%s\n", input.Config.Region)
	fmt.Fprintf(m.Out, "ENV NOMAD_NAMESPACE=%s\n", input.Config.Namespace)
	fmt.Fprintf(m.Out, "ENV NOMAD_HTTP_AUTH=%s\n", input.Config.HttpAuth)

	fmt.Fprintf(m.Out, "ENV NOMAD_CACERT=%s\n", input.Config.CaCert)
	fmt.Fprintf(m.Out, "ENV NOMAD_CAPATH=%s\n", input.Config.CaPath)
	fmt.Fprintf(m.Out, "ENV NOMAD_CLIENT_CERT=%s\n", input.Config.ClientCert)
	fmt.Fprintf(m.Out, "ENV NOMAD_CLIENT_KEY=%s\n", input.Config.ClientKey)
	fmt.Fprintf(m.Out, "ENV NOMAD_TLS_SERVER_NAME=%s\n", input.Config.TLSServerName)
	fmt.Fprintf(m.Out, "ENV NOMAD_SKIP_VERIFY=%s\n", input.Config.SkipVerify)
	fmt.Fprintf(m.Out, "ENV NOMAD_TOKEN=%s\n", input.Config.Token)

	return nil
}
