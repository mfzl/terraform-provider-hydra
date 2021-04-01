module github.com/mfzl/terraform-provider-hydra

go 1.16

replace github.com/rubenv/sql-migrate v0.0.0-20180704111356-ba2c6a7295c59448dbc195cef2f41df5163b3892 => github.com/rubenv/sql-migrate v0.0.0-20191213152630-06338513c237

require (
	github.com/hashicorp/go-hclog v0.0.0-20190109152822-4783caec6f2e // indirect
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mattn/go-isatty v0.0.8 // indirect
	github.com/ory/hydra v1.0.0-rc.6.0.20190103103112-e2b88d211a27
	github.com/ory/hydra-client-go v1.9.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	gopkg.in/resty.v1 v1.11.0 // indirect
)
