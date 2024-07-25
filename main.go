package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v2"
	"github.com/joho/godotenv"
)

// Book struct to hold book data
// struct tags in golang (json)

type Menu struct {
	MenuId   int    `json:"menuId"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Category string `json:"category"`
}

type Order struct {
	OrderId  int    `json:"orderId"`
	Quantity int    `json:"quantity"`
	Date     string `json:"date"`
	Status   string `json:"status"`
}

// Dummy user for example
var user = struct {
	Email    string
	Password string
}{
	Email:    "user@example.com",
	Password: "password123",
}

// UserData represents the user data extracted from the JWT token
type UserData struct {
	Email string
	Role  string
}

// userContextKey is the key used to store user data in the Fiber context
const userContextKey = "admin"

var Menus []Menu   // Slice to store Menu
var Orders []Order // Slice to store Order

func main() {
	app := fiber.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error dot env file")
	}

	// JWT Secret Key


	// Initialize in-memory data
	Menus = append(Menus, Menu{MenuId: 1, Name: "egg", Price: 20, Category: "meat diet"})
	Orders = append(Orders, Order{OrderId: 1, Quantity: 5, Date: "22-12-2025", Status: "doing"})

	// Login route
	app.Post("/login", login(os.Getenv("secretKey")))

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{ // make Token 
		SigningKey: []byte(os.Getenv("secretKey")),
	}))
	
	// Middleware to extract user data from JWT
	app.Use(extractUserFromJWT) // check role form JWT Middleware 

	// Group routes under /book
	MenuGroup := app.Group("/Menu")
	OrdersGroup := app.Group("/Orders")
	// Apply the isAdmin middleware only to the /book routes
	MenuGroup.Use(isAdmin)

	// CRUD
	MenuGroup.Get("/", getMenus)
	MenuGroup.Get("/:id", getMenu)
	MenuGroup.Post("/", createMenu)
	MenuGroup.Put("/:id", updateMenu)
	MenuGroup.Delete("/:id", deleteMenu)

	OrdersGroup.Get("/Orders", getOrders)
	OrdersGroup.Get("/Orders", updateOrder)

	OrdersGroup.Post("/upload", uploadImage)

	// Use the environment variable for the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "Port Error" // Default port if not specified
	}

	// port sever
	app.Listen(":" + port)

}
