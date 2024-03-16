package main

import (
	"fmt"

	"github.com/getlantern/systray"
)

func main() {
	// Initialize systray
	systray.Run(onReady, onExit)
}

func onReady() {
	// Add your items with submenus
	item1 := systray.AddMenuItem("Item 1", "This is item 1")
	subMenu1 := item1.AddSubMenuItem("Subitem 1.1", "This is subitem 1.1")
	subMenu2 := item1.AddSubMenuItem("Subitem 1.2", "This is subitem 1.2")
	subMenu11 := subMenu1.AddSubMenuItem("Sub-subitem 1.1.1", "This is sub-subitem 1.1.1")
	subMenu1.AddSubMenuItem("Sub-subitem 1.1.2", "This is sub-subitem 1.1.2")
	subMenu2.AddSubMenuItem("Sub-subitem 1.2.1", "This is sub-subitem 1.2.1")

	item2 := systray.AddMenuItem("Item 2", "This is item 2")
	subMenu3 := item2.AddSubMenuItem("Subitem 2.1", "This is subitem 2.1")
	subMenu4 := item2.AddSubMenuItem("Subitem 2.2", "This is subitem 2.2")
	subMenu3.AddSubMenuItem("Sub-subitem 2.1.1", "This is sub-subitem 2.1.1")
	subMenu4.AddSubMenuItem("Sub-subitem 2.2.1", "This is sub-subitem 2.2.1")

	// Add a separator
	systray.AddSeparator()

	// Add a quit item
	quitMenuItem := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-quitMenuItem.ClickedCh
		fmt.Print("Quit")
		systray.Quit()
	}()

	// Add click handlers for submenu items
	go func() {
		for {
			select {
			case <-subMenu11.ClickedCh:
				fmt.Println("Subitem 1.1 clicked")
			case <-subMenu2.ClickedCh:
				fmt.Println("Subitem 1.2 clicked")
			case <-subMenu3.ClickedCh:
				fmt.Println("Subitem 2.1 clicked")
			case <-subMenu4.ClickedCh:
				fmt.Println("Subitem 2.2 clicked")
			}
		}
	}()
}

func onExit() {
	// Clean up
}
