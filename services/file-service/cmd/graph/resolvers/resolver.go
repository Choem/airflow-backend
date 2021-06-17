package resolvers

import "github.com/minio/minio-go/v7"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MinioClient *minio.Client
}
