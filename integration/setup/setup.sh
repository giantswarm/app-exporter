#!/bin/bash

go get github.com/giantswarm/apptestctl@v0.1.0
apptestctl bootstrap --kubeconfig="$(kind get kubeconfig)"
