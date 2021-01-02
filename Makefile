terraform-provider-aws:
	rsync -a --exclude='.*' ~/git/hashicorp/terraform-provider-aws/ providers/terraform-provider-aws
	find providers/terraform-provider-aws -type f -print0 | \
	xargs -0 gsed -i 's_github.com/terraform-providers/terraform-provider-aws_github.com/h0tbird/clusterawsadm/providers/terraform-provider-aws_g'