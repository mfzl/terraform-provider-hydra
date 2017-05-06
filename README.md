# Terraform provider for Hydra

[Hydra](https://github.com/ory/hydra) is an open source OAuth2 and OpenID Connect server.

This is a [Terraform](https://terraform.io) provider to create OAuth Clients and Policies required by the services in an 
infrastructure.  


## Available resources

- `hydra_client`
    Manages OAuth2 clients with given configuration

- `hydra_policy`
    Manages Hydra policies

Please have a look at the example configuration in example directory to check which options are available.

For a detailed explanation of each configuration option check hydra [API documentation](http://docs.hydra13.apiary.io/)


## Installation

Check terraform guide on [installtion](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin)


## Provider options

These default environment variables will be used if omitted from provider config block.

    - `HYDRA_CLIENT_ID` 
    - `HYDRA_CLIENT_SECRET` 
    - `HYDRA_CLUSTER_URL` 

