<?php
error_reporting(E_ERROR | E_PARSE);
//foreach (glob("*.csv") as $filename) {
//    echo "$filename size " . filesize($filename) . "\n";



$file = $argv[1];
$bang = explode(".",$file);
$table = "imp" . $bang[0];
$table = substr($table, 0, 63);
// get structure from csv and insert db
ini_set('auto_detect_line_endings',TRUE);
$handle = fopen($file,'r');
// first row, structure
if ( ($data = fgetcsv($handle) ) === FALSE ) {
    echo "";die();
}
$fields = array();
$field_count = 0;
for($i=0;$i<count($data); $i++) {
    $f = strtolower(trim($data[$i]));
    if ($f) {
        // normalize the field name, strip to 20 chars if too long
        $f = substr(preg_replace ('/[^0-9a-z]/', '_', $f), 0, 20);
        $field_count++;
        $fields[] = $f.' VARCHAR(50)';
    }
}

$sql = "CREATE TABLE $table (" . implode(', ', $fields) . ')';
echo $sql . ";\r\n";
//echo $sql . "<br /><br />";
// $db->query($sql);
while ( ($data = fgetcsv($handle) ) !== FALSE ) {
    $fields = array();
    for($i=0;$i<$field_count; $i++) {
        $fields[] = '\''.addslashes($data[$i]).'\'';
    }
    $sql = "Insert into $table values(" . implode(', ', $fields) . ')';
    echo $sql . ";\r\n";
    // $db->query($sql);
}
fclose($handle);
ini_set('auto_detect_line_endings',FALSE);

//}
