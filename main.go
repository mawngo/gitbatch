package main

import "github.com/sitdownrightnow/gitbatch/cmd"

func main() {
	cli := cmd.NewCLI()
	cli.Execute()
}
