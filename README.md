# tools
Beginnings of a command line tool.

## ECR
1. Update the permission policy for ECR repositories

### Commands
#### Update all ECR repo's permission policy
```shell script
go run ecr.go -y myAWSAccount -r myRegion -d 10m -f ecr-access-policy.json
```