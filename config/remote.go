package config

import "os"

//mysql config
var MySql_Password = os.Getenv("MySql_Password")
var MySql_host = os.Getenv("MySql_host")
var MySql_db = os.Getenv("MySql_db")
var MySql_user = os.Getenv("MySql_user")

//update config
var UpdateLink = os.Getenv("UpdateLink")

