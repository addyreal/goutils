package pgsql

var SslModeMap map[string]int = map[string]int{
	"":            0,
	"disable":     1,
	"allow":       2,
	"prefer":      3,
	"require":     4,
	"verify-ca":   5,
	"verify-full": 6,
}

var SslVersionMap map[string]int = map[string]int{
	"":        0,
	"TLSv1":   1,
	"TLSv1.1": 2,
	"TLSv1.2": 3,
	"TLSv1.3": 4,
}
