#!/bin/bash

echo 'BUILDING DEPENDENCIES' &&
npm install &&
echo 'BUILDING JS/CSS' &&
node_modules/.bin/coffee -c -m -o public/js www-src/coffee &&
node_modules/.bin/lessc www-src/less/layout.less public/css/layout.css &&
node_modules/.bin/lessc www-src/less/signin.less public/css/signin.css &&
node_modules/.bin/lessc www-src/less/charts.less public/css/charts.css &&
node_modules/.bin/lessc www-src/less/maptest.less public/css/maptest.css &&

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
