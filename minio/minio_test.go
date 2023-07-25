package minio_try

import (
	"testing"
)

func TestInit(t *testing.T) {
	InitMinio()
}

func TestCheckBucketExists(t *testing.T) {
	Client := InitMinio()

	// UploadNewFiles(Client, "patrick.jpeg")
	UploadNewFiles(Client, "7ec329c4dac04b68ccdd2751c4202976.jpg")
}

func TestGetFile(t *testing.T) {
	Client := InitMinio()
	ReadFile(Client, "7ec329c4dac04b68ccdd2751c4202976.jpg")
}

func TestRmFile(t *testing.T) {
	Client := InitMinio()
	RemoveFile(Client, "nano_chan.jpg")
}
