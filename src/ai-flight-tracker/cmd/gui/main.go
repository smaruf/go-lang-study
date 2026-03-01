// Package main provides a cross-platform native GUI for the Flight Tracker AI server.
//
// Supported platforms: Windows, macOS, Linux (inc. WSL with display), Android, iOS.
//
// Build desktop binary:
//
//	go build -o flight-tracker-gui ./cmd/gui
//
// Run directly:
//
//	go run ./cmd/gui
package main

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/config"
	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/server"
)

const appID = "io.github.smaruf.flight-tracker-ai"

func main() {
	a := app.NewWithID(appID)
	a.SetIcon(theme.InfoIcon())

	w := a.NewWindow("✈ Flight Tracker AI")
	w.Resize(fyne.NewSize(860, 640))
	w.SetMaster()

	ui := newAppUI(w)
	w.SetContent(ui.build())
	w.ShowAndRun()
}

// appUI holds all UI state and widget references.
type appUI struct {
	win      fyne.Window
	cfg      *config.Config
	srv      *server.Server
	running  bool
	logEntry *widget.Entry
}

func newAppUI(win fyne.Window) *appUI {
	return &appUI{
		win: win,
		cfg: config.Load(),
	}
}

// build assembles the tabbed UI and returns the root widget.
func (ui *appUI) build() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Server", theme.ComputerIcon(), ui.buildServerTab()),
		container.NewTabItemWithIcon("Configuration", theme.SettingsIcon(), ui.buildConfigTab()),
		container.NewTabItemWithIcon("Deploy", theme.UploadIcon(), ui.buildDeployTab()),
		container.NewTabItemWithIcon("Logs", theme.DocumentIcon(), ui.buildLogsTab()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	return tabs
}

// ---- Server Tab ----

func (ui *appUI) buildServerTab() fyne.CanvasObject {
	statusLabel := widget.NewLabelWithStyle("● Stopped", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel.Importance = widget.DangerImportance

	urlLabel := widget.NewHyperlink("", nil)

	portEntry := widget.NewEntry()
	portEntry.SetText(strconv.Itoa(ui.cfg.Port))

	var startBtn, stopBtn *widget.Button

	openBrowserBtn := widget.NewButtonWithIcon("Open Web UI in Browser", theme.HomeIcon(), func() {
		openBrowser(fmt.Sprintf("http://localhost:%d", ui.cfg.Port))
	})
	openBrowserBtn.Disable()

	stopBtn = widget.NewButtonWithIcon("Stop Server", theme.MediaStopIcon(), func() {
		if ui.srv != nil {
			ui.srv.Stop() //nolint:errcheck
			ui.srv = nil
		}
		ui.running = false
		statusLabel.SetText("● Stopped")
		statusLabel.Importance = widget.DangerImportance
		urlLabel.SetText("")
		startBtn.Enable()
		stopBtn.Disable()
		openBrowserBtn.Disable()
	})
	stopBtn.Importance = widget.DangerImportance
	stopBtn.Disable()

	startBtn = widget.NewButtonWithIcon("Start Server", theme.MediaPlayIcon(), func() {
		port, err := strconv.Atoi(strings.TrimSpace(portEntry.Text))
		if err != nil || port < 1 || port > 65535 {
			popup := widget.NewModalPopUp(
				widget.NewLabel("Invalid port number (1-65535)"),
				ui.win.Canvas(),
			)
			popup.Show()
			return
		}
		ui.cfg.Port = port

		ui.srv = server.New(ui.cfg)
		if err := ui.srv.Start(); err != nil {
			statusLabel.SetText("● Error: " + err.Error())
			statusLabel.Importance = widget.DangerImportance
			return
		}
		ui.running = true
		webURL := ui.srv.WebURL()
		statusLabel.SetText(fmt.Sprintf("● Running — port %d", port))
		statusLabel.Importance = widget.SuccessImportance
		u, _ := url.Parse(webURL)
		urlLabel.SetURL(u)
		urlLabel.SetText(webURL)
		startBtn.Disable()
		stopBtn.Enable()
		openBrowserBtn.Enable()
		go ui.forwardLogs()
	})
	startBtn.Importance = widget.HighImportance

	card := widget.NewCard("Server Control", "Start or stop the Flight Tracker AI HTTP server",
		container.NewVBox(
			container.NewHBox(widget.NewLabel("Status:"), statusLabel),
			container.NewHBox(widget.NewLabel("URL:"), urlLabel),
			widget.NewSeparator(),
			widget.NewForm(widget.NewFormItem("Port", portEntry)),
			widget.NewSeparator(),
			container.NewHBox(startBtn, stopBtn),
			openBrowserBtn,
		),
	)
	return container.NewPadded(card)
}

// ---- Configuration Tab ----

func (ui *appUI) buildConfigTab() fyne.CanvasObject {
	openSkyURL := widget.NewEntry()
	openSkyURL.SetText(ui.cfg.OpenSkyURL)
	openSkyUser := widget.NewEntry()
	openSkyUser.SetText(ui.cfg.OpenSkyUsername)
	openSkyPass := widget.NewPasswordEntry()
	openSkyPass.SetText(ui.cfg.OpenSkyPassword)
	ollamaURL := widget.NewEntry()
	ollamaURL.SetText(ui.cfg.OllamaURL)
	ollamaModel := widget.NewEntry()
	ollamaModel.SetText(ui.cfg.OllamaModel)
	cacheTTL := widget.NewEntry()
	cacheTTL.SetText(strconv.Itoa(ui.cfg.CacheTTL))
	corsOrigins := widget.NewEntry()
	corsOrigins.SetText(ui.cfg.CORSOrigins)

	applyBtn := widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), func() {
		ui.cfg.OpenSkyURL = openSkyURL.Text
		ui.cfg.OpenSkyUsername = openSkyUser.Text
		ui.cfg.OpenSkyPassword = openSkyPass.Text
		ui.cfg.OllamaURL = ollamaURL.Text
		ui.cfg.OllamaModel = ollamaModel.Text
		if n, err := strconv.Atoi(cacheTTL.Text); err == nil && n > 0 {
			ui.cfg.CacheTTL = n
		}
		ui.cfg.CORSOrigins = corsOrigins.Text
	})
	applyBtn.Importance = widget.HighImportance

	form := widget.NewForm(
		widget.NewFormItem("OpenSky API URL", openSkyURL),
		widget.NewFormItem("OpenSky Username", openSkyUser),
		widget.NewFormItem("OpenSky Password", openSkyPass),
		widget.NewFormItem("Ollama URL", ollamaURL),
		widget.NewFormItem("Ollama Model", ollamaModel),
		widget.NewFormItem("Price Cache TTL (s)", cacheTTL),
		widget.NewFormItem("CORS Origins", corsOrigins),
	)

	card := widget.NewCard("Configuration",
		"Settings take effect the next time the server is started",
		container.NewVBox(form, applyBtn),
	)
	return container.NewPadded(container.NewScroll(card))
}

