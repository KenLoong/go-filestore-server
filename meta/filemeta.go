package meta

import "sort"

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string  //文件的存储路径
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

//更新/新增文件元信息
func UploadFileMeta(meta FileMeta)  {
	fileMetas[meta.FileSha1] = meta
}

//返回文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//获取批量文件元信息
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta,len(fileMetas))
	for _,v := range fileMetas{
		fMetaArray = append(fMetaArray,v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

//删除文件
func RemoveFileMeta(fileSha1 string)  {
	delete(fileMetas,fileSha1)
}

