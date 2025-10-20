# Proyecto Lumen

## 1️⃣ Librerías base que usamos

| Librería              | Uso                                      | Instalación                                     |
| --------------------- | ---------------------------------------- | ----------------------------------------------- |
| **chi**               | Router rápido y minimalista              | `go get github.com/go-chi/chi/v5`               |
| **GORM**              | ORM para manejar DB                      | `go get gorm.io/gorm`                           |
| **GORM MySQL driver** | Conector de GORM a MySQL                 | `go get gorm.io/driver/mysql`                   |
| **godotenv**          | Cargar variables de entorno desde `.env` | `go get github.com/joho/godotenv`               |
| **zap**               | Logging estructurado                     | `go get go.uber.org/zap`                        |
| **validator**         | Validación de structs                    | `go get github.com/go-playground/validator/v10` |
| **viper**             | Manejo flexible de configuración         | `go get github.com/spf13/viper`                 |

---

## 2️⃣ Instalar en tu proyecto nuevo

Desde la raíz de tu proyecto:

```bash
go mod init go-crud-api   # si todavía no lo hiciste
go get github.com/go-chi/chi/v5
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/joho/godotenv
go get go.uber.org/zap
go get github.com/go-playground/validator/v10
go get github.com/spf13/viper
```

## 3️⃣ Configuración básica para MySQL

En tu `.env` podrías tener algo como:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=mi_password
DB_NAME=books_db
PORT=8080
```

Y en Go, la conexión con GORM sería:

```go
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

* Aquí `cfg` lo cargás con `viper` o `.env`.
* Luego podés usar GORM igual que antes (`AutoMigrate`, CRUD, etc.).
