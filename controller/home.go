package controller

import (
	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

// HomeData adalah data untuk template halaman beranda
type HomeData struct {
	Title       string
	Featured    []model.Menu
	Categories  []model.Category
	CartCount   int
}

// HomePage menampilkan halaman beranda
func HomePage(c *fiber.Ctx) error {
	// Ambil cart dari session (sederhana: dari cookie/query)
	cartCount := getCartCount(c)

	data := fiber.Map{
		"Title":      "TasteTech - Pesan Makanan Lezat",
		"Featured":   model.GetFeaturedMenu(),
		"Categories": model.Categories,
		"CartCount":  cartCount,
		"Page":       "home",
	}

	return c.Render("index", data, "layout")
}

// getCartCount mengambil jumlah item di cart dari session/cookie
func getCartCount(c *fiber.Ctx) int {
	// Untuk demo, gunakan query parameter atau default 0
	return 0
}
