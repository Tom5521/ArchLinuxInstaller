package src

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/CmdRunTools/command"
	"github.com/Tom5521/MyGolangTools/file"
	"github.com/gookit/color"
)

var (
	NoPart       *bool = flag.Bool("nopart", false, "Skip the partitionating prosess (not open cfdisk)")
	NoGrub             = flag.Bool("nogrub", false, "Don't install Grub")
	NoWifi             = flag.Bool("nowifi", false, "Don't configure for wifi")
	NoFstab            = flag.Bool("nofstab", false, "Don't generate a fstab for the new system")
	NoKeymap           = flag.Bool("nokeymap", false, "Don't config the keymap for the new system")
	NoReboot           = flag.Bool("noreboot", false, "Don't reboot the system after the prosess")
	NoPasswd           = flag.Bool("nopasswd", false, "Skip the passwd setting")
	NoPacstrap         = flag.Bool("nopacstrap", false, "Skip the pacstrap process")
	NoPacmanConf       = flag.Bool("nopacmanconf", false, "Don't copy the temporal pacman.conf for best performance in the pacstrap")
	NoMount            = flag.Bool("nomount", false, "Don't mount the partitions")
	NoFormat           = flag.Bool("noformat", false, "Don't format the partitions")
)

var (
	// Functions
	sh = func() command.Cmd {
		cmd := command.Cmd{}
		cmd.CustomStd(true, true, true)
		//cmd.RunWithShell(true)
		return cmd
	}()
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

func Error(err string) {
	fmt.Println(color.Red.Render("ERROR:"), err)
}
func Warn(err string) {
	fmt.Println(color.Yellow.Render("Warning:"), err)
}

func Partitioning() {
	if *NoPart {
		return
	}
	var err error
	err = sh.SetAndRun("cfdisk " + conf.GrubInstallDisk)
	if err != nil {
		Error("Error starting cfdisk")
	}
}

func Format() {
	if *NoFormat {
		return
	}
	var err error
	// Set conditionals
	var (
		boot = partitions.Boot.Format && partitions.Boot.Filesystem == "fat32" && partitions.Boot.Partition != ""
		root = partitions.Root.Format && partitions.Root.Partition != ""
		home = partitions.Home.Format && partitions.Home.Partition != ""
		swap = partitions.Swap.Format && partitions.Swap.Partition != ""
	)
	// Format partitions
	sh.CustomStd(false, false, true)
	if boot {
		fmYellow(f("Formatting Boot <%v> fat32", partitions.Boot.Partition))
		err = sh.SetAndRun("mkfs.vfat -F 32 " + partitions.Boot.Partition)
	} else if partitions.Boot.Format {
		fmYellow(f("Formatting Boot <%v> %v", partitions.Boot.Partition, partitions.Boot.Filesystem))
		err = sh.SetAndRun(f("mkfs.%v -F %v", partitions.Boot.Filesystem, partitions.Boot.Partition))
	}
	if err != nil {
		Error("Error formatting boot")
	} else {
		fmGreen("Boot Formatted successfully!")
	}
	if root {
		fmYellow(f("Formatting Root <%v> %v", partitions.Root.Partition, partitions.Root.Filesystem))
		err = sh.SetAndRun(f("mkfs.%v -F %v", partitions.Root.Filesystem, partitions.Root.Partition))
		if err != nil {
			Error("Error formatting Root")
		} else {
			fmGreen("Root formatted successfully!")
		}
	}
	if home {
		fmYellow(f("Formatting Home <%v> %v", partitions.Home.Partition, partitions.Home.Filesystem))
		err = sh.SetAndRun(f("mkfs.%v -F %v", partitions.Home.Filesystem, partitions.Home.Partition))
		if err != nil {
			Error("Error formatting Home")
		} else {
			fmGreen("Home formatted successfully!")
		}
	}
	if swap {
		fmYellow(f("Making swap <%v>", partitions.Swap.Partition))
		err = sh.SetAndRun("mkswap " + partitions.Swap.Partition)
		if err != nil {
			Error("Error making swap")
		} else {
			fmGreen("Swap Maked successfully")
		}
	}
}
func Mount() {
	if *NoMount {
		return
	}
	var err error
	// Mount partitions
	fmYellow("mounting Root...")
	err = sh.SetAndRun(f("mount %v /mnt", partitions.Root.Partition))
	if err != nil {
		Error(f("Error mounting Root <%v>", partitions.Root.Partition))
	} else {
		fmGreen("Root mounted successfully!")
	}
	if conf.Uefi {
		fmYellow("uefi is true")
		if check_efi, _ := file.CheckFile("/mnt/efi"); !check_efi {
			err = os.Mkdir("/mnt/efi", os.ModePerm)
			if err != nil {
				Error("Error making /mnt/efi")
			} else {
				fmGreen("/mnt/efi maked successfully!")
			}
		}
		fmYellow(f("Mounting boot <%v> in </mnt/efi>", partitions.Boot.Partition))
		err = sh.SetAndRun(f("mount %v /mnt/efi", partitions.Boot.Partition))
		if err != nil {
			Error(f("Error mounting <%v> in /mnt/efi", partitions.Boot.Partition))
		} else {
			fmGreen("Boot mounted successfully!")
		}
	} else if check_Boot, _ := file.CheckFile("/mnt/boot"); !check_Boot {
		err = os.Mkdir("/mnt/boot", os.ModePerm)
		if err != nil {
			Error("Error making /mnt/boot")
		} else {
			fmGreen("/mnt/boot maked successfully!")
		}
		fmYellow("Mounting Boot...")
		err = sh.SetAndRun(f("mount %v /mnt/boot", partitions.Boot.Partition))
		if err != nil {
			Error(f("Error mounting Boot <%v> in /mnt/boot", partitions.Boot.Partition))
		} else {
			fmGreen("Boot mounted successfully!")
		}
	}
	if partitions.Home.Partition != "" {
		if checkdir, _ := file.CheckFile("/mnt/home"); !checkdir {
			err = os.Mkdir("/mnt/home", os.ModePerm)
			if err != nil {
				Error("Error making /mnt/home")
			} else {
				fmGreen("/mnt/home maked successfully!")
			}
		}
		fmYellow("Mounting home...")
		err = sh.SetAndRun(f("mount %v /mnt/home", partitions.Home.Partition))
		if err != nil {
			Error(f("Error mounting home <%v> in /mnt/home", partitions.Home.Partition))
		} else {
			fmGreen("Home mounted successfully!")
		}
	}
	if partitions.Swap.Partition != "" {
		fmYellow("Setting Swap...")
		err = sh.SetAndRun(f("swaplabel %v", partitions.Swap.Partition))
		err1 := sh.SetAndRun("swapon")
		if err1 != nil || err != nil {
			Error("Error Setting swap")
		} else {
			fmGreen("Swap setted.")
		}
	}
}

func PacmanConf() {
	if *NoPacmanConf {
		return
	}
	if check_pacman_cfg, _ := file.CheckFile("pacman.conf"); conf.CustomPacmanConfig && !check_pacman_cfg {
		fmYellow("No custom pacman conf found... Creatig a new one...")
		err := file.ReWriteFile("pacman.conf", pacmanConf)
		if err != nil {
			Warn("Error creating new pacman.conf")
		} else {
			fmGreen("pacman.conf created successfully!")
		}
	}
	pacmanfl, _ := file.ReadFileCont("/etc/pacman.conf")
	if conf.CustomPacmanConfig && !strings.Contains(string(pacmanfl), "---YES---THATS---MODIFIED---") {
		err := file.ReWriteFile("/etc/pacman.conf", pacmanConf)
		if err != nil {
			Warn("Error copying pacman.conf file")
		} else {
			fmGreen("pacman.conf pasted.")
		}
	} else {
		Warn("pacman.conf already pasted...")
	}
}

func Wifi() {
	if *NoWifi {
		return
	}
	var err error
	checkAdaptator, err := sh.SetAndOut("ip link")
	if err != nil {
		Error("Error running ip link")
	}
	if wifi.State {
		if strings.Contains(wifi.Adaptator, checkAdaptator) {
			err = sh.SetAndRun("rfkill unblock all")
			if err != nil {
				Error("Error unblocking rfkill")
			} else {
				fmGreen("rfkill unblocked")
			}
			err = sh.SetAndRun(f("ip link set %v up", wifi.Adaptator))
			if err != nil {
				Error("Error setting interface up")
			} else {
				fmGreen(f("<%v> setted up", wifi.Adaptator))
			}
			err = sh.SetAndRun(f("iwctl station %v connect %v --passphrase %v", wifi.Adaptator, wifi.Name, wifi.Password))
			if err != nil {
				Error("Error connecting to the wifi")
			} else {
				fmGreen("Connected to " + wifi.Name)
			}
		} else {
			Warn("Adaptator not found")
		}
		wifi_pkg = "networkmanger iwd"
	}
}

func Pacstrap() {
	if *NoPacstrap {
		return
	}
	if !conf.PacstrapSkip {
		fmYellow("Executing pacstrap...")
		err := sh.SetAndRun(f("pacstrap /mnt base base-devel grub git efibootmgr dialog wpa_supplicant nano linux linux-headers linux-firmware %v %v", wifi_pkg, conf.AdditionalPackages))
		if err != nil {
			Error("Error in pacstrap process")
		} else {
			fmGreen("Pacstrap completed successfully!")
		}
	} else {
		Warn("Skipping pacstrap...")
	}
}

func Fstab() {
	if *NoFstab {
		return
	}
	fmYellow("Generating Fstab...")
	err := sh.SetAndRun("genfstab -pU /mnt >> /mnt/etc/fstab")
	if err != nil {
		Error("Error Creating fstab")
	} else {
		fmGreen("Fstab created successfully!")
	}

}

func Grub() {
	if *NoGrub {
		return
	}
	var err error
	fmYellow("Installing grub...")
	if !conf.Uefi {
		err = sh.SetAndRun(f("echo exit|echo grub-install %v|arch-chroot /mnt", conf.GrubInstallDisk))
	} else {
		err = sh.SetAndRun(f("echo exit|echo grub-install %v --efi-directory /efi|arch-chroot /mnt", conf.GrubInstallDisk))
	}
	if err != nil {
		Error("Error installing grub")
	} else {
		fmGreen("Grub installed successfully!")
	}

	err = sh.SetAndRun("echo exit|echo grub-mkconfig -o /boot/grub/grub.cfg|arch-chroot /mnt")
	if err != nil {
		Error("Error in grub-mkconfig...")
	} else {
		fmGreen("grub-mkconfig maked successfully!")
	}
	if conf.PostInstallCommands != "" {
		err = sh.SetAndRun(conf.PostInstallChrootCommands)
		if err != nil {
			Error("Error in post install commands...")
		}
	}

}

func Keymap() {
	if *NoKeymap {
		return
	}
	sh.UseBashShell(true)
	keys, err := sh.SetAndOut("echo $(localectl list-keymaps)")
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(keys, conf.Keyboard) {
		err := sh.SetAndRun(f("echo KEYMAP=%v > /mnt/etc/vconsole.conf", conf.Keyboard))
		if err != nil {
			Error("Error setting KEYMAP in vconsole.conf")
		}
	} else {
		Warn("Keyboard specification not exist")
	}

}

func ConfigRootPasswd() {
	if *NoPasswd {
		return
	}
	fmYellow("Setting the root password:", conf.Password)
	cmd := exec.Command("chpasswd", "-R", "/mnt")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		Error("Error getting StdinPipe...")
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, fmt.Sprintf("root:%v", conf.Password))
	}()
	err = cmd.Run()
	if err != nil {
		Error("Error setting the password;" + err.Error())
	} else {
		fmGreen("root password setted successfully!")
	}

}

func FinalCmds() {
	var err error
	if conf.PostInstallChrootCommands != "" {
		sh.SetChroot("/mnt")
		err = sh.SetAndRun(conf.PostInstallChrootCommands)
		if err != nil {
			Error("Error in post-install-chroot cmds")
		}
	}
	if conf.ArchChroot {
		err = sh.SetAndRun("arch-chroot /mnt")
		if err != nil {
			Error("Error using arch-chroot.")
		}
	}

}

func Reboot() {
	if *NoReboot {
		return
	}
	var err error
	if !conf.Reboot {
		err = sh.SetAndRun("reboot")
		if err != nil {
			Error("Error in... in... Reboot? wtf")
			for i := 5; i != 0; i-- {
				time.Sleep(1 * time.Second)
				fmt.Print(f("\rRetrying in...%v seconds", i))
				if i == 1 {
					Reboot()
				}
			}
		}
	}
}
