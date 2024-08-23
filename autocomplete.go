package main

import (
	_ "embed"
	"fmt"
)

//go:embed static/bash_autocomplete
var bash_autocomplete string

//go:embed static/zsh_autocomplete
var zsh_autocomplete string

//go:embed static/powershell_autocomplete.ps1
var powershell_autocomplete string

func autocomplete_bash() {
	fmt.Println(bash_autocomplete)
}

func autocomplete_powershell() {
	fmt.Println(powershell_autocomplete)
}

func autocomplete_zsh() {
	fmt.Println(zsh_autocomplete)
}
