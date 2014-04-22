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

then run start.sh.

In the case that you cannot connect to databases you can set the database name though a ENV var called `DATABASE` and it should be in the format of `1.1.1.1:123`