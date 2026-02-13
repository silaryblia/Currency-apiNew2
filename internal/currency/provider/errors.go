package provider

import (
	_ "net"
	_ "net/http"
)

//
//type HTTPError interface {
//	StatusCode() int
//}
//
//func isHTTPError(err error, code int) bool {
//	var httpErr HTTPError
//	if errors.As(err, &httpErr) {
//		return httpErr.StatusCode() == code
//	}
//	return false
//}
//
//func isHTTP5xx(err error) bool {
//	var httpErr HTTPError
//	if errors.As(err, &httpErr) {
//		return httpErr.StatusCode() >= 500
//	}
//	return false
//}
