package models_test

import (
	"ecommerce/models"
	"testing"
)

// ─────────────────────────────────────────
// PRUEBAS UNITARIAS — Producto
// ─────────────────────────────────────────

func TestCrearProductoExitoso(t *testing.T) {
	p, err := models.NewProducto(1, "Laptop", 999.99, 10, "Tecnología")
	if err != nil {
		t.Fatalf("No debería fallar: %v", err)
	}
	if p.Nombre() != "Laptop" {
		t.Errorf("Nombre esperado 'Laptop', got '%s'", p.Nombre())
	}
	if p.Stock() != 10 {
		t.Errorf("Stock esperado 10, got %d", p.Stock())
	}
}

func TestPrecioNegativoError(t *testing.T) {
	p, _ := models.NewProducto(1, "Mouse", 25.00, 10, "Tech")
	err := p.SetPrecio(-5.0)
	if err == nil {
		t.Error("Debería lanzar error con precio negativo")
	}
}

func TestReducirStockCorrecto(t *testing.T) {
	p, _ := models.NewProducto(1, "Teclado", 50.00, 20, "Tech")
	if err := p.ReducirStock(5); err != nil {
		t.Fatalf("No debería fallar: %v", err)
	}
	if p.Stock() != 15 {
		t.Errorf("Stock esperado 15, got %d", p.Stock())
	}
}

func TestReducirStockInsuficienteError(t *testing.T) {
	p, _ := models.NewProducto(1, "Monitor", 300.00, 2, "Tech")
	err := p.ReducirStock(10)
	if err == nil {
		t.Error("Debería lanzar error por stock insuficiente")
	}
}

func TestProductoDigitalNoReduceStock(t *testing.T) {
	p, _ := models.NewProductoDigital(1, "Curso", 29.99, "Edu", "http://url.com")
	stockAntes := p.Stock()
	p.ReducirStock(999)
	if p.Stock() != stockAntes {
		t.Error("Los productos digitales no deben reducir stock")
	}
}

func TestToDictTieneCamposRequeridos(t *testing.T) {
	p, _ := models.NewProducto(1, "Laptop", 999.99, 10, "Tech")
	d := p.ToDict()
	for _, campo := range []string{"id", "nombre", "precio", "stock", "categoria"} {
		if _, ok := d[campo]; !ok {
			t.Errorf("Campo faltante en ToDict: %s", campo)
		}
	}
}

// ─────────────────────────────────────────
// PRUEBAS UNITARIAS — Cliente
// ─────────────────────────────────────────

func TestCrearClienteExitoso(t *testing.T) {
	c, err := models.NewCliente(1, "Ana García", "ana@email.com", "0991234567")
	if err != nil {
		t.Fatalf("No debería fallar: %v", err)
	}
	if c.Email() != "ana@email.com" {
		t.Errorf("Email incorrecto: %s", c.Email())
	}
}

func TestEmailInvalidoError(t *testing.T) {
	_, err := models.NewCliente(1, "Luis", "no-es-email", "")
	if err == nil {
		t.Error("Debería fallar con email inválido")
	}
}

func TestEmailGuardadoEnMinusculas(t *testing.T) {
	c, _ := models.NewCliente(1, "María", "MARIA@EMAIL.COM", "")
	if c.Email() != "maria@email.com" {
		t.Errorf("Email debería estar en minúsculas: %s", c.Email())
	}
}

func TestNombreVacioError(t *testing.T) {
	_, err := models.NewCliente(1, "   ", "test@test.com", "")
	if err == nil {
		t.Error("Debería fallar con nombre vacío")
	}
}

// ─────────────────────────────────────────
// PRUEBAS UNITARIAS — Pedido
// ─────────────────────────────────────────

func TestPedidoIniciaEnPendiente(t *testing.T) {
	p := models.NewPedido(1, 1)
	if p.Estado() != "pendiente" {
		t.Errorf("Estado inicial debería ser 'pendiente', got '%s'", p.Estado())
	}
}

func TestAgregarProductoReduceStock(t *testing.T) {
	prod, _ := models.NewProducto(1, "Laptop", 999.99, 10, "Tech")
	pedido := models.NewPedido(1, 1)
	if err := pedido.AgregarProducto(prod, 2); err != nil {
		t.Fatalf("No debería fallar: %v", err)
	}
	if prod.Stock() != 8 {
		t.Errorf("Stock esperado 8, got %d", prod.Stock())
	}
}

func TestTotalCalculadoCorrecto(t *testing.T) {
	prod, _ := models.NewProducto(1, "Laptop", 999.99, 10, "Tech")
	pedido := models.NewPedido(1, 1)
	pedido.AgregarProducto(prod, 3)
	esperado := 999.99 * 3
	if pedido.Total() != esperado {
		t.Errorf("Total esperado %.2f, got %.2f", esperado, pedido.Total())
	}
}

func TestCambiarEstadoValido(t *testing.T) {
	p := models.NewPedido(1, 1)
	if err := p.CambiarEstado("pagado"); err != nil {
		t.Fatalf("No debería fallar: %v", err)
	}
	if p.Estado() != "pagado" {
		t.Errorf("Estado esperado 'pagado', got '%s'", p.Estado())
	}
}

func TestCambiarEstadoInvalidoError(t *testing.T) {
	p := models.NewPedido(1, 1)
	if err := p.CambiarEstado("estado_inventado"); err == nil {
		t.Error("Debería fallar con estado inválido")
	}
}

func TestCancelarDevuelveStock(t *testing.T) {
	prod, _ := models.NewProducto(1, "Laptop", 999.99, 10, "Tech")
	pedido := models.NewPedido(1, 1)
	pedido.AgregarProducto(prod, 3) // stock = 7
	repo := map[int]*models.Producto{1: prod}
	if err := pedido.Cancelar(repo); err != nil {
		t.Fatalf("No debería fallar: %v", err)
	}
	if prod.Stock() != 10 {
		t.Errorf("Stock debería ser 10 tras cancelar, got %d", prod.Stock())
	}
	if pedido.Estado() != "cancelado" {
		t.Errorf("Estado esperado 'cancelado', got '%s'", pedido.Estado())
	}
}
