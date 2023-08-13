package src

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Tom5521/ArchLinuxInstaller/data"
	"github.com/Tom5521/MyGolangTools/commands"
	"github.com/gookit/color"
)

var (
	conf       = data.GetYamldata()
	partitions = conf.Partitions
	wifi       = conf.Wifi
	sh         = commands.Sh{}
	f          = fmt.Sprintf
	wifi_pkg   string
)

func Format() error {
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
		err = sh.Cmd("mkfs.vfat -F 32" + partitions.Boot.Partition)
		if err != nil {
			return err
		}
	} else {
		err = sh.Cmd(f("mkfs.%v -F %v", partitions.Boot.Filesystem, partitions.Boot.Partition))
		if err != nil {
			return err
		}
	}
	if root {
		err = sh.Cmd(f("mkfs.%v -F %v", partitions.Root.Filesystem, partitions.Root.Partition))
		if err != nil {
			return err
		}
	}
	if home {
		err = sh.Cmd(f("mkfs.%v -F %v", partitions.Home.Filesystem, partitions.Home.Partition))
		if err != nil {
			return err
		}
	}
	if swap {
		err = sh.Cmd("mkswap " + partitions.Swap.Partition)
		if err != nil {
			return err
		}
	}
	return nil
}
func Mount() error {
	var err error
	// Mount partitions
	err = sh.Cmd(f("mount %v /mnt", partitions.Root.Partition))
	if err != nil {
		return err
	}
	if conf.Uefi {
		if !data.CheckDir("/mnt/efi") {
			err = os.Mkdir("/mnt/efi", os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = sh.Cmd(f("mount %v /mnt/efi", partitions.Boot.Partition))
		if err != nil {
			return err
		}
	} else if !data.CheckDir("/mnt/boot") {
		err = os.Mkdir("/mnt/boot", os.ModePerm)
		if err != nil {
			return err
		}
		err = sh.Cmd(f("mount %v /mnt/boot", partitions.Boot.Partition))
		if err != nil {
			return err
		}
	}
	if partitions.Home.Partition != "" {
		if !data.CheckDir("/mnt/home") {
			err = os.Mkdir("/mnt/home", os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = sh.Cmd(f("mount %v /mnt/home", partitions.Home.Partition))
		if err != nil {
			return err
		}
	}
	if partitions.Swap.Partition != "" {
		err = sh.Cmd(f("swaplabel %v", partitions.Swap.Partition))
		if err != nil {
			return err
		}
		err = sh.Cmd("swapon")
		if err != nil {
			return err
		}
	}
	return nil
}

func PacmanConf() error {
	if conf.CustomPacmanConfig && !data.CheckDir("pacman.conf") {
		color.Red.Println("No custom pacman conf found... Downloading...")
		file, err := os.Create("pacman.conf")
		if err != nil {
			color.Red.Println("Error Creating new pacman.conf file")
			return err
		}
		defer file.Close()
		response, err := http.Get("https://raw.githubusercontent.com/Tom5521/ArchLinux-Installer/master/pacman.conf")
		if err != nil {
			color.Red.Println("Error performing request")
			return err
		}
		defer response.Body.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			color.Red.Println("Error copying response to file")
			return err
		}
		color.Green.Println("pacman.conf downloaded successfully")
	}

	pacmanfl, _ := os.ReadFile("/etc/pacman.conf")
	if conf.CustomPacmanConfig && !strings.Contains("---YES---THATS---MODIFIED---", string(pacmanfl)) {
		err := sh.Cmd("cp pacman.conf /etc/")
		if err != nil {
			return err
		}
	} else {
		color.Yellow.Println("pacman.conf already pasted...")
	}
	return nil
}

func Wifi() error {
	var err error
	checkAdaptator, _ := sh.Out("ip link")
	if wifi.State {
		if strings.Contains(wifi.Adaptator, checkAdaptator) {
			err = sh.Cmd("rfkill unblock all")
			if err != nil {
				return err
			}
			err = sh.Cmd(f("ip link set %v up", wifi.Adaptator))
			if err != nil {
				return err
			}
			err = sh.Cmd(f("iwctl station %v connect %v --passphrase %v", wifi.Adaptator, wifi.Name, wifi.Password))
			if err != nil {
				return err
			}
		} else {
			color.Red.Println("Warning: Adaptator not found")
		}
		wifi_pkg = "networkmanger iwd"
	}
	return nil
}

func Pacstrap() error {
	if !conf.PacstrapSkip {
		err := sh.Cmd(f("pacstrap /mnt base base-devel grub git efibootmgr dialog wpa_supplicant nano linux linux-headers linux-firmware %v %v", wifi_pkg, conf.AdditionalPackages))
		if err != nil {
			return err
		}
	} else {
		color.Yellow.Println("Skipping pacstrap...")
	}
	return nil
}

func Fstab() error {
	color.Yellow.Println("Generating Fstab...")
	err := sh.Cmd("genfstab -pU /mnt >> /mnt/etc/fstab")
	if err != nil {
		return err
	}
	return nil
}

func Grub() error {
	var err error
	color.Yellow.Println("Installing grub...")
	if !conf.Uefi {
		err = sh.Cmd(f("echo exit|echo grub-install %v|arch-chroot /mnt", conf.GrubInstallDisk))
		if err != nil {
			return err
		}
	} else {
		err = sh.Cmd(f("echo exit|echo grub-install %v --efi-directory /efi|arch-chroot /mnt", conf.GrubInstallDisk))
		if err != nil {
			return err
		}
	}
	err = sh.Cmd("echo exit|echo grub-mkconfig -o /boot/grub/grub.cfg|arch-chroot /mnt")
	if err != nil {
		return err
	}
	if conf.PostInstallCommands != "" {
		err = sh.Cmd(conf.PostInstallChrootCommands)
		if err != nil {
			return err
		}
	}
	return nil
}

func Keymap() error {
	keys, _ := sh.Out("localectl list-keymaps")
	if strings.Contains(conf.Keyboard, keys) {
		err := sh.Cmd(f("echo exit|echo echo KEYMAP=%v > /mnt/etc/vconsole.conf|arch-chroot /mnt", conf.Keyboard))
		if err != nil {
			return err
		}
	} else {
		color.Red.Println("WARNING:keyboard specification not exist")
	}
	return nil
}
