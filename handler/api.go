package handler

import (
	"artifact_repository/worker"
	"artifact_repository/wss"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func PackagePage(c *gin.Context) {
	c.HTML(http.StatusOK, "package.tmpl", nil)
}

func SyncPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sync.tmpl", nil)
}

func GetPackage(c *gin.Context) {

	var (
		package_name    string
		package_type    string
		package_version string
		with_conf       string = ""
	)
	var wk worker.WorkerHandler
	if in, isExist := c.GetPostForm("package_name"); isExist && in != "" {
		package_name = in
	} else {
		c.HTML(http.StatusBadRequest, "package.tmpl", gin.H{
			"error": errors.New("Please input package name"),
		})
	}
	if in, isExist := c.GetPostForm("package_type"); isExist && in != "" {
		package_type = in
	} else {
		c.HTML(http.StatusBadRequest, "package.tpml", gin.H{
			"error": errors.New("Please input package type"),
		})
	}
	if in, isExist := c.GetPostForm("package_version"); isExist && in != "" {
		package_version = in
	} else {
		c.HTML(http.StatusBadRequest, "package_version", gin.H{
			"error": errors.New("Please input package version"),
		})
	}

	package_name = fmt.Sprintf("%s==%s", package_name, package_version)
	index_url := os.Getenv("internet_index")
	if package_type == "pypi" {
		wk = worker.NewRepositoryWorker(worker.Pypi{})
	}
	wk.DownloadFromIndex("./tmp/", package_name, index_url)
	project_name := strings.Split(package_name, "==")[0]
	wss.DoWhitesourceScan(package_name, project_name, with_conf)
	wss.DoUploadRequest(project_name)

	ch := wss.GenerateProjectReportAsync(project_name)
	_ = wss.GetProcessStatus(ch, project_name)

	report_path := fmt.Sprintf("report/%s", project_name)
	os.Mkdir(report_path, 0755)
	rsp := wss.GetProjectRiskReport(project_name)
	_, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}
	worker.UploadToRepository(wk, os.Getenv("tmp_api"), fmt.Sprintf("./tmp/%s", package_name))
	c.HTML(http.StatusOK, "package.tmpl", gin.H{
		"report":       fmt.Sprintf("http://localhost:8888/report/%s/risk.pdf", strings.Split(package_name, "==")[0]),
		"project_name": project_name,
		"sync_url":     "http://localhost:8888/sync",
	})
}

func SyncPackage(c *gin.Context) {

	var (
		package_name    string
		package_type    string
		package_version string
	)
	var wk worker.WorkerHandler
	if in, isExist := c.GetPostForm("package_name"); isExist && in != "" {
		package_name = in
	} else {
		c.HTML(http.StatusBadRequest, "package.tmpl", gin.H{
			"error": errors.New("Please input package name"),
		})
	}
	if in, isExist := c.GetPostForm("package_type"); isExist && in != "" {
		package_type = in
	} else {
		c.HTML(http.StatusBadRequest, "package.tpml", gin.H{
			"error": errors.New("Please input package type"),
		})
	}
	if in, isExist := c.GetPostForm("package_version"); isExist && in != "" {
		package_version = in
	} else {
		c.HTML(http.StatusBadRequest, "package_version", gin.H{
			"error": errors.New("Please input package version"),
		})
	}

	package_name = fmt.Sprintf("%s==%s", package_name, package_version)
	index_url := os.Getenv("tmp_index")
	if package_type == "pypi" {
		wk = worker.NewRepositoryWorker(worker.Pypi{})
	}
	wk.DownloadFromIndex("./sync_tmp/", package_name, index_url)
	_ = strings.Split(package_name, "==")[0]
	// wss.DoWhitesourceScan(package_name, project_name)
	// wss.DoUploadRequest(project_name)

	worker.UploadToRepository(wk, os.Getenv("prod_api"), fmt.Sprintf("./tmp/%s", package_name))
	c.HTML(http.StatusOK, "sync.tmpl", nil)

}
