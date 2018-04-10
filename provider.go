package fasthttpsession

type Provider interface {
	Init(ProviderConfig) error
	SessionIdIsExist(string) bool
	ReadStore(string) (SessionStore, error)
	Destroy(string) error
	Count() int
}