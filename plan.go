package plan

import (
	"log"
	"backup/util"
)

var Debug = true

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func ExportPlan(site string) []string {
	log.Println("ExportPlan()")
	return nil
}

func ReadPlans() {
	util.DeleteFile("test.123")
}
