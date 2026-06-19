package docs


import _ "embed"

//go:embed intro.txt
var IntroData []byte

//go:embed manual.txt
var ManualData []byte
