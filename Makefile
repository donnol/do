.PHONY:

install_letgo:
	go install ./cmd/letgo

test_letgo_sql2struct:install_letgo
	letgo sql2struct --pkg=sqlparser -f ./cmd/letgo/sqlparser/test.sql -o ./cmd/letgo/sqlparser/test.go

test_letgo_sql2struct_insert:install_letgo
	letgo sql2struct --pkg=sqlparser -f ./cmd/letgo/sqlparser/test.sql insert --amount=3
