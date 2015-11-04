package orm

import (
	"fmt"
	"strings"
)

type Configuration struct {
	Driver   string                 `toml:"driver"`
	Host     string                 `toml:"host"`
	Port     int                    `toml:"port"`
	User     string                 `toml:"user"`
	Password string                 `toml:"password"`
	Name     string                 `toml:"name"`
	MaxOpen  int                    `toml:"max_open"`
	MaxIdle  int                    `toml:"max_idle"`
	Extra    map[string]interface{} `toml:"extra"`
}

func (p *Configuration) Source() string {
	ex := make([]string, 0)
	for k, v := range p.Extra {
		ex = append(ex, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?%s",
		p.Driver,
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Name,
		strings.Join(ex, "&"),
	)
}

//==============================================================================
var modules = make([]string, 0)

func Register(path string) {
	modules = append(modules, path)
}
