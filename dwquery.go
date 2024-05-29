package dwsqlc

type DwQuery struct {
	sql    string
	fields []string
}

func (query *DwQuery) Sql() string {
	return query.sql
}
