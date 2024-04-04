package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// Ruta del directorio de imágenes
const IMAGES_DIRECTORY = "sources"

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Lista para almacenar las imágenes seleccionadas al azar y convertidas a base64
	random_images := []map[string]string{}

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
		imageBytes, err := ioutil.ReadFile(images[idx])
		if err != nil {
			log.Printf("Error al leer la imagen %s: %s\n", images[idx], err)
			continue
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
		random_images = append(random_images, map[string]string{
			"filename": images[idx],
			"base64":   imageBase64,
		})
	}

	// Renderizar la plantilla y pasarle los datos de las imágenes seleccionadas al azar
	renderTemplate(w, "index", random_images)

	// Obtener y mostrar el nombre del host y la dirección IP
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Error al obtener el nombre del host:", err)
		return
	}
	fmt.Fprintf(w, "Nombre del host: %s\n", hostname)
	fmt.Fprintf(w, "Dirección IP: %s\n", r.RemoteAddr)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl + ".html")
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
