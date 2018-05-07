package fasthttpsession

type Provider interface {
	Init(ProviderConfig) error
	MaxLifeTime(int64)
	NeedGC() bool
	GC(int64)
	ReadStore(string) (SessionStore, error)
	Regenerate(string, string) (SessionStore, error)
	Destroy(string) error
	Count() int
}

type ProviderConfig interface {
	Name() string
}
