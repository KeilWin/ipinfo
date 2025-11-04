package handler

import (
	"net/http"

	"github.com/KeilWin/ipinfo/internal/ipinfo/service"
)

func NewIpV6Handler(service service.IpAddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipAddressFromPath := r.PathValue("ipAddress")
		ipAddressInfo := service.GetIpAddress(ipAddressFromPath)
		w.Write([]byte(ipAddressInfo.Value))
	}
}
