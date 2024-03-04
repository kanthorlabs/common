package testdata

var (
	SqliteUri   = "file::memory:?cache=shared"
	PostgresUri = "postgres://postgres:postgres@localhost:2345/postgres"
	RedisUri    = "redis://localhost:6379/0"
	MemoryUri   = "memory://"
	NatsUri     = "nats://localhost:2224"
)
