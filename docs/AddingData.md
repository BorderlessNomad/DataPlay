Adding Datasets
===

There are two entries that you need to fill in for a dataplay dataset to show up. Other than that dataplay datasets are just
plain and simple mysql tables.

##Getting our bearings.

Lets first have a look at the first table we will need to insert stuff into.

```
mysql> SELECT * FROM `index` LIMIT 1;
+------+--------+--------+--------------+----------+-------+
| GUID | Name   | Title  | Notes        | ckan_url | Owner |
+------+--------+--------+--------------+----------+-------+
| gdp  | UK GDP | UK GDP | <h1>GDP</h1> |          |     0 |
+------+--------+--------+--------------+----------+-------+
1 row in set (0.00 sec)
```

and the 2nd table.

```
mysql> SELECT * FROM priv_onlinedata LIMIT 1;
+------+-------------+-----------+------------------------+
| GUID | DatasetGUID | TableName | Defaults               |
+------+-------------+-----------+------------------------+
| gdp  | gdp         | gdp       | {"x":"year","y":"GDP"} |
+------+-------------+-----------+------------------------+
1 row in set (0.00 sec)
```

##Make the table.

This one really does not matter too much how you do it, if you are importing a CSV then there are tools that will turn CSV's into mysql tables.

For this walk though we will make `testtable`.

```
mysql> CREATE TABLE `testtable` (
    -> `year` INT NOT NULL,
    -> `dollers` INT NULL,
    -> `happyness` INT NULL,
    -> PRIMARY KEY (`year`)
    -> )
    -> ENGINE=InnoDB;
Query OK, 0 rows affected (0.01 sec)
```

And then insert some data into it.

```

mysql> INSERT INTO `DataCon`.`testtable` (`year`, `dollers`, `happyness`) VALUES (1999, 8949, 76);
Query OK, 1 row affected (0.00 sec)

```

Then we will inset the needed bits into both tables that store where to find datasets as seen above.

```

mysql> INSERT INTO `DataCon`.`priv_onlinedata` (`GUID`, `DatasetGUID`, `TableName`, `Defaults`) VALUES ('test', 'test', 'testtable', '{}');
Query OK, 1 row affected (0.00 sec)

mysql> INSERT INTO `DataCon`.`index` (`GUID`, `Name`, `Title`, `Notes`) VALUES ('test', 'test', 'A test table', '-');
Query OK, 1 row affected, 1 warning (0.00 sec)

```

Recapping the `testtable` now looks like this:

```

mysql> SELECT * FROM testtable;
+------+---------+-----------+
| year | dollers | happyness |
+------+---------+-----------+
| 1995 |     100 |        99 |
| 1996 |     150 |        99 |
| 1997 |     500 |        99 |
| 1998 |    1000 |        76 |
| 1999 |    8949 |        76 |
+------+---------+-----------+
5 rows in set (0.00 sec)

```

