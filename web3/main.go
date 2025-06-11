package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Struktur data sesuai tabel kriteria
type Kriteria struct {
	Kode  string  `json:"kode"`
	Nama  string  `json:"nama"`
	Bobot float64 `json:"bobot"`
	Tipe  string  `json:"tipe"`
}

// Struktur data sesuai tabel alternatif
type Alternatif struct {
	Kode          string  `json:"kode"`
	Nama          string  `json:"nama"`
	LuasLahan     float64 `json:"luas_lahan"`
	Penghasilan   float64 `json:"penghasilan"`
	HasilPanen    float64 `json:"hasil_panen"`
	LamaUsaha     float64 `json:"lama_usaha_tani"`
	JumlahAnggota float64 `json:"jumlah_anggota_keluarga"`
}

func main() {
	// Koneksi ke database MariaDB/MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/bantuan_petani")
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database:", err)
	}
	defer db.Close()

	// Cek apakah database merespon
	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Database tidak merespon:", err)
	}
	log.Println("✅ Koneksi ke database berhasil.")

	// Endpoint API /kriteria (GET, POST, PUT, DELETE)
	http.HandleFunc("/kriteria", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var k Kriteria
			err := json.NewDecoder(r.Body).Decode(&k)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			_, err = db.Exec("INSERT INTO kriteria (kode, nama, bobot, tipe) VALUES (?, ?, ?, ?)",
				k.Kode, k.Nama, k.Bobot, k.Tipe)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		if r.Method == "PUT" {
			var k Kriteria
			err := json.NewDecoder(r.Body).Decode(&k)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			_, err = db.Exec("UPDATE kriteria SET nama=?, bobot=?, tipe=? WHERE kode=?",
				k.Nama, k.Bobot, k.Tipe, k.Kode)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "DELETE" {
			kode := r.URL.Query().Get("kode")
			if kode == "" {
				http.Error(w, "Kode kriteria kosong", 400)
				return
			}
			_, err := db.Exec("DELETE FROM kriteria WHERE kode = ?", kode)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		// GET
		rows, err := db.Query("SELECT kode, nama, bobot, tipe FROM kriteria")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var hasil []Kriteria
		for rows.Next() {
			var k Kriteria
			err := rows.Scan(&k.Kode, &k.Nama, &k.Bobot, &k.Tipe)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			hasil = append(hasil, k)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hasil)
	})

	// Endpoint API /alternatif (GET, POST, PUT, DELETE)
	http.HandleFunc("/alternatif", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var a Alternatif
			err := json.NewDecoder(r.Body).Decode(&a)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			_, err = db.Exec("INSERT INTO alternatif (kode, nama, luas_lahan, penghasilan, hasil_panen, lama_usaha_tani, jumlah_anggota_keluarga) VALUES (?, ?, ?, ?, ?, ?, ?)",
				a.Kode, a.Nama, a.LuasLahan, a.Penghasilan, a.HasilPanen, a.LamaUsaha, a.JumlahAnggota)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		if r.Method == "PUT" {
			var a Alternatif
			err := json.NewDecoder(r.Body).Decode(&a)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			_, err = db.Exec("UPDATE alternatif SET nama=?, luas_lahan=?, penghasilan=?, hasil_panen=?, lama_usaha_tani=?, jumlah_anggota_keluarga=? WHERE kode=?",
				a.Nama, a.LuasLahan, a.Penghasilan, a.HasilPanen, a.LamaUsaha, a.JumlahAnggota, a.Kode)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "DELETE" {
			kode := r.URL.Query().Get("kode")
			if kode == "" {
				http.Error(w, "Kode alternatif kosong", 400)
				return
			}
			_, err := db.Exec("DELETE FROM alternatif WHERE kode = ?", kode)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		// GET
		rows, err := db.Query("SELECT kode, nama, luas_lahan, penghasilan, hasil_panen, lama_usaha_tani, jumlah_anggota_keluarga FROM alternatif")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var hasil []Alternatif
		for rows.Next() {
			var a Alternatif
			err := rows.Scan(&a.Kode, &a.Nama, &a.LuasLahan, &a.Penghasilan, &a.HasilPanen, &a.LamaUsaha, &a.JumlahAnggota)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			hasil = append(hasil, a)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hasil)
	})

	// Serve file statis dari folder public
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Jalankan server
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Printf("🚀 Server jalan di http://localhost:%s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("❌ Gagal menjalankan server:", err)
	}
}
