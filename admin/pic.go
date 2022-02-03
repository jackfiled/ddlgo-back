package admin

import (
	"ddl/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFileHandler(c *gin.Context) {

	fileNames := make([]string, 0)

	mForm := c.Request.MultipartForm

	for k, _ := range mForm.File {
		// k is the key of file part
		file, fileHeader, err := c.Request.FormFile(k)
		if err != nil {
			c.String(http.StatusBadRequest, "inovke FormFile error:"+err.Error())
			fmt.Println("inovke FormFile error:", err)
			return
		}
		defer file.Close()
		fmt.Printf("the uploaded file: \nname:[%s]\nsize[%d]\nheader[%#v]\n",
			fileHeader.Filename, fileHeader.Size, fileHeader.Header)

		// store uploaded file into local path
		now := time.Now()
		nowString := now.Format("060102150405")
		localFileName := config.UPLOAD_PATH + "/p" + nowString + "_" + fileHeader.Filename

		out, err := os.Create(localFileName)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			fmt.Println(err.Error())
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			fmt.Println(err.Error())
			return
		}
		fileNames = append(fileNames, localFileName)
	}

	c.JSON(200, fileNames)
}
