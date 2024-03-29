//go:build linux

package env

import (
	"encoding/xml"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"server/config"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type linux struct {
}

type Action struct {
	XMLName     xml.Name `xml:"action"`
	Icon        string   `xml:"icon"`
	Name        string   `xml:"name"`
	Uniqueid    string   `xml:"unique-id"`
	Command     string   `xml:"command"`
	Description string   `xml:"description"`
	Patterns    string   `xml:"patterns"`
	Directories string   `xml:"directories"`
	Audiofiles  string   `xml:"audio-files"`
	Imagefiles  string   `xml:"image-files"`
	Otherfiles  string   `xml:"other-files"`
	Textfiles   string   `xml:"text-files"`
	Videofiles  string   `xml:"video-files"`
}

type Actions struct {
	XMLName xml.Name `xml:"actions"`
	Actions []Action `xml:"action"`
}

func (sys *linux) createAction(icon string, name string, ucaId string, command string, description string, patterns string) *Action {
	action := Action{Icon: icon, Name: name,
		Uniqueid: ucaId, Command: command,
		Description: description, Patterns: patterns}
	action.XMLName.Local = "action"
	return &action
}

func (sys *linux) SpecificSetup() {
	log := utils.GetLogger()
	//Setup patch to uca.xml
	ucaDirPath := GetHomeDir() + "/.config/Thunar/uca.xml"

	ucaFile, err := os.Open(ucaDirPath)
	if err != nil {
		ucaFile.Close()
		log.Fatal("Error when trying to open uca.xml: " + err.Error())
	}

	ucaByteValue, _ := io.ReadAll(ucaFile)
	var actions Actions
	err = xml.Unmarshal([]byte(ucaByteValue), &actions)
	if err != nil {
		ucaFile.Close()
		log.Fatal("Error when unmarshaling uca.xml: " + err.Error())
	}

	ucaFile.Close()

	foundIdx := slices.IndexFunc(actions.Actions, func(action Action) bool { return strings.Contains(action.Command, config.App_client_name) })
	if foundIdx == -1 {
		// Client app entry was not found, we will add it at the end of the uca.xml
		log.Info("Menu entries weren't found. They will be added now!")

		command := sys.GetBinDirPath() + "/" + config.App_client_name + " encrypt %f"
		ucaId := strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		action := sys.createAction("ark", "Encrypt source", ucaId, command, "Encrypt source", "*")
		actions.Actions = append(actions.Actions, *action)

		command = sys.GetBinDirPath() + "/" + config.App_client_name + " decrypt %f"
		ucaId = strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		action = sys.createAction("ark", "Decrypt source", ucaId, command, "Decrypt source", "*")
		actions.Actions = append(actions.Actions, *action)

		command = sys.GetBinDirPath() + "/" + config.App_client_name + " addkey %f"
		ucaId = strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		action = sys.createAction("pgp-keys", "Add key", ucaId, command, "Add the key that will be used to encrypt and decrypt", "*")
		actions.Actions = append(actions.Actions, *action)

		updatedUcaBytes, err := xml.MarshalIndent(actions, " ", "	")
		if err != nil {
			log.Fatal("Error when marshaling uca.xml: " + err.Error())
		}

		if err := os.Remove(ucaDirPath); err != nil {
			log.Info("Failed to remove old uca.xml: %v", err)
		}

		ucaFile, err = os.Create(ucaDirPath)
		if err != nil {
			log.Error("Failed to create new uca.xml: %v", err)
		}

		ucaFile.Write(updatedUcaBytes)
		ucaFile.Close()
		restartUI()
	} else {
		log.Info("Menu entries have been found! Nothing to do!")
	}
}

func (sys *linux) GetInterpretor() string {
	//Find the path to the python exec
	cmd := exec.Command("which", "python")
	output, err := cmd.Output()
	log := utils.GetLogger()

	if err != nil {
		log.Fatal("Python not found on the system: ", err)
	}

	interpretor := strings.Replace(string(output), "\n", "", -1)
	interpretor = strings.Replace(interpretor, "\r", "", -1)
	return interpretor
}

func (sys *linux) GetBinDirPath() string {
	return "/usr/bin/"
}

func (sys *linux) ChangeFilePermission(keyPath *string) {
	cmd := exec.Command("chmod", "600", *keyPath)
	_, err := cmd.CombinedOutput()
	log := utils.GetLogger()

	if err != nil {
		log.Error("There was an issue during the change of the key file permmision: ", err)
	}
}

func (sys *linux) GetKeysDirPath() string {
	homePath := GetHomeDir()
	return filepath.Join(homePath, "/."+config.App_generic_name)
}

func GetOsManager() system {
	return new(linux)
}

func restartUI() {
	fileManagerProcessName := "thunar"
	err := exec.Command(fileManagerProcessName, "--quit").Run()
	log := utils.GetLogger()
	if err != nil {
		log.Error("There was an issue during the closing of the thunar process: ", err.Error())
	} else {
		startUI(&fileManagerProcessName)
	}

	log.Info(fileManagerProcessName + " was restarted succesfully!")
}

func startUI(processName *string) {
	err := exec.Command(*processName, "&").Run()
	log := utils.GetLogger()

	if err != nil {
		log.Error("There was an issue starting thunar: ", err.Error())
	}
}
