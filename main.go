package main

import (
	"io/ioutil"

	"os/exec" // for current user
	"strings"

	"github.com/getlantern/systray"
	"github.com/musabiii/parse_1c_v8"
)

func main() {
	// Initialize systray
	systray.Run(onReady, onExit)
}

func onReady() {

	// systray.SetIcon(icon.Data)	// Load your custom icon
	iconBytes, err := ioutil.ReadFile("./icons/1c.ico")
	if err != nil {
		return
	}

	// Set the icon
	systray.SetIcon(iconBytes)

	connections := parse_1c_v8.GetConnections()
	foldersMap := parse_1c_v8.GetFoldersMap(connections)

	itemsMap := make(map[string]*systray.MenuItem)

	fillItemsMap(itemsMap, foldersMap)

	// Add a separator
	systray.AddSeparator()

	// Add a quit item
	quitMenuItem := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-quitMenuItem.ClickedCh
		systray.Quit()
	}()

	handleSubMenuClicks(itemsMap, connections)

}

func onExit() {
	// Clean up
}

func fillItemsMap(itemsMap map[string]*systray.MenuItem, foldersMap map[string][]parse_1c_v8.Connection) {

	for k, v := range foldersMap {
		item1 := systray.AddMenuItem(k, "This is item 1")
		for _, vv := range v {
			item11 := item1.AddSubMenuItem(vv.Name, vv.Connect)
			itemsMap[vv.Name] = item11 // Add to ma

		}

	}
}

func handleSubMenuClicks(itemsMap map[string]*systray.MenuItem, connections []parse_1c_v8.Connection) {
	for title, menuItem := range itemsMap {
		go func(title string, menuItem *systray.MenuItem) {
			for range menuItem.ClickedCh {
				for _, conn := range connections {
					if conn.Name == title {
						runBase(conn.Connect)
					}
				}
			}
		}(title, menuItem)
	}
}

func runBase(connect string) {

	starter1c := "C:\\Program Files\\1cv8\\common\\1cestart.exe"
	basePath, runType := parseConnect(connect)
	cmd := exec.Command(starter1c, "ENTERPRISE", runType, basePath)
	cmd.Output()

}

func parseConnect(connect string) (string, string) {

	parts := strings.SplitN(connect, "=", 2)
	baseType := parts[0]
	var basePath string
	var runType string
	if baseType == "File" {
		basePath = parts[1][1 : len(parts[1])-2]
		runType = "/F" // File
	} else {
		servParts := strings.SplitN(connect, ";", 2)
		serv := servParts[0][6 : len(servParts[0])-1]
		ref := servParts[1][5 : len(servParts[1])-2]
		basePath = serv + "\\" + ref
		runType = "/S" // Server
	}
	return basePath, runType

}
