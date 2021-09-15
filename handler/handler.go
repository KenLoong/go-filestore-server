package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回html文件
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w,"inter error : "+err.Error())
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == "POST" {
		//接收文件流
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Println("failed to get data,err :",err.Error())
			return
		}
		defer file.Close()

		//在window下编程，暂时用这个路径
		stroePath := "F:\\GoProjects\\src\\filestore-server\\tmp\\"
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: stroePath + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//newFile, err := os.Create("/tmp/" + head.Filename) //这是linux系统下的路径
		newFile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Println("failed to create file,err :",err.Error())
			return
		}
		defer newFile.Close()
		//拷贝到新文件中,记录文件大小
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Println("failed to save into file,err :",err.Error())
			return
		}

		//Seek设置下一次读/写的位置。offset为相对偏移量，而whence决定相对位置
		newFile.Seek(0,0)
		//计算哈希值
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UploadFileMeta(fileMeta)

		//重定向
		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
}

//返回上传成功响应
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,"Upload finished!")
}

//根据sha1值获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析表单

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil{
		//500响应码
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

//批量查询文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fileMetas)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

//下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	//从map中获取元信息
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	//一次性读完文件，不建议使用
	data, err := ioutil.ReadAll(f)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//使浏览器知道可以下载文件
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)

}

//更新文件名
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UploadFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除文件
func FileDeleteHandler(w http.ResponseWriter , r *http.Request)  {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fileMeta := meta.GetFileMeta(fileSha1)
	//fmt.Println(fileMeta.Location)
	//fmt.Println(fileMeta.FileName)
	//删除文件
	os.Remove(fileMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}