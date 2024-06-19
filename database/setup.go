package database
import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/althaafka/alk-proj-be.git/models"
)

var DB *gorm.DB;

func Connect() {
	database, err := gorm.Open(mysql.Open("root:qwerty@tcp(localhost:3306)/alkademi?parseTime=true"))
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Article{})

	DB = database
}