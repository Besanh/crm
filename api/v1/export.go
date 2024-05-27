package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/response"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Exported struct {
	exportService service.IExport
}

func NewExport(r *gin.Engine, export service.IExport) {
	handler := &Exported{
		exportService: export,
	}
	Group := r.Group(constants.VERSION + "/crm/export")
	{
		Group.GET("", authMdw.AuthMiddleware(), handler.GetExports)
		Group.GET(":type/:id/download", handler.DownloadExported)
	}
}

func (data *Exported) GetExports(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, result := data.exportService.GetExports(domainUuid, userUuid)
	c.JSON(code, result)
}

func (data *Exported) DownloadExported(c *gin.Context) {
	// domainUuid := c.Query("domain_uuid")
	userUuid := c.Query("user_uuid")
	id := c.Param("id")
	code, fileData := data.exportService.DownloadExport(userUuid, id)
	if code != http.StatusOK {
		c.JSON(response.BadRequestMsg("file is not existed"))
		return
	}
	fileStr, _ := fileData.(string)
	exportValue := strings.Split(fileStr, ";")
	exportFolder := c.Param("type")
	if len(exportValue) > 5 {
		exportFolder = exportValue[len(exportValue)-1]
	}
	fileDir := fmt.Sprintf(constants.EXPORT_DIR+"%s/%s", exportFolder, id)
	file, err := os.Open(fileDir)
	if err != nil {
		c.JSON(response.BadRequestMsg("file is not existed"))
		return
	}
	defer file.Close()
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+id)
	c.Writer.Header().Set("File-Name", id)
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Expires", "0")
	c.Writer.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(response.BadRequestMsg("file is not existed"))
		return
	}
}
