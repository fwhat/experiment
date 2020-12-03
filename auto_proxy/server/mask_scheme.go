package server

import "errors"

func encodeScheme(scheme string) (encode string, err error) {
	if len(scheme) > len(schemeMask) {
		return "", errors.New("not support scheme")
	}

	for i, b := range scheme {
		encode += ByteTo16(byte(b) ^ schemeMask[i])
	}

	return
}

func encodeSchemeWithLen(scheme string) (encode string, err error) {
	encode, err = encodeScheme(scheme)
	if err != nil {
		return
	}

	return ByteTo16(byte(len(encode)/2)^schemeLenMask) + encode, nil
}