// ---- Deploy Tab ----

func (ui *appUI) buildDeployTab() fyne.CanvasObject {
	out := widget.NewMultiLineEntry()
	out.SetMinRowsVisible(10)
	out.Disable()
	out.Wrapping = fyne.TextTruncate

	appendOut := func(msg string) {
		out.SetText(out.Text + msg)
	}

	dockerCmd := func(args ...string) {
		appendOut(fmt.Sprintf("$ docker %s\n", strings.Join(args, " ")))
		go func() {
			b, err := exec.Command("docker", args...).CombinedOutput() //nolint:gosec
			if err != nil {
				appendOut(fmt.Sprintf("error: %v\n%s\n", err, string(b)))
			} else {
				appendOut(string(b))
			}
		}()
	}

	composeCmd := func(args ...string) {
		appendOut(fmt.Sprintf("$ docker-compose %s\n", strings.Join(args, " ")))
		go func() {
			b, err := exec.Command("docker-compose", args...).CombinedOutput() //nolint:gosec
			if err != nil {
				appendOut(fmt.Sprintf("error: %v\n%s\n", err, string(b)))
			} else {
				appendOut(string(b))
			}
		}()
	}

	buildBtn := widget.NewButtonWithIcon("Build Image", theme.StorageIcon(),
		func() { dockerCmd("build", "-t", "flight-tracker-ai", ".") })

	runBtn := widget.NewButtonWithIcon("Run Container", theme.MediaPlayIcon(), func() {
		portMap := fmt.Sprintf("%d:8080", ui.cfg.Port)
		dockerCmd("run", "-d", "--rm", "--name", "flight-tracker-ai", "-p", portMap, "flight-tracker-ai")
	})
	runBtn.Importance = widget.HighImportance

	stopBtn := widget.NewButtonWithIcon("Stop Container", theme.MediaStopIcon(),
		func() { dockerCmd("stop", "flight-tracker-ai") })
	stopBtn.Importance = widget.DangerImportance

	upBtn := widget.NewButtonWithIcon("compose up", theme.MediaPlayIcon(),
		func() { composeCmd("up", "--build", "-d") })
	upBtn.Importance = widget.HighImportance

	downBtn := widget.NewButtonWithIcon("compose down", theme.MediaStopIcon(),
		func() { composeCmd("down") })
	downBtn.Importance = widget.DangerImportance

	clearBtn := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(),
		func() { out.SetText("") })

	card := widget.NewCard("Docker Deploy", "Build and run in containers (Docker must be installed)",
		container.NewVBox(
			widget.NewLabel("Single container:"),
			container.NewGridWithColumns(3, buildBtn, runBtn, stopBtn),
			widget.NewSeparator(),
			widget.NewLabel("docker-compose (app + Ollama AI):"),
			container.NewGridWithColumns(3, upBtn, downBtn, clearBtn),
			widget.NewSeparator(),
			widget.NewLabel("Output:"),
			container.NewScroll(out),
		),
	)
	return container.NewPadded(card)
}

// ---- Logs Tab ----

func (ui *appUI) buildLogsTab() fyne.CanvasObject {
	logEntry := widget.NewMultiLineEntry()
	logEntry.SetMinRowsVisible(22)
	logEntry.Disable()
	logEntry.Wrapping = fyne.TextTruncate
	ui.logEntry = logEntry

	clearBtn := widget.NewButtonWithIcon("Clear Logs", theme.DeleteIcon(),
		func() { logEntry.SetText("") })

	card := widget.NewCard("Server Logs", "Real-time HTTP server output",
		container.NewVBox(
			clearBtn,
			container.NewScroll(logEntry),
		),
	)
	return container.NewPadded(card)
}

// forwardLogs reads from the server's log buffer and appends to the GUI log widget.
func (ui *appUI) forwardLogs() {
	if ui.srv == nil {
		return
	}
	ch := ui.srv.LogBuf.Ch()
	for line := range ch {
		if !ui.running {
			return
		}
		if ui.logEntry != nil {
			ui.logEntry.SetText(ui.logEntry.Text + line)
		}
	}
}

// openBrowser opens a URL in the system's default browser.
func openBrowser(rawURL string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	cmd.Start() //nolint:errcheck
}
