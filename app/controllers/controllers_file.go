package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/repository"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"tanya_dokter_app/config"
	"time"

	"github.com/labstack/echo/v4"
)

// UploadFile godoc
// @Summary File Uploader
// @Description File Uploader
// @Tags File
// @Accept mpfd
// @Param file formData file true "File to upload"
// @Produce json
// @Success 200
// @Router /v1/file [post]
// @Security ApiKeyAuth
// @Security JwtToken
func UploadFile(c echo.Context) error {

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
		// return c.JSON(400, map[string]interface{}{
		// 	"status":  400,
		// 	"message": "Failed to get Asia/Jakarta time. Error: " + err.Error(),
		// 	"error":   err.Error(),
		// })
	}
	userIDData := c.Get("user_id")
	if userIDData == nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Wrong Authorization Token"))
	}
	userID := userIDData.(int)
	// Define the accepted MIME types
	acceptedTypes := []string{
		"image/png",
		"image/jpeg",
		"image/gif",
		"video/quicktime",
		"video/mp4",
		"application/pdf",
		"text/csv",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-excel.sheet.macroenabled.12",
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Get the MIME type of the file
	fileType := file.Header.Get("Content-Type")
	extension := ".jpg"
	if fileType == "image/png" {
		extension = ".png"
	}
	if fileType == "image/jpeg" {
		extension = ".jpg"
	}
	if fileType == "image/gif" {
		extension = ".gif"
	}
	if fileType == "video/quicktime" {
		extension = ".mov"
	}
	if fileType == "video/mp4" {
		extension = ".mov"
	}
	if fileType == "application/pdf" {
		extension = ".pdf"
	}
	if fileType == "application/vnd.ms-excel" {
		extension = ".xls"
	}
	if fileType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		extension = ".xlsx"
	}
	if fileType == "application/vnd.ms-excel.sheet.macroenabled.12" {
		extension = ".et"
	}
	if fileType == "text/csv" {
		extension = ".csv"
	}
	// Check if the MIME type is accepted
	var isAccepted bool
	for _, t := range acceptedTypes {
		if t == fileType {
			isAccepted = true
			break
		}
	}

	if !isAccepted {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"accepted_type": acceptedTypes,
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a destination file
	t := time.Now().In(location)
	time := t.Format("2006-01")
	folder := strconv.Itoa(int(userID)) + "/" + time
	err = os.MkdirAll(config.RootPath()+"/assets/uploads/"+folder, os.ModePerm)
	if err != nil {
		return err
	}
	timestamp := strconv.Itoa(int(t.Unix()))

	// Rename the file to the timestamp with the original file extension
	filePath := filepath.Join(config.RootPath()+"/assets/uploads/", folder, timestamp+extension)

	fmt.Println(filePath)
	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Copy the file to the destination
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	data, err := SaveFileToDatabase(userID, folder+"/"+timestamp+extension, filePath)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}
	data.FullUrl = config.LoadConfig().BaseUrl + "/assets/uploads/" + folder + "/" + timestamp + extension
	// if config.LoadConfig().IsDesktop {
	// 	data.FullUrl = filePath
	// }

	if config.LoadConfig().EnableEncodeID {
		if data.EncodedID == "" {
			data.EncodedID = utils.EndcodeID(int(data.ID))
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "File uploaded successfully",
	})
}

func SaveFileToDatabase(id int, filename, path string) (data models.GlobalFile, err error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}
	t := time.Now().In(location).Unix()
	data = models.GlobalFile{
		Token:    strconv.Itoa(int(t)) + utils.GenerateRandomString(5),
		UserID:   id,
		Filename: filename,
		Path:     path,
	}

	err = repository.SaveFile(&data)
	return
}

