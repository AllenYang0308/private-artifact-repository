package handler

import (
	"artifact_repository/wss"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/signintech/pdft"
	gopdf "github.com/signintech/pdft/minigopdf"
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

	var projectScanInfo wss.ProjectScanInfo

	_ = json.Unmarshal([]byte(rsp), &projectScanInfo)

	reportPath := fmt.Sprintf("report/%s", projectName)
	reportFile := fmt.Sprintf(reportPath + "/alert.json")
	os.Mkdir(reportPath, 0755)
	err := os.WriteFile(reportFile, []byte(rsp), 0644)
	if err != nil {
		panic(err)
	}

	// fmt.Print(rsp)
}

func UpdateRiskReport(projectName string) {

	var ipdf pdft.PDFt

	rsp := wss.GetProjectRiskAlert(projectName)
	rsp, _ = wss.GetPrettyString(rsp)

	var projectScanInfo wss.ProjectScanInfo
	_ = json.Unmarshal([]byte(rsp), &projectScanInfo)

	timestamp := "lastUpload:" + projectScanInfo.ProjectVitals.LastUpdatedDate + " GenReport:" + time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("timestamp: ", timestamp)

	report_file := fmt.Sprintf(
		"report/%s/risk.pdf",
		projectName,
	)
	fmt.Println("report_file: ", report_file)
	err := ipdf.Open(report_file)
	if err != nil {
		fmt.Println("PDF not found")
	}

	ipdf.AddFont("arial", "./angsa.ttf")
	ipdf.SetFont("arial", "", 20)
	ipdf.Insert(timestamp, 1, 300, 10, 100, 100, gopdf.Center|gopdf.Bottom)
	ipdf.Save(report_file)
}
