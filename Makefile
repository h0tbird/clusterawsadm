#------------------------------------------------------------------------------
#
#------------------------------------------------------------------------------

build:
	go build -ldflags="-s -w" .

#------------------------------------------------------------------------------
#
#------------------------------------------------------------------------------

terraform-provider-aws: NAME := terraform-provider-aws
terraform-provider-aws: VERSION := v3.22.0
terraform-provider-aws: TMPDIR := $(shell mktemp -d)
terraform-provider-aws:
	@git clone --depth 1 --branch ${VERSION} https://github.com/hashicorp/${NAME}.git ${TMPDIR}
	@mkdir -p providers/${NAME} && rsync -a --delete --exclude='.*' ${TMPDIR}/ providers/${NAME}
	@rm -rf ${TMPDIR} && find providers/${NAME} -type f -print0 | \
	xargs -0 gsed -i 's_github.com/terraform-providers/${NAME}_github.com/h0tbird/clusterawsadm/providers/${NAME}_g'