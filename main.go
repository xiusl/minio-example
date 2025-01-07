package main

import (
	"os"

	"github.com/xiusl/minio-example/app"
)

func main() {
	minioClient, err := app.NewMinioClient(
		os.Getenv("MINIO_AK"),
		os.Getenv("MINIO_SK"),
		os.Getenv("MINIO_ENDPOINT"),
	)
	if err != nil {
		panic(err)
	}
	uc := app.NewUseCase(minioClient, os.Getenv("MINIO_BUCKET"))
	s := app.NewServer(uc)
	s.Spin()
}

/*
export MINIO_AK=Jf3rIOym21obpRX5FiAv
export MINIO_SK=pa83BWsbki8573y7Q8Wi5vW1emjYFEUkUxi0qqH1
export MINIO_ENDPOINT=127.0.0.1:9000
export MINIO_BUCKET=test
*/
