echo 'BUILDING JS/CSS'
node_modules/.bin/coffee -c -o public/js src/coffee
node_modules/.bin/lessc src/less/layout.less public/css/layout.css
node_modules/.bin/lessc src/less/signin.less public/css/signin.css
node_modules/.bin/lessc src/less/charts.less public/css/charts.css
node_modules/.bin/lessc src/less/maptest.less public/css/maptest.css
echo 'BUILDING GOGRAM'
go build
./datacon

