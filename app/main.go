package main

import (
	"context"
	extractor "github.com/faralaks/go-vk-expander/app/html_builder/html_extractor"
	log "github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
	"os"
	"os/signal"
	"syscall"
)

var config struct {
	App struct {
		Port     string `long:"port" env:"PORT" default:"80" description:"string | Application port"`
		Address  string `long:"address" env:"ADDRESS" default:"0.0.0.0" description:"string | Application web address"`
		Debug    bool   `long:"debug" env:"DEBUG" description:"Debug mode. It provides more output info"`
		InDocker bool   `long:"indocker" env:"INDOCKER" description:"Is this program runs in docker container"`
	} `group:"app" namespace:"app" env-namespace:"VKEXP_APP"`
	OAuth struct {
		ClientId     string `long:"client_id" env:"CLIENT_ID" required:"true"  description:"string | VK OAuth Client ID"`
		Key          string `long:"key" env:"KEY" required:"true"  description:"string | VK OAuth Key"`
		RedirectPath string `long:"redirect_path" env:"REDIRECT_PATH" required:"true"  description:"string | VK OAuth Redirect path like \"/save_token\""`
		RedirectURL  string `long:"redirect_url" env:"REDIRECT_URL" required:"true"  description:"string | VK OAuth full Redirect URL"`
	} `group:"oauth" namespace:"oauth" env-namespace:"VKEXP_OAUTH"`
}
var revision = "unknown"

func main() {
	println("\t ---> Faralaks Vk-Expander starting! Version: " + revision)

	// Load configuration
	flagParser := flags.NewParser(&config, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := flagParser.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			println("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}
	if config.App.InDocker {
		os.Clearenv() // Now only this process has config data
	}

	setupLog(config.App.Debug)
	log.Printf("[DEBUG] Log setup Done!")

	ctx, cancel := context.WithCancel(context.Background())
	go cancelOnInterruptSyscall(cancel)
	run(ctx, "Archive/messages")

	println("\n\t <--- Faralaks Vk-Expander finished!!!")
}

func run(ctx context.Context, path string) {
	_ = extractor.Extract(ctx, path)
}

func setupLog(debug bool) {
	if debug {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}

// cancelOnInterruptSyscall catch signal and invoke graceful termination
func cancelOnInterruptSyscall(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
	log.Printf("[WARN] interrupt signal")
}
