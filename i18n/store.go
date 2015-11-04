package i18n

type Store interface {
	Set(lang, code, msg string) error
	Get(lang, code string) (string, error)
	Loop(lang string, fn func(code, msg string) error) error
}
