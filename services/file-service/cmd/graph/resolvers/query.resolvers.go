package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/choem/airflow-backend/services/file-service/cmd/graph/generated"
	minio "github.com/minio/minio-go/v7"
)

const (
	layoutISO = "2006-01-02T15:04:05"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
func inRange(startDate time.Time, endDate time.Time, dateToMeasure time.Time) bool {
	return startDate.Before(dateToMeasure) && endDate.After(dateToMeasure)
}

func (r *queryResolver) GetPatientLogs(ctx context.Context, patientID int) ([]string, error) {
	logUrls := []string{}

	bucketName := fmt.Sprintf("user-%d", patientID)

	for object := range r.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: "/logs", Recursive: true}) {
		logUrls = append(logUrls, bucketName+object.Key)
	}

	return logUrls, nil
}

func (r *queryResolver) GetActivePatients(ctx context.Context, startDate string, endDate string) ([]int, error) {
	parsedStartDate, err := time.Parse(layoutISO, startDate)
	if err != nil {
		return nil, err
	}

	parsedEndDate, err := time.Parse(layoutISO, endDate)
	if err != nil {
		return nil, err
	}

	if parsedStartDate.After(parsedEndDate) {
		return nil, errors.New("start date cannot be after the end date")
	}

	buckets, err := r.MinioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketsWithNewLogs := []string{}

	for _, bucket := range buckets {
		for object := range r.MinioClient.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{Prefix: "/logs", Recursive: true}) {
			relativeTime := object.LastModified.Add(time.Hour * 2)
			if inRange(parsedStartDate, parsedEndDate, relativeTime) && !contains(bucketsWithNewLogs, bucket.Name) {
				bucketsWithNewLogs = append(bucketsWithNewLogs, bucket.Name)
			}
		}
	}

	patientIds := []int{}

	for _, bucketName := range bucketsWithNewLogs {
		split := strings.Split(bucketName, "-")

		i, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}

		patientIds = append(patientIds, i)
	}

	return patientIds, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
