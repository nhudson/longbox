package api

var (
	requestIDContextKey = contextKey("RequestID")
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}
