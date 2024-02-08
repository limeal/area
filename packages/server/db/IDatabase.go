package db

// IDatabase is an interface that has three methods: Init, Migrate, and Connect.
// @property Init - This is a function that returns a pointer to the database object.
// @property {error} Migrate - This is a function that will be used to migrate the database.
// @property {error} Connect - This is the function that will be called to connect to the database.
type IDatabase interface {
	Init() *IDatabase
	Migrate() error
	Connect() error
}
