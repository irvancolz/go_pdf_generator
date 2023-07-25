package minio_try

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
	endpoint := "192.168.1.80:9000"
	accessKeyID := "3u8B9LQytFwCeVNmtFlq"
	secretAccessKey := "p3TgHQVgx9SBPcndQtwbCPwhbqFrK0lGCWyuAHF9"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return minioClient
}

func UploadNewFiles(minioClient *minio.Client, path string) {
	ctx := context.Background()

	result, errorUpload := minioClient.FPutObject(ctx, "users", "nano_chan.jpg", path, minio.PutObjectOptions{})

	if errorUpload != nil {
		log.Println("failed to upload file to minio :", errorUpload)
		return
	}

	log.Println(result)

	// err := minioClient.MakeBucket(ctx, "mybucket", minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Successfully created mybucket.")

	// isBucketExists, errorCheckBucketExists := minioClient.BucketExists(ctx, "test")

	// if errorCheckBucketExists != nil {
	// 	log.Println("failed to check bucket existing :", errorCheckBucketExists)
	// 	return
	// }

	// log.Println(isBucketExists)

	// buckets, errorGetBucket := minioClient.ListBuckets(ctx)

	// if errorGetBucket != nil {
	// 	log.Println("failed to get buckets list :", errorGetBucket)
	// 	return
	// }

	// for _, bucket := range buckets {
	// 	log.Println(bucket)
	// }
}

func ReadFile(minioClient *minio.Client, path string) error {
	ctx := context.Background()

	errorGetObj := minioClient.FGetObject(ctx, "users", path, path, minio.GetObjectOptions{})
	if errorGetObj != nil {
		log.Println("failed to get file from minio :", errorGetObj)
		return errorGetObj
	}
	return nil
}

func RemoveFile(client *minio.Client, path string) {
	ctx := context.Background()

	deleteProps := minio.RemoveObjectOptions{
		ForceDelete: true,
	}

	errGetFile := ReadFile(client, path)
	if errGetFile != nil {
		return
	}

	errorDelete := client.RemoveObject(ctx, "users", path, deleteProps)

	if errorDelete != nil {
		log.Println("failed to delete file in minio :", errorDelete)
		return
	}
}
