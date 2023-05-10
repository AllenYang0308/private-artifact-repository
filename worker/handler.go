package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Worker interface {
	Download(destination string, packang_name string, index_url string) string
	Sync(target_url string, package_file string) string
	Remove(full_path string) error
}

type WorkerHandler struct {
	worker Worker
}

func (rw WorkerHandler) Download(package_name string, index_url string) {
	rw.worker.Download(
		os.Getenv("package_tmp"),
		package_name,
		index_url,
	)
}

func (rw WorkerHandler) DownloadFromIndex(destination string, package_name string, index_url string) {
	rw.worker.Download(
		destination,
		package_name,
		index_url,
	)
}

func (rw WorkerHandler) Sync(target_url string, package_file string) {
	rw.worker.Sync(
		target_url,
		package_file,
	)
}

func (rw WorkerHandler) Remove(full_path string) error {
	return rw.worker.Remove(full_path)
}

func NewRepositoryWorker(worker Worker) WorkerHandler {
	return WorkerHandler{worker: worker}
}

func UploadToRepository(worker WorkerHandler, target_url string, source_path string) {

	var wg sync.WaitGroup

	files, err := ioutil.ReadDir(source_path)
	if err != nil {
		fmt.Println("讀取目錄錯誤: ", err)
		panic(err)
	}
	for _, file := range files {
		wg.Add(1)
		pkg_name := fmt.Sprintf("%s/%s", source_path, file.Name())
		go func() {
			worker.Sync(target_url, pkg_name)
			wg.Done()
		}()
	}
	wg.Wait()
}
