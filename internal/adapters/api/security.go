package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func (h *Handler) SignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		signatureFromURLString := query.Get("s")

		if signatureFromURLString == "" {
			h.log.Warn("request rejected: missing signature")
			http.Error(w, "Forbidden: signature is missing", http.StatusForbidden)
			return
		}

		signatureFromURL, err := hex.DecodeString(signatureFromURLString)
		if err != nil {
			h.log.Warn("request rejected: invalid signature form")
			http.Error(w, "Forbidden: invalid signature format", http.StatusForbidden)
			return
		}

		query.Del("s")
		baseString := query.Encode()

		mac := hmac.New(sha256.New, []byte(h.cfg.Security.HMACSecretKey))
		mac.Write([]byte(baseString))
		expectedSignature := mac.Sum(nil)

		if !hmac.Equal(signatureFromURL, expectedSignature) {
			h.log.Warn("request rejected: invalid signature")
			http.Error(w, "Forbidden: invalid signature", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