// GetFile godoc
// @Summary Mendapatkan List Files
// @Description Mendapatkan List Files
// @Tags File
// @Accept json
// @Param search query string false "search (string)"
// @Param page query integer false "page (int)"
// @Param limit query integer false "limit (int)"
// @Param token query string false "token (string)"
// @Param company_id query integer false "company_id (int)"
// @Produce json
// @Success 200
// @Router /v1/file [get]
// @Security ApiKeyAuth
// @Security JwtToken
func GetFile(c echo.Context) error {

	userID := c.Get("user_id").(int)

	param := utils.PopulatePaging(c, "token")

	data, err := GetFileControl(userID, param)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}
	return c.JSON(http.StatusOK, data)
}

func GetFileControl(id int, param reqres.ReqPaging) (data reqres.ResPaging, err error) {
	// key := fmt.Sprintf("UCGetFiles%v_%v_%v_%v_%v", param.Order, param.Search, param.Page, id, param.Offset)
	// outputCached, _ := config.RC.Get(key).Result()
	// if err = json.Unmarshal([]byte(outputCached), &data); err == nil {
	// 	return
	// }
	data, err = repository.GetFile(id, param)
	if err != nil {
		return
	}
	// json, err := json.Marshal(&data)
	// if err == nil {
	// 	config.RC.Set(key, json, 30*time.Second).Err()
	// }
	return
}

// UploadMultipleFiles godoc
// @Summary Upload multiple files
// @Description Uploads multiple files
// @Tags File
// @Accept mpfd
// @Param files formData file true "Files to upload"
// @Param width formData string false "width"
// @Param height formData string false "height"
// @Param folder query string false "Folder to store"
// @Param company_id query integer false "company_id (int)"
// @Produce json
// @Success 200
// @Router /v1/file/multi [post]
// @Security JwtToken
func UploadMultipleFiles(c echo.Context) error {

	userid := c.Get("user_id")
	if userid == nil {
		return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Token Authorization salah"))
	}
	userId := userid.(int)

	// Define the accepted MIME types
	acceptedTypes := []string{"image/png", "image/jpeg", "image/gif", "video/quicktime", "video/mp4", "application/pdf"}

	// Get the files from the request
	storefolder := c.QueryParam("folder")
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}
	files := form.File["files"]
	fmt.Println("Files", files)
	// Create a destination folder
	t := time.Now()
	varTime := t.Format("2006-01")
	if storefolder != "" {
		storefolder = storefolder + "/"
	}
	folder := storefolder + strconv.Itoa(int(userId)) + "/" + varTime
	err = os.MkdirAll(config.RootPath()+"/assets/uploads/"+folder, os.ModePerm)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	var uploadedFiles []models.GlobalFile

	for _, file := range files {
		tDetail := time.Now()
		// Get the MIME type of the file
		fileType := file.Header.Get("Content-Type")
		extension := getExtensionForFileType(fileType)

		// Check if the MIME type is accepted
		var isAccepted bool
		for _, t := range acceptedTypes {
			if t == fileType {
				isAccepted = true
				break
			}
		}

		if !isAccepted {
			continue // Skip this file
		}

		// Open the file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Create a destination file
		timestamp := strconv.FormatInt(tDetail.UnixNano(), 10)
		filePath := filepath.Join(config.RootPath()+"/assets/uploads/", folder, timestamp+extension)
		err = saveFile(src, filePath)
		if err != nil {
			return err
		}

		data, err := SaveFileToDatabase(int(userId), folder+"/"+timestamp+extension, filePath)
		if err != nil {
			return c.JSON(utils.ParseHttpError(err))
		}
		data.FullUrl = config.LoadConfig().BaseUrl + "/assets/uploads/" + folder + "/" + timestamp + extension

		if config.LoadConfig().EnableEncodeID {
			if data.EncodedID == "" {
				data.EncodedID = utils.EndcodeID(int(data.ID))
			}
		}
		uploadedFiles = append(uploadedFiles, data)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"data":    uploadedFiles,
		"message": "Upload File Berhasil !",
	})
}

func getExtensionForFileType(fileType string) string {
	switch fileType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "video/quicktime":
		return ".mov"
	case "video/mp4":
		return ".mp4"
	case "application/pdf":
		return ".pdf"
	default:
		return ".unknown"
	}
}

func saveFile(src multipart.File, destPath string) error {
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
