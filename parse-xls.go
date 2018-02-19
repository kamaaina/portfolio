package main

import ("fmt"
	"github.com/tealeg/xlsx")

func main() {
	xlFile, err := xlsx.OpenFile("/home/mwhite/test.xlsx")
	if err != nil {
		panic(err)
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Println(sheet.Name)
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("%s\n", text)
			}
		}
	}
}
