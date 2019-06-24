package monitor

var (
	PacketErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Packet Unrecognized"},
	}
	DialErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Dial To Destination Address Failed"},
	}
	SocksErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Socks Protocol Error"},
	}

	// 记录monitor包出现的问题
	monitorErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Monitor Error"},
	}
	// 记录monitor包出现的warning
	monitorWarn = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Monitor Warning"},
	}
)

func init() {
	PacketErr.Setup()
	DialErr.Setup()
	SocksErr.Setup()
	monitorErr.Setup()
	monitorWarn.Setup()
}
