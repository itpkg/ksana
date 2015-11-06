package atom

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`

	Id       string    `xml:"id"`
	Title    string    `xml:"title"`
	SubTitle string    `xml:"subtitle"`
	Link     string    `xml:"link"`
	Updated  time.Time `xml:"updated"`

	Entries []*Entry

	Xmlns string `xml:"xmlns,attr"`
}

func (p *Feed) Add(e *Entry) {
	p.Entries = append(p.Entries, e)
}

func (p *Feed) Write(w io.Writer, pretty bool) error {
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	if pretty {
		enc.Indent("", "  ")
	}
	return enc.Encode(p)
}

func (p *Feed) Xml(pretty bool) (buf []byte, err error) {

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

type Author struct {
	XMLName xml.Name `xml:"author"`

	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`

	Id      string    `xml:"id"`
	Title   string    `xml:"title"`
	Link    string    `xml:"link"`
	Updated time.Time `xml:"updated"`
	Summary string    `xml:"summary"`

	Author  *Author
	Content *Content
}

func (p *Entry) SetAuthor(name, email string) {
	p.Author = &Author{Name: name, Email: email}
}

func (p *Entry) SetContent(html string) {
	p.Content = &Content{
		Type: "xhtml",
		Div: &Div{
			P:     html,
			Xmlns: "http://www.w3.org/1999/xhtml",
		},
	}
}

type Content struct {
	XMLName xml.Name `xml:"content"`
	Type    string   `xml:"xmlns,attr"`
	Div     *Div
}

type Div struct {
	XMLName xml.Name `xml:"div"`

	P     string `xml:"p"`
	Xmlns string `xml:"xmlns,attr"`
}

//==============================================================================

func NewEntry(id, title, link, summary string, updated time.Time) *Entry {
	return &Entry{
		Id:      id,
		Title:   title,
		Link:    link,
		Summary: summary,
		Updated: updated,
	}
}

func NewFeed(id, title, subtitle, link string) *Feed {
	return &Feed{
		Id:       id,
		Title:    title,
		SubTitle: subtitle,
		Link:     link,
		Updated:  time.Now(),
		Xmlns:    "http://www.w3.org/2005/Atom",
	}
}

func Link(name, title string) string {
	return fmt.Sprintf(`<link href="%s" type="application/atom+xml" rel="alternate" title="%s" />`, name, title)
}
