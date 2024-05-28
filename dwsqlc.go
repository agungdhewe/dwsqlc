package dwsqlc

type DwSqlCommand struct {
	tablename string
}

func New(tablename string) (cmd *DwSqlCommand) {
	cmd = &DwSqlCommand{
		tablename: tablename,
	}

	return cmd
}

func (cmd *DwSqlCommand) Ready() {

}
