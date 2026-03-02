NexShop E-Commerce — Go + Gin
Materia: Programación Orientada a Objetos
Curso: 2-CIB-3A
Integrante: Carlos Paguay
Backend: Go 1.21 + Gin

Objetivo
Sistema de gestión de tienda online desarrollado en Go aplicando los principios de POO mediante structs, embedding (herencia), interfaces y encapsulación con campos privados. Expone una API REST con 11 endpoints y serialización JSON.

Estructura del Proyecto
ecommerce-go/
├── models/
│   ├── models.go        # Unidades 1-3: structs, embedding, encapsulación
│   └── models_test.go   # Pruebas unitarias
├── repository/
│   └── repository.go    # Base de datos en memoria + mutex (concurrencia)
├── handlers/
│   └── handlers.go      # Unidad 4: 11 endpoints REST con Gin
├── static/
│   └── index.html       # Frontend web
├── main.go              # Router Gin + servidor
├── go.mod
└── README.md

Instalación y ejecución en Pop!_OS
bash# 1. Instalar Go
sudo apt install golang-go -y
o se puede descargar desde la app oficial

# 2. Clonar / ir al proyecto
cd gestion

# 3. Descargar dependencias
go mod tidy

# 4. Ejecutar
go run main.go

# 5. Compilar binario
go build -o nexshop .
./nexshop
La API estará en: http://localhost:5000
El sitio web en: http://localhost:5000/static/index.html

Endpoints REST
#MétodoEndpointDescripción1GET/api/productosListar productos2GET/api/productos/:idProducto por ID3POST/api/productosCrear producto4PUT/api/productos/:idActualizar producto5GET/api/clientesListar clientes6POST/api/clientesRegistrar cliente7POST/api/pedidosCrear pedido8GET/api/clientes/:id/pedidosPedidos de cliente9PATCH/api/pedidos/:id/estadoCambiar estado10DELETE/api/pedidos/:idCancelar pedido11GET/api/estadisticasDashboard

Pruebas
bashgo test ./... -v

Equivalencias Go ↔ Python/POO
Concepto POOPythonGoClaseclass Productotype Producto structHerenciaclass Producto(EntidadBase)type Producto struct { EntidadBase } (embedding)Encapsulaciónself.__preciocampo en minúscula precio + getter/setterPolimorfismodef to_dict() sobreescritofunc (p *Producto) ToDict() por tipoConcurrenciathreaded=True en Flasknativa — cada petición es una goroutine

Conceptos POO aplicados

Unidad 1: Structs como clases, campos y métodos, paquetes Go
Unidad 2: Embedding de structs (herencia), polimorfismo con ToDict()
Unidad 3: Encapsulación con campos privados, manejo de errores con error
Unidad 4: API REST con Gin, JSON nativo, concurrencia con goroutines + sync.RWMutex
