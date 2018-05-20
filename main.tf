resource "linode_instance" "test-instance" {
	label = "test-instance"
	type = "g6-nanode-1"
	region = "us-east"
	image = "linode/debian9"
	root_pass = "Test12345"
	swap_size = 512
}

resource "linode_domain" "test-domain" {
	domain = "testdomain.io"
	type = "master"
	soa_email = "test@testdomain.io"
}
