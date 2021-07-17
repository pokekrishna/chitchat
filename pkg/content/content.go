package content

import "context"

const (
	TypeNotSupported = iota
	TypeJSON
)

const KeyAcceptContentType = "AcceptContentType"

type Context struct {
	context.Context
}

func (c *Context) UnsupportedContentType() bool {
	val := c.Value(KeyAcceptContentType)
	if val == nil || val == TypeNotSupported{
		return true
	}
	return false
}

func ContextWithSupportedContentType (parent context.Context, contentType string) context.Context {
	return &Context{context.WithValue(parent, KeyAcceptContentType, contentType)}
}