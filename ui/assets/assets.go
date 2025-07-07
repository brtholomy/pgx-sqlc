package assets

import "embed"

//go:embed css/*
//go:embed js/*
var Assets embed.FS
