package src

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/MyGolangTools/commands"
	"github.com/Tom5521/MyGolangTools/file"
	"github.com/gookit/color"
)

var (
	// Functions
	sh       = commands.Sh{}
	f        = fmt.Sprintf // Set a more comfortable alias for fmt.Sprintf()
	fmRed    = color.Red.Println
	fmGreen  = color.Green.Println
	fmYellow = color.Yellow.Println
	// Data
	conf       = data.GetYamldata()
	partitions = conf.Partitions
	wifi       = conf.Wifi
	wifi_pkg   string
	pacmanConf = `
#Pacman-config-modifyed by Tom5521 ---YES---THATS---MODIFIED---

[options]
HoldPkg     = pacman glibc

Architecture = auto

Color
CheckSpace
VerbosePkgLists
ParallelDownloads = 30
ILoveCandy

SigLevel    = Required DatabaseOptional
LocalFileSigLevel = Optional

[core]
Include = /etc/pacman.d/mirrorlist

[extra]
Include = /etc/pacman.d/mirrorlist

[community]
Include = /etc/pacman.d/mirrorlist

[multilib]
Include = /etc/pacman.d/mirrorlist`
)

func Format() {
	var err error
	// Set conditionals
	var (
		boot = partitions.Boot.Format && partitions.Boot.Filesystem == "fat32" && partitions.Boot.Partition != ""
		root = partitions.Root.Format && partitions.Root.Partition != ""
		home = partitions.Home.Format && partitions.Home.Partition != ""
		swap = partitions.Swap.Format && partitions.Swap.Partition != ""
	)
	// Format partitions
	if boot {
		fmYellow(f("Formatting Boot <%v> fat32", partitions.Boot.Partition))
		err = sh.Cmd("mkfs.vfat -F 32 " + partitions.Boot.Partition)
	} else if partitions.Boot.Format {
		fmYellow(f("Formatting Boot <%v> %v", partitions.Boot.Partition, partitions.Boot.Filesystem))
		err = sh.Cmd(f("mkfs.%v -F %v", partitions.Boot.Filesystem, partitions.Boot.Partition))
	}
	if err != nil {
		fmRed("Error formatting boot")
	} else {
		fmGreen("Boot Formatted successfully!")
	}
	if root {
		fmYellow(f("Formatting Root <%v> %v", partitions.Root.Partition, partitions.Root.Filesystem))
		err = sh.Cmd(f("mkfs.%v -F %v", partitions.Root.Filesystem, partitions.Root.Partition))
		if err != nil {
			fmRed("Error formatting Root")
		} else {
			fmGreen("Root formatted successfully!")
		}
	}
	if home {
		fmYellow(f("Formatting Home <%v> %v", partitions.Home.Partition, partitions.Home.Filesystem))
		err = sh.Cmd(f("mkfs.%v -F %v", partitions.Home.Filesystem, partitions.Home.Partition))
		if err != nil {
			fmRed("Error formatting Home")
		} else {
			fmGreen("Home formatted successfully!")
		}
	}
	if swap {
		fmYellow(f("Making swap <%v>", partitions.Swap.Partition))
		err = sh.Cmd("mkswap " + partitions.Swap.Partition)
		if err != nil {
			fmRed("Error making swap")
		} else {
			fmGreen("Swap Maked successfully")
		}
	}
}
func Mount() {
	var err error
	// Mount partitions
	fmYellow("mounting Root...")
	err = sh.Cmd(f("mount %v /mnt", partitions.Root.Partition))
	if err != nil {
		fmRed(f("Error mounting Root <%v>", partitions.Root.Partition))
	} else {
		fmGreen("Root mounted successfully!")
	}
	if conf.Uefi {
		fmYellow("uefi is true")
		if check_efi, _ := file.CheckFile("/mnt/efi"); !check_efi {
			err = os.Mkdir("/mnt/efi", os.ModePerm)
			if err != nil {
				fmRed("Error making /mnt/efi")
			} else {
				fmGreen("/mnt/efi maked successfully!")
			}
		}
		fmYellow(f("Mounting boot <%v> in </mnt/efi>", partitions.Boot.Partition))
		err = sh.Cmd(f("mount %v /mnt/efi", partitions.Boot.Partition))
		if err != nil {
			fmRed(f("Error mounting <%v> in /mnt/efi", partitions.Boot.Partition))
		} else {
			fmGreen("Boot mounted successfully!")
		}
	} else if check_Boot, _ := file.CheckFile("/mnt/boot"); !check_Boot {
		err = os.Mkdir("/mnt/boot", os.ModePerm)
		if err != nil {
			fmRed("Error making /mnt/boot")
		} else {
			fmGreen("/mnt/boot maked successfully!")
		}
		fmYellow("Mounting Boot...")
		err = sh.Cmd(f("mount %v /mnt/boot", partitions.Boot.Partition))
		if err != nil {
			fmRed(f("Error mounting Boot <%v> in /mnt/boot", partitions.Boot.Partition))
		} else {
			fmGreen("Boot mounted successfully!")
		}
	}
	if partitions.Home.Partition != "" {
		if checkdir, _ := file.CheckFile("/mnt/home"); !checkdir {
			err = os.Mkdir("/mnt/home", os.ModePerm)
			if err != nil {
				fmRed("Error making /mnt/home")
			} else {
				fmGreen("/mnt/home maked successfully!")
			}
		}
		fmYellow("Mounting home...")
		err = sh.Cmd(f("mount %v /mnt/home", partitions.Home.Partition))
		if err != nil {
			fmRed(f("Error mounting home <%v> in /mnt/home", partitions.Home.Partition))
		} else {
			fmGreen("Home mounted successfully!")
		}
	}
	if partitions.Swap.Partition != "" {
		fmYellow("Setting Swap...")
		err = sh.Cmd(f("swaplabel %v", partitions.Swap.Partition))
		err1 := sh.Cmd("swapon")
		if err1 != nil || err != nil {
			fmRed("Error Setting swap")
		} else {
			fmGreen("Swap setted.")
		}
	}
}

