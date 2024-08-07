package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"ypeskov/go-password-manager/internal/config"
	"ypeskov/go-password-manager/internal/logger"
	"ypeskov/go-password-manager/internal/migrations"
	"ypeskov/go-password-manager/internal/server"
)

func main() {
	launchBrowser := flag.Bool("browser", false, "Launch browser on start")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		fmt.Printf(fmt.Sprintf("cannot read config: %s", err))
		panic(err)
	}

	appLogger := logger.New(cfg)

	err = migrations.MakeMigration(appLogger, cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot make migration: %s", err))
	}

	appServer := server.New(cfg, appLogger)

	if *launchBrowser {
		openBrowser(fmt.Sprintf("http://localhost:%s", cfg.Port))
	}

	err = appServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}

/*
openBrowser opens the default web browser at the specified URL.
*/
func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = "open"
		args = append(args, url)
	case "linux":
		cmd = "xdg-open"
		args = append(args, url)
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		return
	}
	fmt.Printf("Opening browser at %s\n", url)
	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}
