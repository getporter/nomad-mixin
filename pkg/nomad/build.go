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

// MixinConfig represents configuration that can be set on the skeletor mixin in porter.yaml
// mixins:
// - skeletor:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`

	ServerAddress string `yaml:"NOMAD_ADDR,omitempty"`
	Region        string `yaml:"NOMAD_REGION,omitempty"`
	Namespace     string `yaml:"NOMAD_NAMESPACE,omitempty"`
	HttpAuth      string `yaml:"NOMAD_HTTP_AUTH,omitempty"`

	CaCert        string `yaml:"NOMAD_CACERT,omitempty"`
	Token         string `yaml:"NOMAD_TOKEN,omitempty"`
	CaPath        string `yaml:"NOMAD_CAPATH,omitempty"`
	ClientCert    string `yaml:"NOMAD_CLIENT_CERT,omitempty"`
	ClientKey     string `yaml:"NOMAD_CLIENT_KEY,omitempty"`
	TLSServerName string `yaml:"NOMAD_TLS_SERVER_NAME,omitempty"`
	SkipVerify    string `yaml:"NOMAD_SKIP_VERIFY,omitempty"`
}

// This is an example. Replace the following with whatever steps are needed to
// install required components into
// const dockerfileLines = `RUN apt-get update && \
// apt-get install gnupg apt-transport-https lsb-release software-properties-common -y && \
// echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ stretch main" | \
//    tee /etc/apt/sources.list.d/azure-cli.list && \
// apt-key --keyring /etc/apt/trusted.gpg.d/Microsoft.gpg adv \
// 	--keyserver packages.microsoft.com \
// 	--recv-keys BC528686B50D79E339D3721CEB3E94ADBE1229CF && \
// apt-get update && apt-get install azure-cli
// `

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
	suppliedClientVersion := input.Config.ClientVersion
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

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	return nil
}
