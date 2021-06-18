package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/choem/airflow-backend/services/file-service/cmd/graph/generated"
	minio "github.com/minio/minio-go/v7"
)

func (r *mutationResolver) CreatePatientLog(ctx context.Context, patientID int, file graphql.Upload) (bool, error) {
	_, err := r.MinioClient.FPutObject(ctx, fmt.Sprintf("user-%d/logs", patientID), file.Filename, file.Filename, minio.PutObjectOptions{
		ContentType: "application/csv",
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SavePatientModel(ctx context.Context, patientID int, file graphql.Upload) (bool, error) {
	_, err := r.MinioClient.FPutObject(ctx, fmt.Sprintf("user-%d/models", patientID), file.Filename, file.Filename, minio.PutObjectOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
