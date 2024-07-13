CURRENT_DIR := $(shell pwd)

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}


exp:
	export DBURL='postgres://macbookpro:1111@localhost:5432/lesson46?sslmode=disable'

mig-up:
	migrate -path migrations -database 'postgres://macbookpro:1111@localhost:5432/lesson46?sslmode=disable' -verbose up

mig-down:
	migrate -path migrations -database ${DBURL} -verbose down

mig-create:
	migrate create -ext sql -dir migrations -seq create_tables2

mig-insert:
	migrate create -ext sql -dir migrations -seq insert_table