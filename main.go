package main

import (
	"os"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/ArchLinuxInstaller/src"
	"github.com/Tom5521/MyGolangTools/commands"
	"github.com/gookit/color"
)

var sh = commands.Sh{}

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
		if os.Args[1] == "newconfig" {
			data.NewYamlFile()
			err := sh.Cmd("vim " + data.Pfilename)
			if err != nil {
				color.Red.Println("Error oppening vim." + err.Error())
			}
		}
	}
}
