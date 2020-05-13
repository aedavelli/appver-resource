package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/aedavelli/appver-resource/models"
	"github.com/aedavelli/appver-resource/utils"
)

func main() {
	var request models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		utils.Fatal("parse error:", err.Error())
	}

	json.NewEncoder(os.Stdout).Encode(generateVersions(request))
	utils.PrintSuccess("Done")
}

func generateVersions(r models.CheckRequest) models.CheckResponse {
	res := utils.GetHttpVersion(r.Source)
	defer res.Body.Close()

	m := utils.ParseResponse(r.Source, res)
	if len(m) == 0 {
		utils.Fatal("Unable to get version from given data")
	}

	scopes := strings.Split(r.Source.VersionField, "::")
	version := m[scopes[len(scopes)-1]]

	if r.Version.Version == "" || r.Version.Version == version {
		return models.CheckResponse{models.AppVersion{version}}
	}

	return models.CheckResponse{r.Version, models.AppVersion{version}}
}
