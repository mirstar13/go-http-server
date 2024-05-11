package methods

import "fmt"

const (
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
	ContentEncoding = "Content-Encoding"
	AcceptEncoding  = "Accept-Encoding"
	UserAgent       = "User-Agent"
)

const (
	TextPlain      = "text/plain"
	AppOctetStream = "application/octet-stream"
	Gzip           = "gzip"
)

var AcceptedEncodings = map[string]string{
	Gzip: "gzip",
}

func FormatContentTypeHeader(contentType string, length int) []string {
	res := make([]string, 0)

	res = append(res, fmt.Sprintf("%s: %s", ContentType, contentType))
	res = append(res, fmt.Sprintf("%s: %d", ContentLength, length))

	return res
}

func FormatHeader(header string, val string) string {
	return fmt.Sprintf("%s: %s", header, val)
}
