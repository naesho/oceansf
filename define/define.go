package define

type (
	Json map[string]interface{}
)

const (
	MemcachePrefix = "oceansf:"
	DefaultCacheTime = 60 * 60
	SessionExpire = 60 * 60
)
