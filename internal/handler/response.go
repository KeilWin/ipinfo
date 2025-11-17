package handler

import (
	"net/http"

	"github.com/KeilWin/ipinfo/internal/entity"
)

type HealthStatus string

const (
	HealthOk = "alive"
)

type ResponseStatus int

const (
	ResponseOk ResponseStatus = iota
	ResponseNotFound
	ResponseBadRequest
	ResponseInternalError
)

type BadResponse struct {
	Code        ResponseStatus `json:"code"`
	Description string         `json:"description"`
}

type OkResponse struct {
	Code ResponseStatus `json:"code"`
	Data any            `json:"data"`
}

type HealthData struct {
	Health HealthStatus `json:"health"`
}

type IpV4Data struct {
	IpAddress        string `json:"ipAddress"`
	IpAddressVersion string `json:"ipAddressVersion"`
	CountryCode      string `json:"countryCode"`
	IpRangeStart     string `json:"ipRangeStart"`
	IpRangeEnd       string `json:"ipRangeEnd"`
	IpRangeQuantity  string `json:"ipRangeQuantity"`
	Status           string `json:"status"`
	StatusUpdatedAt  string `json:"statusUpdatedAt"`
}

type IpV6Data struct {
	IpAddress        string `json:"ipAddress"`
	IpAddressVersion string `json:"ipAddressVersion"`
	CountryCode      string `json:"countryCode"`
	IpRangeStart     string `json:"ipRangeStart"`
	IpRangeEnd       string `json:"ipRangeEnd"`
	IpRangeQuantity  string `json:"ipRangeQuantity"`
	Status           string `json:"status"`
	StatusUpdatedAt  string `json:"statusUpdatedAt"`
}

func NewOkResponse(data any) *OkResponse {
	return &OkResponse{
		Code: ResponseOk,
		Data: data,
	}
}

func NewNotFoundResponse(desription string) *BadResponse {
	return &BadResponse{
		Code:        ResponseNotFound,
		Description: desription,
	}
}

func NewBadRequestResponse(description string) *BadResponse {
	return &BadResponse{
		Code:        ResponseBadRequest,
		Description: description,
	}
}

func NewInternalErrorResponse(description string) *BadResponse {
	return &BadResponse{
		Code:        ResponseInternalError,
		Description: description,
	}
}

func NewHealthData() *HealthData {
	return &HealthData{
		Health: HealthOk,
	}
}

func NewIpV4Data(addr *entity.IpAddressInfo) *IpV4Data {
	return &IpV4Data{
		IpAddress:        addr.IpAddress,
		IpAddressVersion: addr.IpAddressVersion,
		CountryCode:      addr.CountryCode,
		IpRangeStart:     addr.IpRangeStart,
		IpRangeEnd:       addr.IpRangeEnd,
		IpRangeQuantity:  addr.IpRangeQuantity,
		Status:           addr.Status,
		StatusUpdatedAt:  addr.StatusUpdatedAt,
	}
}

func NewIpV6Data(addr *entity.IpAddressInfo) *IpV6Data {
	return &IpV6Data{
		IpAddress:        addr.IpAddress,
		IpAddressVersion: addr.IpAddressVersion,
		CountryCode:      addr.CountryCode,
		IpRangeStart:     addr.IpRangeStart,
		IpRangeEnd:       addr.IpRangeEnd,
		IpRangeQuantity:  addr.IpRangeQuantity,
		Status:           addr.Status,
		StatusUpdatedAt:  addr.StatusUpdatedAt,
	}
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
