#!/usr/bin/env bash

if [[ -z "$(which protoc)" || "$(protoc --version)" != "libprotoc 3.15."* ]]; then
  echo "Generating protobuf requires protoc 3.15.x. Please download and"
  echo "install the platform appropriate Protobuf package for your OS: "
  echo
  echo "  https://github.com/google/protobuf/releases"
  echo
  echo "WARNING: Protobuf changes are not being validated"
  exit 1
fi

set -e
TMP_DIR=$(mktemp -d)
mkdir -p "${TMP_DIR}"/src/github.com/openkruise/kruise
cp -r ./{apis,hack,vendor} "${TMP_DIR}"/src/github.com/openkruise/kruise/

(cd "${TMP_DIR}"/src/github.com/openkruise/kruise; \
    GOPATH=${TMP_DIR} go install github.com/openkruise/kruise/vendor/k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo; \
    PATH=${TMP_DIR}/bin:$PATH GOPATH=${TMP_DIR} protoc --gogo_out=plugins=grpc:. apis/ctrlmesh/proto/ctrlmesh.proto)

cp -f "${TMP_DIR}"/src/github.com/openkruise/kruise/apis/ctrlmesh/proto/*.go apis/ctrlmesh/proto/
