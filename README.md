# SqlCommand Modul
modul ini untuk mempermudah penulisan sintak sql Insert dan Update.

Tambahkan ke modul dengan
	
	go get github.com/agungdhewe/dwsqlc


## Contoh Program

Apabila kita akan melakukan query insert/update ke tabel karyawan yang ada di schema hr

### Struktur Model

Untuk menggunakan modul ini, terlebih dahulu kita buat struktur model yang akan merepresentasikan tabel yang ada di database.
Masing-masing field harus dimapping ke field yang ada di tabel database.
Pada struktur di bawah, `Karyawan`.`Id` akan di mapping ke filed `id_karyawan` yang ada di database.

	type Karyawan struct {
		Id     dwsqlc.DbTypeVarchar `field:"id_karyawan" default:""`
		Nama   dwsqlc.DbTypeVarchar `field:"nama" default:""`
		Alamat dwsqlc.DbTypeVarchar `field:"alamat" default:""`
	}

### Koneksi Database Menggunakan Standard Modul `database/sql`

Buat koneksi ke database

	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=fgodblocal user=fgta password=inipasswordnya")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()



### Inisiasi SqlCommand dengan Model

Contoh disini akan insert/update ke tabel tbl_karyawan yang ada di database. Kita akan menghubungkan tabel tersebut ke struktur Karyawan yang telah dibuat sebelumnya.

	sqlc, err := dwsqlc.New("tbl_karyawan", &Karyawan{})
	if err != nil {
		panic(err.Error())
	}
	sqlc.Connect(conn)


### Insert / Update ke Database

Untuk melakukan query, cukup panggil fungsi `Insert` atau `Update` 

#### Insert
	
	res, err := sqlc.Insert(&Heinv{
		Id:    "240567",
		Nama:   "Agung Nugroho",
		Descr: "ini test insert",
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Data Inserted")

#### Update

	keys := []string{"Id"}
	res, err := sqlc.Update(&Heinv{
		Id:    "240567",
		Nama:   "Agung Nugroho",
		Descr: "descr diupdate",
	}, keys)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Data Updated")

#### Delete

	keys := []string{"Id"}
	res, err := sqlc.Delete(&Heinv{
		Id:    "240567",
	}, keys)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Data Delete")

#### <span style="color:red">*Select*</span>

*Untuk keperluan efisiensi program, sampai versi ini modul ini tidak direncanakan untuk melakukan operasi select.*


### Menggunakan Transaksi

Apabila proses ini akan menggunakan transaksi, terlebih dahulu buat transaksi pada context.

	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		panic(err.Error())
	}

Kemudian, hubungkan dengan transaksi sebelum melakukan query, dan jangan lupa lakukan commit setelah Insert.

	sqlc.SetTransaction(tx)
	res, err := sqlc.Insert(&Heinv{
		Id:    "240567",
		Nama:   "Agung Nugroho",
		Descr: "ini test insert",
	})
	if err != nil {
		panic(err.Error())
	}
	tx.Commit()
	fmt.Println("Data Inserted")

### Menggunakan cara yang lebih komplek

Berikut ini apabila akan menggukan cara yang lebih komplek untuk mendapatkan control yang lebih terhadap parameter dan query SQL.


#### Inisiasi SqlCommand dengan struktur `Relation`

Mereferensikan model di atas dengan tabel di RDBMS.

	rel := dwsqlc.Relation{
		Table:  "karyawan",
		Schema: "hr",
	}

Inisiasi SqlCommand

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
	model.Alamat = "Tangerang Raya" 
	params := sqlc.CreateParameter(query, model)


Atau bisa juga dengan cara berikut

	params := sqlc.CreateParameter(query, &Karyawan{
		Id: "240567",
		Nama: "Agung Nugroho",
		Alamat: "Tangerang Raya", 
	})