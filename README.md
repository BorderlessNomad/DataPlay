![](http://i.imgur.com/esjTHFE.png)

DataPlay is an open-source data analysis and exploration game developed by [PlayGen](http://playgen.com/) as part of the European [CELAR](http://celarcloud.eu) initiative.

The aim of DataPlay, besides taking CELAR for a spin, is to provide a collaborative environment in which non-expert users can play with government data. The system presents different elements of the data to users in a range of visual forms. People are then encouraged to explore this data together. The system also seeks to identify potential correlations between disparate datasets, in order to help people discover hidden patterns.

The back end is written in [Go](http://golang.org/) to provide concurrency for large volume data processing. It also uses [Martini](https://github.com/go-martini/martini) for web routing, [PostgreSQL](http://www.postgresql.org/) with [GORM](https://github.com/jinzhu/gorm) for processing the government data, and [Cassandra](http://cassandra.apache.org/) with [gocql](https://github.com/gocql/gocql) for handling the web data.

The front end is written in [CoffeeScript](http://coffeescript.org/) and [AngularJS](https://angularjs.org/) and utilises the [d3.js](http://d3js.org/), [dc.js](http://dc-js.github.io/dc.js/) and [NVD3.js](http://nvd3.org/) charting packages.

DataPlay alpha contains a rudimentary selection of datasets drawn from [data.gov.uk](http://data.gov.uk/), along with political information taken from the [BBC](http://www.bbc.co.uk/news/), which was extracted and analysed via [python](https://www.python.org/) scripted [import.io](https://import.io/) and Go implemented [embed.ly](http://embed.ly/).

### Landing Page
![](http://i.imgur.com/yJyJ4GC.png)

### Home Page
![](http://i.imgur.com/2vkyTVS.png)

### Overview Screen
![](http://i.imgur.com/N4kCiPG.png)

### Search Page
![](http://i.imgur.com/1ZYsaQb.png)

### Chart Page
![](http://i.imgur.com/cEakHPq.png)

## Installation

TODO: Describe the installation process

## Usage

TODO: Write usage instructions

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## History

TODO: Write history

## Credits

TODO: Write credits

## License

TODO: Write license