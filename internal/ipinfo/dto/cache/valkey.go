package cache

type ValkeyCacheConfig struct {
}

type ValkeyCache struct {
	Cache
}

func (p *ValkeyCache) StartUp() {

}

func (p *ValkeyCache) ShutDown() {

}

func (p *ValkeyCache) AddIpInfo() {

}

func (p *ValkeyCache) GetIpInfo(ipAddress string) {

}

func NewValkeyCache() *ValkeyCache {
	return &ValkeyCache{}
}
