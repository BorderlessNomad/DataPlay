package main

// import (
// 	"fmt"
// )

/**
 * @brief Get the SQL Scheme for a Table
 * @details Almost all of the SQLs support 'information_schema' database which stores metadata about
 * other databases, tables etc.
 *
 * @todo Apply caching to queries which goes to 'information_schema'
 * MySQL has something like innodb_stats_on_metadata=0 which will prevent statistic update upon quering 'information_schema'.
 * Also it won't make 'information_schema' to be stale when changes are made on corresponding metadata.
 *
 * @param string <Table Name>
 * @return <Table Schema>
 */

var TableSchemaStorage = make(map[string][]ColType)

func GetSQLTableSchema(table string) []ColType {
	tableSchema := []TableSchema{}
	schema := make([]ColType, 0)

	if table == "" {
		return schema
	}

	if _, schemaExists := TableSchemaStorage[table]; schemaExists {
		return TableSchemaStorage[table]
	}

	err := DB.Select("column_name, data_type").Where("table_name = ?", table).Find(&tableSchema).Error
	if err != nil {
		return schema
	}

	for _, row := range tableSchema {
		NewCol := ColType{
			Name:    row.ColumnName,
			Sqltype: row.DataType,
		}

		if NewCol.Sqltype == "character varying" {
			NewCol.Sqltype = "varchar"
		} else if NewCol.Sqltype == "numeric" || NewCol.Sqltype == "float" || NewCol.Sqltype == "double" || NewCol.Sqltype == "real" {
			NewCol.Sqltype = "float"
		}

		schema = append(schema, NewCol)
	}

	TableSchemaStorage[table] = schema

	return schema
}
