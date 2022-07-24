GO_COVER_FILE ?= "coverage.out"

test:
	go test ./... --count=1

test-cover:
	[ ! -e $(GO_COVER_FILE) ] || rm $(GO_COVER_FILE)
	go test ./... --count=1 -race -covermode=atomic -coverprofile=$(GO_COVER_FILE)
	go tool cover -func $(GO_COVER_FILE)
