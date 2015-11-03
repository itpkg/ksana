package orm

var mapper = make(map[string]string, 0)

type Mapper struct {
	Driver  string            `toml:"driver"`
	Queries map[string]string `toml:"queries"`
}
