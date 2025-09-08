group "default" {
  targets = ["image"]
}

target "image" {
  target     = "image"
  dockerfile = "Dockerfile"
  tags = [
    "ghcr.io/consciousnss/email-in-tg:dev"
  ]
}
