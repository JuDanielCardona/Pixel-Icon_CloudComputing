package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	IMAGES_DIRECTORY = "images"
	FOLDERS          = "folders"
)

type PageData struct {
	Images   []map[string]string
	Hostname string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/open_folder", openFolderHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	sliceFolder := []map[string]string{}

	folders, err := getFolders(FOLDERS)
	if err != nil {
		log.Println("Error al obtener la lista de imágenes:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, nameFolder := range folders {
		sliceFolder = append(sliceFolder, map[string]string{
			"foldername": nameFolder,
			"base64":     getImageBase64("source/image_source/carpetaimage.png"),
		})
	}

	data := PageData{
		Images:   sliceFolder,
		Hostname: getHostname(),
	}

	renderTemplate(w, "index", data)
}

func openFolderHandler(w http.ResponseWriter, r *http.Request) {
	folderName := r.URL.Query().Get("folder")
	directorio := filepath.Join("folders", folderName)

	randomImages := []map[string]string{}

	images, err := getImagesList(directorio)
	if err != nil {
		log.Println("Error al obtener la lista de imágenes:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		imageBase64 := getImageBase64(image)
		filename := filepath.Base(image)
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		randomImages = append(randomImages, map[string]string{
			"filename": filename,
			"base64":   imageBase64,
		})
	}

	data := PageData{
		Images:   randomImages,
		Hostname: getHostname(),
	}

	renderTemplate(w, "open_folder", data)
}

func renderTemplate(w http.ResponseWriter, tmplPath string, data interface{}) {
	tmpl := filepath.Join("templates", tmplPath) + ".html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println("Error al analizar la plantilla:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		log.Println("Error al ejecutar la plantilla:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getFolders(directory string) ([]string, error) {
	var folders []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != directory {
			folders = append(folders, filepath.Base(path))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func getImagesList(directory string) ([]string, error) {
	var images []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".jpg" || filepath.Ext(info.Name()) == ".png") {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func getImageBase64(imagePath string) string {
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.Printf("Error al leer la imagen %s: %s\n", imagePath, err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(imageBytes)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Error al obtener el nombre del host:", err)
		return ""
	}
	return hostname
}
