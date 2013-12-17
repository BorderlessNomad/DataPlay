<?php
// This tool goes though all the tables and if it can find col's that can be floats of intsm, 
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
$TableAlters = "";
while($row = mysql_fetch_array($q)) {
	$tablename = $row['Name'];
	// okay so we need to scan the table now.
	$TableQuery = mysql_query("SHOW CREATE TABLE `$tablename`;");
	$TableQueryData = mysql_fetch_array($TableQuery);
	$CreateCode = $TableQueryData[1];
	$CreateCodeBits = explode("\n", $CreateCode);

	foreach ($CreateCodeBits as $linenum => $line) {
		if(strstr($line, "  `") && !strstr("FLOAT", $line)) {
			$bitsofline = explode("`", $line);
			// bitsofline[2] is the target.
			$TestSQL = "SELECT `".$bitsofline[1]."` FROM `" . $tablename . "`";
			$TestQuery = mysql_query($TestSQL);
			$NumbersErrywhere = true;
			if($TestQuery) {
				while($TestRow = mysql_fetch_array($TestQuery)) {
					if(!is_numeric($TestRow[0])) {
						$AttemptedFix = str_replace(",", "", $TestRow[0]);
						if(is_numeric($AttemptedFix)) {
							// Generate SQL to fix this issue :)
							echo("UPDATE `" . $tablename . "` SET `".$bitsofline[1]."` = '" . mysql_real_escape_string($AttemptedFix) . "' WHERE `".$bitsofline[1]."` = '" . mysql_real_escape_string($TestRow[0]) . "';\n");
						} else {
							$NumbersErrywhere = false;
							break;
						}
					}
				}
				if($NumbersErrywhere && $linenum != 1) {
					$AlterTableSQL = "ALTER TABLE `" . $tablename . "` CHANGE COLUMN `" . $bitsofline[1] . "` `" . $bitsofline[1] . "` FLOAT NULL DEFAULT NULL AFTER `" .  explode("`", $CreateCodeBits[$linenum - 1])[1]  . "`;";
					$TableAlters = $TableAlters .  $AlterTableSQL . "\n";
					// ALTER TABLE `imp1783c9f20a1573b146fc22e74718235566359c36_448b6173e77013748a2`
					// 		CHANGE COLUMN `amount` `amount` FLOAT NULL DEFAULT NULL AFTER `reference`;
				}
			} else {
				echo("WAT `" . $bitsofline[1] . "` $TestSQL  " . mysql_error() . "\n");
			}
		}
	}
}
echo $TableAlters;