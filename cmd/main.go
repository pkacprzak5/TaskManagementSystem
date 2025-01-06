package main

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkacprzak5/TaskManagementSystem/internal/app"
	"github.com/pkacprzak5/TaskManagementSystem/internal/common"
	"log"
	"os"
	"os/signal"
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

	sqlStorage := app.NewMySQLStorage(cfg)

	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	port := fmt.Sprintf(":%v", common.Envs.Port)

	store := common.NewStore(db)
	api := app.NewAPIServer(port, store)
	err = api.Serve(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
