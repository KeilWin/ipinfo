package dao

import (
	"github.com/KeilWin/ipinfo/internal/dto/database"
	"github.com/KeilWin/ipinfo/internal/entity"
)

type IpAddressRepository interface {
	GetIpAddress(ipAddress string) (*entity.IpAddressInfo, error)
}

type IpAddress struct {
	Db database.Database
}

func (p *IpAddress) GetIpAddress(ipAddress string) (*entity.IpAddressInfo, error) {
	addr, err := p.Db.GetIpInfo(ipAddress)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, nil
	}
	return &entity.IpAddressInfo{
		IpAddress:        ipAddress,
		RirName:          addr.RirName,
		IpAddressVersion: addr.IpAddressVersion,
		CountryCode:      addr.CountryCode,
		IpRangeStart:     addr.IpRangeStart,
		IpRangeEnd:       addr.IpRangeEnd,
		IpRangeQuantity:  addr.IpRangeQuantity,
		Status:           addr.Status,
		StatusUpdatedAt:  addr.StatusUpdatedAt,
	}, nil
}

func NewIpAddressRepository(db database.Database) *IpAddress {
	return &IpAddress{
		Db: db,
	}
}
