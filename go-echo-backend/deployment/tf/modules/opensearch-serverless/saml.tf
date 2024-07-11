data "http" "saml_metadata_document" {
  count = var.saml_metadata_document_url != "" ? 1 : 0
  url = var.saml_metadata_document_url

  request_headers = {
    "Accept" = "application/samlmetadata+xml"
  }
}

resource "local_file" "saml_metadata_document" {
  count = var.saml_metadata_document_url != "" ? 1 : 0
  
  content  = "${data.http.saml_metadata_document[0].body}"
  filename = "${path.module}/${local.saml_metadata_document}"
  
  depends_on = [ data.http.saml_metadata_document ]
}

resource "aws_opensearchserverless_security_config" "default" {
  count = var.saml_metadata_document_url != "" ? 1 : 0

  name = var.saml_provider_name
  type = "saml"
  saml_options {
    metadata = var.saml_metadata_document_url != "" ? file("${path.module}/${local.saml_metadata_document}") : ""
    group_attribute = var.saml_group
  }
  
  depends_on = [ local_file.saml_metadata_document ]
  
}
