package controllers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

func WalkFiles(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var result = map[string]interface{}{}
	var err = filepath.Walk(cc.App.Config.EFSPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			result[path] = map[string]interface{}{
				"is_dir":   info.IsDir(),
				"mod_time": info.ModTime(),
				"mode":     info.Mode(),
				"name":     info.Name(),
				"size":     info.Size(),
			}
			return nil
		})
	if err != nil {
		return err
	}
	return cc.Success(result)
}

func DeleteFile(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var fileName = cc.GetPathParamString("file_name")
	var err = os.Remove(fmt.Sprintf("%s/files/%s", cc.App.Config.EFSPath, fileName))
	if err != nil {
		return err
	}

	return cc.Success("Removed file")
}
