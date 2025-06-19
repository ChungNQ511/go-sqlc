package s3

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

// S3_CONFIG is the configuration for the S3 storage
type S3_CONFIG struct {
	BUCKET            string `mapstructure:"BUCKET"`
	ACCESS_KEY_ID     string `mapstructure:"ACCESS_KEY_ID"`
	ACCESS_KEY_SECRET string `mapstructure:"ACCESS_KEY_SECRET"`
	S3_REGION         string `mapstructure:"S3_REGION"`
	S3_HOST_NAME      string `mapstructure:"S3_HOST_NAME"`

	ROOT_DIR string `mapstructure:"S3_ROOT_DIR"`
	DIR      string `mapstructure:"S3_DIR"`
}

var S3Config S3_CONFIG

// loadS3Config loads the S3 configuration from the environment variables
func loadS3Config(path string) (config S3_CONFIG, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// setup s3 client
func setup(ctx context.Context) (*s3.Client, error) {
	// load config
	loadS3Config(".")

	// setup s3 client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			S3Config.ACCESS_KEY_ID,
			S3Config.ACCESS_KEY_SECRET,
			"",
		)),
		config.WithRegion(S3Config.S3_REGION),
		config.WithBaseEndpoint(fmt.Sprintf("https://%s", S3Config.S3_HOST_NAME)),
	)
	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)
	return s3Client, nil
}

// PresignURL presign url for upload file to s3
func PresignURL(ctx context.Context, prefix string, fileName string) (string, error) {
	s3Client, err := setup(ctx)
	if err != nil {
		return "", err
	}

	// Add RootDir if needed
	if !strings.HasPrefix(prefix, S3Config.ROOT_DIR) {
		prefix = S3Config.ROOT_DIR + prefix
	}

	presignClient := s3.NewPresignClient(s3Client)
	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(S3Config.BUCKET),
		Key:    aws.String(path.Join(prefix, fileName)),
	}, s3.WithPresignExpires(10*time.Minute))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

// GetPresignedURL generates a pre-signed S3 URL to access an object.
func GetPresignedURL(ctx context.Context, key string, durationSeconds int) (string, error) {
	s3Client, err := setup(ctx)
	if err != nil {
		return "", err
	}

	presignClient := s3.NewPresignClient(s3Client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(S3Config.BUCKET),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(durationSeconds)*time.Second)) // 👈 multiply với Second, không phải Minute

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
