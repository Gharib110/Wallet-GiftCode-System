package api

import "net/http"

func (c *ConfigAPI) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Connection") == "Upgrade" && r.Header.Get("Upgrade") == "websocket" {
			// Add CORS headers to allow WebSocket handshake from any origin to
			// prevent abnormal closure
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
		}
		next.ServeHTTP(w, r)
	})
}
