#!/usr/bin/env bash

workingDirName=${PWD##*/}
projectName=${workingDirName}

#if [[ $1 != "" ]]; then
#    projectName=$1
#fi

echo "running tests for project: $projectName"

echo "mkdir -p /go/src/gitlab.zitcom.dk/smartweb/proj"
echo "mkdir -p /go/src/_/builds/smartweb/proj/$projectName"
echo "cp -r ${CI_PROJECT_DIR} /go/src/gitlab.zitcom.dk/smartweb/proj/$projectName"
echo "ln -s /go/src/gitlab.zitcom.dk/smartweb/proj/$projectName /go/src/_/builds/smartweb/proj/$projectName"
echo "cd /go/src/gitlab.zitcom.dk/smartweb/proj/$projectName"

mkdir -p /go/src/gitlab.zitcom.dk/smartweb/proj
mkdir -p /go/src/_/builds/smartweb/proj/php-dependency-checker
cp -r ${CI_PROJECT_DIR} /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker
ln -s /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker /go/src/_/builds/smartweb/proj/php-dependency-checker
cd /go/src/gitlab.zitcom.dk/smartweb/proj/php-dependency-checker

make dep
make test
