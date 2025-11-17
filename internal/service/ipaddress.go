package service

import (
	"github.com/KeilWin/ipinfo/internal/dao"
	"github.com/KeilWin/ipinfo/internal/entity"
)

type IpAddressService interface {
	GetIpAddress(ipAddress string) (*entity.IpAddressInfo, error)
}

type IpAddress struct {
	Repository dao.IpAddressRepository
}

func (p *IpAddress) GetIpAddress(ipAddress string) (*entity.IpAddressInfo, error) {
	return p.Repository.GetIpAddress(ipAddress)
}

func NewIpAddress(repository dao.IpAddressRepository) *IpAddress {
	return &IpAddress{
		Repository: repository,
	}
}
