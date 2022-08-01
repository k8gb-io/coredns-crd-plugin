package gateway

type Opts struct {
	annotation string
	apex       string
	hostmaster string
	ttlLow     uint32
	ttlHigh    uint32
	zones      []string
}

var (
	ttlLowDefault     = uint32(60)
	ttlHighDefault    = uint32(3600)
	defaultApex       = "dns"
	defaultHostmaster = "hostmaster"
)

func NewGatewayOpts(annotation, apex string, ttlLow, ttlHigh uint32, zones []string) Opts {
	opts := Opts{
		apex:       defaultApex,
		ttlLow:     ttlLowDefault,
		ttlHigh:    ttlHighDefault,
		hostmaster: defaultHostmaster,
	}
	if len(apex) != 0 {
		opts.apex = apex
	}
	if ttlLow != 0 {
		opts.ttlLow = ttlLow
	}
	if ttlHigh != 0 {
		opts.ttlHigh = ttlHigh
	}
	opts.annotation = annotation
	opts.zones = zones
	return opts
}
