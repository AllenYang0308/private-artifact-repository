package handler

import (
	"artifact_repository/wss"
	"encoding/json"
	"fmt"
	"os"
)

func GetPackageReport(packageName string, projectName string, withConf string) {
	fmt.Println(withConf)
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
