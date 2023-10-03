package nomad

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseInstallAction(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *Nomad
		wantErr assert.ErrorAssertionFunc
	}{
		{name: "test", args: args{
			file: "testdata/install.yaml",
		}, want: &Nomad{Jobs: []Job{
			{
				Path:             "nomad/path",
				Dispatch:         "",
				IdPrefixTemplate: "",
				Payload:          "",
				Meta:             map[string]string(nil),
				Stop:             "",
				Purge:            false,
				Address:          "123",
				Region:           "us-east-1",
				Namespace:        "default",
				CaCert:           "cert",
				CaPath:           "path",
				ClientCert:       "cert",
				ClientKey:        "key",
				TlsServerName:    "name",
				TlsSkipVerify:    true,
				Token:            "token"},
			{
				Path:             "",
				Dispatch:         "dispatch",
				IdPrefixTemplate: "prefix",
				Payload:          "payload",
				Meta:             map[string]string{"budget": "200"},
				Stop:             "",
				Purge:            false,
				Address:          "123",
				Region:           "us-east-1",
				Namespace:        "default",
				CaCert:           "cert",
				CaPath:           "path",
				ClientCert:       "cert",
				ClientKey:        "key",
				TlsServerName:    "name",
				TlsSkipVerify:    true,
				Token:            "token"},
			{
				Path:             "",
				Dispatch:         "",
				IdPrefixTemplate: "",
				Payload:          "",
				Meta:             map[string]string(nil),
				Stop:             "stop",
				Purge:            true,
				Address:          "123",
				Region:           "us-east-1",
				Namespace:        "default",
				CaCert:           "cert",
				CaPath:           "path",
				ClientCert:       "cert",
				ClientKey:        "key",
				TlsServerName:    "name",
				TlsSkipVerify:    false,
				Token:            "token"},
		},
		}, wantErr: assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTestMixin(t).Mixin
			// cast m.In to Buffer
			buffer := bytes.Buffer{}
			bytes, err := os.ReadFile(tt.args.file)
			assert.NoError(t, err)
			buffer.Write(bytes)
			m.In = &buffer

			got, err := parseInstallAction(m)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
