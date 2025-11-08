package handler

import (
	"net/http"

	"github.com/KeilWin/ipinfo/internal/service"
)

func NewIpV4Handler(service service.IpAddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipAddressFromPath := r.PathValue("ipAddress")
		ipAddressInfo := service.GetIpAddress(ipAddressFromPath)
		w.Write([]byte(ipAddressInfo.Value))
	}
}
