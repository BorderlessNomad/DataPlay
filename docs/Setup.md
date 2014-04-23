Setup
===

To setup your very own dataplay. You will need the following:

	* Go 1.2 (Higher versions *proablly* work)
	* MySQL 5.5.30 (Or somthing compatible)
	* NPM (You can get it from installing node.js)
	* Redis 2.8.6

You can start by running `layout.sql` from the root directory into a new database.

Then start a Redis server up and note down the host name.

Then goto the root of the dataplay directory and run "npm install" This will install a bunch of stuff needed to build the front end of dataplay.

then run `go get` in the src directory to fetch the deps.

then run start.sh.

In the case that you cannot connect to databases you can set the database name though a ENV var called `DATABASE` and it should be in the format of `1.1.1.1:123`
like wise with the redis host can be changed with `redishost` env var.

You can then start loading in your own mysql tables (use just varchar, int and float for now) and then add them into both `priv_index` and `priv_onlinedata` to make sure they are searchable. there are many tools in the `tool` folder to help you out with lots of things like data filitering.