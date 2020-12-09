package main // import "github.com/lwahlmeier/goloader"

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = stimlog.GetLogger()
var version string
var config *viper.Viper

var labelCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "goloader_label_total",
	Help: "The total number non-stun packets received",
})

func main() {
	var err error

	config = viper.New()
	config.SetConfigName("config.yaml")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")
	config.AddConfigPath("/")
	config.AddConfigPath("/config/")

	if version == "" || version == "latest" {
		version = "unknown"
	}

	var cmd = &cobra.Command{
		Use:   "goloader",
		Short: "launch goloader service",
		Long:  "launch goloader service",
		Run:   cMain,
	}

	cmd.PersistentFlags().String("loglevel", "info", "level to show logs at (warn, info, debug, trace)")
	config.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))

	cmd.PersistentFlags().String("metricsAddress", "0.0.0.0:8080", "The ip:port to listen for prometheus metrics requests on, only accepts a single address ('127.0.0.1:8080')")
	config.BindPFlag("metricsAddress", cmd.PersistentFlags().Lookup("metricsAddress"))

	err = cmd.Execute()
	CheckError(err)
}

func cMain(cmd *cobra.Command, args []string) {
	err := config.ReadInConfig()
	if err != nil {
		log.Warn("Got Error reading configfile: {}", err)
	}

	if config.GetBool("version") {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}
	var ll stimlog.Level
	switch strings.ToLower(config.GetString("loglevel")) {
	case "info":
		ll = stimlog.InfoLevel
	case "warn":
		ll = stimlog.WarnLevel
	case "debug":
		ll = stimlog.DebugLevel
	case "trace":
		ll = stimlog.TraceLevel
	}
	stimlog.GetLoggerConfig().SetLevel(ll)
	metricsAddress := config.GetString("metricsAddress")
	http.Handle("/metrics", promhttp.Handler())

	go ConfigReLoader()
	go CpuWatcher()
	go MemoryWatcher()
	go SimpleCounterWatcher()
	err = http.ListenAndServe(metricsAddress, nil)
	CheckError(err)
}

func ConfigReLoader() {
	for {
		time.Sleep(time.Second * 5)
		log.Debug("Reading Config File")
		err := config.ReadInConfig()
		if err != nil {
			log.Warn("Got Error reading configfile: {}", err)
		}
	}
}

func CheckError(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal("Fatal Error, Exiting!:{}", err)
	}
}
