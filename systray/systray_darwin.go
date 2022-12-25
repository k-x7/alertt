//go:build darwin

package systray

/*
#cgo darwin CFLAGS: -I $GOPATH/src/github.com/fyne-io/systray
#cgo darwin CFLAGS: -DDARWIN -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdbool.h>
#include "systray.h"
void setInternalLoop(bool);
*/
import "C"

import "fyne.io/systray"

func Run(onReady, onExit func()) {
	systray.Run(onReady, onExit)
}

func SetupTray(iconTray []byte, title, desc string) {
	systray.SetIcon(iconTray)
	systray.SetTitle(title)
	systray.SetTooltip(desc)
	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func Quit() {
	systray.Quit()
}
