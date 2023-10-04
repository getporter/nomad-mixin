package nomad

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_invalidRun(t *testing.T) {

	tests := []struct {
		name    string
		job     Job
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "invalid job",
			job: Job{
				Path:             "path",
				Dispatch:         "dispatch",
				IdPrefixTemplate: "",
				Payload:          "",
				Meta:             nil,
				Stop:             "stop",
				Purge:            false,
				Address:          "",
				Region:           "",
				Namespace:        "",
				CaCert:           "",
				CaPath:           "",
				ClientCert:       "",
				ClientKey:        "",
				TlsServerName:    "",
				TlsSkipVerify:    false,
				Token:            "",
			}, wantErr: assert.Error},
		{
			name: "invalid job",
			job: Job{
				Path:             "path",
				Dispatch:         "dispatch",
				IdPrefixTemplate: "",
				Payload:          "",
				Meta:             nil,
				Stop:             "",
				Purge:            false,
				Address:          "",
				Region:           "",
				Namespace:        "",
				CaCert:           "",
				CaPath:           "",
				ClientCert:       "",
				ClientKey:        "",
				TlsServerName:    "",
				TlsSkipVerify:    false,
				Token:            "",
			}, wantErr: assert.Error},
		{
			name: "invalid job",
			job: Job{
				Path:             "",
				Dispatch:         "dispatch",
				IdPrefixTemplate: "",
				Payload:          "",
				Meta:             nil,
				Stop:             "stop",
				Purge:            false,
				Address:          "",
				Region:           "",
				Namespace:        "",
				CaCert:           "",
				CaPath:           "",
				ClientCert:       "",
				ClientKey:        "",
				TlsServerName:    "",
				TlsSkipVerify:    false,
				Token:            "",
			}, wantErr: assert.Error},
		{
			name: "invalid job",
			job: Job{
				Path:             "path",
				Dispatch:         "",
				IdPrefixTemplate: "",
				Payload:          "",
				Meta:             nil,
				Stop:             "stop",
				Purge:            false,
				Address:          "",
				Region:           "",
				Namespace:        "",
				CaCert:           "",
				CaPath:           "",
				ClientCert:       "",
				ClientKey:        "",
				TlsServerName:    "",
				TlsSkipVerify:    false,
				Token:            "",
			}, wantErr: assert.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateJob(tt.job), fmt.Sprintf("validateJob(%v)", tt.job))
		})
	}
}
