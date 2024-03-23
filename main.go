package main

import (
	"fmt"
	"io/ioutil"

	"bufio"
	"os"
	"os/exec"
	"os/user" // for current user
	"strings"

	"github.com/getlantern/systray"
)

type Connection struct {
	Connect                  string            `json:"connect"`
	Name                     string            `json:"name"`
	ID                       string            `json:"id"`
	OrderInList              int               `json:"order_in_list"`
	Folder                   string            `json:"folder"`
	OrderInTree              int               `json:"order_in_tree"`
	External                 int               `json:"external"`
	ClientConnectionSpeed    string            `json:"client_connection_speed,omitempty"`
	App                      string            `json:"app,omitempty"`
	WA                       int               `json:"wa,omitempty"`
	Version                  string            `json:"version,omitempty"`
	DisableLocalSpeechToText int               `json:"disable_local_speech_to_text,omitempty"`
	DefaultVersion           string            `json:"default_version,omitempty"`
	DefaultApp               string            `json:"default_app,omitempty"`
	menuitem                 *systray.MenuItem `json:""`
}

func main() {
	// Initialize systray
	systray.Run(onReady, onExit)
}

func onReady() {

	// systray.SetIcon(icon.Data)	// Load your custom icon
	iconBytes, err := ioutil.ReadFile("./icons/1c.ico")
	if err != nil {
		fmt.Println("Failed to load icon:", err)
		return
	}

	// Set the icon
	systray.SetIcon(iconBytes)

	foldersMap, connections := getFoldersMap() // Get folders map

	itemsMap := make(map[string]*systray.MenuItem)

	menuItems := []*systray.MenuItem{}

	for k, v := range foldersMap {
		fmt.Println(k)
		item1 := systray.AddMenuItem(k, "This is item 1")
		for _, vv := range v {
			// fmt.Println(vv.Connect)
			fmt.Printf("vv: %v\n", vv.Name)
			item11 := item1.AddSubMenuItem(vv.Name, vv.Connect)
			menuItems = append(menuItems, item11)
			itemsMap[vv.Name] = item11 // Add to ma
			vv.menuitem = item11       // Add to vv

		}

	}

	// Add a separator
	systray.AddSeparator()

	// Add a quit item
	quitMenuItem := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-quitMenuItem.ClickedCh
		fmt.Print("Quit")
		systray.Quit()
	}()

	// Handle clicks for submenu items
	// handleSubMenuClicks([](*systray.MenuItem){subMenu1, subMenu2, subMenu3, subMenu4})
	handleSubMenuClicks2(itemsMap, connections)

}

func onExit() {
	// Clean up
}

func handleSubMenuClicks(subMenus []*systray.MenuItem) {
	for _, subMenu := range subMenus {
		var sm systray.MenuItem = *subMenu
		go func(subMenu *systray.MenuItem, Title string) {
			for {
				select {
				case <-subMenu.ClickedCh:
					// sm := *subMenu["tooltip"]
					// tt := (*subMenu).tooltip
					// tt := sm.tooltip\
					// var sm systray.MenuItem

					fmt.Printf("Submenu item %s clicked\n", subMenu, sm)
				}
			}
		}(subMenu, "")
	}
}

func handleSubMenuClicks2(itemsMap map[string]*systray.MenuItem, connections []Connection) {

	for k, v := range itemsMap {
		go func(subMenu *systray.MenuItem, Title string) {
			for {
				select {
				case <-subMenu.ClickedCh:
					// find item in connections with Name = Title
					for _, vv := range connections {
						if vv.Name == Title {
							starter1c := "C:\\Program Files\\1cv8\\common\\1cestart.exe"
							parts := strings.SplitN(vv.Connect, "=", 2)
							baseType := parts[0]

							var basePath string
							var runType string

							if baseType == "File" {

								basePath = parts[1][1 : len(parts[1])-2]
								runType = "/F" // File

							} else {
								servParts := strings.SplitN(vv.Connect, ";", 2)
								serv := servParts[0][6 : len(servParts[0])-1]
								ref := servParts[1][5 : len(servParts[1])-2]
								basePath = serv + "\\" + ref
								runType = "/S" // Server
							}

							cmd := exec.Command(starter1c, "ENTERPRISE", runType, basePath)
							cmd.Output()
						}
					}
					// subMenu.SetTitle(vv.Connect)
					fmt.Printf("Submenu item %s clicked\n", Title)
				}
			}
		}(v, k)

	}
}

func getFoldersMap() (map[string][]Connection, []Connection) {

	foldersMap := make(map[string][]Connection)
	var currentConnection Connection
	var connections []Connection

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return foldersMap, connections
	}

	// Open the file
	file, err := os.Open(currentUser.HomeDir + "\\AppData\\Roaming\\1C\\1CEStart\\ibases.v8i")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return foldersMap, connections
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a map to store unique folders

	// Iterate over each line
	for scanner.Scan() {
		line := scanner.Text()

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentConnection.ID != "" {
				connections = append(connections, currentConnection)
				foldersMap[currentConnection.Folder] = append(foldersMap[currentConnection.Folder], currentConnection)
			}
			currentConnection = Connection{}
			line = "Name=" + line[1:len(line)-1]
			// continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Connect":
			currentConnection.Connect = value
		case "ID":
			currentConnection.ID = value
		case "Name":
			currentConnection.Name = value
		case "OrderInList":
			fmt.Sscanf(value, "%d", &currentConnection.OrderInList)
		case "Folder":
			currentConnection.Folder = value
		case "OrderInTree":
			fmt.Sscanf(value, "%d", &currentConnection.OrderInTree)
		case "External":
			fmt.Sscanf(value, "%d", &currentConnection.External)
		case "ClientConnectionSpeed":
			currentConnection.ClientConnectionSpeed = value
		case "App":
			currentConnection.App = value
		case "WA":
			fmt.Sscanf(value, "%d", &currentConnection.WA)
		case "Version":
			currentConnection.Version = value
		case "DisableLocalSpeechToText":
			fmt.Sscanf(value, "%d", &currentConnection.DisableLocalSpeechToText)
		case "DefaultVersion":
			currentConnection.DefaultVersion = value
		case "DefaultApp":
			currentConnection.DefaultApp = value
		}
	}

	if currentConnection.ID != "" {
		connections = append(connections, currentConnection)
		foldersMap[currentConnection.Folder] = append(foldersMap[currentConnection.Folder], currentConnection)
	}

	return foldersMap, connections

}
