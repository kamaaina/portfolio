package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/leekchan/accounting"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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
	ThreeMonthYeild  float64
	ThreeMonthYeildN float64
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

func main() {
	funds := make(map[string][]fund)
	getFunds("/home/mwhite/port.csv", funds, "V", false)
	getFunds("/home/mwhite/portF.csv", funds, "F", false)
	getFunds("/home/mwhite/portRM.csv", funds, "RM", true)
	getFunds("/home/mwhite/portRC.csv", funds, "RC", true)
	getFunds("/home/mwhite/portira.csv", funds, "IRA", true)
	getLMFunds("/home/mwhite/portssp.csv", funds) // SSP and CAP

	normalizeYields(funds)

	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	fmt.Printf("retirement: %s\n", ac.FormatMoney(retirement))
	fmt.Printf("non-retirement: %s\n", ac.FormatMoney(nonRetirement))
	fmt.Printf("total: %s\n", ac.FormatMoney(nonRetirement+retirement))
	//fmt.Println(funds)
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
		f.ThreeMonthYeild = getFloat64FromString(record[9])
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

func normalizeYields(fundMap map[string][]fund) {
	for key, acct := range fundMap {
		tmp := strings.Split(key, ":") // name:sum
		s, _ := strconv.ParseFloat(tmp[1], 32)
		sum := float64(s)

		for i := 0; i < len(acct); i++ {
			pct := acct[i].Total / sum
			acct[i].YTDN = pct * acct[i].YTD
			acct[i].ThreeMonthYeildN = pct * acct[i].ThreeMonthYeild
			acct[i].OneYearYieldN = pct * acct[i].OneYearYield
			acct[i].ThreeYearYieldN = pct * acct[i].ThreeYearYield
			acct[i].FiveYearYieldN = pct * acct[i].FiveYearYield
			acct[i].Allocation.CashN = pct * acct[i].Allocation.Cash
			acct[i].Allocation.DomesticN = pct * acct[i].Allocation.Domestic
			acct[i].Allocation.InternationalN = pct * acct[i].Allocation.International
			acct[i].Allocation.BondN = pct * acct[i].Allocation.Bond
			acct[i].Allocation.OtherN = pct * acct[i].Allocation.Other
		}
	}
}

func getIntFromString(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func getFloat64FromString(s string) float64 {
	var smod = s
	if strings.Contains(smod, "%") {
		smod = strings.Replace(smod, "%", "", -1)
	}
	if strings.Contains(smod, "$") {
		smod = strings.Replace(smod, "$", "", -1)
	}
	if strings.Contains(smod, ",") {
		smod = strings.Replace(smod, ",", "", -1)
	}
	val, _ := strconv.ParseFloat(smod, 32)
	return val
}
