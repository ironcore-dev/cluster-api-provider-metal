#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

./bin/kustomize build $1 | ./bin/envsubst
