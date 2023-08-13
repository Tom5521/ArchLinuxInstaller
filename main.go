package main

import (
	"fmt"
	"os"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/ArchLinuxInstaller/src"
	"github.com/Tom5521/MyGolangTools/commands"
	"github.com/gookit/color"
)

var (
	conf = data.GetYamldata()
	sh   = commands.Sh{}
	f    = fmt.Sprintf
	err  error
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err = src.Wifi()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.PacmanConf()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.Format()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.Mount()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.Pacstrap()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.Fstab()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.Grub()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = src.Keymap()
			if err != nil {
				fmt.Println(err)
				return
			}
			if conf.PostInstallChrootCommands != "" {
				err = sh.Cmd(f("echo exit|echo %v|arch-chroot /mnt", conf.PostInstallChrootCommands))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			color.Red.Println("please use the command \"passwd --root /mnt\" to set the root password before rebooting.")
			if conf.ArchChroot {
				err = sh.Cmd("arch-chroot /mnt")
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if conf.Reboot {
				err = sh.Cmd("reboot")
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
