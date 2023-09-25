package nomad

import (
	"fmt"
	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/pkg/errors"
)

type Nomad struct {
	Runs []Run `yaml:"runs"`
}

type Run struct {
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

	Token string `yaml:"token"`
}

func (m *Mixin) DoAction(action *Nomad) error {
	for _, run := range action.Runs {
		if len(map[interface{}]int{run.Path: 1, run.Dispatch: 1, run.Stop: 1}) > 1 {
			return fmt.Errorf("unexpected nomad run format, expected only one out of path, dispatch, stop for a single run")
		}
		// config will respect user global ADDR, REGION, NAMESPACE, HTTP_AUTH
		// and the mTLS environment variables passed in via porter.yaml
		config := getConfig(&run)
		client, err := api.NewClient(config)
		if err != nil {
			return fmt.Errorf("unable to create nomad client: %w", err)
		}
		job, err := parseJob(&run)
		if err != nil {
			return fmt.Errorf("unable to parse job: %w", err)
		}
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

func parseJob(run *Run) (*api.Job, error) {
	hclFile := run.Path
	job, err := jobspec.ParseFile(hclFile)
	return job, err
}

func getConfig(run *Run) *api.Config {
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

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
}

// MarshalYAML converts the action back to a YAML representation
// install:
//
//	skeletor:
//	  ...
func (a Action) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{a.Name: a.Steps}, nil
}

// MakeSteps builds a slice of Step for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - skeletor: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	// Go doesn't have generics, nothing to see here...
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

type Step struct {
	Instruction `yaml:"nomad"`
}

// Actions is a set of actions, and the steps, passed from Porter.
type Actions []Action

// UnmarshalYAML takes chunks of a porter.yaml file associated with this mixin
// and populates it on the current action set.
// install:
//
//	skeletor:
//	  ...
//	skeletor:
//	  ...
//
// upgrade:
//
//	skeletor:
//	  ...
func (a *Actions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, Action{})
	if err != nil {
		return err
	}

	for actionName, action := range results {
		for _, result := range action {
			s := result.(*[]Step)
			*a = append(*a, Action{
				Name:  actionName,
				Steps: *s,
			})
		}
	}
	return nil
}

var _ builder.HasOrderedArguments = Instruction{}
var _ builder.ExecutableStep = Instruction{}
var _ builder.StepWithOutputs = Instruction{}

type Instruction struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	WorkingDir  string   `yaml:"dir,omitempty"`
	Arguments   []string `yaml:"arguments,omitempty"`

	// Useful when the CLI you are calling wants some arguments to come after flags
	// Arguments are passed first, then Flags, then SuffixArguments.
	SuffixArguments []string `yaml:"suffix-arguments,omitempty"`

	Flags          builder.Flags `yaml:"flags,omitempty"`
	Outputs        []Output      `yaml:"outputs,omitempty"`
	SuppressOutput bool          `yaml:"suppress-output,omitempty"`

	// Allow the user to ignore some errors
	// Adds the ignoreError functionality from the exec mixin
	// https://release-v1.porter.sh/mixins/exec/#ignore-error
	builder.IgnoreErrorHandler `yaml:"ignoreError,omitempty"`
}

func (s Instruction) GetCommand() string {
	return "nomad"
}

func (s Instruction) GetWorkingDir() string {
	return s.WorkingDir
}

func (s Instruction) GetArguments() []string {
	return s.Arguments
}

func (s Instruction) GetSuffixArguments() []string {
	return s.SuffixArguments
}

func (s Instruction) GetFlags() builder.Flags {
	return s.Flags
}

func (s Instruction) SuppressesOutput() bool {
	return s.SuppressOutput
}

func (s Instruction) GetOutputs() []builder.Output {
	// Go doesn't have generics, nothing to see here...
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}

var _ builder.OutputJsonPath = Output{}
var _ builder.OutputFile = Output{}
var _ builder.OutputRegex = Output{}

type Output struct {
	Name string `yaml:"name"`

	// See https://porter.sh/mixins/exec/#outputs
	// TODO: If your mixin doesn't support these output types, you can remove these and the interface assertions above, and from #/definitions/outputs in schema.json
	JsonPath string `yaml:"jsonPath,omitempty"`
	FilePath string `yaml:"path,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}

func (o Output) GetFilePath() string {
	return o.FilePath
}

func (o Output) GetRegex() string {
	return o.Regex
}
