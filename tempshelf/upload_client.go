package tempshelf

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

// S3Client returns cert configed s3 client with Manifest
func S3Client(manifest Manifest) *s3.S3 {
    cert := credentials.NewStaticCredentials(
        manifest.Meta.Token,
        manifest.Meta.Secret,
        "")

    s3cli := s3.New(session.New(), &aws.Config{
        Credentials: cert,
        Region:      aws.String(manifest.Meta.Region),
    })

    return s3cli
}
