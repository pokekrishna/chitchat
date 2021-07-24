package content

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	TypeNotSupported = "text/plain; charset=utf-8"
	TypeJSON = "application/json" // RFC8259
)

const KeyAcceptContentType = "AcceptContentType"

var (
	ErrContentContextNotFound = errors.New("content.Context Not Found in http.Request.Context tree")
)

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
	if c.ContentType() == TypeNotSupported {
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

// TODO: try to return content.Context in a crafty or idiomatic way...
// TODO: it will help the user to access methods like UnsupportedContentType

// ExtractContentType takes in a request and tries to get the content type set in
// the request context.
//
// If the context tree has content.Context then the returned context type is
// either a valid content type or TypeNotSupported, otherwise it returns an error
// ErrContentContextNotFound indicating there was some error populating the
// Context. In such cases caller is expected to not set any response headers
// explicitly
func ExtractContentType(r *http.Request) (string, error){
	contentType := r.Context().Value(KeyAcceptContentType)
	switch contentType :=  contentType.(type){
	case string:
		return contentType, nil
	default:
		return "", ErrContentContextNotFound
	}
}