package update

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/minio/selfupdate"
	"github.com/projectdiscovery/gologger"
)

// References:
// https://github.com/minio/selfupdate
// https://github.com/Ciyfly/Argo/blob/main/pkg/updateself/update.go
var (
	ZIP                    = "zip"
	PROJECT_NAME           = "CLITemplate"
	REPOSITORY_RELEASE_API = "https://api.github.com/repos/N0el4kLs/CLITemplate/releases/latest"
)

type LastVersionInfo struct {
	TagName string   `json:"tag_name"`
	Assets  []Assets `json:"assets"`
	Body    string   `json:"body"`
}
type Assets struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func Update(crtVersion string) {
	// 1. get remote version
	remoteLatest, err := remoteVersion()
	if err != nil {
		gologger.Fatal().Msgf("get remote version error: %s\n", err)
	}

	// 2. compare version
	if crtVersion < remoteLatest.TagName {
		gologger.Info().
			Label("Update").
			Msgf("Found new version, try to update to latest version %s ...\n", remoteLatest.TagName)

		// 3. download latest version
		if err = remoteLatest.download(); err != nil {
			gologger.Fatal().Msgf("Update to latest version error: %s\n", err)
		}
	} else {
		gologger.Info().
			Label("Update").
			Msgf("Current version %s is the latest version \n", crtVersion)
	}
}

func remoteVersion() (LastVersionInfo, error) {
	var lvi = LastVersionInfo{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", REPOSITORY_RELEASE_API, nil)
	if err != nil {
		return lvi, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return lvi, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	jsonErr := json.Unmarshal(body, &lvi)
	if jsonErr != nil {
		gologger.Error().Msgf("parser github api json error: %s\n", err)
		return lvi, err
	}
	return lvi, nil
}

func (l LastVersionInfo) download() error {
	// find corresponding achievement
	var (
		downloadUrl string
		name        string
	)
	currentOsArch := getOsArch()
	for _, assest := range l.Assets {
		if strings.Contains(assest.Name, currentOsArch) {
			downloadUrl = assest.BrowserDownloadURL
			name = assest.Name
			break
		}
	}

	// download binary
	gologger.Debug().Msgf("downloadUrl: %s\n", downloadUrl)
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	downloadedFileName := name
	downloadedFile, err := os.Create(downloadedFileName)
	if err != nil {
		return err
	}
	if _, err = io.Copy(downloadedFile, resp.Body); err != nil {
		return err
	}
	downloadedFile.Close()

	// unzip achievement
	if err = unzipFile(downloadedFileName); err != nil {
		return err
	}
	bin, err := os.Open(PROJECT_NAME)
	if err != nil {
		return err
	}
	defer bin.Close()

	if err = selfupdate.Apply(bin, selfupdate.Options{}); err != nil {
		if err = selfupdate.RollbackError(err); err != nil {
			gologger.Fatal().
				Label("updater").
				Msgf("rollback of update failed got %v ,pls reinstall \n", err)
		}
		return err
	}

	// clean cache
	if err = cleanCache(downloadedFileName); err != nil {
		return err
	}

	return nil
}

func getOsArch() string {
	var os, arch string
	if runtime.GOOS == "darwin" {
		os = "macOS"
	} else {
		os = runtime.GOOS
	}
	if runtime.GOARCH == "arm" {
		arch = "arm64"
	} else {
		arch = runtime.GOARCH
	}
	return os + "_" + arch
}

func unzipFile(name string) error {
	if strings.ToLower(filepath.Ext(name)[1:]) == ZIP {
		zipReader, err := zip.OpenReader(name)
		if err != nil {
			return err
		}

		for _, tmpF := range zipReader.File {
			filePath := tmpF.Name
			if filePath != "CLITemplate" {
				continue
			}
			zippedFile, err := tmpF.Open()
			if err != nil {
				return err
			}
			defer zippedFile.Close()

			if tmpF.FileInfo().IsDir() {
				if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
					return fmt.Errorf("UncompressZip makedir error: %w", err)
				}
				continue
			}

			targetFile, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("UncompressZip create downloadedFile error: %w", err)
			}
			defer targetFile.Close()

			srcFile, err := tmpF.Open()
			if err != nil {
				return fmt.Errorf("UncompressZip open downloadedFile error: %w", err)
			}
			defer srcFile.Close()

			if _, err := io.Copy(targetFile, srcFile); err != nil {
				return fmt.Errorf("UncompressZip copy downloadedFile error: %w", err)
			}

			gologger.Debug().Msgf("UncompressZip success: %s \n", filePath)
		}
	}
	return nil
}

func cleanCache(n string) error {
	neededClean := []string{
		n,
		PROJECT_NAME,
	}
	for _, f := range neededClean {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
