module github.com/donnol/do/test

go 1.18

require (
	github.com/donnol/do v0.76.0
	gorm.io/driver/sqlite v1.4.4
	gorm.io/gorm v1.24.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/smallnest/chanx v1.2.0 // indirect
)

replace github.com/donnol/do v0.76.0 => ../
