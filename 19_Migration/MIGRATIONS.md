# Migration's Commands
```bash
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_users_table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_events_table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_attendees_table
```