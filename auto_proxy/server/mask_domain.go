package server

func EncodeDomain(domain string) (encode string) {
	for i, b := range []byte(domain) {
		encode += ByteTo16(b ^ domainMask[i])
	}

	return encode
}

func encodeDomainWithLen(domain string) (encode string, err error) {
	encode = EncodeDomain(domain)

	return ByteTo16(byte(len(encode)/2)^domainLenMask) + encode, nil
}
