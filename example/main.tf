provider "hydra" {
    admin_url = ""
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

output "hydra_client_id" {
    value = "${hydra_client.main.id}" 
}
