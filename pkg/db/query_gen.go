package db

func genCreateTableQuery(tableName string, headers []string, columnTypes []string) string {
	var createTableQuery string
	createTableQuery = "CREATE TABLE " + tableName + " ("
	for i, header := range headers {
		createTableQuery += "'" + header + "' " + columnTypes[i]
		if i != len(headers)-1 {
			createTableQuery += ", "
		}
	}
	createTableQuery += ");"
	return createTableQuery
}

func genInsertQuery(tableName string, headers []string) string {
	var insertQuery string
	insertQuery = "INSERT INTO " + tableName + " ("
	for i, header := range headers {
		insertQuery += "'" + header + "'"
		if i != len(headers)-1 {
			insertQuery += ", "
		}
	}
	insertQuery += ") VALUES ("
	for i := range headers {
		insertQuery += "?"
		if i != len(headers)-1 {
			insertQuery += ", "
		}
	}
	insertQuery += ");"
	return insertQuery
}
