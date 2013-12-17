<?php
// This file is written in PHP just to make things go a tad faster.
mysql_connect("10.0.0.2", "root", "");
mysql_select_db("DataCon");

$q = mysql_query("SHOW TABLE STATUS FROM `DataCon`");
/*
CREATE TABLE `impd8acd9a08a2be08d59dad18f5368d409546aba6f_85531314a1abfd8e544` (
  `department_family` varchar(50) DEFAULT NULL,
  `entity` varchar(50) DEFAULT NULL,
  `date` varchar(50) DEFAULT NULL,
  `expense_type` varchar(50) DEFAULT NULL,
  `expense_area` varchar(50) DEFAULT NULL,
  `supplier` varchar(50) DEFAULT NULL,
  `transaction_number` varchar(50) DEFAULT NULL,
  `amount` varchar(50) DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1
*/
while($row = mysql_fetch_array($q)) {
	$tablename = $row['Name'];
	// okay so we need to scan the table now.
	$TableQuery = mysql_query("SHOW CREATE TABLE `$tablename`;");
	$TableQueryData = mysql_fetch_array($TableQuery);
	$CreateCode = $TableQueryData[1];
	$CreateCodeBits = explode("\n", $CreateCode);
	$Blanking_Query = "DELETE FROM `$tablename` WHERE 1 = 1 AND ";
	foreach ($CreateCodeBits as $linenum => $line) {
		if(strstr($line, "  `")) {
			$bitsofline = explode("`", $line);
			$Blanking_Query = $Blanking_Query . "`" .  $bitsofline[1] . "` = '' AND ";
		}
	}
	$Blanking_Query = $Blanking_Query . " 2 = 2;";
	echo($Blanking_Query . " \n\n");
}