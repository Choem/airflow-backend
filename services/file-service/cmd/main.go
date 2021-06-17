package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/choem/airflow-backend/services/file-service/cmd/graph/generated"
	"github.com/choem/airflow-backend/services/file-service/cmd/graph/resolvers"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	port := os.Getenv("PORT")
	endpointPrefix := os.Getenv("ENDPOINT_PREFIX")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSLString := os.Getenv("MINIO_USE_SSL")

	useSSL, useSSLConvErr := strconv.ParseBool(useSSLString)
	if useSSLConvErr != nil {
		useSSL = false
	}

	// Initialize minio client object
	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal(err)
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers:  &resolvers.Resolver{MinioClient: minioClient},
			Directives: generated.DirectiveRoot{},
			Complexity: generated.ComplexityRoot{},
		},
	))

	http.Handle(endpointPrefix+"/playground", playground.Handler("GraphQL playground", endpointPrefix+"/query"))
	http.Handle(endpointPrefix+"/query", server)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
