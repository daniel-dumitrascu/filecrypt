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

func (sys *linux) GetHomeDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	if len(homePath) < 1 {
		log.Fatalln("The user home directory path is not valid: " + string(homePath))
	}
	return homePath
}

func (sys *linux) GetContextAppPath() string {
	homePath := sys.GetHomeDir()
	return homePath + "/bin/context-app"
}

func (sys *linux) SpecificSetup() {
	//Setup patch to uca.xml
	homePath := sys.GetHomeDir()
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

	foundIdx := slices.IndexFunc(actions.Actions, func(action Action) bool { return strings.Contains(action.Command, "context-app") })
	if foundIdx == -1 {
		//context-app entry was not found
		//we will add it at the end of the uca.xml
		ucaId := strconv.FormatInt(time.Now().UnixMicro(), 10) + "-" + strconv.Itoa(rand.Intn(5)+1)
		command := sys.GetContextAppPath() + " %f"
		action := Action{Icon: "ark", Name: "Encrypt/Decrypt source",
			Uniqueid: ucaId, Command: command,
			Description: "Encrypt/Decrypt source", Patterns: "*"}
		action.XMLName.Local = "action"

		actions.Actions = append(actions.Actions, action)
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

func GetOsManager() system {
	return new(linux)
}
