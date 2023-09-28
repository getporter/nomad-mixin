package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/pkg/errors"
)

type Nomad struct {
	Jobs []Job `yaml:"jobs"`
}

type Job struct {
	// regular job run
	Path string `yaml:"path,omitempty"`

	// dispatch job run
	Dispatch         string            `yaml:"dispatch,omitempty"`
	IdPrefixTemplate string            `yaml:"idPrefixTemplate"`
	Payload          string            `yaml:"payload"` //todo make PayloadPath?
	Meta             map[string]string `yaml:"meta,omitempty"`

	// job stop
	Stop  string `yaml:"stop,omitempty"`
	Purge bool   `yaml:"purge,omitempty"`

	// nomad client config
	Address   string `yaml:"address"`
	Region    string `yaml:"region"`
	Namespace string `yaml:"namespace"`

	CaCert        string `yaml:"caCert"`
	CaPath        string `yaml:"caPath"`
	ClientCert    string `yaml:"clientCert"`
	ClientKey     string `yaml:"clientKey"`
	TlsServerName string `yaml:"tlsServerName"`
	TlsSkipVerify bool   `yaml:"tlsSkipVerify"`

	// acl
	Token string `yaml:"token"`
}

func (m *Mixin) DoAction(action *Nomad) error {
	for _, run := range action.Jobs {
		err := validateRun(run)
		if err != nil {
			return err
		}

		// config will respect user global ADDR, REGION, NAMESPACE, HTTP_AUTH
		// and the TLS environment variables passed in via mixin in porter.yaml
		config := getConfig(&run)
		client, err := api.NewClient(config)
		if err != nil {
			return fmt.Errorf("unable to create nomad client: %w", err)
		}

		job, err := parseJob(&run)
		if err != nil {
			return fmt.Errorf("unable to parse job: %w", err)
		}

		// do the actual nomad job run
		if run.Path != "" {
			// regular run
			jobRegResp, _, err := client.Jobs().Register(job, nil)
			if err != nil {
				return fmt.Errorf("unable to register job: %w", err)
			}
			fmt.Fprintf(m.Out, "Job registration succesful\nEvaluation ID: %s", jobRegResp.EvalID) //todo
		} else if run.Dispatch != "" {
			// dispatch run
			jobDispResp, _, err := client.Jobs().Dispatch(run.Dispatch, run.Meta, []byte(run.Payload), run.IdPrefixTemplate, nil)
			if err != nil {
				return fmt.Errorf("unable to dispatch job: %w", err)
			}
			fmt.Fprintf(m.Out, "Job dispatched successfully\nDispatched Job ID: %s\nEvaluation ID: %s", jobDispResp.DispatchedJobID, jobDispResp.EvalID) //todo
		} else if run.Stop != "" {
			// stop run
			resp, _, err := client.Jobs().Deregister(run.Stop, run.Purge, nil)
			if err != nil {
				return fmt.Errorf("unable to stop job: %w", err)
			}
			fmt.Fprintf(m.Out, "Job stopped successfully\n%s", resp)
		} else {
			return errors.Errorf("unknown nomad run format, should either specify path, dispatch, stop")
		}
		fmt.Fprintf(m.Out, "\n")
	}
	return nil
}

func validateRun(run Job) error {
	// use below map as a set
	check := map[string]int{run.Path: 1, run.Dispatch: 1, run.Stop: 1}
	delete(check, "")
	if len(check) > 1 {
		return fmt.Errorf("unexpected nomad run format, expected only 1 out of path, dispatch, stop for a single run, but got %d", len(check))
	}
	return nil
}

func parseJob(run *Job) (*api.Job, error) {
	hclFile := run.Path
	job, err := jobspec.ParseFile(hclFile)
	return job, err
}

func getConfig(run *Job) *api.Config {
	config := api.DefaultConfig() // respects env variables

	// overwrite env variables if provided in each run
	if run.Address != "" {
		config.Address = run.Address
	}
	if run.Region != "" {
		config.Region = run.Region
	}
	if run.Namespace != "" {
		config.Namespace = run.Namespace
	}
	if run.CaCert != "" {
		config.TLSConfig.CACert = run.CaCert
	}
	if run.CaPath != "" {
		config.TLSConfig.CAPath = run.CaPath
	}
	if run.ClientCert != "" {
		config.TLSConfig.ClientCert = run.ClientCert
	}
	if run.ClientKey != "" {
		config.TLSConfig.ClientKey = run.ClientKey
	}
	if run.TlsServerName != "" {
		config.TLSConfig.TLSServerName = run.TlsServerName
	}
	if run.TlsSkipVerify {
		config.TLSConfig.Insecure = run.TlsSkipVerify
	}
	if run.Token != "" {
		config.SecretID = run.Token
	}
	return config
}
