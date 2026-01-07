package templates

import (
	_ "embed"
)

//go:embed assets/logo.png
var LogoPNG []byte

//go:embed assets/logo_text.png
var LogoTextPNG []byte

//go:embed assets/logo.ico
var LogoICO []byte
