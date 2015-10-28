package ksana

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	"github.com/robfig/cron"
)

type BaseEngine struct {
	Db  *gorm.DB       `inject:""`
	Dao *Dao           `inject:""`
	Cfg *Configuration `inject:""`
}

func (p *BaseEngine) Router() {
}

func (p *BaseEngine) Migrate() error {
	db := p.Db

	db.AutoMigrate(&Setting{}, &Locale{}, &User{}, &Log{}, &Role{}, &Permission{})
	db.Model(&Locale{}).AddUniqueIndex("idx_base_locales_lang_type_key", "lang", "type", "key")
	db.Model(&User{}).AddUniqueIndex("idx_base_users_uid_provider", "uid", "provider")
	db.Model(&Role{}).AddUniqueIndex("idx_base_roles_name_resource_type_id", "name", "resource_type", "resource_id")
	db.Model(&Permission{}).AddUniqueIndex("idx_base_permissions_role_user", "role_id", "user_id")
	return nil
}

func (p *BaseEngine) Cron() map[string]func() {
	return map[string]func(){
		"0 0 3 * * *": func() {
			//todo
			log.Println(time.Now())
		},
	}
}

func (p *BaseEngine) Worker() {

	workers.Process("email",
		func(message *workers.Msg) {

		},
		p.Cfg.Workers["email"])
}

func (p *BaseEngine) Seed() error {
	db := p.Db
	//--------------administrator-------------
	admin_e := "root@localhost.localdomain"

	if !p.Dao.IsEmailUserExist(db, admin_e) {
		admin_u, err := p.Dao.CreateEmailUser(db, "Admin", admin_e, "changeme")
		if err != nil {
			return err
		}
		role_a := Role{Name: "admin"}
		role_r := Role{Name: "root"}
		db.Create(&role_a)
		db.Create(&role_r)

		begin := time.Now()
		end := begin.AddDate(10, 0, 0)
		db.Create(&Permission{
			RoleID:   role_a.ID,
			UserID:   admin_u.ID,
			StartUp:  begin,
			ShutDown: end,
		})
		db.Create(&Permission{
			RoleID:   role_r.ID,
			UserID:   admin_u.ID,
			StartUp:  begin,
			ShutDown: end,
		})
		db.Create(&Log{
			UserID:  admin_u.ID,
			Message: "Init.",
		})
	}
	//--------------locales-------------------
	root := fmt.Sprintf("%s/locales", PkgRoot(p))
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, f := range files {
		fn := fmt.Sprintf("%s/%s", root, f.Name())
		log.Printf("Find locale file %s", fn)
		ss := strings.Split(f.Name(), ".")
		if len(ss) != 3 {
			return errors.New(fmt.Sprintf("bad locale file name %s", f.Name))
		}
		items := make(map[string]string, 0)
		if _, err := toml.DecodeFile(fn, &items); err != nil {
			return err
		}
		for k, v := range items {
			var cn int
			db.Model(Locale{}).Where(&Locale{Type: ss[0], Lang: ss[1], Key: k}).Count(&cn)
			if cn == 0 {
				db.Create(&Locale{
					Type: ss[0],
					Lang: ss[1],
					Key:  k,
					Val:  v,
				})
			}
		}

	}

	return nil
}

func (p *BaseEngine) Deploy() {
}

