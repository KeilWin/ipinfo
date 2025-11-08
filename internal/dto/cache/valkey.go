package cache

type ValkeyCacheConfig struct {
}

type ValkeyCache struct {
	Cache

	Config *CacheConfig
}

func (p *ValkeyCache) StartUp() error {
	return nil
}

func (p *ValkeyCache) ShutDown() error {
	return nil
}

func (p *ValkeyCache) AddIpInfo() {

}

func (p *ValkeyCache) GetIpInfo(ipAddress string) {

}

func NewValkeyCache(config *CacheConfig) *ValkeyCache {
	return &ValkeyCache{
		Config: config,
	}
}
