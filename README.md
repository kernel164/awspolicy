# awspolicy

## e.g Run terraform and generate aws policy

```
export AWSPOLICY_CAPTURE_TYPE=terraform-cli 
awspolicy init
awspolicy apply
```

## OR 

```
ln -s /usr/bin/awspolicy terraform
./terraform init
./terraform apply
```

## Sample Output

```
$ AWSPOLICY_CAPTURE_TYPE=terraform-cli awspolicy plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.aws_region.current: Refreshing state...
aws_ecr_repository.default: Refreshing state... [id=test]
data.aws_caller_identity.current: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # aws_ecr_repository.default will be created
  + resource "aws_ecr_repository" "default" {
      + arn                  = (known after apply)
      + id                   = (known after apply)
      + image_tag_mutability = "MUTABLE"
      + name                 = "test"
      + registry_id          = (known after apply)
      + repository_url       = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.

==> Policy Document:

{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:DescribeAccountAttributes",
        "ecr:DescribeRepositories",
        "sts:GetCallerIdentity"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```