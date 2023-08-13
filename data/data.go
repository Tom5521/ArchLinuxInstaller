package data

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
)

var (
	Yellow   = color.FgYellow.Render
	Red      = color.FgRed.Render
	filename = "config.yml"
)

func CheckDir(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

type yamlfile struct {
	CustomPacmanConfig bool   `yaml:"custom_pacman_config"`
	Keyboard           string `yaml:"keyboard"`
	Wifi               struct {
		State     bool   `yaml:"state"`
		Name      string `yaml:"name"`
		Adaptator string `yaml:"adaptator"`
		Password  string `yaml:"password"`
	} `yaml:"wifi"`
	Partitions struct {
		Boot struct {
			Partition  string `yaml:"partition"`
			Format     bool   `yaml:"format"`
			Filesystem string `yaml:"filesystem"`
		} `yaml:"boot"`
		Root struct {
			Partition  string `yaml:"partition"`
			Format     bool   `yaml:"format"`
			Filesystem string `yaml:"filesystem"`
		} `yaml:"root"`
		Home struct {
			Partition  string `yaml:"partition"`
			Format     bool   `yaml:"format"`
			Filesystem string `yaml:"filesystem"`
		} `yaml:"home"`
		Swap struct {
			Partition string `yaml:"partition"`
			Format    bool   `yaml:"format"`
		} `yaml:"swap"`
	} `yaml:"partitions"`
	GrubInstallDisk           string `yaml:"grub_install_disk"`
	PacstrapSkip              bool   `yaml:"pacstrap_skip"`
	AdditionalPackages        string `yaml:"additional_packages"`
	Uefi                      bool   `yaml:"uefi"`
	ArchChroot                bool   `yaml:"arch-chroot"`
	PostInstallCommands       string `yaml:"post_install_commands"`
	PostInstallChrootCommands string `yaml:"post_install_chroot_commands"`
	Reboot                    bool   `yaml:"reboot"`
}

func GetYamldata() yamlfile {
	yamldata := yamlfile{}
	if !CheckDir(filename) {
		fmt.Printf(Red(filename+" not found...") + Yellow("Creating a new one...\n"))
		NewYamlFile()
		if CheckDir(filename) {
			color.Green.Println("config file created!!!")
		}
		os.Exit(0)
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		color.Red.Println("Error reading " + filename)
	}
	err = yaml.Unmarshal(file, &yamldata)
	if err != nil {
		color.Red.Println("Error Unmarshalling the data")
	}
	return yamldata
}

func NewYamlFile() {
	yamlstruct := yamlfile{}
	file, err := os.Create(filename)
	if err != nil {
		color.Red.Println("Error creating " + filename)
		return
	}
	defer file.Close()
	data, err := yaml.Marshal(yamlstruct)
	if err != nil {
		color.Red.Println("Error Marshalling config file")
		return
	}
	_, err = file.WriteString(string(data))
	if err != nil {
		color.Red.Println("Error writing the data in the new yml file")
		return
	}
}
