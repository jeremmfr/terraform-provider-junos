default: install

.PHONY: install testacc testacc_srx testacc_router testacc_switch
# Install to use dev_overrides in provider_installation of Terraform
install:
	go install
# Run acceptance tests
testacc:
	cd junos ; TF_ACC=1 go test -v --timeout 0 -coverprofile=../coverage.out $(TESTARGS)
	go tool cover -html=coverage.out
testacc/srx:
	cd junos ; TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../coverage_srx.out $(TESTARGS)
	go tool cover -html=coverage_srx.out
testacc/router:
	cd junos ; TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../coverage_router.out $(TESTARGS)
	go tool cover -html=coverage_router.out
testacc/switch:
	cd junos ; TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../coverage_switch.out $(TESTARGS)
	go tool cover -html=coverage_switch.out
