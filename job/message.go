package job

import(
	"time"

	"github.com/itpkg/ksana/utils"
)

type Message struct{
	Id string `json:"id"`
	Body []byte `json:"body"`
	Created time.Time `json:"created"`
}

func (p *Message) Parse(args...interface{}) error{
	return utils.FromJson(p.Body, &args)
}

//==============================================================================

func NewMessage(args ... interface{}) (*Message, error){
	buf, err := utils.ToJson(args)
	if err!=nil{
		return nil, err
	}
	return &Message{
		Id: utils.Uuid(),
		Body: buf,
		Created: time.Now(),
	}, nil
}
