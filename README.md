![](http://i.imgur.com/esjTHFE.png)

DataPlay is an open-source data analysis and exploration game developed by [PlayGen](http://playgen.com/) as part of the European [CELAR](http://celarcloud.eu) elastic cloud computing initiative.

The aim of DataPlay, besides taking the CELAR architecture for a spin, is to provide a collaborative online environment in which non-expert users can play with government data, enabling people to quickly discover patterns within it. The system presents many different elements of government data to users in a range of visual forms. Users are then encouraged to explore and comment upon this data. The system also seeks out potential correlations between disparate datasets to show to the user, in order to help them discover hidden meaning witihin the data.

The back end of the application is written in [Golang](http://golang.org/) to provide concurrency when dealing with large volume data processing. It also utilises [Martini](https://github.com/go-martini/martini) for web routing, [PostGreSQL](http://www.postgresql.org/) coupled with [Gorm](https://github.com/jinzhu/gorm) for processing the government data, and [Cassandra](http://cassandra.apache.org/) in conjunction with [gocql](https://github.com/gocql/gocql) for handling the web data.

The front end is written in [Coffeescript](http://coffeescript.org/) and [AngularJS](https://angularjs.org/) and utilises the [d3js](http://d3js.org/), [dc.js](http://dc-js.github.io/dc.js/) and [NVD3](http://nvd3.org/) charting packages.

DataPlay alpha contains a rudimentary selection of government datasets drawn from [data.gov.uk](http://data.gov.uk/), along with political information taken from the [BBC](http://www.bbc.co.uk/news/) news archive, which was extracted and analysed via [python](https://www.python.org/) scripted [import.io](https://import.io/) in conjunction with golang implemented [embed.ly](http://embed.ly/).


#Screenshots

Landing Page
![](http://i.imgur.com/yJyJ4GC.png)

Home Page
![](http://i.imgur.com/2vkyTVS.png)

Overview Screen
![](http://i.imgur.com/N4kCiPG.png)

Search Page
![](http://i.imgur.com/1ZYsaQb.png)

Chart Page
![](http://i.imgur.com/cEakHPq.png)