package diskx

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"

	"github.com/zengineDev/x/configx"
)

type Disk struct {
	D *s3.S3
}

var (
	instance *Disk
)

var once sync.Once

func GetDiskConnection() *Disk {
	conf := configx.GetConfig()

	once.Do(func() {
		s3Config := &aws.Config{
			Credentials: credentials.NewStaticCredentials(conf.Disks.S3.Key, conf.Disks.S3.Secret, ""),
			Endpoint:    aws.String(conf.Disks.S3.Endpoint),
			Region:      aws.String("us-east-1"),
		}

		newSession := session.New(s3Config)
		s3Client := s3.New(newSession)

		instance = &Disk{D: s3Client}

	})

	return instance

}
