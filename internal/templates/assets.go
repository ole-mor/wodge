package templates

import (
	_ "embed"
)

//go:embed ../../logo.png
var LogoPNG []byte

//go:embed ../../logo_text.png
var LogoTextPNG []byte
