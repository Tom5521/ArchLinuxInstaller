package main

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

func format() {
	// Set conditionals
	var (
		boot = partitions.Boot.Format && partitions.Boot.Filesystem == "fat32" && partitions.Boot.Partition != ""
		root = partitions.Root.Format && partitions.Root.Partition != ""
		home = partitions.Home.Format && partitions.Home.Partition != ""
		swap = partitions.Swap.Format && partitions.Swap.Partition != ""
	)
	// Format partitions
	if boot {
		sh.Cmd("mkfs.vfat -F 32" + partitions.Boot.Partition)
	} else {
		sh.Cmd(f("mkfs.%v -F %v", partitions.Boot.Filesystem, partitions.Boot.Partition))
	}
	if root {
		sh.Cmd(f("mkfs.%v -F %v", partitions.Root.Filesystem, partitions.Root.Partition))
	}
	if home {
		sh.Cmd(f("mkfs.%v -F %v", partitions.Home.Filesystem, partitions.Home.Partition))
	}
	if swap {
		sh.Cmd("mkswap " + partitions.Swap.Partition)
	}
}
func mount() {
	// Mount partitions
	sh.Cmd(f("mount %v /mnt", partitions.Root.Partition))
	if conf.Uefi {
		if !data.CheckDir("/mnt/efi") {
			os.Mkdir("/mnt/efi", os.ModePerm)
		}
		sh.Cmd(f("mount %v /mnt/efi", partitions.Boot.Partition))
	} else if !data.CheckDir("/mnt/boot") {
		os.Mkdir("/mnt/boot", os.ModePerm)
		sh.Cmd(f("mount %v /mnt/boot", partitions.Boot.Partition))
	}
	if partitions.Home.Partition != "" {
		if !data.CheckDir("/mnt/home") {
			os.Mkdir("/mnt/home", os.ModePerm)
		}
		sh.Cmd(f("mount %v /mnt/home", partitions.Home.Partition))
	}
	if partitions.Swap.Partition != "" {
		sh.Cmd(f("swaplabel %v", partitions.Swap.Partition))
		sh.Cmd("swapon")
	}
}

func main() {
	if conf.CustomPacmanConfig && !data.CheckDir("pacman.conf") {
		color.Red.Println("No custom pacman conf found... Downloading...")
		file, err := os.Create("pacman.conf")
		if err != nil {
			color.Red.Println("Error Creating new pacman.conf file")
			return
		}
		defer file.Close()
		response, err := http.Get("https://raw.githubusercontent.com/Tom5521/ArchLinux-Installer/master/pacman.conf")
		if err != nil {
			color.Red.Println("Error performing request")
			return
		}
		defer response.Body.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			color.Red.Println("Error copying response to file")
			return
		}
		color.Green.Println("pacman.conf downloaded successfully")
	}
	checkAdaptator, _ := sh.Out("ip link")
	if wifi.State {
		if strings.Contains(wifi.Adaptator, checkAdaptator) {
			sh.Cmd("rfkill unblock all")
			sh.Cmd(f("ip link set %v up", wifi.Adaptator))
			sh.Cmd(f("iwctl station %v connect %v --passphrase %v", wifi.Adaptator, wifi.Name, wifi.Password))
		} else {
			color.Red.Println("Warning: Adaptator not found")
		}
		wifi_pkg = "networkmanger iwd"
	}
	pacmanfl, _ := os.ReadFile("/etc/pacman.conf")
	if conf.CustomPacmanConfig && !strings.Contains("---YES---THATS---MODIFIED---", string(pacmanfl)) {
		sh.Cmd("cp pacman.conf /etc/")
	} else {
		color.Yellow.Println("pacman.conf already pasted...")
	}
	format()
	mount()
	if !conf.PacstrapSkip {
		sh.Cmd(f("pacstrap /mnt base base-devel grub git efibootmgr dialog wpa_supplicant nano linux linux-headers linux-firmware %v %v", wifi_pkg, conf.AdditionalPackages))
	} else {
		color.Yellow.Println("Skipping pacstrap...")
	}
	color.Yellow.Println("Generating Fstab...")
	sh.Cmd("genfstab -pU /mnt >> /mnt/etc/fstab")
	color.Yellow.Println("Installing grub...")
	if !conf.Uefi {
		sh.Cmd(f("echo exit|echo grub-install %v|arch-chroot /mnt", conf.GrubInstallDisk))
	} else {
		sh.Cmd(f("echo exit|echo grub-install %v --efi-directory /efi|arch-chroot /mnt", conf.GrubInstallDisk))
	}
	sh.Cmd("echo exit|echo grub-mkconfig -o /boot/grub/grub.cfg|arch-chroot /mnt")
	if conf.PostInstallCommands != "" {
		sh.Cmd(conf.PostInstallChrootCommands)
	}
	keys, _ := sh.Out("localectl list-keymaps")
	if strings.Contains(conf.Keyboard, keys) {
		sh.Cmd(f("echo exit|echo echo KEYMAP=%v > /mnt/etc/vconsole.conf|arch-chroot /mnt", conf.Keyboard))
	} else {
		color.Red.Println("WARNING:keyboard specification not exist")
	}
	if conf.PostInstallChrootCommands != "" {
		sh.Cmd(f("echo exit|echo %v|arch-chroot /mnt", conf.PostInstallChrootCommands))
	}
	color.Red.Println("please use the command \"passwd --root /mnt\" to set the root password before rebooting.")
	if conf.ArchChroot {
		sh.Cmd("arch-chroot /mnt")
	}
	if conf.Reboot {
		sh.Cmd("reboot")
	}

}
