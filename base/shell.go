package base

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	ks "github.com/itpkg/ksana"
	"github.com/jrallison/go-workers"
	"github.com/robfig/cron"
)

func (p *BaseEngine) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "assets",
			Aliases: []string{"a"},
			Usage:   "assets operations",
			Subcommands: []cli.Command{
				{
					Name:    "build",
					Aliases: []string{"b"},
					Usage:   "compile all the assets.",
					Flags: []cli.Flag{
						ks.KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						//todo
					},
				},
			},
		},
		{
			Name:    "routes",
			Aliases: []string{"ro"},
			Usage:   "Print out all defined routes in match order, with names",
			Action: func(c *cli.Context) {
				router := mux.NewRouter()
				if err := ks.LoopEngine(func(en ks.Engine) error {
					en.Mount(router)
					return nil
				}); err != nil {
					log.Fatal(err)
				}
				//todo

			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the ksana server",
			Flags: []cli.Flag{
				ks.KSANA_ENV,
			},
			Action: func(c *cli.Context) {
				app, err := ks.New(c)
				if err != nil {
					log.Fatal(err)
				}
				router := mux.NewRouter()
				if err = ks.LoopEngine(func(en ks.Engine) error {
					en.Mount(router)
					return nil
				}); err != nil {
					log.Fatal(err)
				}

				http.Handle("/", router)

				if err = http.ListenAndServe(fmt.Sprintf(":%d", app.Cfg.Http.Port), router); err != nil {
					log.Fatal(err)
				}

			},
		},
		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "starting background job",
			Flags: []cli.Flag{
				ks.KSANA_ENV,
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
				if _, err := ks.New(c); err != nil {
					log.Fatal(err)
				}
				if c.Bool("dispatcher") {
					log.Println("start cron job...")
					cn := cron.New()
					if err := ks.LoopEngine(func(en ks.Engine) error {
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

				if err := ks.LoopEngine(func(en ks.Engine) error {
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
						ks.KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						_, err := ks.New(c)
						if err != nil {
							log.Fatal(err)
						}
						if err = ks.LoopEngine(func(en ks.Engine) error {
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
						ks.KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						_, err := ks.New(c)
						if err != nil {
							log.Fatal(err)
						}
						if err = ks.LoopEngine(func(en ks.Engine) error {
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
						ks.KSANA_ENV,
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
						if err := ks.Mkdirs(fmt.Sprintf("config/%s", env)); err != nil {
							log.Fatal(err)
						}

						if _, err := os.Stat(fn); err == nil {
							log.Fatalf("%s already exist!", fn)
						}
						buf, err := ks.RandomBytes(512)
						if err != nil {
							log.Fatal(err)
						}

						cfg := ks.Configuration{Secrets: buf,
							Http: ks.HttpCfg{
								Host: "localhost",
								Port: 8080,
							},
							Database: ks.DatabaseCfg{
								Dialect: db,
								Pool: ks.PoolCfg{
									MaxIdle: 6,
									MaxOpen: 180,
								},
							},
							Redis: ks.RedisCfg{
								Host: "localhost",
								Port: 6379,
								Db:   0,
								Pool: ks.PoolCfg{
									MaxIdle: 4,
									MaxOpen: 120,
								},
							},
							Elasticsearch: ks.ElasticsearchCfg{
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
						ks.KSANA_ENV,
					},
					Action: func(c *cli.Context) {
						cfg, err := ks.Load(c)
						if err != nil {
							log.Fatal(err)
						}

						if err := ks.Mkdirs(fmt.Sprintf("config/%s", cfg.Env)); err != nil {
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
						ks.KSANA_ENV,
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

						if err := ks.Mkdirs(root); err != nil {
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
