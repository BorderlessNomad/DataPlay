<?php
// This is a tool that is used to link what ckandownloader makes and build a index on those names,
// it is presuming that the tables have already been imported into the system just not linked
// into the index table.

// This file is written in PHP just to make things go a tad faster.

function What($input) {
	$realinput = $input;
	$input = str_replace("imp", "", $input);
	$input = str_replace("_", "", $input);
	$file = file_get_contents("./ckanindex.log");
	$file = explode("\n", $file);
	foreach ($file as $linenumber => $line) {
		$bits_of_line = explode(" ", $line);
		// Example
		// http://data.gov.uk/dataset/1-50-000-scale-gazetteer | http://www.ordnancesurvey.co.uk/oswebsite/products/50kgazetteer_open/ => ./data/85f998f715dc0e587880eefb5de3f87236add4c0_58beefbd5af5426de6c78c02e6663851de12866a
		$bitname = str_replace("./data/","",str_replace("_", "", $bits_of_line[4]));
		// echo(substr($bitname, 0, 59) . " vs $input\n");
		if($input == substr($bitname, 0, 59)) {
			echo ("Ooooh Ohhh found it!!!");
			$BrokenDataSetName = str_replace("http://data.gov.uk/dataset/", "http://data.gov.uk//dataset/", $bits_of_line[0]);
			$q2 = mysql_query("SELECT * FROM `index_old` WHERE ckan_url = '" . mysql_real_escape_string($BrokenDataSetName) . "' LIMIT 1");
			$d2 = mysql_fetch_array($q2);
			mysql_query("INSERT INTO `DataCon`.`index` (`GUID`, `Name`, `Title`, `Notes`) VALUES ('$input', '".mysql_real_escape_string($d2[1])."', '".mysql_real_escape_string($d2[2])."', '".mysql_real_escape_string($d2[3])."')");
			mysql_query("INSERT INTO `DataCon`.`priv_onlinedata` (`GUID`, `DatasetGUID`, `TableName`, `Defaults`) VALUES ('$input', '$input', '$realinput', '{}');");
			mysql_query("INSERT INTO `DataCon`.`priv_cache` (`cid`, `contents`) VALUES ('QCheck::$input', '".mysql_real_escape_string("{\"Amount\":3,\"Reques\":\"$input\"}")."');");
			// var_dump($bitname);
		}
		// echo("'$input' ".strlen($input)."\n");
	}
}

mysql_connect("10.0.0.2", "root", "");
mysql_select_db("DataCon");

$q = mysql_query("SHOW TABLE STATUS FROM `DataCon`");

while($row = mysql_fetch_array($q)) {
	$tablename = $row['Name'];
	// okay so we need to scan the table now.
	What($tablename);
}