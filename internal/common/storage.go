package common

const (
	AIpRangesTable = "ip_ranges_a"
	BIpRangesTable = "ip_ranges_b"
)

type Storage interface {
	StartUp() error
	ShutDown() error

	GetCurrentIpRangesName() (string, error)
	CopyToIpRangesFromArray(table string, ip_ranges []IpRange) error
}

type IpRange struct {
	Rir         int
	CountryCode string
	VersionIp   int
	StartIp     string
	EndIp       string
	Status      int
	CreatedAt   string
}
