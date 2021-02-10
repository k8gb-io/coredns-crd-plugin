package test

import (
	"fmt"
	"strings"
	"sort"
	"testing"
	"github.com/gruntwork-io/terratest/modules/shell"
)

func Dig(t *testing.T, dnsServer string, dnsPort int, dnsName string) ([]string, error) {
        port := fmt.Sprintf("-p%v", dnsPort)
        dnsServer = fmt.Sprintf("@%s", dnsServer)

        digApp := shell.Command{
                Command: "dig",
                Args:    []string{port, dnsServer, dnsName, "+short"},
        }

        digAppOut := shell.RunCommandAndGetOutput(t, digApp)
        digAppSlice := strings.Split(digAppOut, "\n")

        sort.Strings(digAppSlice)

        return digAppSlice, nil
}
