package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
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

	prefix := "/logs"

	for object := range r.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		logUrls = append(logUrls, strings.Split(object.Key, prefix+"/")[1])
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
			if inRange(parsedStartDate, parsedEndDate, object.LastModified) && !contains(bucketsWithNewLogs, bucket.Name) {
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

func (r *queryResolver) GetPatientModelDownloadURL(ctx context.Context, patientID int) (*string, error) {
	var bucketName = fmt.Sprintf("user-%d", patientID)

	var latestModel *minio.ObjectInfo = nil
	for object := range r.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: "models", Recursive: true}) {
		if latestModel == nil || object.LastModified.After(latestModel.LastModified) {
			latestModel = &object
		}
	}

	if latestModel != nil {
		url, err := r.MinioClient.PresignedGetObject(ctx, bucketName, latestModel.Key, time.Minute*1, url.Values{})
		if err != nil {
			return nil, err
		}

		var downloadURL string = url.String()

		// Parse admin credential query param
		accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
		var downloadURLQueryParams = strings.Split(downloadURL, fmt.Sprintf("&X-Amz-Credential=%s", accessKeyID))
		downloadURL = downloadURLQueryParams[0] + downloadURLQueryParams[1]

		if strings.Contains(downloadURL, "minio:9000") {
			var downloadURLParts = strings.Split(downloadURL, "minio:9000")
			downloadURL = downloadURLParts[0] + "localhost:9000" + downloadURLParts[1]
		}

		return &downloadURL, nil
	}

	return nil, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
