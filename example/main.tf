provider "hydra" {
    # all options maybe omitted, they will be taken
    # from environment variables listed in README
    client_id = ""
    client_secret = ""
    cluster_url = ""
    skip_tls_verify = false 
}

resource "hydra_client" "main" {
    client_id = "mainclient"
    name = "main clients"
    response_types = ["id_token", "code", "token"]
    grant_types = ["authorization_code", "client_credentials"]
    owner = "org.hydra.com"

    redirect_uris = ["https://localhost:8080/callback", "blz://login"]
    scope = ["hydra"]
}

resource "hydra_policy" "main" {
    description = "One policy to rule them all."
    subjects = ["clients:${hydra_client.main.id}", "users:<[peter|ken]>", "users:maria"]
    resources = [
        "resources:articles:<.*>",
        "resources:printer"
    ],
    actions = ["delete", "<[create|update]>"]
    effect = "allow"

    condition {
        name = "remoteIP"
        type = "CIDRCondition"
        options {
            cidr = "192.168.0.2/16"
        }
    }
}

output "hydra_client_id" {
    value = "${hydra_client.main.id}" 
}
