package fasthttpsession

type Provider interface {
	Init(int64, ProviderConfig) error
	NeedGC() bool
	GC()
	ReadStore(string) (SessionStore, error)
	Regenerate(string, string) (SessionStore, error)
	Destroy(string) error
	Count() int
}

type ProviderConfig interface {
	Name() string
}
