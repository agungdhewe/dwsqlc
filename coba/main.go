package main

import (
	"reflect"

	"github.com/agungdhewe/dwsqlc"
)

type Heinv struct {
	Id  dwsqlc.DbTypeVarchar `field:"heinv_id" default:""`
	Art dwsqlc.DbTypeVarchar `field:"heinv_art" default:""`
	Mat dwsqlc.DbTypeVarchar `field:"heinv_mat" default:""`
}

func main() {

	rel := dwsqlc.Relation{
		Table:  "heinv",
		Schema: "temp",
	}

	sqlc, err := dwsqlc.New(rel, reflect.TypeOf(Heinv{}))
	if err != nil {
		panic(err.Error())
	}

	query := sqlc.CreateInsertQuery()
	//fmt.Println(query.Sql())

	// untuk update
	//keys := []string{"Id"}
	//sqlc.CreateUpdateQuery(keys)

	var model *Heinv

	// loop 1
	model = sqlc.GetModel().(*Heinv)
	model.Id = "FFF"
	model.Art = "234"
	model.Mat = dwsqlc.DbValueNull
	sqlc.CreateParameter(query, model)

	// loop 2
	model = sqlc.GetModel().(*Heinv)
	model.Id = "CC"
	sqlc.CreateParameter(query, model)
}
