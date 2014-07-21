#!/bin/bash

echo 'BUILDING DEPENDENCIES' &&
npm install -g &&
echo 'BUILDING JS/CSS' &&
grunt &&

if [ ! -f bin/public/lib/openlayers/build/OpenLayers.js ]; then
	mkdir -p bin/public/lib/dependencies/js/ &&
	pushd bin/public/lib/openlayers/build &&
	python build.py &&
	cp -r OpenLayers.js ../../dependencies/js/ &&
	popd
fi

echo 'BUILDING GOGRAM'
oldgo=$GOPATH
if [[ "$OSTYPE" == "msys" ]]; then
	GOPATH=$oldgo";"$(pwd -W)
else
	GOPATH=$oldgo:$(pwd)
fi
export GOPATH
project=dataplay
go get -v $project &&
go install -v $project &&
cd bin &&
./$project $@
export GOPATH=$oldgo
