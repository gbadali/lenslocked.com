package rand

// Box 11.3. Wrapping packages Wrapping packages is a pattern used
// by many Go developers as a way to take a more general package,
// like the crypto/ rand package, and simplify it for their
//  application by handling a lot of the application-specific logic
//  inside of their wrapping. In our particular case, creating a
// custom rand package means we can individually test the package
//  and then easily use it throughout our code without worrying about
//  the details of how we generate a random string.

import (
	"crypto/rand"
	"encoding/base64"
)

// Set this up so we don't accidentaly create it too small
const RememberTokenBytes = 32

// RememberToken is a helper function designed to generate
// remember tokens of a predetermined byte size.
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// Bytes will help us generate n random bytes, or will
// return an error if there was one.   This uses the
// crypto/rand package so it is safe to use with things
// like remember tokens.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byte slice of size nBytes and then
// return a string that is base64 URL encoded version
// of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
