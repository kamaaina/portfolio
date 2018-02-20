package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/leekchan/accounting"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type fund struct {
	Ticker           string
	Name             string
	Rating           int
	ExpenseRatio     float64
	Shares           float64
	Price            float64
	Total            float64
	YTD              float64
	YTDN             float64
	ThreeMonthYield  float64
	ThreeMonthYieldN float64
	OneYearYield     float64
	OneYearYieldN    float64
	ThreeYearYield   float64
	ThreeYearYieldN  float64
	FiveYearYield    float64
	FiveYearYieldN   float64
	Allocation       allocation
}

type allocation struct {
	Cash           float64
	CashN          float64
	Domestic       float64
	DomesticN      float64
	International  float64
	InternationalN float64
	Bond           float64
	BondN          float64
	Other          float64
	OtherN         float64
}

// global
var retirement float64
var nonRetirement float64

const acctInsert = "INSERT INTO accounts(name, is_retirement) VALUES(?, ?)"
const fundInsert = "INSERT INTO fund(ticker, name, morningstar_rating, expense_ratio, shares, price, account_id) VALUES(?, ?, ?, ?, ?, ?, ?)"
const assetAllocationInsert = "INSERT INTO asset_allocation(fund_id, cash, domestic, international, bonds, other) VALUES(?, ?, ?, ?, ?, ?)"
const perfInsert = "INSERT INTO performance(fund_id, ytd, three_month, one_year, three_year, five_year) VALUES(?, ?, ?, ?, ?, ?)"

func main() {
	funds := make(map[string][]fund)
	getFunds("/home/mwhite/port.csv", funds, "V", false)
	//getStocks("/Users/mike/mike/port_data/portHEI.csv", funds, "HEI")
	getFunds("/home/mwhite/portF.csv", funds, "F", false)
	getFunds("/home/mwhite/portRM.csv", funds, "RM", true)
	//getFunds("/Users/mike/mike/port_data/portTIAA.csv", funds, "TIAA", true)
	getFunds("/home/mwhite/portRC.csv", funds, "RC", true)
	getFunds("/home/mwhite/portira.csv", funds, "IRA", true)
	getLMFunds("/home/mwhite/portssp.csv", funds) // SSP and CAP

	db, err := sql.Open("mysql", fmt.Sprintf("mike:mike@tcp(%s:3306)/portfolio", "192.168.2.41"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	time := time.Now()
	timeStr := fmt.Sprintf(time.Format("2006-01-02 15:04:05"))
	sql := fmt.Sprintf("INSERT INTO `summary`(`key`, `value`) VALUES('date', '%s')", timeStr)
	_, err = db.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
	normalizeYields(funds, db)

	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	fmt.Printf("retirement: %s\n", ac.FormatMoney(retirement))
	fmt.Printf("non-retirement: %s\n", ac.FormatMoney(nonRetirement))
	fmt.Printf("total: %s\n", ac.FormatMoney(nonRetirement+retirement))
	//fmt.Println(funds)
}

func getStocks(filename string, fundMap map[string][]fund, name string) {
	f, _ := os.Open(filename)
	r := csv.NewReader(bufio.NewReader(f))
	isFirstLine := true

	fundList := make([]fund, 0)
	var sum float64

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// ignore first line
		if isFirstLine {
			isFirstLine = false
			continue
		}

		// check if we have a fund
		if len(record[0]) <= 0 {
			continue
		}

		f := fund{}
		f.Ticker = record[0]
		f.Name = record[1]
		f.Shares = getFloat64FromString(record[2])
		f.Price = getFloat64FromString(record[3])
		f.Total = f.Shares * f.Price

		a := allocation{}
		a.Cash = 0.0
		a.Domestic = 100.0
		a.International = 0.0
		a.Bond = 0.0
		a.Other = 0.0
		f.Allocation = a
		sum += f.Shares * f.Price

		fundList = append(fundList, f)
	}
	key := fmt.Sprintf("%s:%f", name, sum)
	nonRetirement += sum
	fundMap[key] = fundList
}

func getLMFunds(filename string, fundMap map[string][]fund) {
	f, _ := os.Open(filename)
	r := csv.NewReader(bufio.NewReader(f))
	isFirstLine := true
	var name = "CAP"

	fundList := make([]fund, 0)
	var sum float64
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// ignore first line
		if isFirstLine {
			isFirstLine = false
			continue
		}

		// check if we have a fund
		if len(record[0]) <= 0 {
			continue
		} else if record[0] == "SSP" {
			key := fmt.Sprintf("%s:%f", name, sum)
			retirement += sum
			fundMap[key] = fundList
			fundList = fundList[:0] // delete items in the slice
			name = "SSP"
			sum = 0
			continue
		}

		f := fund{}
		f.Name = record[0]
		f.Total = getFloat64FromString(record[7])

		a := allocation{}
		a.Cash = getFloat64FromString(record[5])
		a.Domestic = getFloat64FromString(record[2])
		a.International = getFloat64FromString(record[3])
		a.Bond = getFloat64FromString(record[4])
		a.Other = getFloat64FromString(record[6])
		f.Allocation = a
		sum += f.Total

		fundList = append(fundList, f)
	}
	key := fmt.Sprintf("%s:%f", name, sum)
	retirement += sum
	fundMap[key] = fundList
}

