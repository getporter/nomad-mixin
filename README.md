# Nomad Mixin for Porter

This is a [Porter](https://porter.sh) mixin for interacting with [Nomad](https://www.nomadproject.io/).

## Installation

todo

## Mixin Syntax

In your porter.yaml file, add `nomad` as a mixin:

```yaml
mixins:
  - nomad:
```

To use the mixin in an Install/Upgrade/Uninstall step, add the "jobs" block to your porter.yaml file (using install here
an example):

```yaml
install:
  - nomad:
      jobs:
        - path: nomad/pytechco-redis.nomad.hcl
        - path: nomad/pytechco-web.nomad.hcl
        - path: nomad/pytechco-setup.nomad.hcl
        - dispatch: pytechco-setup
          meta:
            budget: "200"
        - path: nomad/pytechco-employee.nomad.hcl
```

The job step supports the following 3 actions which can be specified as a list item in the `jobs:` block:

- `path`: runs a Nomad job file using the given Nomad hcl file.
- `dispatch`: dispatches a parametrized Nomad job using the given job name. Note that this requires a previously registered Nomad batch job.
    - `meta`: (optional): a map of key/value pairs that will be passed as metadata to the Nomad job.
    - `idPrefixTemplate` (optional): a prefix added to dispatched job IDs.
    - `payload` (optional): base64 encoded string containng the payload to pass to the Nomad job. Limited to 65536
      bytes.
- `stop`: stops a Nomad job using the given job name.
    - `purge`: (optional): if set to true, the job will be purged immediately.

Each of the 3 actions above can also take in the following options, which can be passed as arguments to each individual
Nomad run, overriding any global Nomad parameters set earlier:

- `address` (string): the address of the Nomad server to connect to.
- `region` (string): the region of the Nomad server to connect to.
- `namespace` (string): the namespace of the Nomad server to connect to.
- `caCert` (string): the path to the CA cert file to use for TLS.
- `caPath` (string): the path to a directory of CA cert files to use for TLS.
- `clientCert` (string): the path to the client cert file to use for TLS.
- `clientKey` (string): the path to the client key file to use for TLS.
- `tlsServerName` (string): the server name to use as the SNI host when connecting via TLS.
- `tlsSkipVerify` (bool): disables TLS host verification.
- `token` (string): the ACL token to use when connecting to Nomad.

If you don't want to specify these settings for each individual job, you can specify them globally in the porter.yaml
file as porter parameters, which will create corresponding environment variable inside the bundle. For example, the
sample-bundle in this repo specifies a porter parameter named `ip_address` which sets the `NOMAD_ADDR` environment variable. 
Thus, there is no need to specify the address in each individual Nomad job run, since the Nomad mixin will automatically 
use the `NOMAD_ADDR` if set. The mixin will automatically use the following environment variables if set:

- `NOMAD_ADDR`
- `NOMAD_REGION`
- `NOMAD_NAMESPACE`
- `NOMAD_HTTP_AUTH`
- `NOMAD_CACERT`
- `NOMAD_CAPATH`
- `NOMAD_CLIENT_CERT`
- `NOMAD_CLIENT_KEY`
- `NOMAD_TLS_SERVER_NAME`
- `NOMAD_SKIP_VERIFY`
- `NOMAD_TOKEN`

For more information about these environment variables, see the Nomad [documentation](https://developer.hashicorp.com/nomad/docs/commands#environment-variables).
Note that you can still override these environment variables for each individual job run by specifying the corresponding arguments in the run itself.

Additionally, the mixin also supports porter outputs. To specify an output for a job, add the `outputs` block to your
job, with a
list of `name/key` pairs to your run like so

```yaml
nomad:
  jobs:
    - path: nomad/pytechco-redis.nomad.hcl
      outputs:
        - name: redis-evalId
          key: evalId
```

You can also specify a top-level `outputs` block in the porter.yaml, to specify the outputs of the bundle as a whole

```yaml
outputs:
  - name: redis-evalId
    type: string
  - name: web-evalId
    type: string
```

This output can then be used in subsequent steps in your porter.yaml file using the syntax ${
bundle.outputs.redis-evalId }, see
the porter [documentation](https://porter.sh/wiring/#wiring-outputs) for more details.

Currently, the only supported output key is `evalId`, which is the ID of the Nomad evaluation that was created when the
job was run.
Note that not all Nomad jobs will create an evaluation, so this output will only be available for jobs that do.

## Example Bundle

See the [sample-bundle](./sample-bundle) directory for a bundle that uses the Nomad mixin to run the official
Nomad [tutorial](https://developer.hashicorp.com/nomad/tutorials/get-started/gs-deploy-job).
The install step runs the various jobs from the tutorial, and the uninstall step stops and purges the jobs. Remember to
provide the ip_address porter param when running the bundle (refer to the official Nomad tutorial for
how to start a development sever agent). The bundle also specifies 2 porter outputs which gets used with the exec mixin.