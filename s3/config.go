package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func Config(region, endpoint string) *aws.Config {
	if os.Getenv("ENV") == "development" {
		return &aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ACCESS_KEY"),
				os.Getenv("SECRET_KEY"),
				"",
			),
			S3ForcePathStyle: aws.Bool(true),
			Region:           aws.String(region),
			Endpoint:         aws.String(endpoint),
		}
	}
	return &aws.Config{
		Region: aws.String(region),
	}
}
