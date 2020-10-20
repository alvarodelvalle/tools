// Package main does this....
package cmd

import (
	_ "context"
	_ "flag"
	"fmt"
	_ "fmt"
	_ "os"
	"time"
	_ "time"

	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/request"
	_ "github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"
)

var (
	//flags
	registryId string
	region     string
	file       string
	timeout    time.Duration

	ecrCmd = &cobra.Command{
		Use:     "ecr-policy",
		Aliases: []string{"ecr", "ecrpolicy"},
		Short:   "Update ECR Policy",
		Long: `Update an AWS ECR Repository with the given permissions policy. 
Use this to allow cross-accounts like 'dev' to access the ECR in the root account`,
		Run: func(cmd *cobra.Command, args []string) {
			updateEcrPermissionsPolicy()
		},
	}
)

// Usage:
// #Update repository with permission policy
// go run ecr-policy.go -y myAWSAccount -r myRegion -d 10m -f ecr-access-policy.json
func updateEcrPermissionsPolicy() {
	fmt.Println("updated ecr policy")
}
