package utils

// contextKey は context に値を格納するための専用のキータイプです。
type ContextKey struct{}

// cookieKey は Cookie 値を context に格納するための専用のキーです。
var CookieKey = ContextKey{}