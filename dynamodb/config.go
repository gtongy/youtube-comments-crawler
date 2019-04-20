package dynamodb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
)

// Config is get aws settings
func Config(region, endpoint string) *aws.Config {
	if os.Getenv("ENV") == "development" {
		return &aws.Config{
			Region:   aws.String(region),
			Endpoint: aws.String(endpoint),
		}
	}
	return &aws.Config{
		Region: aws.String(region),
	}
}
