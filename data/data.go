package data

import (
	"fmt"
	"os"

	"github.com/Tom5521/MyGolangTools/file"
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
)

var (
	Yellow    = color.FgYellow.Render
	Red       = color.FgRed.Render
	Pfilename = filename
	filename  = "config.yml"
)

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
	Password                  string `yaml:"passwd"`
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
	if check, _ := file.CheckFile(filename); !check {
		fmt.Printf(Red(filename+" not found...") + Yellow("Creating a new one...\n"))
		NewYamlFile()
		if newcheck, _ := file.CheckFile(filename); newcheck {
			color.Green.Println("config file created!!!")
		}
		os.Exit(0)
	}
	file, err := file.ReadFileCont(filename)
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
	data, err := yaml.Marshal(yamlstruct)
	if err != nil {
		color.Red.Println("Error Marshalling config file")
		return
	}
	err = file.ReWriteFile(filename, string(data))
	if err != nil {
		color.Red.Println("Error writing the data in the new yml file")
		return
	}
}
