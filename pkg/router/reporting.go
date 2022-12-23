package router

import (
	"encoding/json"
	"github.com/g-portal/metadata-server/pkg/sources"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func CloudInitReport(c *gin.Context) {
	ip := sources.GetServer(c.Request)
	if ip == nil {
		log.Printf("Failed to get remote address: %v", sources.ErrFailedGetRemoteAddress)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	cloudInitReport := &sources.CloudInitReport{}
	err = json.Unmarshal(bodyBytes, &cloudInitReport)
	if err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	report := cloudInitReport.ToReportMessage()
	report.IP = ip

	if sourceList, ok := c.MustGet("datasources").([]sources.Source); ok {
		for _, source := range sourceList {
			err = source.ReportLog(report)
			if err != nil {
				log.Printf("Failed to report log to datasource %s: %v", source.Type(), err)

				continue
			}
		}
	} else {
		_ = c.AbortWithError(http.StatusInternalServerError, sources.ErrNoDatasourceFound)

		return
	}

	c.Status(http.StatusOK)
}
