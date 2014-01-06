@echo off
echo BUILDING DEPENDENCIES
call npm install
echo BUILDING JS / CSS
call node_modules\.bin\coffee -c -o public\js src\coffee && call node_modules\.bin\lessc src\less\layout.less public\css\layout.css && call node_modules\.bin\lessc src\less\maptest.less public\css\maptest.css && call node_modules\.bin\lessc src\less\charts.less public\css\charts.css && call node_modules\.bin\lessc src\less\signin.less public\css\signin.css
echo BUILDING GOGRAM
go build && datacon