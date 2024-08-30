module processor

go 1.22.3

replace models => ../models

require (
	github.com/disintegration/imaging v1.6.2
	github.com/google/uuid v1.6.0
	github.com/streadway/amqp v1.1.0
	models v0.0.0
)

require golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