func PacmanConf() {
	if check_pacman_cfg, _ := file.CheckFile("pacman.conf"); conf.CustomPacmanConfig && !check_pacman_cfg {
		fmYellow("No custom pacman conf found... Creatig a new one...")
		err := file.ReWriteFile("pacman.conf", pacmanConf)
		if err != nil {
			fmRed("Error creating new pacman.conf")
		} else {
			fmGreen("pacman.conf created successfully!")
		}
	}
	pacmanfl, _ := file.ReadFileCont("/etc/pacman.conf")
	if conf.CustomPacmanConfig && !strings.Contains("---YES---THATS---MODIFIED---", string(pacmanfl)) {
		err := file.ReWriteFile("/etc/pacman.conf", pacmanConf)
		if err != nil {
			fmRed("Error copying pacman.conf file")
		} else {
			fmGreen("pacman.conf pasted.")
		}
	} else {
		fmYellow("pacman.conf already pasted...")
	}
}

func Wifi() {
	var err error
	checkAdaptator, _ := sh.Out("ip link")
	if wifi.State {
		if strings.Contains(wifi.Adaptator, checkAdaptator) {
			err = sh.Cmd("rfkill unblock all")
			if err != nil {
				fmRed("Error unblocking rfkill")
			} else {
				fmGreen("rfkill unblocked")
			}
			err = sh.Cmd(f("ip link set %v up", wifi.Adaptator))
			if err != nil {
				fmRed("Error setting interface up")
			} else {
				fmGreen(f("<%v> setted up", wifi.Adaptator))
			}
			err = sh.Cmd(f("iwctl station %v connect %v --passphrase %v", wifi.Adaptator, wifi.Name, wifi.Password))
			if err != nil {
				fmRed("Error connecting to the wifi")
			} else {
				fmGreen("Connected to " + wifi.Name)
			}
		} else {
			fmRed("Warning: Adaptator not found")
		}
		wifi_pkg = "networkmanger iwd"
	}
}

func Pacstrap() {
	if !conf.PacstrapSkip {
		fmYellow("Executing pacstrap...")
		err := sh.Cmd(f("pacstrap /mnt base base-devel grub git efibootmgr dialog wpa_supplicant nano linux linux-headers linux-firmware %v %v", wifi_pkg, conf.AdditionalPackages))
		if err != nil {
			fmRed("Error in pacstrap process")
		} else {
			fmGreen("Pacstrap completed successfully!")
		}
	} else {
		fmYellow("Skipping pacstrap...")
	}
}

func Fstab() {
	fmYellow("Generating Fstab...")
	err := sh.Cmd("genfstab -pU /mnt >> /mnt/etc/fstab")
	if err != nil {
		fmRed("Error Creating fstab")
	} else {
		fmGreen("Fstab created successfully!")
	}

}

func Grub() {
	var err error
	fmYellow("Installing grub...")
	if !conf.Uefi {
		err = sh.Cmd(f("echo exit|echo grub-install %v|arch-chroot /mnt", conf.GrubInstallDisk))
	} else {
		err = sh.Cmd(f("echo exit|echo grub-install %v --efi-directory /efi|arch-chroot /mnt", conf.GrubInstallDisk))
	}
	if err != nil {
		fmRed("Error installing grub")
	} else {
		fmGreen("Grub installed successfully!")
	}

	err = sh.Cmd("echo exit|echo grub-mkconfig -o /boot/grub/grub.cfg|arch-chroot /mnt")
	if err != nil {
		fmRed("Error in grub-mkconfig...")
	} else {
		fmGreen("grub-mkconfig maked successfully!")
	}
	if conf.PostInstallCommands != "" {
		err = sh.Cmd(conf.PostInstallChrootCommands)
		if err != nil {
			fmRed("Error in post install commands...")
		}
	}

}

func Keymap() {
	keys, _ := sh.Out("localectl list-keymaps")
	if strings.Contains(conf.Keyboard, keys) {
		err := sh.Cmd(f("echo exit|echo echo KEYMAP=%v > /mnt/etc/vconsole.conf|arch-chroot /mnt", conf.Keyboard))
		if err != nil {
			fmRed("Error setting KEYMAP in vconsole.conf")
		}
	} else {
		fmRed("WARNING:keyboard specification not exist")
	}
}

func FinalCmds() {
	var err error
	if conf.PostInstallChrootCommands != "" {
		err = sh.Cmd(f("echo exit|echo %v|arch-chroot /mnt", conf.PostInstallChrootCommands))
		if err != nil {
			fmRed("Error in post-install-chroot cmds")
		}
	}
	fmRed("please use the command \"passwd --root /mnt\" to set the root password before rebooting.")
	if conf.ArchChroot {
		err = sh.Cmd("arch-chroot /mnt")
		if err != nil {
			fmRed("Error using arch-chroot.")
		}
	}
	if conf.Reboot {
		err = sh.Cmd("reboot")
		if err != nil {
			fmRed("Error in... in... Reboot? wtf")
		}
	}
}
