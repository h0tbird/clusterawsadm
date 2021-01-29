#------------------------------------------------------------------------------
# -s Omit the symbol table and debug information.
# -w Omit the DWARF symbol table.
#------------------------------------------------------------------------------

build:
	go build -ldflags="-s -w" .

#------------------------------------------------------------------------------
# Builds and runs against the local checkout of terrago.
#------------------------------------------------------------------------------

run-local-terrago:
	@go mod edit -replace github.com/h0tbird/terrago=../terrago
	@go run .
	@go mod edit -dropreplace github.com/h0tbird/terrago

#------------------------------------------------------------------------------
# Creates a local copy of the upstream provider module and replaces all the
# upstream references with references to the local copy using 'go mod edit'.
#------------------------------------------------------------------------------

terraform-provider-aws: NAME := terraform-provider-aws
terraform-provider-aws: VERSION := v3.23.0
terraform-provider-aws: TMPDIR := $(shell mktemp -d)
terraform-provider-aws:
	@git clone --depth 1 --branch ${VERSION} https://github.com/hashicorp/${NAME}.git ${TMPDIR}
	@mkdir -p providers/${NAME} && rsync -a --delete --exclude='.*' ${TMPDIR}/ providers/${NAME}
	@go mod edit -replace github.com/terraform-providers/${NAME}=./providers/${NAME} && rm -rf ${TMPDIR}
	# TODO: Automate rewriting the provider's source code using go/parser go/ast go/printer.
	# Delete unreferenced entries in DataSourcesMap and ResourcesMap so that the code is not
	# referenced and the linker's DCE (Dead Code Elimination) strips out the unused functions.
