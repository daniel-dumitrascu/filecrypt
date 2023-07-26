//go:build linux

package env

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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
	//Setup patch to uca.xml
	homePath := GetHomeDir()
	ucaDirPath := homePath + "/.config/Thunar/uca.xml"

	ucaFile, err := os.Open(ucaDirPath)
	if err != nil {
		ucaFile.Close()
		log.Fatalln("Error when trying to open uca.xml: " + err.Error())
	}

	ucaByteValue, _ := ioutil.ReadAll(ucaFile)
	var actions Actions
	err = xml.Unmarshal([]byte(ucaByteValue), &actions)
	if err != nil {
		ucaFile.Close()
		log.Fatalln("Error when unmarshaling uca.xml: " + err.Error())
	}

	ucaFile.Close()

	foundIdx := slices.IndexFunc(actions.Actions, func(action Action) bool { return strings.Contains(action.Command, app_client_name) })
	if foundIdx == -1 {
		// Client app entry was not found, we will add it at the end of the uca.xml

		command := sys.GetBinDirPath() + app_client_name + " encrypt %f"
		ucaId := strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		action := sys.createAction("ark", "Encrypt source", ucaId, command, "Encrypt source", "*")
		actions.Actions = append(actions.Actions, *action)

		command = sys.GetBinDirPath() + app_client_name + " decrypt %f"
		ucaId = strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		action = sys.createAction("ark", "Decrypt source", ucaId, command, "Decrypt source", "*")
		actions.Actions = append(actions.Actions, *action)

		command = sys.GetBinDirPath() + app_client_name + " addkey %f"
		ucaId = strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		action = sys.createAction("pgp-keys", "Add key", ucaId, command, "Add the key that will be used to encrypt and decrypt", "*")
		actions.Actions = append(actions.Actions, *action)

		updatedUcaBytes, err := xml.MarshalIndent(actions, " ", "	")
		if err != nil {
			log.Fatalln("Error when marshaling uca.xml: " + err.Error())
		}

		if err := os.Remove(ucaDirPath); err != nil {
			log.Printf("Failed to remove old uca.xml: %v", err)
		}

		ucaFile, err = os.Create(ucaDirPath)
		if err != nil {
			log.Printf("Failed to create new uca.xml: %v", err)
		}

		ucaFile.Write(updatedUcaBytes)
		ucaFile.Close()
	}
}

func (sys *linux) GetInterpretor() string {
	//Find the path to the python exec
	cmd := exec.Command("which", "python")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Python not found on the system", err)
	}

	interpretor := strings.Replace(string(output), "\n", "", -1)
	interpretor = strings.Replace(interpretor, "\r", "", -1)
	return interpretor
}

func GetOsManager() system {
	return new(linux)
}
