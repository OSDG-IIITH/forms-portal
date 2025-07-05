package utils

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	Production bool

	Host   string
	Port   string
	Domain string

	DatabaseUri string
	FrontendUrl string

	CasBaseUrl    string
	CasServiceUrl string

	SessionSecret string
}

func defaultConfig() config {
	return config{
		Production: false,

		Host:   "127.0.0.1",
		Port:   "8647",
		Domain: "localhost",

		DatabaseUri: "postgres://user:pass@localhost:5432/forms",
		FrontendUrl: "http://localhost:8648",

		CasBaseUrl:    "https://login.iiit.ac.in/cas",
		CasServiceUrl: "http://localhost:8647/api/auth/login/callback",

		SessionSecret: "quis-custodiet-ipsos-custodes",
	}
}

func LoadConfig() {
	c := defaultConfig()

	prod, ok := os.LookupEnv("PRODUCTION")
	if ok && (prod == "true" || prod == "1") {
		c.Production = true
	}
	host, ok := os.LookupEnv("FORMS_SERVER_HOST")
	if ok {
		c.Host = host
	}
	port, ok := os.LookupEnv("FORMS_SERVER_PORT")
	if ok {
		c.Port = port
	}
	domain, ok := os.LookupEnv("FORMS_APP_DOMAIN")
	if ok {
		c.Domain = domain
	}

	databaseUri, ok := os.LookupEnv("FORMS_DATABASE_URI")
	if ok {
		c.DatabaseUri = databaseUri
	}
	frontendUrl, ok := os.LookupEnv("FORMS_FRONTEND_URL")
	if ok {
		c.FrontendUrl = frontendUrl
	}

	casBaseUrl, ok := os.LookupEnv("FORMS_CAS_BASE_URL")
	if ok {
		c.CasBaseUrl = casBaseUrl
	}
	casServiceUrl, ok := os.LookupEnv("FORMS_CAS_SERVICE_URL")
	if ok {
		c.CasServiceUrl = casServiceUrl
	}

	sessionSecret, ok := os.LookupEnv("FORMS_SESSION_SECRET")
	if ok {
		c.SessionSecret = sessionSecret
	}

	Config = c
}
