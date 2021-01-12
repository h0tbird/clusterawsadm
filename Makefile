REPO := github.com/h0tbird/clusterawsadm

#------------------------------------------------------------------------------
# -s Omit the symbol table and debug information.
# -w Omit the DWARF symbol table.
#------------------------------------------------------------------------------

build:
	go build -ldflags="-s -w" .

#------------------------------------------------------------------------------
# Creates a local copy of the upstream provider module and replaces all the
# upstream references with references to the local copy.
#------------------------------------------------------------------------------

terraform-provider-aws: NAME := terraform-provider-aws
terraform-provider-aws: VERSION := v3.22.0
terraform-provider-aws: TMPDIR := $(shell mktemp -d)
terraform-provider-aws:
	@git clone --depth 1 --branch ${VERSION} https://github.com/hashicorp/${NAME}.git ${TMPDIR}
	@mkdir -p providers/${NAME} && rsync -a --delete --exclude='.*' ${TMPDIR}/ providers/${NAME}
	@go mod edit -replace ${REPO}/providers/${NAME}=./providers/${NAME} && rm -rf ${TMPDIR}
	@find providers/${NAME} -type f -print0 | \
	xargs -0 gsed -i 's_github.com/terraform-providers/${NAME}_${REPO}/providers/${NAME}_g'
	# TODO: Automate rewriting the provider's source code using go/parser go/ast go/printer.
	# Delete unreferenced entries in DataSourcesMap and ResourcesMap so that the code is not
	# referenced and the linker's DCE (Dead Code Elimination) strips out the unused functions.
