package orm

import (
	"os"
)

var (
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

func Sql10() string {
	query, err := os.ReadFile(openSql("sql10.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralStatsQuery() string {
	query, err := os.ReadFile(openSql("sql11.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralTrendQuery() string {
	query, err := os.ReadFile(openSql("sql12.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralTrendClicksQuery() string {
	query, err := os.ReadFile(openSql("sql13.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralRecentClicksQuery() string {
	query, err := os.ReadFile(openSql("sql14.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralRecentEarningsQuery() string {
	query, err := os.ReadFile(openSql("sql15.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateListQuery() string {
	query, err := os.ReadFile(openSql("sql16.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAllReferralListQuery() string {
	query, err := os.ReadFile(openSql("sql17.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateReferralListQuery() string {
	query, err := os.ReadFile(openSql("sql18.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateReferralListWithNameQuery() string {
	query, err := os.ReadFile(openSql("sql19.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetBookingListQuery() string {
	query, err := os.ReadFile(openSql("sql20.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralDetailsByIdQuery() string {
	query, err := os.ReadFile(openSql("sql21.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetReferralBookingDetailsQuery() string {
	query, err := os.ReadFile(openSql("sql22.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateByCodeQuery() string {
	query, err := os.ReadFile(openSql("sql23.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateInfoQuery() string {
	query, err := os.ReadFile(openSql("sql24.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetUserInfoWithUserIdQuery() string {
	query, err := os.ReadFile(openSql("sql25.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetUserInfoWithUserNameQuery() string {
	query, err := os.ReadFile(openSql("sql26.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetUserInfoWithUserEmailQuery() string {
	query, err := os.ReadFile(openSql("sql27.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateInfoWithEntityNameQuery() string {
	query, err := os.ReadFile(openSql("sql28.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetAffiliateInfoWithReferralCodeQuery() string {
	query, err := os.ReadFile(openSql("sql29.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}

func GetUserInfoWithAuthQuery() string {
	query, err := os.ReadFile(openSql("sql30.sql"))
	if err != nil {
		panic(err)
	}
	return string(query)
}
