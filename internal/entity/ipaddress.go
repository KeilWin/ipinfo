package entity

type IpAddressVersionType string
type RirType string
type StatusType string

const (
	IpAddressV4 IpAddressVersionType = "ipv4"
	IpAddressV6 IpAddressVersionType = "ipv6"
)

const (
	Apnic   RirType = "apnic"
	Arin    RirType = "arin"
	Iana    RirType = "iana"
	Lacnic  RirType = "lacnic"
	Ripencc RirType = "ripencc"
)

const (
	Allocated StatusType = "allocated"
	Assigned  StatusType = "assigned"
)

type IpAddress struct {
	Rir     RirType
	Country string
	Value   string
	Start   string
	Count   string
	Version IpAddressVersionType
	Date    string
	Status  StatusType
}

func NewIpAddress() *IpAddress {
	return &IpAddress{}
}
