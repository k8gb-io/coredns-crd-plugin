package main

import (
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func main() {
	fh, err := os.Create("geoip.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	writer, _ := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "Test-IP-DB",
			RecordSize:   24,
			IPVersion:    4,
		},
	)
	_, absaSDCNet, _ := net.ParseCIDR("192.200.1.0/24")
	_, absa270Net, _ := net.ParseCIDR("192.200.2.0/24")

	absaSDCData := mmdbtype.Map{
		"datacenter": mmdbtype.String("site1"),
	}

	absa270Data := mmdbtype.Map{
		"datacenter": mmdbtype.String("site2"),
	}

	if err := writer.InsertFunc(absaSDCNet, inserter.TopLevelMergeWith(absaSDCData)); err != nil {
		log.Fatal(err)
	}
	if err := writer.InsertFunc(absa270Net, inserter.TopLevelMergeWith(absa270Data)); err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
