package orm

type Migration struct {
	Id   string   `toml:"-"`
	Up   []string `toml:"up"`
	Down []string `toml:"down"`
}
