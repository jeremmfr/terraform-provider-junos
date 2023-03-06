default: install

.PHONY: install testacc testacc_srx testacc_router testacc_switch testunit cleanout

# Install to use dev_overrides in provider_installation of Terraform
install:
	go install

# Run acceptance tests
testacc:
	cd internal/providersdk ; TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk.out $(TESTARGS)
	cd internal/providerfwk ; TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk.out $(TESTARGS)
	go tool cover -html=coverage_sdk.out
	go tool cover -html=coverage_fwk.out
testacc/srx:
	cd internal/providersdk ; TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk_srx.out $(TESTARGS)
	cd internal/providerfwk ; TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_srx.out $(TESTARGS)
	go tool cover -html=coverage_sdk_srx.out
	go tool cover -html=coverage_fwk_srx.out
testacc/router:
	cd internal/providersdk ; TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk_router.out $(TESTARGS)
	cd internal/providerfwk ; TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_router.out $(TESTARGS)
	go tool cover -html=coverage_sdk_router.out
	go tool cover -html=coverage_fwk_router.out
testacc/switch:
	cd internal/providersdk ; TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk_switch.out $(TESTARGS)
	cd internal/providerfwk ; TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_switch.out $(TESTARGS)
	go tool cover -html=coverage_sdk_switch.out
	go tool cover -html=coverage_fwk_switch.out

# Run unit tests
testunit: 
	go test -race -v -coverprofile=coverage_unit.out ./...
	go tool cover -html=coverage_unit.out

# Cleanup out files from tests
cleanout:
	find . -maxdepth 1 -name "*.out" -type f -delete
