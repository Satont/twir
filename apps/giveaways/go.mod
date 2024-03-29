module github.com/twirapp/twir/apps/giveaways

go 1.21.5

toolchain go1.22.0

replace (
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/twirapp/twir/libs/gomodels => ../../libs/gomodels
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/twirapp/twir/libs/integrations => ../../libs/integrations
	github.com/twirapp/twir/libs/logger => ../../libs/logger
	github.com/twirapp/twir/libs/sentry => ../../libs/sentry
	github.com/twirapp/twir/libs/twitch => ../../libs/twitch
	github.com/twirapp/twir/libs/types => ../../libs/types
	github.com/twirapp/twir/libs/uptrace => ../../libs/uptrace
)

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/nicklaw5/helix/v2 v2.28.1
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/logger v0.0.0-20240318150924-e0edef040053
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/gomodels v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	github.com/twirapp/twir/libs/logger v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/sentry v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/twitch v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/uptrace v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.21.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.8
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/getsentry/sentry-go v0.26.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nats.go v1.33.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.31.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	github.com/samber/slog-common v0.15.0 // indirect
	github.com/samber/slog-multi v1.0.2 // indirect
	github.com/samber/slog-sentry/v2 v2.4.0 // indirect
	github.com/samber/slog-zerolog/v2 v2.2.0 // indirect
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5 // indirect
	github.com/satont/twir/libs/types v0.0.0-20240316203219-21491ee1343c // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/uptrace/uptrace-go v1.21.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.48.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.46.1 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
