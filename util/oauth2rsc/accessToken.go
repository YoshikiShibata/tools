// Copyright © 2017 Yoshiki Shibata. All rights reserved.

package oauth2rsc

import (
	"fmt"
	"net/http"
	"strings"
)

// GetAccessToken extract the beare token from the request.
func GetAccessToken(r *http.Request) (string, error) {
	// The OAuth Bearer Token Usage specification (RFC 6750) defines
	// three different methods for passing bearer tokens to the protected
	// resource: the HTTP Authorization header, inside a form-encoded POST
	// body, and as a query parameter.
	//
	// Note that the HTTP specifications says that the "Authorization" header
	// keyword is itself not case sensitive so http.CanonicalHeaderKey() is
	// used as a key for the header map.
	auths := r.Header[http.CanonicalHeaderKey("Authorization")]
	if len(auths) != 0 && containsBearer(auths[0]) {
		// Note that the token value itself is a case sensitive, so the
		// original string is sliced.
		return auths[0][len("bearer "):], nil
	}

	// Here we have to support form-encoded body and query parameter
	// for bearer token, but not implemented yet
	return "", fmt.Errorf("Not Implemented Yet")
}

func containsBearer(auth string) bool {
	// The OAuth Bearer Token specification says that when the token is passed
	// as an HTTP Authorization header, the value of the header consists of
	// the keyword "Bearer", followed by a single space, and followed by the
	// token value itself.
	if auth == "" {
		return false
	}
	// The OAuth specification says that the "Bearer" keyword
	// is not case sensitive.
	return strings.HasPrefix(strings.ToLower(auth), "bearer ")
}
