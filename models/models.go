package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ─────────────────────────────────────────
// ENTIDAD BASE  (equivale a la clase padre)
// ─────────────────────────────────────────

type EntidadBase struct {
	id     int
	nombre string
}

func NewEntidadBase(id int, nombre string) (*EntidadBase, error) {
	if strings.TrimSpace(nombre) == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	return &EntidadBase{id: id, nombre: strings.TrimSpace(nombre)}, nil
}

func (e *EntidadBase) ID() int        { return e.id }
func (e *EntidadBase) Nombre() string { return e.nombre }

func (e *EntidadBase) SetNombre(nombre string) error {
	if strings.TrimSpace(nombre) == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	e.nombre = strings.TrimSpace(nombre)
	return nil
}

func (e *EntidadBase) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     e.id,
		"nombre": e.nombre,
	}
}

// ─────────────────────────────────────────
// PRODUCTO  (hereda EntidadBase via embedding)
// ─────────────────────────────────────────

type Producto struct {
	EntidadBase // Embedding = herencia en Go
	precio      float64
	stock       int
	Categoria   string
	Activo      bool
	Tipo        string // "fisico" o "digital"
	URLDescarga string
}

func NewProducto(id int, nombre string, precio float64, stock int, categoria string) (*Producto, error) {
	base, err := NewEntidadBase(id, nombre)
	if err != nil {
		return nil, err
	}
	if precio < 0 {
		return nil, errors.New("el precio no puede ser negativo")
	}
	return &Producto{
		EntidadBase: *base,
		precio:      precio,
		stock:       stock,
		Categoria:   categoria,
		Activo:      true,
		Tipo:        "fisico",
	}, nil
}

func NewProductoDigital(id int, nombre string, precio float64, categoria, url string) (*Producto, error) {
	p, err := NewProducto(id, nombre, precio, 9999, categoria)
	if err != nil {
		return nil, err
	}
	p.Tipo = "digital"
	p.URLDescarga = url
	return p, nil
}

func (p *Producto) Precio() float64 { return p.precio }
func (p *Producto) Stock() int      { return p.stock }

func (p *Producto) SetPrecio(v float64) error {
	if v < 0 {
		return errors.New("el precio no puede ser negativo")
	}
	p.precio = v
	return nil
}

// ReducirStock → encapsulación: valida antes de modificar
func (p *Producto) ReducirStock(cantidad int) error {
	if p.Tipo == "digital" {
		return nil // los digitales no se agotan
	}
	if cantidad > p.stock {
		return fmt.Errorf("stock insuficiente. Disponible: %d", p.stock)
	}
	p.stock -= cantidad
	return nil
}

func (p *Producto) AumentarStock(cantidad int) error {
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser positiva")
	}
	p.stock += cantidad
	return nil
}

// ToDict → polimorfismo: cada modelo lo sobreescribe (como to_dict en Python)
func (p *Producto) ToDict() map[string]interface{} {
	d := p.EntidadBase.ToMap()
	d["precio"] = p.precio
	d["stock"] = p.stock
	d["categoria"] = p.Categoria
	d["activo"] = p.Activo
	d["tipo"] = p.Tipo
	if p.Tipo == "digital" {
		d["url_descarga"] = p.URLDescarga
	}
	return d
}

// ─────────────────────────────────────────
// CLIENTE
// ─────────────────────────────────────────

type Cliente struct {
	EntidadBase
	email       string
	Telefono    string
	Direcciones []string
	Activo      bool
}

var emailRegex = regexp.MustCompile(`^[\w\.\-]+@[\w\.\-]+\.\w{2,}$`)

func NewCliente(id int, nombre, email, telefono string) (*Cliente, error) {
	base, err := NewEntidadBase(id, nombre)
	if err != nil {
		return nil, err
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("email inválido: %s", email)
	}
	return &Cliente{
		EntidadBase: *base,
		email:       email,
		Telefono:    telefono,
		Direcciones: []string{},
		Activo:      true,
	}, nil
}

func (c *Cliente) Email() string { return c.email }

