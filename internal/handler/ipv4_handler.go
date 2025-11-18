package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/KeilWin/ipinfo/internal/service"
)

func NewIpV4Handler(service service.IpAddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipAddressFromPath := r.PathValue("ipAddress")
		if parsedAddr := net.ParseIP(ipAddressFromPath).To4(); parsedAddr == nil {
			badResponse := NewBadRequestResponse("invalid ipv4 address")
			res, err := json.Marshal(badResponse)
			if err != nil {
				WriteInternalServerError(w)
				return
			}
			w.Write(res)
			return
		}
		ipAddressInfo, err := service.GetIpAddress(ipAddressFromPath)
		if err != nil {
			slog.Error("can't get ip address info", "err", err)
			badResponse := NewInternalErrorResponse("can't get ip address info")
			res, err := json.Marshal(badResponse)
			if err != nil {
				WriteInternalServerError(w)
				return
			}
			w.Write(res)
			return
		}
		if ipAddressInfo == nil {
			notFoundResponse := NewNotFoundResponse(fmt.Sprintf("ip address '%s' not found", ipAddressFromPath))
			res, err := json.Marshal(notFoundResponse)
			if err != nil {
				WriteInternalServerError(w)
				return
			}
			w.Write(res)
			return
		}
		ipV4Data := NewIpV4Data(ipAddressInfo)
		okResponse := NewOkResponse(ipV4Data)
		res, err := json.Marshal(okResponse)
		if err != nil {
			WriteInternalServerError(w)
			return
		}
		w.Write(res)
	}
}
