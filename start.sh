#!/bin/bash

echo 'BUILDING DEPENDENCIES'
npm install
echo 'BUILDING JS/CSS'
node_modules/.bin/coffee -c -o public/js src/coffee
node_modules/.bin/lessc src/less/layout.less public/css/layout.css
node_modules/.bin/lessc src/less/signin.less public/css/signin.css
node_modules/.bin/lessc src/less/charts.less public/css/charts.css
node_modules/.bin/lessc src/less/maptest.less public/css/maptest.css
cp public/lib/openlayers/lib/OpenLayers.js public/lib/dependencies/js/OpenLayers.js
echo 'BUILDING GOGRAM'
go build
./datacon
