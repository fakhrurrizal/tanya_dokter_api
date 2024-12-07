package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                     string
	AppKey                      string
	BaseUrl                     string
	FrontEndUrl                 string
	Environtment                string
	DatabaseUsername            string
	DatabasePassword            string
	DatabaseHost                string
	DatabasePort                string
	DatabaseName                string
	DatabasePlannerName         string
	PathDB                      string
	CacheURL                    string
	CachePassword               string
	LoggerLevel                 string
	ContextTimeout              int
	Port                        string
	GoogleClientID              string
	GoogleClientSecret          string
	POSFrontendUrl              string
	BOFrontendUrl               string
	EnableCronJob               bool
	EnableConcurrent            bool
	EnableCSRF                  bool
	EnableDatabaseAutomigration bool
	EnableSaas                  bool
	EnableAPIKey                bool
	APIKey                      string
	IsDesktop                   bool
	GoldAPIUrl                  string
	// GoldAPIKey                  string
	OpenExchangeRatesUrl        string
	APIGeolocationAPIKey        string
	RunLocalDatabaseVia         string
	EnableEncodeID              bool
	EnableIDDuplicationHandling bool
	MailMailer                  string
	MailHost                    string
	MailPort                    int
	MailUsername                string
	MailPassword                string
	MailEncryption              string
}

func LoadConfig() (config *Config) {

	fmt.Println("Root Path :", RootPath())
	if err := godotenv.Load(RootPath() + `/.env`); err != nil {
		fmt.Println("error loading .env file:", err)

	}

	appName := os.Getenv("APP_NAME")
	appKey := os.Getenv("APP_KEY")
	baseurl := os.Getenv("BASE_URL")
	frontendurl := os.Getenv("FRONT_END_URL")
	environment := strings.ToUpper(os.Getenv("ENVIRONMENT"))
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseName := os.Getenv("DATABASE_NAME")
	databasePlannerName := os.Getenv("DATABASE_PLANNER_NAME")
	PathDB := os.Getenv("PATH_DB")
	cacheURL := os.Getenv("CACHE_URL")
	cachePassword := os.Getenv("CACHE_PASSWORD")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	port := os.Getenv("PORT")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	mailMailer := os.Getenv("MAIL_MAILER")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailEncryption := os.Getenv("MAIL_ENCRYPTION")
	posFrontendUrl := os.Getenv("POS_FRONT_END_URL")
	boFrontendUrl := os.Getenv("BO_FRONT_END_URL")
	enableCronJob, _ := strconv.ParseBool(os.Getenv("ENABLE_CRONJOB"))
	enableConcurrent, _ := strconv.ParseBool(os.Getenv("ENABLE_CONCURRENT"))
	enableCSRF, _ := strconv.ParseBool(os.Getenv("ENABLE_CSRF"))
	enableDatabaseAutomigration, _ := strconv.ParseBool(os.Getenv("ENABLE_DATABASE_AUTOMIGRATION"))
	enableSaas, _ := strconv.ParseBool(os.Getenv("ENABLE_SAAS"))
	enableApiKey, _ := strconv.ParseBool(os.Getenv("ENABLE_API_KEY"))
	goldAPIUrl := os.Getenv("GOLDAPI_URL")
	openExchangeRatesUrl := os.Getenv("OPENEXCHAGERATES_URL")
	apiGeolocationAPIKey := os.Getenv("APIGEOLOCATION_API_KEY")
	runLocalDatabaseVia := strings.ToUpper(os.Getenv("RUN_LOCAL_DATABASE_VIA"))
	enableEncodeID, _ := strconv.ParseBool(os.Getenv("ENABLE_ENCODE_ID"))
	apiKey := os.Getenv("API_KEY")
	enableIDDuplicationHandling, _ := strconv.ParseBool(os.Getenv("ENABLE_ID_UPLICATION_HANDLING"))

	return &Config{
		AppName:                     appName,
		AppKey:                      appKey,
		BaseUrl:                     baseurl,
		FrontEndUrl:                 frontendurl,
		Environtment:                environment,
		DatabaseUsername:            databaseUsername,
		DatabasePassword:            databasePassword,
		DatabaseHost:                databaseHost,
		DatabasePort:                databasePort,
		DatabaseName:                databaseName,
		DatabasePlannerName:         databasePlannerName,
		PathDB:                      PathDB,
		CacheURL:                    cacheURL,
		CachePassword:               cachePassword,
		LoggerLevel:                 loggerLevel,
		ContextTimeout:              contextTimeout,
		Port:                        port,
		GoogleClientID:              googleClientID,
		GoogleClientSecret:          googleClientSecret,
		POSFrontendUrl:              posFrontendUrl,
		MailMailer:                  mailMailer,
		MailHost:                    mailHost,
		MailUsername:                mailUsername,
		MailPassword:                mailPassword,
		MailPort:                    mailPort,
		MailEncryption:              mailEncryption,
		BOFrontendUrl:               boFrontendUrl,
		EnableCronJob:               enableCronJob,
		EnableConcurrent:            enableConcurrent,
		EnableCSRF:                  enableCSRF,
		EnableDatabaseAutomigration: enableDatabaseAutomigration,
		EnableSaas:                  enableSaas,
		GoldAPIUrl:                  goldAPIUrl,
		// GoldAPIKey:                  goldAPIKey,
		OpenExchangeRatesUrl:        openExchangeRatesUrl,
		APIGeolocationAPIKey:        apiGeolocationAPIKey,
		RunLocalDatabaseVia:         runLocalDatabaseVia,
		EnableEncodeID:              enableEncodeID,
		EnableAPIKey:                enableApiKey,
		APIKey:                      apiKey,
		EnableIDDuplicationHandling: enableIDDuplicationHandling,
	}
}

func RootPath() string {
	projectDirName := os.Getenv("DIR_NAME")
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	return string(rootPath)
}
