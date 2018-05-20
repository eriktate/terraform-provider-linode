reset:
	rm -f terraform.tfstate
	rm -f terraform.tfstate.backup
	go build -o terraform-provider-linode
	terraform init
