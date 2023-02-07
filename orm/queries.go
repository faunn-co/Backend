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

func Sql2() string {
	query, err := os.ReadFile(openSql("sql2.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func Sql3() string {
	query, err := os.ReadFile(openSql("sql3.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func Sql4() string {
	query, err := os.ReadFile(openSql("sql4.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}
