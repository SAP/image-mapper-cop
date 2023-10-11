#!/usr/bin/env bash

set -eo pipefail

export GOROOT=$(go env GOROOT)

BASEDIR=$(realpath $(dirname "$0")/..)
TEMPDIR=$BASEDIR/tmp/gen
trap 'rm -rf "$TEMPDIR"' EXIT
mkdir -p "$TEMPDIR"

mkdir -p "$TEMPDIR"/apis/operator.kyma-project.io
ln -s "$BASEDIR"/api/v1alpha1 "$TEMPDIR"/apis/operator.kyma-project.io/v1alpha1

"$BASEDIR"/bin/client-gen \
  --clientset-name versioned \
  --input-base "" \
  --input github.com/sap/image-mapper-cop/tmp/gen/apis/operator.kyma-project.io/v1alpha1 \
  --go-header-file "$BASEDIR"/hack/boilerplate.go.txt \
  --output-package github.com/sap/image-mapper-cop/pkg/client/clientset \
  --output-base "$TEMPDIR"/pkg/client \
  --plural-exceptions ImageMapper:imagemappers

"$BASEDIR"/bin/lister-gen \
  --input-dirs github.com/sap/image-mapper-cop/tmp/gen/apis/operator.kyma-project.io/v1alpha1 \
  --go-header-file "$BASEDIR"/hack/boilerplate.go.txt \
  --output-package github.com/sap/image-mapper-cop/pkg/client/listers \
  --output-base "$TEMPDIR"/pkg/client \
  --plural-exceptions ImageMapper:imagemappers

"$BASEDIR"/bin/informer-gen \
  --input-dirs github.com/sap/image-mapper-cop/tmp/gen/apis/operator.kyma-project.io/v1alpha1 \
  --versioned-clientset-package github.com/sap/image-mapper-cop/pkg/client/clientset/versioned \
  --listers-package github.com/sap/image-mapper-cop/pkg/client/listers \
  --go-header-file "$BASEDIR"/hack/boilerplate.go.txt \
  --output-package github.com/sap/image-mapper-cop/pkg/client/informers \
  --output-base "$TEMPDIR"/pkg/client \
  --plural-exceptions ImageMapper:imagemappers

find "$TEMPDIR"/pkg/client -name "*.go" -exec \
  perl -pi -e "s#github\.com/sap/image-mapper-cop/tmp/gen/apis/operator\.kyma-project\.io/v1alpha1#github.com/sap/image-mapper-cop/api/v1alpha1#g" \
  {} +

rm -rf "$BASEDIR"/pkg/client
mv "$TEMPDIR"/pkg/client/github.com/sap/image-mapper-cop/pkg/client "$BASEDIR"/pkg

cd "$BASEDIR"
go fmt ./pkg/client/...
go vet ./pkg/client/...