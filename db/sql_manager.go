package db

var sqlContainer = make(map[string]string)

// RegisterSQL registers a sql string in memory.
func RegisterSQL(name string, sql string) {
	sqlContainer[name] = sql
}

// GetSQL gets a stored sql string by name.
func GetSQL(name string) string {
	return sqlContainer[name]
}
