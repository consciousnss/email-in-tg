group "default" {
  targets = ["image"]
}

target "image" {
  target     = "image"
  dockerfile = "Dockerfile"
  tags = [
    "ghcr.io/un1uckyyy/email-in-tg:dev"
  ]
}
