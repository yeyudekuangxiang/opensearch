package sdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"sort"
	"strings"
)

const (
	HeaderSeparator = "\n"
)

func ShaHmac1(source, secret string) string {
	key := []byte(secret)
	hm := hmac.New(sha1.New, key)
	hm.Write([]byte(source))
	signedBytes := hm.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
func BuildStringToSign(request Request) (stringToSign string) {

	headers := request.GetHeaders()

	stringToSignBuilder := bytes.Buffer{}
	stringToSignBuilder.WriteString(request.GetMethod())
	stringToSignBuilder.WriteString(HeaderSeparator)

	// append header keys for sign
	appendIfContain(headers, &stringToSignBuilder, "Content-MD5", HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Content-Type", HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Accept", HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Date", HeaderSeparator)

	// sort and append headers witch starts with 'x-acs-'
	var acsHeaders []string
	for key := range headers {
		if strings.HasPrefix(strings.ToLower(key), "x-opensearch-") {
			acsHeaders = append(acsHeaders, key)
		}
	}
	sort.Strings(acsHeaders)
	for _, key := range acsHeaders {
		stringToSignBuilder.WriteString(strings.ToLower(key) + ":" + headers[key])
		stringToSignBuilder.WriteString(HeaderSeparator)
	}

	// append query params
	stringToSignBuilder.WriteString(request.BuildQueries())
	stringToSign = stringToSignBuilder.String()
	return
}
func appendIfContain(sourceMap map[string]string, target *bytes.Buffer, key, separator string) {
	if _, contain := sourceMap[key]; contain {
		target.WriteString(sourceMap[key])
		target.WriteString(separator)
	}
}
func GetMD5Base64(bytes []byte) (base64Value string) {
	md5Ctx := md5.New()
	md5Ctx.Write(bytes)
	md5Value := md5Ctx.Sum(nil)
	base64Value = base64.StdEncoding.EncodeToString(md5Value)
	return
}
