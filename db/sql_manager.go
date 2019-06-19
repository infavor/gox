package db

var sqlContainer = make(map[string]string)

func RegisterSQL(name string, sql string) {
	sqlContainer[name] = sql
}

func GetSQL(name string) string {
	return sqlContainer[name]
}
