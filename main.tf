resource "linode_instance" "test_instance" {
	label = "test-instance"
	type = "g6-nanode-1"
	region = "us-west"
	image = "linode/debian9"
	root_pass = "Test12345"
	swap_size = 512
}

resource "linode_disk" "test_disk" {
	linode_id = "${linode_instance.test_instance.id}"
	label = "test-disk"
	size = 128
}

resource "linode_domain" "test_domain" {
	domain = "testdomain.io"
	type = "master"
	soa_email = "test@testdomain.io"
}
