package main

// import (
// 	"fmt"
// )

var TableSchemaStorage = make(map[string][]ColType)

/**
 * @brief Get the SQL Scheme for a Table
 * @details Almost all of the SQLs support 'information_schema' database which stores metadata about
 * other databases, tables etc.
 *
 * @param string <Table Name>
 * @return <Table Schema>
 */
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
