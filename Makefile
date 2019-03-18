migrate-latest:

		migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -source file://build/migrations up

migrate-rollback:

		migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -source file://build/migrations down

tests:

		go test ./...

coveralls-ci:

		go get golang.org/x/tools/cmd/cover
		go get github.com/mattn/goveralls
		go test -v -covermode=count -coverprofile=coverage.out ${HOME}/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci -repotoken ${COVERALLS_REPO_TOKEN}