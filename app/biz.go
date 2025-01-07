package app

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type UseCase struct {
	minioClient *minio.Client
	bucket      string

	imageKeys []string
}

func NewUseCase(minioClient *minio.Client, bucket string) *UseCase {
	return &UseCase{
		minioClient: minioClient,
		imageKeys:   []string{},
		bucket:      bucket,
	}
}

func (uc *UseCase) PresignedPostPolicy(ctx context.Context, key string, contentType string) (string, map[string]string, error) {
	var err error
	p := minio.NewPostPolicy()
	_ = p.SetExpires(time.Now().Add(60 * time.Second))
	_ = p.SetBucket(uc.bucket)
	_ = p.SetKey(key)
	if contentType != "" {
		_ = p.SetContentType(contentType)
	}

	res, fd, err := uc.minioClient.PresignedPostPolicy(ctx, p)
	if err != nil {
		return "", nil, err
	}

	uc.imageKeys = append(uc.imageKeys, key)
	return res.String(), fd, nil

}
func (uc *UseCase) ListImages(ctx context.Context) []string {
	arr := make([]string, len(uc.imageKeys))
	for i, key := range uc.imageKeys {
		if strings.HasPrefix(key, "protect") {
			val := url.Values{}
			res, err := uc.minioClient.PresignedGetObject(ctx, uc.bucket, key, time.Second*60, val)
			if err != nil {
				arr[i] = uc.minioClient.EndpointURL().String() + "/" + uc.bucket + "/" + key
			} else {
				arr[i] = res.String()
			}
		} else {
			arr[i] = uc.minioClient.EndpointURL().String() + "/" + uc.bucket + "/" + key
		}
	}
	return arr
}

func (uc *UseCase) PresignedGetObject(ctx context.Context, key string) (string, error) {
	val := url.Values{}
	res, err := uc.minioClient.PresignedGetObject(ctx, uc.bucket, key, time.Second*60, val)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
