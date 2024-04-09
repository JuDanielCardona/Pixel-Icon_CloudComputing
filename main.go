package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Ruta del directorio de imágenes
const IMAGES_DIRECTORY = "images"
const FOLDERS = "folders"

type PageData struct {
	Images   []map[string]string
	Hostname string
}

func main() {

	// Manejador para archivos estáticos en el directorio 'static'
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandlerOficial)
	http.HandleFunc("/open_folder", openFolderHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Lista para almacenar las imágenes seleccionadas al azar y convertidas a base64
	randomImages := []map[string]string{}

	// Obtener lista de imágenes del directorio
	images, err := getImagesList(IMAGES_DIRECTORY)
	if err != nil {
		log.Println("Error al obtener la lista de imágenes:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Semilla para la generación de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Seleccionar 4 imágenes al azar
	for _, idx := range rand.Perm(len(images))[:4] {
		imageBytes, err := os.ReadFile(images[idx])
		if err != nil {
			log.Printf("Error al leer la imagen %s: %s\n", images[idx], err)
			continue
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		// Obtener solo el nombre del archivo sin la extensión ni la ruta
		filename := filepath.Base(images[idx])
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		randomImages = append(randomImages, map[string]string{
			"filename": filename,
			"base64":   imageBase64,
		})
	}

	// Obtener el nombre del host
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Error al obtener el nombre del host:", err)
		return
	}

	// Datos de la página a pasar a la plantilla
	data := PageData{
		Images:   randomImages,
		Hostname: hostname,
	}

	// Renderizar la plantilla y pasarle los datos de las imágenes seleccionadas al azar y el nombre del host
	renderTemplate(w, "index", data)
}

func indexHandlerOficial(w http.ResponseWriter, r *http.Request) {
	// Datos de la página a pasar a la plantilla
	sliceFolder := []map[string]string{}

	folders, err := getFolders(FOLDERS)
	if err != nil {
		log.Println("Error al obtener la lista de imágenes:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	imageBytes, err := os.ReadFile("source/image_source/carpetaimage.png")
	if err != nil {
		log.Printf("Error al leer la imagen %s: %s\n", "capetaimage.png", err)

	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	for _, nameFolder := range folders {
		println(nameFolder + "  ")
		sliceFolder = append(sliceFolder, map[string]string{
			"foldername": nameFolder,
			"base64":     imageBase64,
		})
	}

	// Obtener el nombre del host
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Error al obtener el nombre del host:", err)
		return
	}

	data := PageData{
		Images:   sliceFolder,
		Hostname: hostname,
	}

	// Renderizar la plantilla y pasarle los datos de las imágenes seleccionadas al azar y el nombre del host
	renderTemplate(w, "indexO", data)
}

func openFolderHandler(w http.ResponseWriter, r *http.Request) {
	// Lista para almacenar las imágenes seleccionadas al azar y convertidas a base64

	folderName := r.URL.Query().Get("folder")
	directorio := filepath.Join("folders", folderName)

	print(directorio)

	randomImages := []map[string]string{}

	// Obtener lista de imágenes del directorio
	images, err := getImagesList(directorio)
	if err != nil {
		log.Println("Error al obtener la lista de imágenes:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Semilla para la generación de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Seleccionar 4 imágenes al azar
	for _, image := range images {
		imageBytes, err := os.ReadFile(image)
		if err != nil {
			log.Printf("Error al leer la imagen %s: %s\n", image, err)
			continue
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		// Obtener solo el nombre del archivo sin la extensión ni la ruta
		filename := filepath.Base(image)
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		randomImages = append(randomImages, map[string]string{
			"filename": filename,
			"base64":   imageBase64,
		})
	}

	// Obtener el nombre del host
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Error al obtener el nombre del host:", err)
		return
	}

	// Datos de la página a pasar a la plantilla
	data := PageData{
		Images:   randomImages,
		Hostname: hostname,
	}

	// Renderizar la plantilla y pasarle los datos de las imágenes seleccionadas al azar y el nombre del host
	renderTemplate(w, "open_folder", data)
}

func renderTemplate(w http.ResponseWriter, tmplPath string, data interface{}) {
	// Construye la ruta completa al archivo de la plantilla
	tmpl := filepath.Join("templates", tmplPath) + ".html"

	// Analiza el archivo de la plantilla
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println("Error al analizar la plantilla:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ejecuta la plantilla con los datos proporcionados
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
