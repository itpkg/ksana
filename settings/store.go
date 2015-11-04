package settings

type Store interface {
	Set(key string, val interface{}, encrypt bool) error
	Get(key string, val interface{}) error
}
