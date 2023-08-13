package minio_try

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
	endpoint := "localhost:9000"
	accessKeyID := "hi6mFwCKqNgw6AmGGVmX"
	secretAccessKey := "vMRTTEEUnRouKatIeJz2nWPEk4lnYzsQLmn15bZk"
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

	result, errorUpload := minioClient.FPutObject(ctx, "testing", "nano_chan.jpg", path, minio.PutObjectOptions{})

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

func GetFilePath(client *minio.Client, name string) {
	bucketName := "testing"
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": [
					"arn:aws:s3:::` + bucketName + `/*"
				]
			}
		]
	}`
	c := context.Background()
	errCreateBucket := client.MakeBucket(c, bucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
	if errCreateBucket != nil {
		log.Println("failed to create new bucket :", errCreateBucket)
		return
	}

	errSetBucketPolicy := client.SetBucketPolicy(c, bucketName, policy)
	if errSetBucketPolicy != nil {
		log.Println(errSetBucketPolicy)
		return
	}
}

func CheckObjExist(client *minio.Client, filename string) {
	bucketName := "testing"
	c := context.Background()
	objStat, errGetStat := client.StatObject(c, bucketName, filename, minio.StatObjectOptions{})
	if errGetStat != nil {
		log.Println("failed to get obj stat :", errGetStat)
		return
	}

	log.Println(objStat)
}
