// UNIDAD 4 - Servicios Web REST con Gin
// Concurrencia nativa en Go (goroutines por petición)
// Serialización JSON con encoding/json

package handlers

import (
	"ecommerce/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ── Utilidades ──

func ok(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"mensaje": msg, "datos": data})
}

func created(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"mensaje": msg, "datos": data})
}

func errResp(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

func paramID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResp(c, http.StatusBadRequest, "ID inválido")
		return 0, false
	}
	return id, true
}

// ═══════════════════════════════════════════════════
// SERVICIO 1: GET /api/productos
// ═══════════════════════════════════════════════════
func ListarProductos(c *gin.Context) {
	categoria := c.Query("categoria")
	lista := repository.Global.ListarProductos(categoria)
	result := make([]map[string]interface{}, len(lista))
	for i, p := range lista {
		result[i] = p.ToDict()
	}
	ok(c, strconv.Itoa(len(result))+" productos encontrados", result)
}

// ═══════════════════════════════════════════════════
// SERVICIO 2: GET /api/productos/:id
// ═══════════════════════════════════════════════════
func ObtenerProducto(c *gin.Context) {
	id, valid := paramID(c)
	if !valid {
		return
	}
	p, err := repository.Global.ObtenerProducto(id)
	if err != nil {
		errResp(c, http.StatusNotFound, err.Error())
		return
	}
	ok(c, "OK", p.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 3: POST /api/productos
// Body: {"nombre":"...","precio":0.0,"stock":0,"categoria":"..."}
// ═══════════════════════════════════════════════════
func CrearProducto(c *gin.Context) {
	var body struct {
		Nombre    string  `json:"nombre" binding:"required"`
		Precio    float64 `json:"precio" binding:"required"`
		Stock     int     `json:"stock"`
		Categoria string  `json:"categoria" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errResp(c, http.StatusBadRequest, "Campos requeridos: nombre, precio, stock, categoria")
		return
	}
	p, err := repository.Global.AgregarProducto(body.Nombre, body.Precio, body.Stock, body.Categoria)
	if err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}
	created(c, "Producto creado exitosamente", p.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 4: PUT /api/productos/:id
// Body: {"nombre":"...","precio":0.0,"stock":0}
// ═══════════════════════════════════════════════════
func ActualizarProducto(c *gin.Context) {
	id, valid := paramID(c)
	if !valid {
		return
	}
	p, err := repository.Global.ObtenerProducto(id)
	if err != nil {
		errResp(c, http.StatusNotFound, err.Error())
		return
	}

	var body struct {
		Nombre *string  `json:"nombre"`
		Precio *float64 `json:"precio"`
		Stock  *int     `json:"stock"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errResp(c, http.StatusBadRequest, "JSON inválido")
		return
	}

	if body.Nombre != nil {
		if err := p.SetNombre(*body.Nombre); err != nil {
			errResp(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if body.Precio != nil {
		if err := p.SetPrecio(*body.Precio); err != nil {
			errResp(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if body.Stock != nil {
		diff := *body.Stock - p.Stock()
		if diff > 0 {
			p.AumentarStock(diff)
		}
	}

	ok(c, "Producto actualizado", p.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 5: GET /api/clientes
// ═══════════════════════════════════════════════════
func ListarClientes(c *gin.Context) {
	lista := repository.Global.ListarClientes()
	result := make([]map[string]interface{}, len(lista))
	for i, cl := range lista {
		result[i] = cl.ToDict()
	}
	ok(c, strconv.Itoa(len(result))+" clientes", result)
}

// ═══════════════════════════════════════════════════
// SERVICIO 6: POST /api/clientes
// Body: {"nombre":"...","email":"...","telefono":"..."}
// ═══════════════════════════════════════════════════
func RegistrarCliente(c *gin.Context) {
	var body struct {
		Nombre   string `json:"nombre" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Telefono string `json:"telefono"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errResp(c, http.StatusBadRequest, "Campos requeridos: nombre, email")
		return
	}
	cl, err := repository.Global.AgregarCliente(body.Nombre, body.Email, body.Telefono)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "el email ya está registrado" {
			code = http.StatusConflict
		}
		errResp(c, code, err.Error())
		return
	}
	created(c, "Cliente registrado", cl.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 7: POST /api/pedidos
// Body: {"cliente_id":1,"items":[{"producto_id":1,"cantidad":2}]}
// ═══════════════════════════════════════════════════
func CrearPedido(c *gin.Context) {
	var body struct {
		ClienteID int `json:"cliente_id" binding:"required"`
		Items     []struct {
			ProductoID int `json:"producto_id"`
			Cantidad   int `json:"cantidad"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errResp(c, http.StatusBadRequest, "Campos requeridos: cliente_id, items")
		return
	}
	if len(body.Items) == 0 {
		errResp(c, http.StatusBadRequest, "El pedido debe tener al menos un producto")
		return
	}

	pedido, err := repository.Global.CrearPedido(body.ClienteID)
	if err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}

	for _, item := range body.Items {
		prod, err := repository.Global.ObtenerProducto(item.ProductoID)
		if err != nil {
			errResp(c, http.StatusNotFound, "Producto "+strconv.Itoa(item.ProductoID)+" no encontrado")
			return
		}
		if err := pedido.AgregarProducto(prod, item.Cantidad); err != nil {
			errResp(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	created(c, "Pedido creado exitosamente", pedido.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 8: GET /api/clientes/:id/pedidos
// ═══════════════════════════════════════════════════
func PedidosPorCliente(c *gin.Context) {
	id, valid := paramID(c)
	if !valid {
		return
	}
	if _, err := repository.Global.ObtenerCliente(id); err != nil {
		errResp(c, http.StatusNotFound, err.Error())
		return
	}
	lista := repository.Global.PedidosPorCliente(id)
	result := make([]map[string]interface{}, len(lista))
	for i, p := range lista {
		result[i] = p.ToDict()
	}
	ok(c, strconv.Itoa(len(result))+" pedidos encontrados", result)
}

// ═══════════════════════════════════════════════════
// SERVICIO 9: PATCH /api/pedidos/:id/estado
// Body: {"estado":"pagado"}
// ═══════════════════════════════════════════════════
func CambiarEstadoPedido(c *gin.Context) {
	id, valid := paramID(c)
	if !valid {
		return
	}
	p, err := repository.Global.ObtenerPedido(id)
	if err != nil {
		errResp(c, http.StatusNotFound, err.Error())
		return
	}

	var body struct {
		Estado string `json:"estado" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errResp(c, http.StatusBadRequest, "Se requiere el campo 'estado'")
		return
	}
	if err := p.CambiarEstado(body.Estado); err != nil {
		errResp(c, http.StatusBadRequest, err.Error())
		return
	}
	ok(c, "Estado cambiado a '"+body.Estado+"'", p.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 10: DELETE /api/pedidos/:id
// ═══════════════════════════════════════════════════
func CancelarPedido(c *gin.Context) {
	id, valid := paramID(c)
	if !valid {
		return
	}
	p, err := repository.Global.ObtenerPedido(id)
	if err != nil {
		errResp(c, http.StatusNotFound, err.Error())
		return
	}
	if err := p.Cancelar(repository.Global.Productos); err != nil {
		errResp(c, http.StatusForbidden, err.Error())
		return
	}
	ok(c, "Pedido cancelado. Stock devuelto.", p.ToDict())
}

// ═══════════════════════════════════════════════════
// SERVICIO 11: GET /api/estadisticas
// ═══════════════════════════════════════════════════
func Estadisticas(c *gin.Context) {
	db := repository.Global
	db.Mu().RLock()
	defer db.Mu().RUnlock()

	var ingresos float64
	porEstado := map[string]int{}

	for _, p := range db.Pedidos {
		porEstado[p.Estado()]++
		if p.Estado() == "pagado" || p.Estado() == "enviado" || p.Estado() == "entregado" {
			ingresos += p.Total()
		}
	}

	ok(c, "OK", gin.H{
		"total_productos":      len(db.Productos),
		"total_clientes":       len(db.Clientes),
		"total_pedidos":        len(db.Pedidos),
		"ingresos_confirmados": ingresos,
		"pedidos_por_estado":   porEstado,
	})
}
