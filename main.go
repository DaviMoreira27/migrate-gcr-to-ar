package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	images := getAllGcrDevImages()
	fmt.Println("Lista de imagens:")
	for _, img := range images {
		if (strings.Contains(img, "staging")) {
			imgTag, errorTag := getImageTag(img)
			if (errorTag != nil) {
				panic(errorTag)
			}
			transferFromGcrToArtifactRegistry(img, imgTag)
		}
	}
}

func getAllGcrDevImages() []string {
	cmd := exec.Command("gcloud", "container", "images", "list")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erro ao executar comando:", err)
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	return lines
}

func getImageTag(image string) (string, error) {
	getLastTag := exec.Command("gcloud", "container", "images", "list-tags", image, "--sort-by=~TIMESTAMP", "--limit=1", "--format=value(TAGS)")

	output, err := getLastTag.CombinedOutput()
	if err != nil {
		fmt.Println("Erro ao executar comando:", err.Error())
		fmt.Println("Saída do comando:", string(output))
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func transferFromGcrToArtifactRegistry(image string, tag string) {
	imagePlusTag := fmt.Sprintf("%s:%s", image, tag)
	currentDate := time.Now().Unix()

	targetImage := fmt.Sprintf("us-docker.pkg.dev/PROJECT/gcr.io/%s:%d", getImageName(image), currentDate)

	fmt.Println(fmt.Sprintln("gcloud", "container", "images", "add-tag", imagePlusTag, targetImage))

	transferImage := exec.Command("gcloud", "container", "images", "add-tag", imagePlusTag, targetImage)

	output, err := transferImage.CombinedOutput()
	if err != nil {
		fmt.Println("Erro ao transferir imagem:", err)
		fmt.Println("Saída do comando:", string(output))
		return
	}

	fmt.Println("Transferência concluída com sucesso! \n \n")
}

func getImageName(imageURL string) string {
	parts := strings.Split(imageURL, "/")
	return parts[len(parts)-1]
}
