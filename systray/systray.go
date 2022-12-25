//go:build !darwin

package systray

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
