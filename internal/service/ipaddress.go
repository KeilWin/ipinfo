package service

import (
	"github.com/KeilWin/ipinfo/internal/dao"
	"github.com/KeilWin/ipinfo/internal/entity"
)

type IpAddressService interface {
	GetIpAddress(ipAddress string) *entity.IpAddress
}

type IpAddress struct {
	Repository dao.IpAddressRepository
}

func (p *IpAddress) GetIpAddress(ipAddress string) *entity.IpAddress {
	return p.Repository.GetIpAddress(ipAddress)
}

func NewIpAddress(repository dao.IpAddressRepository) *IpAddress {
	return &IpAddress{
		Repository: repository,
	}
}
