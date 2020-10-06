#!/bin/bash

curl -L https://github.com/giantswarm/apptestctl/releases/download/v0.2.0/apptestctl-v0.2.0-linux-amd64.tar.gz > ./apptestctl.tar.gz
tar xzvf apptestctl.tar.gz
chmod u+x apptestctl-v0.2.0-darwin-amd64/apptestctl
mv apptestctl-v0.2.0-darwin-amd64/apptestctl .

./apptestctl bootstrap --kubeconfig="$(kind get kubeconfig)"
