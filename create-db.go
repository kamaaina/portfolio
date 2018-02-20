package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("mike:mike@tcp(%s:3306)/portfolio", "192.168.2.41"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	
	fmt.Print("Dropping tables...")
	tables := []string{"asset_allocation", "performance", "fund", "accounts", "summary", "totals"}
	for _, table := range tables {
		dropTable(table, db)
	}
	fmt.Println("done")

	fmt.Print("Creating tables...")
	createTable("accounts", db)
	createTable("summary", db)
	createTable("totals", db)
	createTable("fund", db)
	createTable("performance", db)
	createTable("asset_allocation", db)
	fmt.Println("done")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func dropTable(table string, db *sql.DB) {
	sql := fmt.Sprintf("drop table if exists `%s`", table)
	_, err := db.Exec(sql)
	checkErr(err)
}

func createTable(table string, db *sql.DB) {
	var sql string
	switch table {
	case "accounts":
		sql = "CREATE TABLE `accounts` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(45) NOT NULL, `is_retirement` TINYINT DEFAULT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	case "asset_allocation":
		sql = "CREATE TABLE `asset_allocation` (`id` int(11) NOT NULL AUTO_INCREMENT, `fund_id` int(11) DEFAULT NULL, `cash` double DEFAULT NULL, `domestic` double DEFAULT NULL, `international` double DEFAULT NULL, `bonds` double DEFAULT NULL, `other` double DEFAULT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`fund_id`) REFERENCES `fund` (`id`) ON DELETE CASCADE ON UPDATE CASCADE) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	case "fund":
		sql = "CREATE TABLE `fund` (`id` int(11) NOT NULL AUTO_INCREMENT, `ticker` varchar(5) DEFAULT NULL, `name` varchar(64) NOT NULL, `morningstar_rating` tinyint(1) DEFAULT NULL, `expense_ratio` double DEFAULT NULL, `shares` double DEFAULT NULL, `price` double DEFAULT NULL, `account_id` int(11) DEFAULT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;"
	case "performance":
		sql = "CREATE TABLE `performance` (`id` int(11) NOT NULL AUTO_INCREMENT, `fund_id` int(11) DEFAULT NULL, `ytd` double DEFAULT NULL, `three_month` double DEFAULT NULL, `one_year` double DEFAULT NULL, `three_year` double DEFAULT NULL, `five_year` double DEFAULT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`fund_id`) REFERENCES `fund` (`id`) ON DELETE CASCADE ON UPDATE CASCADE) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	case "summary":
		sql = "CREATE TABLE `summary` (`id` int(11) NOT NULL AUTO_INCREMENT, `key` varchar(45) DEFAULT NULL, `value` varchar(45) DEFAULT NULL, PRIMARY KEY (`id`), KEY `id_sum_key_index` (`key`)) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;"
	case "totals":
		sql = "CREATE TABLE `totals` (`id` int(11) NOT NULL AUTO_INCREMENT, `date` date DEFAULT NULL, `key` int(11) DEFAULT NULL, `value` double DEFAULT NULL, PRIMARY KEY (`id`), KEY `date_index` (`date`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	}
	_, err := db.Exec(sql)
	checkErr(err)
}
