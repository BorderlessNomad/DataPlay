call node_modules\.bin\coffee -c -o public\js src\coffee
call node_modules\.bin\lessc src\less\layout.less public\css\layout.css
call node_modules\.bin\lessc src\less\charts.less public\css\charts.css
call node_modules\.bin\lessc src\less\signin.less public\css\signin.css
call go build
call datacon
