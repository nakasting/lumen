package databse

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectSQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("../lumen.db"))
}
