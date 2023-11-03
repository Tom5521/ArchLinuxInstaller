package main

import (
	"fmt"
	"os"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/ArchLinuxInstaller/src"
	"github.com/Tom5521/MyGolangTools/commands"
)

var sh = commands.Sh{}

var HelpStr = `
Usage:
#[bin] [argument] -[option]

Arguments:
- help				Print this text

- install			Run all the nesesary functions to install completely Arch Linux 

- pacstrap			Only runs the pacstrap functions

- grub				Only installs Grub

- newconfig			Creates a new config overwriting the original

- part				Open cfdisk to partitionate the install disk

- passwd			Only changes the password of the new root 

- mount				Only mounts the disks in her routes

Options:
Argument options will be applied before the config fil

-nopasswd			Skip the passwd set

-nopacstrap			Skip the pacstrap prosess

-nopart				Skip the partitionating prosess (not open cfdisk)

-nogrub				Don't install Grub

-noformat			Don't format the partitions

-nomount			Don't mount the partitions

-nopacmanconf			Don't paste custom pacman config

-nowifi				Don't configure for wifi 

-nofstab			Don't generate a fstab for the new system

-nokeymap			Don't config the keymap for the new system

`

func PrintHelp() {
	fmt.Print(HelpStr)
}

func main() {
	if len(os.Args) == 0 {
		fmt.Println("Not enough arguments")
		return
	}
	switch os.Args[1] {
	case "help":
		PrintHelp()
	case "passwd":
		src.ConfigRootPasswd()
	case "part":
		src.Partitioning()
	case "install":
		src.Partitioning()
		src.Wifi()
		src.PacmanConf()
		src.Format()
		src.Mount()
		src.Pacstrap()
		src.Fstab()
		src.Grub()
		src.Keymap()
		src.ConfigRootPasswd()
		src.FinalCmds()
	case "pacstrap":
		src.Wifi()
		src.PacmanConf()
		src.Mount()
		src.Pacstrap()
	case "grub":
		src.Mount()
		src.Fstab()
		src.Grub()
	case "mount":
		src.Mount()
	case "newconfig":
		data.NewYamlFile()
		err := sh.Cmd("vim " + data.Pfilename)
		if err != nil {
			src.Error("Error oppening vim.\n" + err.Error())
		}
	}
}
