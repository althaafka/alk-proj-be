package database
import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/althaafka/alk-proj-be.git/models"
)

var DB *gorm.DB;

func Connect() {
	dsn := "host=localhost user=postgres password=qwerty dbname=alkademi port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&models.User{}, &models.Article{})

	DB = database
}