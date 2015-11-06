package sitemap

import (
	"encoding/xml"
	"io"
	"time"
)

type Frequency string

const (
	Always  Frequency = "always"
	Hourly  Frequency = "hourly"
	Dialy   Frequency = "dialy"
	Weekly  Frequency = "weekly"
	Monthly Frequency = "monthly"
	Yearly  Frequency = "yearly"
	Never   Frequency = "never"
)

type Item struct {
	XMLName xml.Name `xml:"url"`

	Url       string    `xml:"loc"`
	Updated   time.Time `xml:"lastmod"`
	Frequency Frequency `xml:"changefreq"`
	Priority  float32   `xml:"priority"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`

	Items []*Item
}

func (p *Sitemap) Add(url string, updated time.Time, frequency Frequency, priority float32) {
	p.Items = append(p.Items, &Item{
		Url:       url,
		Updated:   updated,
		Frequency: frequency,
		Priority:  priority,
	})
}

func (p *Sitemap) Write(w io.Writer, pretty bool) error {
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	if pretty {
		enc.Indent("", "  ")
	}
	return enc.Encode(p)
}

func (p *Sitemap) Xml(pretty bool) (buf []byte, err error) {

	if pretty {
		buf, err = xml.MarshalIndent(p, "  ", "    ")
	} else {
		buf, err = xml.Marshal(p)
	}

	if err == nil {
		buf = append([]byte(xml.Header), buf...)
	}

	return
}

//==============================================================================
func New() *Sitemap {
	return &Sitemap{
		Items: make([]*Item, 0),
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
}
