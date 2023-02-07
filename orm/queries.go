package orm

import (
	"os"
)

const (
	DIR = "orm/queries/"
)

func openSql(fileName string) string {
	return DIR + fileName
}

func Sql1() string {
	query, err := os.ReadFile(openSql("sql1.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}
