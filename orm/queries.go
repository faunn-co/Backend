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

func Sql5() string {
	query, err := os.ReadFile(openSql("sql5.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func Sql6() string {
	query, err := os.ReadFile(openSql("sql6.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func Sql7() string {
	query, err := os.ReadFile(openSql("sql7.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func Sql8() string {
	query, err := os.ReadFile(openSql("sql8.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func Sql9() string {
	query, err := os.ReadFile(openSql("sql9.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}
