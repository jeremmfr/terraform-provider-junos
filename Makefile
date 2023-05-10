default: install

.PHONY: install testacc testacc_srx testacc_router testacc_switch testunit cleanout changemd

# Install to use dev_overrides in provider_installation of Terraform
install:
	go install

# Run acceptance tests
testacc:
	cd internal/providerfwk ; TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk.out $(TESTARGS)
	go tool cover -html=coverage_fwk.out
	cd internal/providersdk ; TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk.out $(TESTARGS)
	go tool cover -html=coverage_sdk.out
testacc/srx:
	cd internal/providerfwk ; TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_srx.out $(TESTARGS)
	go tool cover -html=coverage_fwk_srx.out
	cd internal/providersdk ; TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk_srx.out $(TESTARGS)
	go tool cover -html=coverage_sdk_srx.out
testacc/router:
	cd internal/providerfwk ; TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_router.out $(TESTARGS)
	go tool cover -html=coverage_fwk_router.out
	cd internal/providersdk ; TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk_router.out $(TESTARGS)
	go tool cover -html=coverage_sdk_router.out
testacc/switch:
	cd internal/providerfwk ; TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_switch.out $(TESTARGS)
	go tool cover -html=coverage_fwk_switch.out
	cd internal/providersdk ; TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_sdk_switch.out $(TESTARGS)
	go tool cover -html=coverage_sdk_switch.out

# Run unit tests
testunit: 
	go test -race -v -coverprofile=coverage_unit.out ./...
	go tool cover -html=coverage_unit.out

# Cleanup out files from tests
cleanout:
	find . -maxdepth 1 -name "*.out" -type f -delete

changemd:
	cp .changes/.template.md .changes/$(shell git rev-parse --abbrev-ref HEAD).md
