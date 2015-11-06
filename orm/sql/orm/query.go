package orm

type Query struct {
	table string
	where Where
}

func (p *Query) From(table string) *Query {
	p.table = table
	return p
}

func (p *Query) Where(w *Where) *Qu {
	p.where = w
}

func (p *Query) Select(columns ...string) {
}

func (p *Query) Count() {
}

func (p *Query) Insert(map[string]interface{}) {
}

func (p *Query) Update(map[string]interface{}) {
}
