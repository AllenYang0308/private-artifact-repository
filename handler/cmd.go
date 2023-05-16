package handler

import (
	"artifact_repository/wss"
	"encoding/json"
	"fmt"
	"os"
)

func GetPackageReport(packageName string, projectName string, withConf string) {
	wss.DoWhitesourceScan(packageName, projectName, withConf)
	wss.DoUploadRequest(projectName)

	ch := wss.GenerateProjectReportAsync(projectName)
	_ = wss.GetProcessStatus(ch, projectName)

	reportPath := fmt.Sprintf("report/%s", projectName)
	os.Mkdir(reportPath, 0755)
	rsp := wss.GetProjectRiskReport(projectName)
	_, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}
}

func GetProjectAlert(projectName string) {
	rsp := wss.GetProjectRiskAlert(projectName)
	rsp, _ = wss.GetPrettyString(rsp)

	reportPath := fmt.Sprintf("report/%s", projectName)
	reportFile := fmt.Sprintf(reportPath + "/alert.json")
	os.Mkdir(reportPath, 0755)
	err := os.WriteFile(reportFile, []byte(rsp), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Print(rsp)
}
