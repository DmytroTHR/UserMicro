package configs

import (
	"os"
)

var PG_HOST = getStringParameter("PG_HOST", "localhost")
var PG_PORT = getStringParameter("PG_PORT", "5444")
var POSTGRES_DB = getStringParameter("POSTGRES_DB", "scooterdb")
var POSTGRES_USER = getStringParameter("POSTGRES_USER", "scooteradmin")
var POSTGRES_PASSWORD = getStringParameter("POSTGRES_PASSWORD", "Megascooter!")
var GRPC_PORT = getStringParameter("USERS_GRPC_PORT", "5555")
var TOKEN_SECRET = getStringParameter("TOKEN_SECRET", "SomeSuperSECRETpassword123!#$")

func getStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}