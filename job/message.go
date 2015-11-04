package job

import (
	"fmt"
	"time"

	"github.com/itpkg/ksana/utils"
)

type Message struct {
	Id      string    `json:"id"`
	Title   string    `json:"title"`
	Body    []byte    `json:"body"`
	Created time.Time `json:"created"`
}

func (p *Message) String() string {
	return fmt.Sprintf("[%s] %s", p.Id, p.Title)
}
func (p *Message) Parse(args ...interface{}) error {
	return utils.FromJson(p.Body, &args)
}

//==============================================================================

func NewMessage(args ...interface{}) (*Message, error) {
	buf, err := utils.ToJson(args)
	if err != nil {
		return nil, err
	}
	return &Message{
		Id:      utils.Uuid(),
		Body:    buf,
		Created: time.Now(),
	}, nil
}
