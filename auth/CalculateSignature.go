package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func CalculateSignature(r *http.Request, parsed *ParsedAuth, secretKey, amzDate, payloadHash string) (string, error) {
	canonicalHeaders := ""

	signedHeaders := strings.Split(parsed.SignedHeaders, ";")

	for _, header := range signedHeaders {
		header = strings.ToLower(header)
		var value string

		if header == "host" {
			value = r.Host
		} else {
			vals := r.Header.Values(header)
			if len(vals) == 0 {
				return "", fmt.Errorf("missing signed header: %s", header)
			}
			value = strings.Join(vals, ",")
		}

		canonicalHeaders += header + ":" + normalizeHeaderValue(value) + "\n"
	}

	canonicalRequest := strings.Join([]string{
		r.Method,
		normalizePath(r.URL.Path),
		canonicalQueryString(r.URL.Query()),
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
	return signature, nil
}

func canonicalQueryString(query url.Values) string {
	var keys []string
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var result []string
	for _, key := range keys {
		values := query[key]
		sort.Strings(values)
		for _, value := range values {
			encodedKey := strings.ReplaceAll(url.QueryEscape(key), "+", "%20")
			encodedValue := strings.ReplaceAll(url.QueryEscape(value), "+", "%20")
			result = append(result, encodedKey+"="+encodedValue)
		}
	}
	return strings.Join(result, "&")
}

func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}

func normalizeHeaderValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
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