func getFunds(filename string, fundMap map[string][]fund, name string, isRetirement bool) {
	f, _ := os.Open(filename)
	r := csv.NewReader(bufio.NewReader(f))
	isFirstLine := true

	fundList := make([]fund, 0)
	var sum float64

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// ignore first line
		if isFirstLine {
			isFirstLine = false
			continue
		}

		// check if we have a fund
		if len(record[0]) <= 0 {
			continue
		}

		f := fund{}
		f.Ticker = record[0]
		f.Name = record[1]
		f.Rating = getIntFromString(record[2])
		f.ExpenseRatio = getFloat64FromString(record[3])
		f.Shares = getFloat64FromString(record[4])
		f.Price = getFloat64FromString(record[5])
		f.Total = f.Shares * f.Price
		f.YTD = getFloat64FromString(record[7])
		f.ThreeMonthYield = getFloat64FromString(record[9])
		f.OneYearYield = getFloat64FromString(record[11])
		f.ThreeYearYield = getFloat64FromString(record[13])
		f.FiveYearYield = getFloat64FromString(record[15])

		a := allocation{}
		a.Cash = getFloat64FromString(record[17])
		a.Domestic = getFloat64FromString(record[18])
		a.International = getFloat64FromString(record[19])
		a.Bond = getFloat64FromString(record[20])
		a.Other = getFloat64FromString(record[21])
		f.Allocation = a
		sum += f.Shares * f.Price

		fundList = append(fundList, f)
	}
	key := fmt.Sprintf("%s:%f", name, sum)
	if isRetirement {
		retirement += sum
	} else {
		nonRetirement += sum
	}
	fundMap[key] = fundList
}

func normalizeYields(fundMap map[string][]fund, db *sql.DB) {
	for key, acct := range fundMap {
		tmp := strings.Split(key, ":") // name:sum
		s, _ := strconv.ParseFloat(tmp[1], 32)

		// FIXME: check for dups!

		// insert into accounts
		acctIns, err := db.Prepare(acctInsert)
		if err != nil {
			panic(err.Error())
		}
		defer acctIns.Close()
		result, e := acctIns.Exec(tmp[0], false) // FIXME - hardcode to false
		if e != nil {
			panic(e.Error())
		}
		var acctID int64
		var fundID int64
		acctID, err = result.LastInsertId()
		if err != nil {
			panic(err.Error())
		}

		// fund
		fundIns, fundErr := db.Prepare(fundInsert)
		if fundErr != nil {
			panic(fundErr.Error())
		}
		defer fundIns.Close()

		// asset allocation
		assetAllocIns, aaErr := db.Prepare(assetAllocationInsert)
		if aaErr != nil {
			panic(aaErr.Error())
		}
		defer assetAllocIns.Close()

		// performance
		perfIns, perfErr := db.Prepare(perfInsert)
		if perfErr != nil {
			panic(perfErr.Error())
		}
		defer perfIns.Close()

		sum := float64(s)

		var ytd float64
		var threeMonth float64
		var oneYear float64
		var threeYear float64
		var fiveYear float64
		var fundRes sql.Result
		for i := 0; i < len(acct); i++ {
			pct := acct[i].Total / sum

			acct[i].YTDN = pct * acct[i].YTD
			acct[i].ThreeMonthYieldN = pct * acct[i].ThreeMonthYield
			acct[i].OneYearYieldN = pct * acct[i].OneYearYield
			acct[i].ThreeYearYieldN = pct * acct[i].ThreeYearYield
			acct[i].FiveYearYieldN = pct * acct[i].FiveYearYield
			acct[i].Allocation.CashN = pct * acct[i].Allocation.Cash
			acct[i].Allocation.DomesticN = pct * acct[i].Allocation.Domestic
			acct[i].Allocation.InternationalN = pct * acct[i].Allocation.International
			acct[i].Allocation.BondN = pct * acct[i].Allocation.Bond
			acct[i].Allocation.OtherN = pct * acct[i].Allocation.Other
			ytd += acct[i].YTDN
			threeMonth += acct[i].ThreeMonthYieldN
			oneYear += acct[i].OneYearYieldN
			threeYear += acct[i].ThreeYearYieldN
			fiveYear += acct[i].FiveYearYieldN
			fundRes, err = fundIns.Exec(acct[i].Ticker, acct[i].Name, acct[i].Rating, acct[i].ExpenseRatio, acct[i].Shares, acct[i].Price, acctID)
			if err != nil {
				panic(err.Error())
			}
			fundID, err = fundRes.LastInsertId()
			if err != nil {
				panic(err.Error())
			}

			// asset allocation
			_, err = assetAllocIns.Exec(fundID, acct[i].Allocation.Cash, acct[i].Allocation.Domestic, acct[i].Allocation.International, acct[i].Allocation.Bond, acct[i].Allocation.Other)
			if err != nil {
				panic(err.Error())
			}

			// performance
			_, err = perfIns.Exec(fundID, acct[i].YTD, acct[i].ThreeMonthYield, acct[i].OneYearYield, acct[i].ThreeYearYield, acct[i].FiveYearYield)
		}
		if tmp[0] != "SSP" && tmp[0] != "CAP" && tmp[0] != "HEI" {
			fmt.Printf("%s:\tYTD %.2f%s", tmp[0], ytd, "%")
			fmt.Printf("\t3mo %.2f%s", threeMonth, "%")
			fmt.Printf("\t1yr %.2f%s", oneYear, "%")
			fmt.Printf("\t3yr %.2f%s", threeYear, "%")
			fmt.Printf("\t5yr %.2f%s\n", fiveYear, "%")
		}
	}
}

func getIntFromString(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func getFloat64FromString(s string) float64 {
	var smod = strings.TrimSpace(s)
	smod = strings.Replace(smod, "%", "", -1)
	smod = strings.Replace(smod, "$", "", -1)
	smod = strings.Replace(smod, ",", "", -1)
	val, _ := strconv.ParseFloat(smod, 32)
	return val
}

func dbWrite(sql string, db *sql.DB) {

}
