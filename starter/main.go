package main

func main() {
	// TODO: HTTP server with POST /events handler — validate, persist, ack.
	// TODO: Persistence layer — your choice. Justify it in the Decisions section.
	// TODO: Forwarder — read persisted events and POST them to http://localhost:9000/events,
	//       preserving per-user_id ordering and surviving downstream failures.
	// TODO: Graceful shutdown on SIGTERM/SIGINT — drain in-flight work, don't lose events.
}
