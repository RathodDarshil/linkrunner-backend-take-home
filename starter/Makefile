.PHONY: run test chaos-up chaos-down

run:
	go run .

test:
	go test ./...

chaos-up:
	go run ./chaos/receiver.go

chaos-down:
	@pkill -f "chaos/receiver" || true
