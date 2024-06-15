package dwsqlc

type DwQuery struct {
	sql    string
	fields []string
}

// Sql returns the SQL query string of the DwQuery.
//
// It does not take any parameters.
// It returns a string.
func (query *DwQuery) Sql() string {
	return query.sql
}
