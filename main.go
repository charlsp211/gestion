package main

import (
	"ecommerce/handlers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ── CORS ──
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// ── Frontend ──
	r.StaticFS("/static", http.Dir("/home/carlos/gestion/static"))
	r.GET("/", func(c *gin.Context) {
		c.File("/home/carlos/gestion/static/index.html")
	})

	// ── API Routes ──
	api := r.Group("/api")
	{
		api.GET("/productos", handlers.ListarProductos)
		api.GET("/productos/:id", handlers.ObtenerProducto)
		api.POST("/productos", handlers.CrearProducto)
		api.PUT("/productos/:id", handlers.ActualizarProducto)
		api.GET("/clientes", handlers.ListarClientes)
		api.POST("/clientes", handlers.RegistrarCliente)
		api.GET("/clientes/:id/pedidos", handlers.PedidosPorCliente)
		api.POST("/pedidos", handlers.CrearPedido)
		api.PATCH("/pedidos/:id/estado", handlers.CambiarEstadoPedido)
		api.DELETE("/pedidos/:id", handlers.CancelarPedido)
		api.GET("/estadisticas", handlers.Estadisticas)
	}

	fmt.Println("════════════════════════════════════════")
	fmt.Println("  NexShop E-Commerce API — Go + Gin")
	fmt.Println("  http://localhost:5000")
	fmt.Println("════════════════════════════════════════")

	r.Run(":5000")
}
