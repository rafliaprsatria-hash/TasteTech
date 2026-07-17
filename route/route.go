package route

import (
	"TasteTech/controller"
	"TasteTech/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mendefinisikan semua routes aplikasi
func SetupRoutes(app *fiber.App) {
	// ===== Halaman Web Publik (SSR) =====
	app.Get("/", controller.HomePage)
	app.Get("/menu", controller.MenuPage)
	app.Get("/promo", controller.PromoPage)
	app.Get("/keranjang", controller.CartPage)
	app.Get("/checkout", controller.CheckoutPage)
	app.Get("/pesanan/:id", controller.OrderDetailPage)

	// ===== REST API Publik =====
	api := app.Group("/api")

	// Menu API
	api.Get("/menu", controller.GetMenuAPI)
	api.Get("/menu/:id", controller.GetMenuDetailAPI)

	// Cart API
	api.Get("/cart", controller.GetCartAPI)
	api.Post("/cart/add", controller.AddToCartAPI)
	api.Put("/cart/update", controller.UpdateCartAPI)
	api.Delete("/cart/remove/:id", controller.RemoveFromCartAPI)

	// Order API
	api.Post("/order/checkout", controller.ProcessOrder)

	// ===== ADMIN ROUTES =====
	// Login (tidak butuh auth)
	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.Redirect("/admin/login")
	})
	app.Get("/admin/login", controller.AdminLoginPage)
	app.Post("/admin/login", controller.AdminLogin)
	app.Get("/admin/logout", controller.AdminLogout)

	// Protected admin routes (butuh auth)
	admin := app.Group("/admin", middleware.AdminAuth)

	// Dashboard
	admin.Get("/dashboard", controller.AdminDashboard)

	// Menu CRUD
	admin.Get("/menu", controller.AdminMenuList)
	admin.Get("/menu/tambah", controller.AdminMenuCreatePage)
	admin.Post("/menu/tambah", controller.AdminMenuCreate)
	admin.Get("/menu/:id/edit", controller.AdminMenuEditPage)
	admin.Post("/menu/:id/edit", controller.AdminMenuUpdate)
	admin.Delete("/menu/:id", controller.AdminMenuDelete)
	admin.Patch("/menu/:id/toggle", controller.AdminMenuToggle)

	// Order Management
	admin.Get("/pesanan", controller.AdminOrderList)
	admin.Post("/pesanan/:id/status", controller.AdminOrderUpdateStatus)

	// Promo Management
	admin.Get("/promo", controller.AdminPromoList)
	admin.Post("/promo/:id/toggle", controller.AdminPromoToggle)
	admin.Post("/promo/:id/set", controller.AdminPromoSet)
}
