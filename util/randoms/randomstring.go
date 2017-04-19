// Copyright © 2017 Yoshiki Shibata. All rights reserved.

package randoms

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomString generates a random string with the specified length.
func RandomString(length int) string {
	if length <= 0 {
		panic(fmt.Errorf("length is %d, but want a positive value", length))
	}

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
