#!/usr/bin/env bash

mkdir -p /go/src/gitlab.zitcom.dk/smartweb/proj
mkdir -p /go/src/_/builds/smartweb/proj/php-dependency-checker
cp -r ${CI_PROJECT_DIR} /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker
ln -s /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker /go/src/_/builds/smartweb/proj/php-dependency-checker
cd /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker
make build
mv /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/bin/php-dependency-checker ${CI_PROJECT_DIR}
