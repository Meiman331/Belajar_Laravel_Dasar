package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Product struct {
	ID    int
	Name  string
	Price int
}

type productStore struct {
	mu       sync.Mutex
	nextID   int
	products []Product
}

func newStore() *productStore {
	return &productStore{
		nextID:   1,
		products: []Product{},
	}
}

func (s *productStore) Add(name string, price int) Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := Product{ID: s.nextID, Name: name, Price: price}
	s.nextID++
	s.products = append(s.products, p)
	return p
}

func (s *productStore) List() []Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Product, len(s.products))
	copy(out, s.products)
	return out
}

var tpl = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="id">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Daftar Produk</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 24px; }
    table { border-collapse: collapse; width: 100%; margin-top: 16px; }
    th, td { border: 1px solid #ddd; padding: 8px; }
    th { background: #f5f5f5; }
    .container { max-width: 900px; }
    .form-row { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
    input[type=text], input[type=number] { padding: 8px; width: 260px; }
    button { padding: 8px 14px; cursor: pointer; }
  </style>
</head>
<body>
<div class="container">
  <h1>Daftar Produk</h1>

  <form method="POST" action="/products">
    <div class="form-row">
      <label>Nama Produk</label>
      <input type="text" name="name" required>

      <label>Harga</label>
      <input type="number" name="price" min="0" step="1" required>

      <button type="submit">Tambah Produk</button>
    </div>
  </form>

  <table>
    <thead>
      <tr>
        <th>ID</th>
        <th>Nama</th>
        <th>Harga</th>
      </tr>
    </thead>
    <tbody>
    {{range .Products}}
      <tr>
        <td>{{.ID}}</td>
        <td>{{.Name}}</td>
        <td>{{.Price}}</td>
      </tr>
    {{else}}
      <tr>
        <td colspan="3">Belum ada produk</td>
      </tr>
    {{end}}
    </tbody>
  </table>
</div>
</body>
</html>`))

func main() {
	store := newStore()

	handlers := http.NewServeMux()

	// Tambahan: Store instance ada di main() dan dipakai oleh semua request.
	handlers.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	})

	handlers.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			name := r.FormValue("name")
			priceStr := r.FormValue("price")
			price, err := strconv.Atoi(priceStr)
			if err != nil {
				http.Error(w, "Harga tidak valid", http.StatusBadRequest)
				return
			}

			store.Add(name, price)
			// redirect agar halaman menampilkan daftar terbaru
			http.Redirect(w, r, "/products", http.StatusSeeOther)
			return
		default:
			// GET /products (render daftar)
			products := store.List()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tpl.Execute(w, map[string]any{"Products": products}); err != nil {
				log.Println(err)
			}
		}
	})

	addr := ":8080"
	log.Printf("Server berjalan di %s", addr)
	if err := http.ListenAndServe(addr, handlers); err != nil {
		log.Fatal(err)
	}
}
