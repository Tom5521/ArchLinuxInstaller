package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/ArchLinuxInstaller/src"
	"github.com/Tom5521/CmdRunTools/command"
)

func main() {
	var sh = func() command.Cmd {
		cmd := command.Cmd{}
		cmd.CustomStd(true, true, true)
		return cmd
	}()

	if len(os.Args) == 0 {
		fmt.Println("Not enough arguments.")
		return
	}
	versionFlag := flag.Bool("version", false, "Show the version of the binary")
	helpFlag := flag.Bool("help", false, "Print the help message")

	// Only individual functions
	installFlag := flag.Bool("install", false, "Run all the nesesary functions to install completely Arch Linux")
	pacstrapFlag := flag.Bool("pacstrap", false, "Only runs the pacstrap functions")
	grubFlag := flag.Bool("grub", false, "Only installs Grub")
	passwdFlag := flag.Bool("passwd", false, "Only changes the password of the new root")
	partFlag := flag.Bool("part", false, "Only changes the password of the new root")
	mountFlag := flag.Bool("mount", false, "Only mounts the disks in her routes")
	newConfigFlag := flag.Bool("newconfig", false, "Creates a new config overwriting the original")

	flag.Parse()

	if *newConfigFlag {
		data.NewYamlFile()
		err := sh.SetAndRun("vim " + data.Pfilename)
		if err != nil {
			src.Error("Error oppening vim.\n" + err.Error())
		}
		return
	}

	parserFlags := []bool{*installFlag, *pacstrapFlag, *grubFlag, *passwdFlag, *partFlag, *mountFlag}
	catchBadFlags := func(flags []bool) bool {
		var trueValues int
		for _, parser := range flags {
			if parser {
				trueValues++
			}
		}
		return trueValues > 1
	}

	if catchBadFlags(parserFlags) {
		fmt.Println("There are arguments that cannot be used with others!")
		flag.PrintDefaults()
		return
	}

	if *versionFlag {
		fmt.Println("...")
		return
	}
	if *helpFlag {
		fmt.Println(`
Usage:
#[bin] [argument] -[option]`)
		flag.PrintDefaults()
		return
	}
	if *pacstrapFlag {
		src.Wifi()
		src.PacmanConf()
		src.Mount()
		src.Pacstrap()
	}
	if *installFlag {
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
		src.Reboot()
	}
	if *passwdFlag {
		src.ConfigRootPasswd()
	}
	if *partFlag {
		src.Partitioning()
	}
	if *grubFlag {
		src.Mount()
		src.Fstab()
		src.Grub()
	}
	if *mountFlag {
		src.Mount()
	}
}
