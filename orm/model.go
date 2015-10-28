package orm

type Entity interface {
	TableName() string
}
