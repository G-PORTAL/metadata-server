package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

const ForbiddenResponse = `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>403 - Forbidden</title>
 </head>
 <body>
  <h1>403 - Forbidden</h1>
 </body>
</html>`

func ForbiddenRequest(c *gin.Context) {
	c.Render(http.StatusForbidden, render.Data{
		Data:        []byte(ForbiddenResponse),
		ContentType: "text/html",
	})
}
