package content

import (
	"context"
	"strings"
)

const (
	TypeNotSupported = "text/plain; charset=utf-8"
	TypeJSON = "application/json" // RFC8259
)

const KeyAcceptContentType = "AcceptContentType"

type Context struct {
	context.Context
}

func (c *Context) ContentType() string {
	val := c.Value(KeyAcceptContentType)
	switch val := val.(type) {
	case string:
		return val
	default:
		return TypeNotSupported
	}
}

func (c *Context) UnsupportedContentType() bool {
	val := c.ContentType()
	if val == TypeNotSupported {
		return true
	}
	return false
}


func ContextWithSupportedType(parent context.Context, contentType string) *Context {
	return &Context{context.WithValue(parent, KeyAcceptContentType, ValidateType(contentType))}
}

// ValidateType returns enums by checking whether the content type 't' is
// supported or not for response
func ValidateType(t string) string {
	if len(t) < 1 {
		return TypeNotSupported
	}
	t = strings.ToLower(t)
	switch t{
	case TypeJSON:
		return TypeJSON
	default:
		return TypeNotSupported
	}
}
