package wrr

type groups []*group

// parseGroups create slice of groups in order as they are defined in Endpoint
func parseGroups(labels map[string]string) (g groups, err error) {
	filter := make(map[string]*group, 0)
	for k, v := range labels {
		pg, weight, err := parseGroup(k)
		if !weight {
			continue
		}
		if err != nil {
			return g, err
		}
		if filter[pg.String()] == nil {
			filter[pg.String()] = pg
			g = append(g, pg)
		}
		filter[pg.String()].IPs = append(filter[pg.String()].IPs, v)
	}
	return g, err
}

func (g groups) pdf() (pdf []int) {
	for _, v := range g {
		pdf = append(pdf, v.weight)
	}
	return pdf
}

func (g groups) shuffle(vec []int) {
	for i, v := range vec {
		g[v], g[i] = g[i], g[v]
	}
}

// asSlice converts groups to array of IP address
// Function respects order of groups
func (g groups) asSlice() (arr []string) {
	for _, v := range g {
		arr = append(arr, v.IPs...)
	}
	return arr
}
