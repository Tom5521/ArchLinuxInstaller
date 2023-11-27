# ArchLinux Installer (written in go)

## Usage

First execute the binary with the `newconfig` arg, like this to create a empty config file
```
./[binary] [arg] [options]
```
After the configuration of the file execute `--install` arg to init the installation

You can see more info using `[binary] --help`

The help message:
```
Usage:
#[bin] [argument] -[option]
  -grub
        Only installs Grub
  -help
        Print the help message
  -install
        Run all the nesesary functions to install completely Arch Linux
  -mount
        Only mounts the disks in her routes
  -newconfig
        Creates a new config overwriting the original
  -noformat
        Don't format the partitions
  -nofstab
        Don't generate a fstab for the new system
  -nogrub
        Don't install Grub
  -nokeymap
        Don't config the keymap for the new system
  -nomount
        Don't mount the partitions
  -nopacmanconf
        Don't copy the temporal pacman.conf for best performance in the pacstrap
  -nopacstrap
        Skip the pacstrap process
  -nopart
        Skip the partitionating prosess (not open cfdisk)
  -nopasswd
        Skip the passwd setting
  -noreboot
        Don't reboot the system after the prosess
  -nowifi
        Don't configure for wifi
  -pacstrap
        Only runs the pacstrap functions
  -part
        Only changes the password of the new root
  -passwd
        Only changes the password of the new root
  -version
        Show the version of the binary

```

### Configuration

The YAML file is a data structure used to configure the installation parameters of the script. Below are the available configuration options:

`custom_pacman_config`: A boolean value indicating whether to use a custom pacman configuration. The default value is "false". 
`keyboard`: Configures the keyboard layout for the installation. The default value is an empty string.

#### Wifi

- `wifi`: An object specifying the wireless network configuration. It consists of the following four properties: 
1. `state`: A boolean value indicating whether to enable the wireless connection. The default value is "false". 
2. `name`: The name of the wireless network to connect to. The default value is an empty string. 
3. `adapter`: The wireless network adapter to use for the connection. The default value is an empty string. 
4. `password`: The password for the wireless network. The default value is an empty string.

#### Partitions

`partitions`: An object specifying the partitions to be used for the installation. It consists of the following properties:

1.  `boot`: An object specifying the boot partition. It consists of the following properties:

*   `partition`: The partition device. The default value is "".
*   `format`: A boolean value indicating whether to format the partition. The default value is "".
*   `filesystem`: The file system to be used for the partition. The default value is "".

2.  `root`: An object specifying the root partition. It has the same properties as the boot partition.
3.  `home`: An object specifying the home folder partition. It has the same properties as the boot partition.
4.  `swap`: An object specifying the swap partition. It has the same properties as the boot partition, except it does not have the "filesystem" property.

#### Extra Configs

`passwd`: Set the root password to be used for the new installation

`grub_install_disk`: The storage device where the GRUB bootloader will be installed. The default value is "".

`pacstrap_skip`: A boolean value indicating whether to skip the installation of Arch Linux basic packages. The default value is "false".

`additional_packages`: A string specifying the additional packages to install. The default value is an empty string.

`uefi`: A boolean value indicating whether the system uses UEFI instead of BIOS. The default value is "false".

`arch-chroot`: A boolean value indicating whether to run the script in the Arch Linux chroot environment. The default value is "false".

`post_install_commands`: A string specifying the commands to be executed after the package installation. The default value is an empty string.

`post_install_chroot_commands`: A string specifying the commands to be executed after the package installation in the Arch Linux chroot environment. The default value is an empty string.

`reboot`: Allows specifying whether to restart the system after installation. If this field is omitted, the default configuration will be used.
