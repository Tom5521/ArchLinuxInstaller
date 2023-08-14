package main

import (
	"os"

	"github.com/Tom5521/ArchLinuxInstaller/src"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			src.Wifi()
			src.PacmanConf()
			src.Format()
			src.Mount()
			src.Pacstrap()
			src.Fstab()
			src.Grub()
			src.Keymap()
			src.FinalCmds()
		}
		if os.Args[1] == "pacstrap" {
			src.Wifi()
			src.PacmanConf()
			src.Mount()
			src.Pacstrap()
		}
		if os.Args[1] == "grub" {
			src.Mount()
			src.Fstab()
			src.Grub()
		}
	}
}
