package app

import "net/http"

func getAliveBytes() []byte {
	return []byte{0x61, 0x6c, 0x69, 0x76, 0x65}
}

func initHandler(handler *http.ServeMux, appCfg *IpInfoAppConfig) {
	handler.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(getAliveBytes())
	}))
}

func NewAppHandler(appCfg *IpInfoAppConfig) *http.ServeMux {
	handler := http.NewServeMux()
	initHandler(handler, appCfg)
	return handler
}
