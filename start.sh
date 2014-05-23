#!/bin/bash

echo 'BUILDING DEPENDENCIES' &&
npm install &&
npm install -g grunt-cli
echo 'BUILDING JS/CSS' &&
grunt &&

if [ ! -f public/lib/openlayers/build/OpenLayers.js ]; then
	cd public/lib/openlayers/build &&
	python build.py &&
	mkdir ../../../../public/ &&
	mkdir ../../../../public/lib/ &&
	mkdir ../../../../public/lib/dependencies/ &&
	mkdir ../../../../public/lib/dependencies/js/ &&
	cp -r OpenLayers.js ../../../../public/lib/dependencies/js/
	cd ../../../../
fi

echo 'BUILDING GOGRAM' &&
go build -o datacon &&
./datacon
