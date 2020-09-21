package main

import (
	"context"
	_ "context"
	"flag"
	_ "flag"
	_ "fmt"
	"github.com/aws/aws-sdk-go/service/ecr"
	_ "os"
	"time"
	_ "time"

	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/service/ecr"
)

// Usage:
// #Update repository with permission policy
// go run ecr.go -r myRegion -d 10m < ecr-access-policy.json
func main()  {
	var region string
	var timeout time.Duration

	flag.StringVar(&region, "r", "", "Bucket name.")
	flag.DurationVar(&timeout, "d", 0, "Upload timeout.")
	flag.Parse()


	sess := session.Must(session.NewSession())

	svc := ecr.New(sess)

	ctx := context.Background()

	var cancelFn func ()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}
	
	i := ecr.DescribeRepositoriesInput{
		MaxResults:      nil,
		NextToken:       nil,
		RegistryId:      nil,
		RepositoryNames: nil,
	}
	r, err := svc.DescribeRepositories()
}