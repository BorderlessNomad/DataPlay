![](http://i.imgur.com/l7RwTtQ.png)

### Overview
DataPlay is an open-source data analysis and exploration game developed by [PlayGen](http://playgen.com/) as part of the EU's [CELAR](http://celarcloud.eu) initiative.

The aim of DataPlay, besides taking CELAR for a spin, is to provide a collaborative environment in which non-expert users get to "play" with government data. The system presents the user with a range of elements of the data, displayed in a variety of visual forms. People are then encouraged to explore this data together. The system also seeks to identify potential correlations between disparate datasets, in order to help users discover hidden patterns within the data.

### Architecture
The back end is written in [Go](http://golang.org/), to provide concurrency for large volume data processing. There is a multiple master/frontend architecture which relies on [HAProxy](http://www.haproxy.org/) for its Load-balancing capabilities. The backend also utilises [Martini](https://github.com/go-martini/martini) for parametric API routing, number of [PostgreSQL](http://www.postgresql.org/) replicated and load balanced using [pgpool-II](http://www.pgpool.net/mediawiki/index.php/Main_Page) with [GORM](https://github.com/jinzhu/gorm) for facilitating communication between back end and database, [Cassandra](http://cassandra.apache.org/) coupled with [gocql](https://github.com/gocql/gocql) for data obtained via scraping of 3rd party news sources. [Redis](http://redis.io/) for storing monitoring and session related data.

The front end is written in [CoffeeScript](http://coffeescript.org/) on top of [AngularJS](https://angularjs.org/) and makes extensive use of the libraries such as [D3.js](http://d3js.org/), [dc.js](http://dc-js.github.io/dc.js/) and [NVD3.js](http://nvd3.org/) for presenting data in the form of various charts. The user interface is created using [Bootstrap](http://getbootstrap.com/), [Bootswatch](https://bootswatch.com/) and [Font Awesome](http://fontawesome.io/).

DataPlay (Beta) contains a rudimentary selection of datasets drawn from [DATA.GOV.UK](http://data.gov.uk/) & [LONDON DATASTORE](http://data.london.gov.uk/), along with political information taken from the [BBC](http://www.bbc.co.uk/news/), which was extracted and analysed via [import.io](https://import.io/), [kimono](https://www.kimonolabs.com/) and [embed.ly](http://embed.ly/).

##Screens
### Landing Page
![](http://i.imgur.com/lMcdMYL.png)

### Home Page
![](http://i.imgur.com/81T9n1k.png)

### Activity Monitor
![](http://i.imgur.com/U1gm66j.png)

### Search Page
![](http://i.imgur.com/RbI9YbX.png)

### Chart Page
![](http://i.imgur.com/LP5c5C3.png)

## Installation

1. Install Ubuntu & Node.js
2. Install all necessary dependencies `npm install`

*Note*: Refer [`tools/deployment/base.sh`](tools/deployment/base.sh) for base system config and libs.

### Production:

1. HAProxy Load Balancer [`tools/deployment/loadbalancer/haproxy.sh`](tools/deployment/loadbalancer/haproxy.sh)
2. Gamification instances [`tools/deployment/app/frontend.sh`](tools/deployment/app/frontend.sh)
3. Computation/API instances [`tools/deployment/app/master.sh`](tools/deployment/app/master.sh)
4. PostgreSQL DB instance [`tools/deployment/db/postgresql.sh`](tools/deployment/db/postgresql.sh)
5. Cassandra DB instance [`tools/deployment/db/cassandra.sh`](tools/deployment/db/cassandra.sh)
6. Redis instance [`tools/deployment/db/redis.sh`](tools/deployment/db/redis.sh)

### Monitoring:

1. API response time monitoring [`tools/deployment/monitoring/api.sh`](tools/deployment/monitoring/api.sh)
2. HAProxy API for dynamic scaling [`tools/deployment/loadbalancer/api/`](tools/deployment/loadbalancer/api/)

## Usage

### Development:

1. Run back end & API server using `./start.sh`
2. Install PostgreSQL and import data
3. Run front end `cd www-src && npm install && grunt serve`

### Staging:

1. Run back end & API server using `./start.sh`
2. Install PostgreSQL and import data
3. Deploy & run front end in `cd www-src && npm install && grunt serve:dist`

### Production:

1. Deploy HAProxy Server & DataPlay HAProxy API (written in Node.js)
2. Deploy number of required master nodes
	a. Run back end and API server on each master nodes using `./start.sh`
	b. Send add `master` node requests to DataPlay HAProxy API via cURL
3. Deploy number of required master nodes
	a. Install Nginx to serve data & set appropriate path for `www-src` directory
	b. Send add `gamification` node requests to DataPlay HAProxy API via cURL
4. Install pgpool-II on CentOS server (for best compatibility) & DataPlay PGPOOL API (written in Node.js)
5. Deploy number of required PostgreSQL nodes
	a. Install PostgreSQL along with pgpool-II client plugin
	b. Send add node request to DataPlay PGPOOL API via cURL
6. Create an A Record for required domain and point it to HAProxy Sever IP

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## History

v1.1.1: 	API metric probes are now on each master nodes.

v1.1.0: 	Support for PostgreSQL cluster using pgpool-II

v1.0.2: 	Added API health monitoring probes.

v1.0.1: 	System-wide environment variables

v1.0.0: 	CELAR compatible scripts

v0.9.1: 	Deployment Scripts

v0.9.0: 	Update to AngularJS v1.3

v0.8.9: 	Use NVD3 for Correlated Charts

v0.8.5: 	Added Correlated Charts

v0.8.0: 	Added Related Charts

v0.7.5: 	Added DC.js for Charts

v0.7.2: 	Use Bootstrap for HTML

v0.7.0: 	Migrated frontend to AngularJS

v0.6.1: 	Moved Sessions to Redis

v0.6.0: 	Use GORM for PostgreSQL

v0.5.0: 	Go 1.1 with embedded HTML views

## Authors

Mayur Ahir [mayur@playgen.com]

Jack Cannon [jack@playgen.com]

Lex Robinson [lex@playgen.com]

## License

**Copyright (C) 2014 PlayGen LTD**

This program is free software: you can redistribute it and/or modify it under the terms of the **GNU General Public License version 3.0** as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the [GNU General Public License](LICENSE.md) for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see [LICENSE](LICENSE.md).
