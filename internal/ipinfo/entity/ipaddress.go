package domain

type IpAddressVersion string

const (
	IpAddressV4 IpAddressVersion = "ipv4"
	IpAddressV6 IpAddressVersion = "ipv6"
)

type IpAddress struct {
}

func (p *IpAddress) ParseRow() {

}

func (p *IpAddress) CheckValid() {

}
