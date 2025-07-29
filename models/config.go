package models

type Config struct {
	DB     PostGresConfig
	Port   string
	JWT    string
	SignUp SignUpConfig
	MQ     MQ
	S3     S3Client
}

type PostGresConfig struct {
	Username string
	Password string
	URL      string
	Port     string
	Host     string
	Name     string
}

type SignUpConfig struct {
	ORG_ADD_KEY   string
	ORG_ADD_KEY_1 string
	ORG_ADD_KEY_2 string
	ORG_ADD_KEY_3 string
}

type MQ struct {
	Username string
	Password string
	Host     string
	Port     string
	AmazonMQ string
}

type S3Client struct {
	AccessKey string
	SecretKey string
	Region    string
}
