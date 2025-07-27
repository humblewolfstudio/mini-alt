package auth

import (
	"errors"
	"strings"
)

type ParsedAuth struct {
	AccessKeyID   string
	Date          string
	Region        string
	Service       string
	SignedHeaders string
	Signature     string
}

func ParseAuthorizationHeader(authHeader string) (*ParsedAuth, error) {
	const prefix = "AWS4-HMAC-SHA256 "

	if !strings.HasPrefix(authHeader, prefix) {
		return nil, errors.New("invalid auth header prefix")
	}

	authHeader = strings.TrimPrefix(authHeader, prefix)

	parts := strings.Split(authHeader, ", ")
	if len(parts) != 3 {
		return nil, errors.New("malformed auth header")
	}

	var parsed ParsedAuth

	for _, part := range parts {
		if strings.HasPrefix(part, "Credential=") {
			cred := strings.TrimPrefix(part, "Credential=")
			credParts := strings.Split(cred, "/")
			if len(credParts) != 5 {
				return nil, errors.New("malformed credential scope")
			}
			parsed.AccessKeyID = credParts[0]
			parsed.Date = credParts[1]
			parsed.Region = credParts[2]
			parsed.Service = credParts[3]
		} else if strings.HasPrefix(part, "SignedHeaders=") {
			parsed.SignedHeaders = strings.TrimPrefix(part, "SignedHeaders=")
		} else if strings.HasPrefix(part, "Signature=") {
			parsed.Signature = strings.TrimPrefix(part, "Signature=")
		}
	}

	if parsed.AccessKeyID == "" || parsed.Signature == "" {
		return nil, errors.New("missing fields in auth header")
	}

	return &parsed, nil
}
