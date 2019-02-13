#!/usr/bin/env bash

workingDirName=${PWD##*/}
projectName=${workingDirName}

if [[ $1 != "" ]]; then
    projectName=$1
fi

echo "building project: $projectName"

mkdir -p /go/src/gitlab.zitcom.dk/smartweb/proj
mkdir -p /go/src/_/builds/smartweb/proj/${projectName}
cp -r ${CI_PROJECT_DIR} /go/src/gitlab.zitcom.dk/smartweb/proj/${projectName}
ln -s /go/src/gitlab.zitcom.dk/smartweb/proj/${projectName} /go/src/_/builds/smartweb/proj/${projectName}
cd /go/src/gitlab.zitcom.dk/smartweb/proj/${projectName}

make build
mv /go/src/gitlab.zitcom.dk/smartweb/proj/${projectName}/bin/${projectName} ${CI_PROJECT_DIR}
