package ksana

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
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

type BaseEngine struct {
}

func (p *BaseEngine) Router() {
}

func (p *BaseEngine) Migrate() error {
	return nil
}

func (p *BaseEngine) Job() {
}

func (p *BaseEngine) Deploy() {
}

func (p *BaseEngine) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "openssl",
			Aliases: []string{"ssl"},
			Usage:   "Openssl certs generate",
			Flags: []cli.Flag{
				KSANA_ENV,
				cli.StringFlag{
					Name:  "name",
					Value: "whoami",
					Usage: "certificate name",
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
					if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.cert.pem", root, name), cert, 0644); err != nil {
						log.Fatal(err)
					}
					if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.key.pem", root, name), key, 0600); err != nil {
						log.Fatal(err)
					}
					log.Println("Done.")
				} else {
					log.Fatal(err)
				}

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

func init() {
	var en Engine
	en = &BaseEngine{}
	Use(en)
}
