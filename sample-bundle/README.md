# Sample Nomad Bundle

This directory contains an example bundle that uses the nomad mixin. To use this bundle first you need a nomad cluster
up and running. To get a quick one up and running locally,
you can run (you might need to run as root depending on your environment):

```bash
nomad agent -dev \
  -bind 0.0.0.0 \
  -network-interface='{{ GetDefaultInterfaces | attr "name" }}'
```

Make note of the IP address of the nomad server, it should get output in the terminal after running the command.
This bundle declares a required parameter `ip_address` in the porter.yaml. There are many ways to pass a parameter to a
bundle, but the simplest is just specify it when running the bundle:

Run `porter install --param ip_address=<YOUR_CLUSTER_IP>` in this directory containing the porter.yaml. This will
install the nomad jobs in the nomad/ subdirectory into your cluster.
You can run `porter uninstall` to remove the jobs from the cluster. The parameter is saved from the previous run, so it's
not necessary to specify it again when running porter uninstall.