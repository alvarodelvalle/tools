// Package main does this....
package main

import (
	"context"
	_ "context"
	"encoding/json"
	"flag"
	_ "flag"
	"fmt"
	_ "fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"io/ioutil"
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
// go run ecr.go -y myAWSAccount -r myRegion -d 10m -f ecr-access-policy.json
func main() {
	var registryId, region, file string
	var timeout time.Duration

	flag.StringVar(&registryId, "y", "", "AWS Account ID")
	flag.StringVar(&region, "r", "", "Region in AWS")
	flag.DurationVar(&timeout, "d", 0, "Upload timeout.")
	flag.StringVar(&file, "f", "", "Access Policy File (JSON)")
	flag.Parse()

	sess := session.Must(session.NewSession())
	svc := ecr.New(sess, aws.NewConfig().WithRegion(region))
	ctx := context.Background()

	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	describeRepositoriesOutput, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{RegistryId: &registryId})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.
			//Error to get the Code and Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	bytes, err := json.Marshal(describeRepositoriesOutput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Describing Repositories: %v\n", string(bytes))

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	policyText := string(data)
	var repos []string

	for _,  v := range describeRepositoriesOutput.Repositories{
		repos = append(repos, *v.RepositoryName)
	}

	for _, v := range repos {
		repoName:= v

		_, err := svc.SetRepositoryPolicy(&ecr.SetRepositoryPolicyInput{
			PolicyText:     &policyText,
			RegistryId:     &registryId,
			RepositoryName: &repoName,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ecr.ErrCodeServerException:
					fmt.Println(ecr.ErrCodeServerException, aerr.Error())
				case ecr.ErrCodeInvalidParameterException:
					fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
				case ecr.ErrCodeRepositoryNotFoundException:
					fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.
				//Error to get the Code and Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		getRepositoryPolicyOutput, err := svc.GetRepositoryPolicy(&ecr.GetRepositoryPolicyInput{
			RegistryId:     &registryId,
			RepositoryName: &repoName,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ecr.ErrCodeServerException:
					fmt.Println(ecr.ErrCodeServerException, aerr.Error())
				case ecr.ErrCodeInvalidParameterException:
					fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
				case ecr.ErrCodeRepositoryNotFoundException:
					fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.
				//Error to get the Code and Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		if *getRepositoryPolicyOutput.PolicyText == "" {
			err := fmt.Errorf("Policy for repo %v is empty, should be \n%v", repoName, policyText)
			panic(err)
		} else {
			fmt.Printf("Updated repo %v with policy %v\n", repoName, *getRepositoryPolicyOutput.PolicyText)
		}
	}
}
