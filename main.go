package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ujum/ftran/pkg/data"
	"github.com/ujum/ftran/pkg/ftran"
	"path/filepath"
)

var (
	targetDir = "result"
	oneDir    = true
	sourceDir = ""
	resultLog = &bytes.Buffer{}
)

func main() {
	a := app.NewWithID("Ftran")
	mainWindow := a.NewWindow("Ftran")
	mainWindow.CenterOnScreen()
	mainWindow.SetMaster()

	progressBarInfinite := widget.NewProgressBarInfinite()
	progressBarInfinite.Hide()
	progressBarInfinite.Stop()

	reportButton := widget.NewButton("Download report", func() {
		if resultLog.Len() == 0 {
			dialog.ShowInformation("Info", "Log is empty", mainWindow)
			return
		}
		dialog.ShowFileSave(func(closer fyne.URIWriteCloser, err error) {
			if closer != nil {
				defer closer.Close()
				_, err := closer.Write(resultLog.Bytes())
				if err != nil {
					dialog.ShowError(err, mainWindow)
				}
			}
		}, mainWindow)
	})
	reportButton.Hide()
	reportButton.Disable()

	doneLabel := widget.NewLabel("Done!")
	doneLabel.Hide()

	runButton := createRunButton(mainWindow, progressBarInfinite, doneLabel, reportButton)
	folderLabel := widget.NewLabel("Source folder: ")
	oneDirCheck := widget.NewCheck("Transfer to the same extension directory", func(res bool) {
		oneDir = res
	})
	oneDirCheck.Checked = true

	mainWindow.SetContent(container.NewScroll(container.NewVBox(
		oneDirCheck,
		folderLabel,
		widget.NewButton("Select folder", func() {
			folderOpen := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if uri != nil {
					sourceDir = uri.Path()
					runButton.Enable()
					folderLabel.SetText("Source folder: " + sourceDir)
				}
			}, mainWindow)
			folderOpen.Resize(fyne.NewSize(600, 400))
			folderOpen.Show()
		}),
		widget.NewSeparator(),
		runButton,
		doneLabel,
		reportButton,
		progressBarInfinite,
	)))
	mainWindow.Resize(fyne.NewSize(600, 400))
	mainWindow.ShowAndRun()
}

func createRunButton(mainWindow fyne.Window, progressBarInfinite *widget.ProgressBarInfinite, doneLabel *widget.Label, reportButton *widget.Button) *widget.Button {
	runButton := widget.NewButton("Run", func() {
		if sourceDir == "" {
			dialog.ShowInformation("Warning!", "Please, choose the source folder", mainWindow)
			return
		}
		dialog.ShowConfirm("Are you sure?", "Do you want to transfer files from folder: "+sourceDir+" ?", func(confirm bool) {
			if confirm {
				progressBarInfinite.Show()
				progressBarInfinite.Start()
				doneLabel.Hide()
				go func() {
					resultLog.Reset()
					fullTargetDir := filepath.Join(filepath.Dir(sourceDir), targetDir)
					resourceLogs := make(chan *data.ResourceLog)
					go func() {
						err := ftran.Run(&ftran.Options{
							SourceDir:  sourceDir,
							TargetDir:  fullTargetDir,
							SameExtDir: oneDir,
						}, resourceLogs)
						if err != nil {
							dialog.ShowError(err, mainWindow)
							return
						}
						resultFolderContent := "<a href=\"" + fullTargetDir + "\">Open result folder</a>"
						fyne.CurrentApp().SendNotification(fyne.NewNotification("Done", resultFolderContent))
					}()
					for res := range resourceLogs {
						if res.Error == nil {
							resultLog.WriteString(fmt.Sprintf("%s --> %s\n", res.Source, res.Target))
						} else {
							resultLog.WriteString(fmt.Sprintf("%v, %s --> %s\n", res.Error, res.Source, res.Target))
						}
					}
					reportButton.Enable()
					reportButton.Show()
					progressBarInfinite.Stop()
					progressBarInfinite.Hide()
					doneLabel.Show()
				}()
			}
		}, mainWindow)
	})
	runButton.Disable()
	return runButton
}
