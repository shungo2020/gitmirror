package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func ExecutableDir() string {
	exePath, _ := os.Executable()
	exePath, _ = filepath.Abs(exePath)
	return path.Dir(exePath)
}

type GitInfo struct {
	Name string
	Url  string
}

func pathExist(path string) bool {
	log.Debugf("check path: %s", path)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

func checkGitCloned(storePath string, gitInfo *GitInfo) bool {
	repoPath := filepath.Join(storePath, gitInfo.Name)
	if !pathExist(repoPath) {
		return false
	}

	gitPath := filepath.Join(repoPath, ".git")
	if !pathExist(gitPath) {
		return false
	}
	return true
}

func doClone(exePath string, storePath string, gitInfo *GitInfo) {
	log.Info("start clone git: %s", gitInfo.Url)
	cmdPath := filepath.Join(exePath, "clone.bat")
	repoPath := filepath.Join(storePath, gitInfo.Name)
	os.RemoveAll(repoPath)
	cmd := exec.Command(cmdPath, storePath, gitInfo.Name, gitInfo.Url)
	log.Debugf("run cmd: %s", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Errorf("command error %s", err.Error())
		return
	}
	log.Info("git: %s cloned", gitInfo.Url)
}

func doPull(exePath string, storePath string, gitInfo *GitInfo) {
	log.Info("start pull git: %s", gitInfo.Url)
	cmdPath := filepath.Join(exePath, "pull.bat")
	repoPath := filepath.Join(storePath, gitInfo.Name)
	os.RemoveAll(repoPath)
	cmd := exec.Command(cmdPath, storePath, gitInfo.Name, gitInfo.Url)
	log.Debugf("run cmd: %s", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Errorf("command error %s", err.Error())
		return
	}
	log.Info("git: %s pull", gitInfo.Url)
}

func main() {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{FullTimestamp: true, DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"}
	log.SetFormatter(formatter)
	log.SetReportCaller(true)

	exePath := ExecutableDir()
	storePath := filepath.Join(exePath, "data")
	os.Mkdir(storePath, 0666)
	var gitInfos []GitInfo
	confPath := filepath.Join(exePath, "git.json")
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Warnf("read conf error: %v", err)
		return
	}
	err = json.Unmarshal(bytes, &gitInfos)
	if err != nil {
		log.Warnf("json conf error: %v", err)
		return
	}

	for _, gitInfo := range gitInfos {
		log.Debugf("mirror: %s %s", gitInfo.Name, gitInfo.Url)
		if checkGitCloned(storePath, &gitInfo) {
			doPull(exePath, storePath, &gitInfo)
		} else {
			doClone(exePath, storePath, &gitInfo)
		}
	}

}
