package cmd

type ContextKey struct {
	name string
}

func (c *ContextKey) String() string {
	return c.name
}

func CtxKey(s string) *ContextKey {
	return &ContextKey{
		name: s,
	}
}

var (
	ContextKeyID = CtxKey("id")
)
