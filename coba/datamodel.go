package main

// untuk testing pembuatan datamodel

type Toko struct {
	Id     string `field:"toko_id"`
	Nama   string `field:"name"`
	Alamat string `field:"alamat"`
	Kota   string `field:"kota"`
	Aktif  bool   `field:"isdisabled"`
	Luas   int    `field:"luas"`
}
