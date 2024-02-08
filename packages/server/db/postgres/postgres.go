package postgres

import (
	"area-server/db/postgres/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// `PSDatabase` is a struct with 6 fields.
// @property {string} Host - The hostname of the database server.
// @property {int} Port - The port number of the PostgreSQL server.
// @property {string} User - The user to connect to the database with.
// @property {string} Password - The password for the user you created in the previous step.
// @property {string} Database - The name of the database you want to connect to.
// @property {string} SSLMode - This is the SSL mode to use when connecting to the database.
type PSDatabase struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

var DB *gorm.DB

// Init the database
func Init() *PSDatabase {

	user, present := os.LookupEnv("POSTGRES_USER")
	if !present {
		panic("POSTGRES_USER not set")
	}
	password, present := os.LookupEnv("POSTGRES_PASSWORD")
	if !present {
		panic("POSTGRES_PASSWORD not set")
	}
	database, present := os.LookupEnv("POSTGRES_DB")
	if !present {
		panic("POSTGRES_DB not set")
	}
	sslmode, present := os.LookupEnv("POSTGRES_SSLMODE")
	if !present {
		panic("POSTGRES_SSLMODE not set")
	}

	host := "postgres"
	if os.Getenv("POSTGRES_HOST") != "" {
		host = os.Getenv("POSTGRES_HOST")
	}

	return &PSDatabase{
		Host:     host,
		Port:     5432,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslmode,
	}
}

// Connect to the database
func (d *PSDatabase) Connect() error {
	fmt.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(d.String()), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database!")
	DB = db
	return nil
}

// Returning a string with the database connection information.
func (d *PSDatabase) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)
}

// Dropping the tables and then creating them again.
func (d *PSDatabase) Migrate() error {
	fmt.Println("Dropping tables...")
	if DB.Migrator().DropTable(&models.Account{}, &models.Authorization{}, &models.Applet{}, &models.Area{}) != nil {
		panic("Failed to drop tables")
	}
	fmt.Println("Creating tables...")
	if DB.AutoMigrate(&models.Account{}, &models.Authorization{}, &models.Applet{}, &models.Area{}, &models.Area{}) != nil {
		panic("Failed to migrate databases")
	}
	return nil
}

// Checking if the database is ok.
func (d *PSDatabase) OK() bool {

	if DB == nil {
		return false
	}

	return DB.Migrator().HasTable(&models.Area{}) &&
		DB.Migrator().HasTable(&models.Applet{}) &&
		DB.Migrator().HasTable(&models.Authorization{}) &&
		DB.Migrator().HasTable(&models.Account{})
}
