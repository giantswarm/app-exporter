#!/bin/bash

curl -L https://github.com/giantswarm/apptestctl/releases/download/v0.3.0/apptestctl-v0.3.0-linux-amd64.tar.gz > /tmp/apptestctl.tar.gz
cd /tmp
tar xzvf apptestctl.tar.gz
chmod u+x /tmp/apptestctl-v0.3.0-linux-amd64/apptestctl
sudo mv /tmp/apptestctl-v0.3.0-linux-amd64/apptestctl /usr/local/bin

apptestctl bootstrap --kubeconfig="$(kind get kubeconfig)"
