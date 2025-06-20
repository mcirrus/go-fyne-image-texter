package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	dataURL   = "https://example.com/data.json"
	uploadURL = "https://example.com/upload"
)

type TextItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BuildItem struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Text        []TextItem `json:"text"`
}

func fetchBuilds(url string) ([]BuildItem, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var builds []BuildItem
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return nil, err
	}
	return builds, nil
}

func loadImageFromURL(url string) (fyne.CanvasObject, error) {
	resp, err := http.Get(url)
	if err != nil {
		return widget.NewLabel(""), err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return widget.NewLabel(""), err
	}
	cimg := canvas.NewImageFromImage(img)
	cimg.SetMinSize(fyne.NewSize(80, 80))
	cimg.FillMode = canvas.ImageFillContain
	return cimg, nil
}

func startRecording(file string) (*exec.Cmd, error) {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"-y", "-f", "avfoundation", "-i", ":0", file}
	case "windows":
		args = []string{"-y", "-f", "dshow", "-i", "audio=Microphone", file}
	default:
		args = []string{"-y", "-f", "alsa", "-i", "default", file}
	}
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, cmd.Start()
}

func stopRecording(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Signal(os.Interrupt)
}

func uploadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("open file:", err)
		return
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		log.Println("writer:", err)
		return
	}
	if _, err := io.Copy(part, file); err != nil {
		log.Println("copy:", err)
		return
	}
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, uploadURL, body)
	if err != nil {
		log.Println("request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("upload:", err)
		return
	}
	resp.Body.Close()
	log.Println("uploaded", path)
}

func createRow(item BuildItem) fyne.CanvasObject {
	id := widget.NewLabel(item.ID)
	imageObj, err := loadImageFromURL(item.URL)
	if err != nil {
		log.Println("image:", err)
	}
	title := widget.NewLabelWithStyle(item.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel(item.Description)
	textCol := container.NewVBox(title, desc)
	recordBtn := widget.NewButton("Record", nil)
	var cmd *exec.Cmd
	var file string
	recordBtn.OnTapped = func() {
		if cmd == nil {
			file = fmt.Sprintf("record_%s.wav", item.ID)
			var err error
			cmd, err = startRecording(file)
			if err != nil {
				log.Println("record start:", err)
				cmd = nil
				return
			}
			recordBtn.SetText("Stop")
		} else {
			if err := stopRecording(cmd); err != nil {
				log.Println("record stop:", err)
			}
			cmd = nil
			recordBtn.SetText("Record")
			go uploadFile(file)
		}
	}
	return container.NewGridWithColumns(4, id, imageObj, textCol, recordBtn)
}

func main() {
	a := app.New()
	w := a.NewWindow("Image Texter")

	items, err := fetchBuilds(dataURL)
	if err != nil {
		log.Fatal(err)
	}
	rows := make([]fyne.CanvasObject, 0, len(items))
	for _, item := range items {
		rows = append(rows, createRow(item))
	}
	content := container.NewVBox(rows...)
	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