func (p *BaseEngine) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:        "assets",
			Aliases:     []string{"ass"},
			Usage:       "assets operations",
			Subcommands: []cli.Command{},
		},

		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "starting background job",
			Flags: []cli.Flag{
				KSANA_ENV,
				cli.IntFlag{
					Name:  "port, p",
					Value: 10001,
					Usage: "stats will be available at http://localhost:PORT/stats",
				},
				cli.BoolFlag{
					Name:  "dispatcher, d",
					Usage: "starting with dispatcher",
				},
			},
			Action: func(c *cli.Context) {
				if _, err := New(c); err != nil {
					log.Fatal(err)
				}
				if c.Bool("dispatcher") {
					log.Println("start cron job...")
					cn := cron.New()
					if err := LoopEngine(func(en Engine) error {
						for k, v := range en.Cron() {
							cn.AddFunc(k, v)
						}
						return nil
					}); err != nil {
						log.Fatal(err)
					}
					cn.Start()
				}
				go workers.StatsServer(c.Int("port"))

				if err := LoopEngine(func(en Engine) error {
					en.Worker()
					return nil
				}); err != nil {
					log.Fatal(err)
				}

				workers.Run()
			},
		},

		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "seed",
					Aliases: []string{"s"},
					Usage:   "load the seed data into database",
					Flags: []cli.Flag{
						KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						_, err := New(c)
						if err != nil {
							log.Fatal(err)
						}
						if err = LoopEngine(func(en Engine) error {
							return en.Seed()
						}); err != nil {
							log.Fatal(err)
						}
						log.Println("Done.")

					},
				},
				{
					Name:    "migrate",
					Aliases: []string{"m"},
					Usage:   "migrate the database",
					Flags: []cli.Flag{
						KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						_, err := New(c)
						if err != nil {
							log.Fatal(err)
						}
						if err = LoopEngine(func(en Engine) error {
							return en.Migrate()
						}); err != nil {
							log.Fatal(err)
						}
						log.Println("Done.")
					},
				},
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate new code",
			Subcommands: []cli.Command{
				{
					Name:    "settings",
					Usage:   "generate a new settings.toml",
					Aliases: []string{"s"},
					Flags: []cli.Flag{
						KSANA_ENV,
						cli.StringFlag{
							Name:   "database, d",
							Value:  "postgresql",
							Usage:  "Preconfigure for selected database (options: mysql/postgresql/sqlite3)",
							EnvVar: "KSANA_ENV",
						},
					},
					Action: func(c *cli.Context) {
						env := c.String("environment")
						db := c.String("database")
						fn := fmt.Sprintf("config/%s/settings.toml", env)
						if err := Mkdirs(fmt.Sprintf("config/%s", env)); err != nil {
							log.Fatal(err)
						}

						if _, err := os.Stat(fn); err == nil {
							log.Fatalf("%s already exist!", fn)
						}
						buf, err := RandomBytes(512)
						if err != nil {
							log.Fatal(err)
						}

						cfg := Configuration{Secrets: buf,
							Http: HttpCfg{
								Host: "localhost",
								Port: 8080,
							},
							Database: DatabaseCfg{
								Dialect: db,
								Pool: PoolCfg{
									MaxIdle: 6,
									MaxOpen: 180,
								},
							},
							Redis: RedisCfg{
								Host: "localhost",
								Port: 6379,
								Db:   0,
								Pool: PoolCfg{
									MaxIdle: 4,
									MaxOpen: 120,
								},
							},
							Elasticsearch: ElasticsearchCfg{
								Host: "localhost",
								Port: 9200,
							},
							Workers: map[string]int{"email": 2, "default": 12},
						}

						switch db {
						case "postgresql":
							cfg.Database.Dialect = "postgres"
							cfg.Database.Url = fmt.Sprintf("user=postgres dbname=itpkg_%s sslmode=disable", env)
						case "mysql":
							cfg.Database.Url = fmt.Sprintf("root:@/itpkg_%s?charset=utf8&parseTime=True&loc=Local", env)
						case "sqlite3":
							cfg.Database.Url = fmt.Sprintf("tmp/itpkg_%s.db", env)
						default:
							log.Fatalf("Unsupport database %s", db)
						}

						log.Printf("Generate file %s\n", fn)
						if err := cfg.Store(fn); err != nil {
							log.Fatal(err)
						}
						log.Println("Done.")
					},
				},
				{
					Name:    "nginx",
					Aliases: []string{"n"},
					Usage:   "generage a new nginx config file",
					Flags: []cli.Flag{
						KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						cfg, err := Load(c)
						if err != nil {
							log.Fatal(err)
						}

						if err := Mkdirs(fmt.Sprintf("config/%s", cfg.Env)); err != nil {
							log.Fatal(err)
						}

						fn := fmt.Sprintf("config/%s/nginx.conf", cfg.Env)
						if _, err := os.Stat(fn); err == nil {
							log.Fatalf("%s already exist!", fn)
						}

						const nginx = `
upstream {{ .Host }}.conf {
  server http://localhost:{{ .Port }} fail_timeout=0;
}

server {
  listen 443;

  ssl  on;
  ssl_certificate  ssl/{{ .Host }}.crt.pem;
  ssl_certificate_key  ssl/{{ .Host }}.key.pem;
  ssl_session_timeout  5m;
  ssl_protocols  SSLv2 SSLv3 TLSv1;
  ssl_ciphers  RC4:HIGH:!aNULL:!MD5;
  ssl_prefer_server_ciphers  on;

  client_max_body_size 4G;
  keepalive_timeout 10;

  server_name {{ .Host }};

  root /var/www/{{ .Host }}/current/public;
  try_files $uri $uri/index.html @ksana.conf;

  location @ksana.conf {
    proxy_set_header X-Forwarded-Proto https;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $http_host;
    proxy_set_header  X-Real-IP $remote_addr;
    proxy_redirect off;
    proxy_pass http://{{ .Host }}.conf;
    # limit_req zone=one;
    access_log log/{{ .Host }}.access.log;
    error_log log/{{ .Host }}.error.log;
  }
  
  location ~* \.(?:css|js|html|jpg|jpeg|gif|png|ico)$ {
    gzip_static on;
    expires max;
    add_header Cache-Control public;
  }


  location = /50x.html {
    root html;
  }

  location = /404.html {
    root html;
  }

  location @503 {
    error_page 405 = /system/maintenance.html;
    if (-f $document_root/system/maintenance.html) {
      rewrite ^(.*)$ /system/maintenance.html break;
    }
    rewrite ^(.*)$ /503.html break;
  }

  if ($request_method !~ ^(GET|HEAD|PUT|PATCH|POST|DELETE|OPTIONS)$ ){
    return 405;
  }

  if (-f $document_root/system/maintenance.html) {
    return 503;
  }

  location ~ \.(php|jsp|asp)$ {
    return 405;
  }

}
`

						log.Printf("Generate file %s", fn)
						fd, err := os.Create(fn)
						if err != nil {
							log.Fatal(err)
						}
						defer fd.Close()
						tpl := template.Must(template.New("").Parse(nginx))
						err = tpl.Execute(fd, cfg.Http)
						if err != nil {
							log.Fatal(err)
						}

						log.Println("Done.")
					},
				},
				{
					Name:    "certificate",
					Aliases: []string{"c"},
					Usage:   "generate certificate files",
					Flags: []cli.Flag{
						KSANA_ENV,
						cli.StringFlag{
							Name:  "name",
							Value: "whoami",
							Usage: "name",
						},
						cli.StringFlag{
							Name:  "hosts",
							Value: "",
							Usage: "comma-separated hostnames and IPs to generate a certificate for",
						},
						cli.IntFlag{
							Name:  "years",
							Value: 10,
							Usage: "duration(years) that certificate is valid for",
						},

						cli.BoolFlag{
							Name:  "ca",
							Usage: "whether this cert should be its own Certificate Authority",
						},

						cli.IntFlag{
							Name:  "bits",
							Value: 2048,
							Usage: "size of RSA key to generate. ",
						},
					},
					Action: func(c *cli.Context) {
						env := c.String("environment")
						name := c.String("name")
						root := fmt.Sprintf("config/%s/ssl", env)

						if err := Mkdirs(root); err != nil {
							log.Fatal(err)
						}
						if cert, key, err := generate_certificate(name, c.String("hosts"), c.Bool("ca"), c.Int("bits"), c.Int("years")); err == nil {
							fn := fmt.Sprintf("%s/%s.crt.pem", root, name)
							log.Printf("Generate file %s\n", fn)
							if err := ioutil.WriteFile(fn, cert, 0644); err != nil {
								log.Fatal(err)
							}
							log.Printf("Verify: openssl x509 -noout -modulus -in %s | openssl md5", fn)

							fn = fmt.Sprintf("%s/%s.key.pem", root, name)
							log.Printf("Generate file %s\n", fn)
							if err := ioutil.WriteFile(fn, key, 0600); err != nil {
								log.Fatal(err)
							}
							log.Printf("Verify: openssl rsa -noout -modulus -in %s | openssl md5", fn)

							log.Println("Done.")
						} else {
							log.Fatal(err)
						}

					},
				},
			},
		},
	}
}

