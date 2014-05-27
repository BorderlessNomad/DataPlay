#!/bin/bash

echo 'BUILDING DEPENDENCIES' &&
npm install &&
echo 'BUILDING JS/CSS' &&
grunt &&

if [ ! -f public/lib/openlayers/build/OpenLayers.js ]; then
	mkdir -p public/lib/dependencies/js/ &&
	pushd public/lib/openlayers/build &&
	python build.py &&
	cp -r OpenLayers.js ../../dependencies/js/ &&
	popd
fi

echo 'BUILDING GOGRAM'
oldgo=$GOPATH
export GOPATH=$(pwd)
go get dataplay &&
go install dataplay &&
(./bin/dataplay || ./bin/dataplay.exe)
export GOPATH=$oldgo
