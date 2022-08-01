package wrr

import (
	"fmt"
	"strconv"
	"strings"
)

type group struct {
	region string
	weight int
	IPs    []string
}

func (g group) String() string {
	return g.region + strconv.Itoa(g.weight)
}

// parseGroup parses "weight-eu-0-50" into group
func parseGroup(s string) (g *group, isweight bool, err error) {
	g = &group{}
	if !strings.HasPrefix(s, "weight") {
		return g, false, err
	}
	splits := strings.Split(s, "-")
	if len(splits) != 4 {
		return g, true, fmt.Errorf("invalid label: %s", s)
	}
	if splits[0] != "weight" {
		return g, true, fmt.Errorf("invalid label: %s", s)
	}
	g.region = splits[1]
	g.weight, err = strconv.Atoi(splits[3])
	return g, true, err
}