func generate_certificate(name, hosts string, ca bool, size int, years int) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, size)
	notBefore := time.Now()
	notAfter := notBefore.AddDate(years, 0, 0)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	tpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"OPS"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	tpl.DNSNames = strings.Split(hosts, ",")

	if name == "root" {
		tpl.IsCA = true
		tpl.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	var certOut bytes.Buffer
	err = pem.Encode(
		&certOut,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var keyOut bytes.Buffer
	err = pem.Encode(
		&keyOut,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return certOut.Bytes(), keyOut.Bytes(), nil
}

//==============================================================================
type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `sql:"not null"`
	UpdatedAt time.Time `sql:"not null"`
}
type Setting struct {
	Model
	Key  string `sql:"size:255;not null;unique"`
	Flag bool   `sql:"not null"`
	Val  []byte `sql:"not null"`
}

func (p Setting) TableName() string {
	return "base_settings"
}

type Locale struct {
	Model
	Lang string `sql:"size:5;not null;index"`
	Key  string `sql:"size:255;not null;index"`
	Val  string `sql:"type:TEXT;not null"`
	Type string `sql:"size:8;not null;index"`
}

func (p Locale) TableName() string {
	return "base_locales"
}

type User struct {
	Model
	Username    string `sql:"size:255;not null;index"`
	Email       string `sql:"size:255;not null;unique"`
	Uid         string `sql:"size:255;not null;index"`
	Provider    string `sql:"size:8;not null;index"`
	Password    string `sql:"size:255"`
	Details     []byte
	Logs        []Log
	Permissions []Permission
}

func (p User) TableName() string {
	return "base_users"
}

type Log struct {
	ID        uint `gorm:"primary_key"`
	User      User
	UserID    uint      `sql:"not null"`
	Message   string    `sql:"size:255;not null"`
	CreatedAt time.Time `sql:"not null"`
}

func (p Log) TableName() string {
	return "base_logs"
}

type Role struct {
	Model

	Name         string `sql:"size:255;not null;index"`
	ResourceType string `sql:"size:255;not null;index;default:''"`
	ResourceID   uint   `sql:"not null;index;default:0"`
	Permissions  []Permission
}

func (p Role) TableName() string {
	return "base_roles"
}

type Permission struct {
	Model

	User     User
	UserID   uint `sql:"not null"`
	Role     Role
	RoleID   uint      `sql:"not null"`
	StartUp  time.Time `sql:"not null;type:DATE;default:'9999-12-31'"`
	ShutDown time.Time `sql:"not null;type:DATE;default:'2015-10-27'"`
}

func (p Permission) TableName() string {
	return "base_permissions"
}

//==============================================================================
type Dao struct {
	Aes  *Aes  `inject:""`
	Hmac *Hmac `inject:""`
}

func (p *Dao) Get(db *gorm.DB, key string, val interface{}) {
}
func (p *Dao) Set(db *gorm.DB, key string, val []byte, encrypt bool) {
}
func (p *Dao) IsEmailUserExist(db *gorm.DB, email string) bool {
	var cn int
	db.Model(User{}).Where(&User{Email: email, Provider: "email"}).Count(&cn)
	return cn > 0
}
func (p *Dao) CreateEmailUser(db *gorm.DB, username, email, password string) (*User, error) {
	pwd, err := Ssha512([]byte(password), 8)
	if err != nil {
		return nil, err
	}
	user := User{
		Username: username,
		Email:    email,
		Uid:      Uuid(),
		Password: pwd,
		Provider: "email",
	}
	db.Create(&user)
	return &user, nil
}

//==============================================================================

func init() {
	var en Engine
	en = &BaseEngine{}
	Use(en)
}
