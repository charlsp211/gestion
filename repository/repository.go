// UNIDAD 2 - Servicios y repositorios
// Base de datos en memoria (equivalente a base_datos.py en Python)

package repository

import (
	"ecommerce/models"
	"errors"
	"sync"
)

// DB es el repositorio central con mutex para concurrencia (Unidad 4)
type DB struct {
	mu           sync.RWMutex
	Productos    map[int]*models.Producto
	Clientes     map[int]*models.Cliente
	Pedidos      map[int]*models.Pedido
	lastProdID   int
	lastClientID int
	lastPedidoID int
}

func NewDB() *DB {
	db := &DB{
		Productos: make(map[int]*models.Producto),
		Clientes:  make(map[int]*models.Cliente),
		Pedidos:   make(map[int]*models.Pedido),
	}
	db.cargarDatosPrueba()
	return db
}

func (db *DB) cargarDatosPrueba() {
	// Productos físicos
	db.AgregarProducto("Laptop HP 15", 799.99, 10, "Tecnología")
	db.AgregarProducto("Mouse Inalámbrico", 25.50, 50, "Tecnología")
	db.AgregarProducto("Teclado Mecánico", 89.99, 30, "Tecnología")
	db.AgregarProducto("Monitor 24 pulgadas", 349.99, 15, "Tecnología")
	db.AgregarProducto("Mochila Laptop", 45.00, 40, "Accesorios")
	db.AgregarProducto("Audífonos Bluetooth", 65.00, 25, "Audio")

	// Producto digital
	db.lastProdID++
	p, _ := models.NewProductoDigital(
		db.lastProdID,
		"Curso Python Avanzado",
		29.99,
		"Educación",
		"https://mitienda.com/descargas/python-avanzado",
	)
	db.Productos[p.ID()] = p

	// Clientes
	db.AgregarCliente("Ana García", "ana@email.com", "0991234567")
	db.AgregarCliente("Luis Pérez", "luis@email.com", "0987654321")
}

// ── Generadores de ID ──

func (db *DB) nextProdID() int   { db.lastProdID++; return db.lastProdID }
func (db *DB) nextClientID() int { db.lastClientID++; return db.lastClientID }
func (db *DB) nextPedidoID() int { db.lastPedidoID++; return db.lastPedidoID }

// ── CRUD Productos ──

func (db *DB) AgregarProducto(nombre string, precio float64, stock int, categoria string) (*models.Producto, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	p, err := models.NewProducto(db.nextProdID(), nombre, precio, stock, categoria)
	if err != nil {
		return nil, err
	}
	db.Productos[p.ID()] = p
	return p, nil
}

func (db *DB) ObtenerProducto(id int) (*models.Producto, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	p, ok := db.Productos[id]
	if !ok {
		return nil, errors.New("producto no encontrado")
	}
	return p, nil
}

func (db *DB) ListarProductos(categoria string) []*models.Producto {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var lista []*models.Producto
	for _, p := range db.Productos {
		if !p.Activo {
			continue
		}
		if categoria != "" && p.Categoria != categoria {
			continue
		}
		lista = append(lista, p)
	}
	return lista
}

// ── CRUD Clientes ──

func (db *DB) AgregarCliente(nombre, email, telefono string) (*models.Cliente, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	// Verificar email único
	for _, c := range db.Clientes {
		if c.Email() == email {
			return nil, errors.New("el email ya está registrado")
		}
	}
	c, err := models.NewCliente(db.nextClientID(), nombre, email, telefono)
	if err != nil {
		return nil, err
	}
	db.Clientes[c.ID()] = c
	return c, nil
}

func (db *DB) ObtenerCliente(id int) (*models.Cliente, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	c, ok := db.Clientes[id]
	if !ok {
		return nil, errors.New("cliente no encontrado")
	}
	return c, nil
}

func (db *DB) ListarClientes() []*models.Cliente {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var lista []*models.Cliente
	for _, c := range db.Clientes {
		if c.Activo {
			lista = append(lista, c)
		}
	}
	return lista
}

// ── CRUD Pedidos ──

func (db *DB) CrearPedido(clienteID int) (*models.Pedido, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Clientes[clienteID]; !ok {
		return nil, errors.New("cliente no encontrado")
	}
	p := models.NewPedido(db.nextPedidoID(), clienteID)
	db.Pedidos[p.ID()] = p
	return p, nil
}

func (db *DB) ObtenerPedido(id int) (*models.Pedido, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	p, ok := db.Pedidos[id]
	if !ok {
		return nil, errors.New("pedido no encontrado")
	}
	return p, nil
}

func (db *DB) PedidosPorCliente(clienteID int) []*models.Pedido {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var lista []*models.Pedido
	for _, p := range db.Pedidos {
		if p.ClienteID == clienteID {
			lista = append(lista, p)
		}
	}
	return lista
}

// Mu expone el mutex para lecturas externas (handlers)
func (db *DB) Mu() *sync.RWMutex { return &db.mu }

// Instancia global
var Global = NewDB()
