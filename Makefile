default: install

.PHONY: install testacc testacc_srx testacc_router testacc_switch testunit cleanout changemd

# Install to use dev_overrides in provider_installation of Terraform
install:
	go install

# Run acceptance tests
testacc/srx:
	cd internal/providerfwk ; TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_srx.out $(TESTARGS)
	go tool cover -html=coverage_fwk_srx.out
testacc/upgradestate/srx:
	cd internal/providerfwk ; TF_CLI_CONFIG_FILE= TESTACC_UPGRADE_STATE=1 TESTACC_SRX=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_upg_srx.out -run "TestAccUpgradeState" $(TESTARGS)
	go tool cover -html=coverage_fwk_upg_srx.out
testacc/router:
	cd internal/providerfwk ; TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_router.out $(TESTARGS)
	go tool cover -html=coverage_fwk_router.out
testacc/upgradestate/router:
	cd internal/providerfwk ; TF_CLI_CONFIG_FILE= TESTACC_UPGRADE_STATE=1 TESTACC_ROUTER=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_upg_router.out -run "TestAccUpgradeState" $(TESTARGS)
	go tool cover -html=coverage_fwk_upg_router.out
testacc/switch:
	cd internal/providerfwk ; TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_switch.out $(TESTARGS)
	go tool cover -html=coverage_fwk_switch.out
testacc/upgradestate/switch:
	cd internal/providerfwk ; TF_CLI_CONFIG_FILE= TESTACC_UPGRADE_STATE=1 TESTACC_SWITCH=1 TF_ACC=1 go test -v --timeout 0 -coverprofile=../../coverage_fwk_upg_switch.out -run "TestAccUpgradeState" $(TESTARGS)
	go tool cover -html=coverage_fwk_upg_switch.out

# Run unit tests
testunit: 
	go test -race -v -coverprofile=coverage_unit.out ./...
	go tool cover -html=coverage_unit.out

# Cleanup out files from tests
cleanout:
	find . -maxdepth 1 -name "*.out" -type f -delete

changemd:
	cp .changes/.template.md .changes/$(shell git rev-parse --abbrev-ref HEAD).md
