package orm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itpkg/ksana/utils"
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

func (p *Configuration) Load(env string) error {
	return utils.FromToml(fmt.Sprintf("config/%s/database.toml", env), p)
}

func (p *Configuration) Create() (string, []string) {
	switch p.Driver {
	case "postgres":
		return "psql", []string{"-U", p.User, "-h", p.Host, "-p", strconv.Itoa(p.Port), "-c", fmt.Sprintf("CREATE DATABASE %s WITH ENCODING='UTF8'", p.Name)}
	default:
		return "echo", []string{"unknown database."}
	}
}

func (p *Configuration) Drop() (string, []string) {
	switch p.Driver {
	case "postgres":
		return "psql", []string{"-U", p.User, "-h", p.Host, "-p", strconv.Itoa(p.Port), "-c", fmt.Sprintf("DROP DATABASE %s", p.Name)}
	default:
		return "echo", []string{"unknown database."}
	}
}

func (p *Configuration) Connect() (string, []string) {
	switch p.Driver {
	case "postgres":
		return "psql", []string{"-U", p.User, "-d", p.Name, "-h", p.Host, "-p", strconv.Itoa(p.Port)}
	default:
		return "echo", []string{"unknown database."}
	}
}

//==============================================================================
var modules = make([]string, 0)

func Register(path string) {
	modules = append(modules, path)
}
