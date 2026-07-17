package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"TasteTech/model"

	"github.com/gofiber/fiber/v2"
)

const cartCookieName = "TasteTech_cart"

// getCartFromCookie mengambil data cart dari cookie
func getCartFromCookie(c *fiber.Ctx) *model.Cart {
	cart := &model.Cart{Items: []model.CartItem{}}
	cookieVal := c.Cookies(cartCookieName)
	if cookieVal == "" {
		return cart
	}
	if err := json.Unmarshal([]byte(cookieVal), cart); err != nil {
		return &model.Cart{Items: []model.CartItem{}}
	}
	return cart
}

// saveCartToCookie menyimpan cart ke cookie
func saveCartToCookie(c *fiber.Ctx, cart *model.Cart) {
	data, err := json.Marshal(cart)
	if err != nil {
		return
	}
	c.Cookie(&fiber.Cookie{
		Name:    cartCookieName,
		Value:   string(data),
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})
}

// CartPage menampilkan halaman keranjang belanja
func CartPage(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)

	data := fiber.Map{
		"Title":       "Keranjang - TasteTech",
		"Cart":        cart,
		"CartCount":   cart.ItemCount,
		"DeliveryFee": 5000.0,
		"GrandTotal":  cart.Total + 5000.0,
		"Page":        "cart",
	}

	return c.Render("cart", data, "layout")
}

// AddToCartAPI menambahkan item ke keranjang
func AddToCartAPI(c *fiber.Ctx) error {
	type AddRequest struct {
		MenuID   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	}

	var req AddRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	menu := model.GetMenuByID(req.MenuID)
	if menu == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Menu tidak ditemukan",
		})
	}

	if !menu.IsAvailable {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Menu sedang tidak tersedia",
		})
	}

	cart := getCartFromCookie(c)
	cart.AddToCart(menu, req.Quantity)
	saveCartToCookie(c, cart)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    fmt.Sprintf("%s berhasil ditambahkan ke keranjang", menu.Name),
		"cart_count": cart.ItemCount,
		"cart_total": cart.Total,
	})
}

// UpdateCartAPI mengupdate kuantitas item di keranjang
func UpdateCartAPI(c *fiber.Ctx) error {
	type UpdateRequest struct {
		MenuID   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid",
		})
	}

	cart := getCartFromCookie(c)
	cart.UpdateQuantity(req.MenuID, req.Quantity)
	saveCartToCookie(c, cart)

	return c.JSON(fiber.Map{
		"success":    true,
		"cart":       cart,
		"cart_count": cart.ItemCount,
		"cart_total": cart.Total,
	})
}

// RemoveFromCartAPI menghapus item dari keranjang
func RemoveFromCartAPI(c *fiber.Ctx) error {
	idStr := c.Params("id")
	menuID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	cart := getCartFromCookie(c)
	cart.RemoveItem(menuID)
	saveCartToCookie(c, cart)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Item berhasil dihapus dari keranjang",
		"cart_count": cart.ItemCount,
		"cart_total": cart.Total,
	})
}

// GetCartAPI mengembalikan isi keranjang dalam format JSON
func GetCartAPI(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	deliveryFee := 5000.0
	grandTotal := cart.Total + deliveryFee

	return c.JSON(fiber.Map{
		"success":      true,
		"cart":         cart,
		"delivery_fee": deliveryFee,
		"grand_total":  grandTotal,
	})
}

// CheckoutPage menampilkan halaman checkout
func CheckoutPage(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	if cart.IsEmpty() {
		return c.Redirect("/keranjang")
	}

	deliveryFee := 5000.0
	data := fiber.Map{
		"Title":       "Checkout - TasteTech",
		"Cart":        cart,
		"CartCount":   cart.ItemCount,
		"DeliveryFee": deliveryFee,
		"GrandTotal":  cart.Total + deliveryFee,
		"Page":        "checkout",
	}

	return c.Render("checkout", data, "layout")
}

// ProcessOrder memproses pesanan
func ProcessOrder(c *fiber.Ctx) error {
	cart := getCartFromCookie(c)
	if cart.IsEmpty() {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Keranjang kosong",
		})
	}

	var req model.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if req.CustomerName == "" || req.Phone == "" || req.Address == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Nama, telepon, dan alamat wajib diisi",
		})
	}

	deliveryFee := 5000.0
	orderID := fmt.Sprintf("GF-%d", time.Now().UnixNano()%1000000)

	order := &model.Order{
		ID:            orderID,
		CustomerName:  req.CustomerName,
		Phone:         req.Phone,
		Address:       req.Address,
		Notes:         req.Notes,
		Items:         cart.Items,
		Total:         cart.Total,
		DeliveryFee:   deliveryFee,
		GrandTotal:    cart.Total + deliveryFee,
		PaymentMethod: req.PaymentMethod,
		Status:        model.StatusPending,
		CreatedAt:     time.Now(),
	}

	model.SaveOrder(order)

	// Kosongkan cart setelah order berhasil
	emptyCart := &model.Cart{Items: []model.CartItem{}}
	saveCartToCookie(c, emptyCart)

	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Pesanan berhasil dibuat!",
		"order_id": orderID,
		"order":    order,
	})
}

// OrderDetailPage menampilkan halaman detail pesanan
func OrderDetailPage(c *fiber.Ctx) error {
	orderID := c.Params("id")
	order := model.GetOrderByID(orderID)

	if order == nil {
		return c.Status(404).Render("404", fiber.Map{
			"Title": "Pesanan tidak ditemukan",
		}, "layout")
	}

	data := fiber.Map{
		"Title":  fmt.Sprintf("Pesanan %s - TasteTech", orderID),
		"Order":  order,
		"Page":   "order",
		"CartCount": 0,
	}

	return c.Render("order", data, "layout")
}
