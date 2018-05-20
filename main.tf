resource "linode_domain" "test-domain" {
	domain = "testdomain.io"
	type = "master"
	soa_email = "test@testdomain.io"
}
