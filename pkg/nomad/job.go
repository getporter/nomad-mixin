package nomad

import (
	"fmt"
	"os"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec2"
)

type Nomad struct {
	Jobs []Job `yaml:"jobs"`
}

type NomadOutput struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type Job struct {
	// outputs to be used in next steps of porter bundle
	Outputs []NomadOutput `yaml:"outputs"`

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
		err := validateJob(run)
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

		evalId := ""
		// we are guaranteed to enter one of the below branches since we validated the job
		if run.Path != "" {
			// service job
			evalId, err = serviceJob(client, &run, m)
		} else if run.Dispatch != "" {
			// dispatch run
			evalId, err = dispatchJob(client, &run, m)
		} else if run.Stop != "" {
			// stop run
			evalId, err = stopJob(client, &run, m)
		}
		if err != nil {
			return err
		}

		err = handleOutputs(&run, m, evalId)
		if err != nil {
			return err
		}
		fmt.Fprintf(m.Out, "\n")
	}
	return nil
}

// handleOutputs writes the evalId to the outputs
func handleOutputs(run *Job, m *Mixin, evalId string) error {
	for _, output := range run.Outputs {
		err := m.WriteMixinOutputToFile(output.Name, []byte(evalId))
		if err != nil {
			return err
		}
	}
	return nil
}

// stopJob stops a nomad job
func stopJob(client *api.Client, run *Job, m *Mixin) (string, error) {
	evalId, _, err := client.Jobs().Deregister(run.Stop, run.Purge, nil)
	if err != nil {
		return "", fmt.Errorf("unable to stop job: %w", err)
	}
	fmt.Fprintf(m.Out, "Job stopped successfully\n")
	if evalId != "" {
		fmt.Fprintf(m.Out, "Job stop response: %s\n", evalId)
	}
	return evalId, nil
}

// dispatchJob dispatches a nomad job
func dispatchJob(client *api.Client, run *Job, m *Mixin) (string, error) {
	jobDispResp, _, err := client.Jobs().Dispatch(run.Dispatch, run.Meta, []byte(run.Payload), run.IdPrefixTemplate, nil)
	if err != nil {
		return "", fmt.Errorf("unable to dispatch job: %w", err)
	}
	fmt.Fprintf(m.Out, "Job dispatched successfully\n")
	if jobDispResp.EvalID != "" {
		fmt.Fprintf(m.Out, "Job evaluation ID: %s\n", jobDispResp.EvalID)
	}
	if jobDispResp.DispatchedJobID != "" {
		fmt.Fprintf(m.Out, "Job dispatched ID: %s\n", jobDispResp.DispatchedJobID)
	}
	return jobDispResp.EvalID, nil
}

// serviceJob registers a job with nomad
func serviceJob(client *api.Client, run *Job, m *Mixin) (string, error) {
	job, err := parseJob(run.Path)
	if err != nil {
		return "", fmt.Errorf("unable to parse job: %w", err)
	}
	jobRegResp, _, err := client.Jobs().Register(job, nil)
	if err != nil {
		return "", fmt.Errorf("unable to register job: %w", err)
	}
	fmt.Fprintf(m.Out, "Job registration succesful\n")
	if jobRegResp.EvalID != "" {
		fmt.Fprintf(m.Out, "Job evaluation ID: %s\n", jobRegResp.EvalID)
	}
	if jobRegResp.Warnings != "" {
		fmt.Fprintf(m.Out, "Job registration warnings: %s\n", jobRegResp.Warnings)
	}
	return jobRegResp.EvalID, nil
}

// validateJob checks that exactly one of path, dispatch or stop is specified
func validateJob(run Job) error {
	// use below map as a set
	check := map[string]int{run.Path: 1, run.Dispatch: 1, run.Stop: 1}
	delete(check, "")
	if len(check) != 1 {
		return fmt.Errorf("unexpected nomad run format, expected exactly 1 out of path, dispatch, stop for a single run, but got %d", len(check))
	}
	return nil
}

// parseJob parses a nomad job file
func parseJob(filepath string) (*api.Job, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	job, err := jobspec2.Parse(filepath, f)
	return job, err
}

// getConfig returns a nomad client config with the user provided env variables
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