func (c *Cliente) SetEmail(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email inválido: %s", email)
	}
	c.email = email
	return nil
}

func (c *Cliente) AgregarDireccion(dir string) {
	for _, d := range c.Direcciones {
		if d == dir {
			return
		}
	}
	c.Direcciones = append(c.Direcciones, dir)
}

func (c *Cliente) ToDict() map[string]interface{} {
	d := c.EntidadBase.ToMap()
	d["email"] = c.email
	d["telefono"] = c.Telefono
	d["direcciones"] = c.Direcciones
	d["activo"] = c.Activo
	return d
}

// ─────────────────────────────────────────
// PEDIDO
// ─────────────────────────────────────────

var EstadosValidos = []string{"pendiente", "pagado", "enviado", "entregado", "cancelado"}

type LineaPedido struct {
	Producto       *Producto
	Cantidad       int
	PrecioUnitario float64
}

func (l *LineaPedido) Subtotal() float64 {
	return float64(l.Cantidad) * l.PrecioUnitario
}

func (l *LineaPedido) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"producto_id":     l.Producto.ID(),
		"producto_nombre": l.Producto.Nombre(),
		"cantidad":        l.Cantidad,
		"precio_unitario": l.PrecioUnitario,
		"subtotal":        l.Subtotal(),
	}
}

type Pedido struct {
	EntidadBase
	ClienteID          int
	estado             string
	Lineas             []*LineaPedido
	FechaCreacion      string
	FechaActualizacion string
}

func NewPedido(id, clienteID int) *Pedido {
	now := time.Now().Format(time.RFC3339)
	return &Pedido{
		EntidadBase:        EntidadBase{id: id, nombre: fmt.Sprintf("Pedido-%d", id)},
		ClienteID:          clienteID,
		estado:             "pendiente",
		Lineas:             []*LineaPedido{},
		FechaCreacion:      now,
		FechaActualizacion: now,
	}
}

func (p *Pedido) Estado() string { return p.estado }

func (p *Pedido) Total() float64 {
	var total float64
	for _, l := range p.Lineas {
		total += l.Subtotal()
	}
	return total
}

func (p *Pedido) AgregarProducto(producto *Producto, cantidad int) error {
	if p.estado != "pendiente" {
		return errors.New("solo se pueden modificar pedidos en estado 'pendiente'")
	}
	if err := producto.ReducirStock(cantidad); err != nil {
		return err
	}
	p.Lineas = append(p.Lineas, &LineaPedido{
		Producto:       producto,
		Cantidad:       cantidad,
		PrecioUnitario: producto.Precio(),
	})
	p.actualizarFecha()
	return nil
}

func (p *Pedido) CambiarEstado(nuevo string) error {
	for _, e := range EstadosValidos {
		if e == nuevo {
			p.estado = nuevo
			p.actualizarFecha()
			return nil
		}
	}
	return fmt.Errorf("estado inválido. Opciones: %v", EstadosValidos)
}

func (p *Pedido) Cancelar(productosRepo map[int]*Producto) error {
	if p.estado == "entregado" || p.estado == "cancelado" {
		return fmt.Errorf("no se puede cancelar un pedido '%s'", p.estado)
	}
	for _, l := range p.Lineas {
		if prod, ok := productosRepo[l.Producto.ID()]; ok {
			prod.AumentarStock(l.Cantidad)
		}
	}
	p.estado = "cancelado"
	p.actualizarFecha()
	return nil
}

func (p *Pedido) actualizarFecha() {
	p.FechaActualizacion = time.Now().Format(time.RFC3339)
}

func (p *Pedido) ToDict() map[string]interface{} {
	lineas := make([]map[string]interface{}, len(p.Lineas))
	for i, l := range p.Lineas {
		lineas[i] = l.ToDict()
	}
	return map[string]interface{}{
		"id":                  p.ID(),
		"nombre":              p.Nombre(),
		"cliente_id":          p.ClienteID,
		"estado":              p.estado,
		"lineas":              lineas,
		"total":               p.Total(),
		"fecha_creacion":      p.FechaCreacion,
		"fecha_actualizacion": p.FechaActualizacion,
	}
}
