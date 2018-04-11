package fasthttpsession

type Provider interface {
	Init(ProviderConfig) error
	GC(int64)
	SessionIdIsExist(string) bool
	ReadStore(string) (SessionStore, error)
	Regenerate(string, string) (SessionStore, error)
	Destroy(string) error
	Count() int
}