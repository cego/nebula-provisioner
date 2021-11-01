package webapp

import "embed"

//go:embed dist/*
var Webapp embed.FS
