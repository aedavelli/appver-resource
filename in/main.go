package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aedavelli/appver-resource/models"
	"github.com/aedavelli/appver-resource/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.Fatal("usage: " + os.Args[0] + " <destination>")
	}

	destination := os.Args[1]

	utils.Print("creating destination dir " + destination)
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		utils.Fatal("creating destination", err)
	}

	var request models.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		utils.Fatal("reading request", err)
	}

	var inVersion = request.Version

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version:  inVersion,
		Metadata: generateMeta(request.Source, destination),
	})

	utils.PrintSuccess("Done")
}

func generateMeta(s models.Source, destination string) models.Metadata {
	res := utils.GetHttpVersion(s)
	defer res.Body.Close()

	m := utils.ParseResponse(s, res)

	metadata := models.Metadata{}

	for key, val := range m {
		metaval := val
		if s.Redact {
			metaval = "***********"
		}
		metadata = append(metadata,
			models.MetadataField{Name: key, Value: metaval})
		createFile(destination, key, val)
	}
	return metadata
}

func createFile(destination, filename, prop string) {
	output := filepath.Join(destination, filename)
	utils.Print("creating output file " + output)
	file, err := os.Create(output)
	if err != nil {
		utils.Fatal("creating output file "+output, err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%s", prop)

	err = w.Flush()

	if err != nil {
		utils.Fatal("writing output file"+output, err)
	}
}
