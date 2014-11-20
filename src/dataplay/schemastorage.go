package main

import (
	"fmt"
	"strings"
)

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

var OnlineDataSchemaStorageGuid = make(map[string]OnlineData)

/**
 * @brief Converts GUID ('friendly' name) into actual table inside database
 *
 * @param string GUID
 * @param http http.ResponseWriter
 *
 * @return string output, error
 */
func GetOnlineDataByGuid(guid string) (OnlineData, error) {
	guid = strings.ToLower(strings.Trim(guid, " "))
	if guid == "" || guid == "No Record Found!" {
		return OnlineData{}, fmt.Errorf("Invalid GUID.")
	}

	if _, schemaExists := OnlineDataSchemaStorageGuid[guid]; schemaExists {
		return OnlineDataSchemaStorageGuid[guid], nil
	}

	onlineData := OnlineData{}
	err := DB.Select("tablename").Where("guid = ?", guid).Find(&onlineData).Error
	if err != nil {
		return OnlineData{}, err
	}

	OnlineDataSchemaStorageGuid[guid] = onlineData

	return onlineData, nil
}

var OnlineDataSchemaStorageTablename = make(map[string]OnlineData)

func GetOnlineDataByTablename(tablename string) (OnlineData, error) {
	tablename = strings.ToLower(strings.Trim(tablename, " "))
	if tablename == "" || tablename == "No Record Found!" {
		return OnlineData{}, fmt.Errorf("Invalid Tablename.")
	}

	if _, schemaExists := OnlineDataSchemaStorageTablename[tablename]; schemaExists {
		return OnlineDataSchemaStorageTablename[tablename], nil
	}

	onlineData := OnlineData{}
	err := DB.Where("LOWER(tablename) = ?", tablename).Find(&onlineData).Error
	if err != nil {
		return OnlineData{}, err
	}

	OnlineDataSchemaStorageTablename[tablename] = onlineData

	return onlineData, nil
}

var IndexSchemaStorage = make(map[string]Index)

func GetTableIndex(guid string) (Index, error) {
	guid = strings.ToLower(strings.Trim(guid, " "))
	if guid == "" || guid == "No Record Found!" {
		return Index{}, fmt.Errorf("Invalid GUID")
	}

	if _, schemaExists := IndexSchemaStorage[guid]; schemaExists {
		return IndexSchemaStorage[guid], nil
	}

	index := Index{}
	err := DB.Where("LOWER(guid) = ?", guid).Find(&index).Error
	if err != nil {
		return Index{}, err
	}

	IndexSchemaStorage[guid] = index

	return index, nil
}
