package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func CalculateSignature(r *http.Request, parsed *ParsedAuth, secretKey, amzDate, payloadHash string) string {
	canonicalHeaders := ""
	for _, header := range strings.Split(parsed.SignedHeaders, ";") {
		header = strings.ToLower(header)
		value := r.Header.Get(header)
		if header == "host" {
			value = r.Host
		}
		canonicalHeaders += header + ":" + normalizeHeaderValue(value) + "\n"
	}

	canonicalRequest := strings.Join([]string{
		r.Method,
		r.URL.EscapedPath(),
		r.URL.RawQuery,
		canonicalHeaders,
		parsed.SignedHeaders,
		payloadHash,
	}, "\n")

	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", parsed.Date, parsed.Region, parsed.Service)
	hashedCanonical := sha256Hex([]byte(canonicalRequest))

	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hashedCanonical,
	}, "\n")

	signingKey := getSignatureKey(secretKey, parsed.Date, parsed.Region, parsed.Service)

	signature := hmacSHA256Hex(signingKey, stringToSign)
	return signature
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func hmacSHA256Hex(key []byte, data string) string {
	return hex.EncodeToString(hmacSHA256(key, data))
}

func sha256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func getSignatureKey(secret, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}

func normalizeHeaderValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
