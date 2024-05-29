# SqlCommand Modul
modul ini untuk mempermudah penulisan sintak SQL Insert, Update dan Delete.

## Contoh Program

Apabila kita akan melakukan query insert/update ke tabel karyawan yang ada di schema hr

Kita siapkan dulu data modelnya

	type Karyawan struct {
		Id     dwsqlc.DbTypeVarchar `field:"id_karyawan" default:""`
		Nama   dwsqlc.DbTypeVarchar `field:"nama" default:""`
		Alamat dwsqlc.DbTypeVarchar `field:"alamat" default:""`
	}

Selanjutnya kita siapkan parameter relasi, untuk mereferensikan model di atas dengan tabel di RDBMS.

	rel := dwsqlc.Relation{
		Table:  "karyawan",
		Schema: "hr",
	}

Kemudian kita siapkan untuk querynya

	sqlc, err := dwsqlc.New(rel, reflect.TypeOf(Karyawan{}))
	if err != nil {
		panic(err.Error())
	}

Untuk mebuat query Insert, kita menggunakan fungsi `CreateInsertQuery`

	query := sqlc.CreateInsertQuery()
	fmt.Println(query.Sql())

atau jika hanya memilih untuk field tertentu

	query := sqlc.CreateInsertQuery("Nama", "Alamat")
	fmt.Println(query.Sql())

untuk update, kita menggunakan fungsi `CreateUpdateQuery` dengan memberikan parameter keys

	keys := []string{"Id"}
	sqlc.CreateUpdateQuery(keys)

atau untuk field-field tertentu

	keys := []string{"Id"}
	sqlc.CreateUpdateQuery(keys, "Nama", "Alamat")

Untuk membuat parameter query, laukan sebagai berikut

	model = sqlc.GetModel().(*Karyawan)
	model.Id = "240567"
	model.Nama = "Agung Nugroho"
	model.Alamat = dwsqlc.DbValueNull
	params := sqlc.CreateParameter(query, model)