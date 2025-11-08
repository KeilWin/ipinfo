package dao

import (
	"github.com/KeilWin/ipinfo/internal/dto/database"
	"github.com/KeilWin/ipinfo/internal/entity"
)

type IpAddressRepository interface {
	GetIpAddress(ipAddress string) *entity.IpAddress
}

type IpAddress struct {
	Db database.Database
}

func (p *IpAddress) GetIpAddress(ipAddress string) *entity.IpAddress {
	addr := p.Db.GetIpInfo(ipAddress)
	return &entity.IpAddress{
		Value: addr,
	}
}

func NewIpAddress(db database.Database) *IpAddress {
	return &IpAddress{
		Db: db,
	}
}
