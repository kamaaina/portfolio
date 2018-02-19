package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
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
	ExpenseRatio     float32
	Shares           float32
	Price            float32
	Total            float32
	YTD              float32
	YTDN             float32
	ThreeMonthYeild  float32
	ThreeMonthYeildN float32
	OneYearYield     float32
	OneYearYieldN    float32
	ThreeYearYield   float32
	ThreeYearYieldN  float32
	FiveYearYield    float32
	FiveYearYieldN   float32
	Allocation       allocation
}

type allocation struct {
	Cash           float32
	CashN          float32
	Domestic       float32
	DomesticN      float32
	International  float32
	InternationalN float32
	Bond           float32
	BondN          float32
	Other          float32
	OtherN         float32
}

func main() {
	funds := make(map[string][]fund)
	getFunds("/home/mwhite/port.csv", funds, "V")
	getFunds("/home/mwhite/portRM.csv", funds, "RM")
	normalizeYields(funds)
	fmt.Println(funds)
}

func getFunds(filename string, fundMap map[string][]fund, name string) {
	f, _ := os.Open(filename)
	r := csv.NewReader(bufio.NewReader(f))
	isFirstLine := true

	fundList := make([]fund, 0)
	var sum float32
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
		f.ExpenseRatio = getFloat32FromString(record[3])
		f.Shares = getFloat32FromString(record[4])
		f.Price = getFloat32FromString(record[5])
		f.Total = f.Shares * f.Price
		f.YTD = getFloat32FromString(record[7])
		f.ThreeMonthYeild = getFloat32FromString(record[9])
		f.OneYearYield = getFloat32FromString(record[11])
		f.ThreeYearYield = getFloat32FromString(record[13])
		f.FiveYearYield = getFloat32FromString(record[15])

		a := allocation{}
		a.Cash = getFloat32FromString(record[17])
		a.Domestic = getFloat32FromString(record[18])
		a.International = getFloat32FromString(record[19])
		a.Bond = getFloat32FromString(record[20])
		a.Other = getFloat32FromString(record[21])
		f.Allocation = a
		sum += f.Shares * f.Price

		fundList = append(fundList, f)
	}
	key := fmt.Sprintf("%s:%f", name, sum)
	fundMap[key] = fundList
}

func normalizeYields(fundMap map[string][]fund) {
	for key, acct := range fundMap {
		tmp := strings.Split(key, ":") // name:sum
		s, _ := strconv.ParseFloat(tmp[1], 32)
		sum := float32(s)

		for i:=0; i<len(acct); i++ {
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

func getFloat32FromString(s string) float32 {
	var smod string = s
	if strings.Contains(smod, "%") {
		smod = strings.Replace(s, "%", "", -1)
	}
	if strings.Contains(smod, "$") {
		smod = strings.Replace(s, "$", "", -1)
	}
	val, _ := strconv.ParseFloat(smod, 32)
	return float32(val)
}
