# Nomad Mixin for Porter

This is a [Porter](https://porter.sh) mixin for interacting with [Nomad](https://www.nomadproject.io/).

## Installation

todo

## Mixin Syntax

In your porter.yaml file, add the nomad as a mixin:

```yaml
mixins:
  - nomad:
```

It supports the following global settings, which will apply to all your nomad actions, unless overridden in the action
itself:

```yaml
mixins:
  - nomad:
      NOMAD_ADDR: ""
      NOMAD_REGION: ""
      NOMAD_NAMESPACE: ""
      NOMAD_HTTP_AUTH: ""

      # tls
      NOMAD_CACERT: ""
      NOMAD_CAPATH: ""
      NOMAD_CLIENT_CERT: ""
      NOMAD_CLIENT_KEY: ""
      NOMAD_TLS_SERVER_NAME: ""
      NOMAD_SKIP_VERIFY: ""

      # acl
      NOMAD_TOKEN: ""
```

To use the mixin in an Install/Upgrade/Uninstall step, add the "nomad:jobs" to your porter.yaml file (using install here
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

The job step supports the following 3 actions which can be specifed as a list item in the `jobs:` list:

- `path`: runs a nomad job file using the given nomad hcl file.
- `dispatch`: dispatches a parametrized nomad job using the given job name.
    - `meta`: (optional): a map of key/value pairs that will be passed as metadata to the nomad job.
    - `idPrefixTemplate` (optional): a prefix added to dispathd job IDs.
    - `payload` (optional): base64 encoded string containg the payload to pass to the nomad job. Limited to 65536 bytes.
- `stop`: stops a nomad job using the given job name.
    - `purge`: (optional): a boolean value that determines whether to purge the job immediately.

Each of the 3 actions above can also take in the following options, which can be passed as arguments to each individual
nomad run, overriding any global nomad parameters set earlier:

- `address` (string): the address of the nomad server to connect to.
- `region` (string): the region of the nomad server to connect to.
- `namespace` (string): the namespace of the nomad server to connect to.
- `caCert` (string): the path to the CA cert file to use for TLS.
- `caPath` (string): the path to a directory of CA cert files to use for TLS.
- `clientCert` (string): the path to the client cert file to use for TLS.
- `clientKey` (string): the path to the client key file to use for TLS.
- `tlsServerName` (string): the server name to use as the SNI host when connecting via TLS.
- `skipVerify` (bool): disables TLS host verification.
- `token` (string): the ACL token to use when connecting to Nomad.

## Example Bundle

See the [sample-bundle](./sample-bundle) directory for a bundle that uses the nomad mixin to run the official
nomad [tutorial](https://developer.hashicorp.com/nomad/tutorials/get-started/gs-deploy-job).
The install step runs the various jobs from the tutorial, and the uninstall step stops and purges the jobs. Remember to
set NOMAD_ADDR to the address of your nomad server before running the bundle (refer to the official nomad tutorial for 
how to start a development sever agent).
