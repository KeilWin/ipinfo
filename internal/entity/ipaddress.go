package entity

type IpAddressInfo struct {
	IpAddress        string `json:"ipAddress"`
	RirName          string `json:"rirName"`
	IpAddressVersion string `json:"ipAddressVersion"`
	CountryCode      string `json:"countryCode"`
	IpRangeStart     string `json:"ipRangeStart"`
	IpRangeEnd       string `json:"ipRangeEnd"`
	IpRangeQuantity  string `json:"ipRangeQuantity"`
	Status           string `json:"status"`
	StatusUpdatedAt  string `json:"statusUpdatedAt"`
}

func NewIpAddressInfo() *IpAddressInfo {
	return &IpAddressInfo{}
}
