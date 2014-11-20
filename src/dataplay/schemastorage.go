package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
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

var RealTableStorage = make(map[string]OnlineData)

/**
 * @brief Converts GUID ('friendly' name) into actual table inside database
 *
 * @param string GUID
 * @param http http.ResponseWriter
 *
 * @return string output, error
 */
func GetRealTableName(guid string) (out string, e error) {
	if guid == "" || guid == "No Record Found!" {
		return "", fmt.Errorf("Invalid tablename")
	}

	if _, schemaExists := RealTableStorage[guid]; schemaExists {
		return RealTableStorage[guid].Tablename, nil
	}

	data := OnlineData{}
	err := DB.Select("tablename").Where("guid = ?", guid).Find(&data).Error
	if err != nil && err != gorm.RecordNotFound {
		return "", fmt.Errorf("Database query failed (OnlineData)")
	} else if err == gorm.RecordNotFound {
		return "", fmt.Errorf("Could not find table")
	}

	RealTableStorage[guid] = data

	return data.Tablename, nil
}

var IndexSchemaStorage = make(map[string]Index)

func GetTableIndex(guid string) (Index, error) {
	guid = strings.ToLower(strings.Trim(guid, " "))
	if guid == "" || guid == "No Record Found!" {
		return Index{}, fmt.Errorf("Invalid tablename")
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
