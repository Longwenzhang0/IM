package ctrl

import (
	"IM/util"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	err := os.Mkdir("./mnt", os.ModePerm)
	if err != nil {
		log.Printf("make directory ./mnt failed,the detailed error is %s: ", err.Error())
		return
	}
}

// 上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	UploadLocal(w, r)
	//UploadOss(w,r)
}

// 1. 存储位置 ./mnt/,确保已经创建完毕；init创建
// 2. url的格式为：/mnt/xxx.png	确保网络能访问./mnt/
func UploadLocal(writer http.ResponseWriter, request *http.Request) {
	srcfile, head, err := request.FormFile("file")
	if err != nil {
		// 发生错误时返回failed响应
		util.RespFail(writer, err.Error())
	}

	// 创建一个新的文件
	suffix := ".png"
	// 如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	// 如果前端指定了filetype，就以type为准;这个值是前端填充的
	//  formdata.append("filetype",".mp3");
	filetype := request.FormValue("filetype")
	if len(filetype) > 0 {
		// 后缀判空
		suffix = filetype
	}
	// 随机生成名称
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	// 创建新文件
	dstfile, err := os.Create("./mnt/" + filename)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	// 将源文件copy到新文件
	_, err = io.Copy(dstfile, srcfile)
	if err != nil {
		util.RespFail(writer, err.Error())
	}
	// 将新文件路径转换为url地址；
	url := "/mnt/" + filename
	// 响应到前端
	util.RespOk(writer, url, "")

}

// 即将删掉,定期更新
const (
	AccessKeyId     = "5p2RZKnrUanMuQw9"
	AccessKeySecret = "bsNmjU8Au08axedV40TRPCS5XIFAkK"
	EndPoint        = "oss-cn-shenzhen.aliyuncs.com"
	Bucket          = "winliondev"
)

// 权限设置为公共读状态，阿里后台页面设置？oss: object storage Service 对象存储
// 需要安装
func UploadOss(writer http.ResponseWriter,
	request *http.Request) {
	//todo 获得上传的文件
	srcfile, head, err := request.FormFile("file")
	if err != nil {
		// 发生错误时返回failed响应
		util.RespFail(writer, err.Error())
	}
	//todo 获得文件后缀.png/.mp3
	suffix := ".png"
	// 如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	// 如果前端指定了filetype，就以type为准;这个值是前端填充的
	//  formdata.append("filetype",".mp3");
	filetype := request.FormValue("filetype")
	if len(filetype) > 0 {
		// 后缀判空
		suffix = filetype
	}
	//todo 初始化ossclient
	client, err := oss.New(EndPoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	//todo 获得bucket
	bucket, err := client.Bucket(Bucket)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	//todo 设置文件名称
	filename := fmt.Sprintf("mnt/%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)

	//todo 通过bucket上传
	err = bucket.PutObject(filename, srcfile)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	//todo 获得url地址
	url := "http://" + Bucket + "." + EndPoint + "/" + filename

	//todo 响应到前端
	util.RespOk(writer, url, "")
}
