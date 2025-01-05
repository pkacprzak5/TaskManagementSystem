package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
)

func main() {
	cfg := mysql.Config{
		User:                 common.Envs.DBUser,
		Passwd:               common.Envs.DBPassword,
		Addr:                 common.Envs.DBAddress,
		DBName:               common.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

}
