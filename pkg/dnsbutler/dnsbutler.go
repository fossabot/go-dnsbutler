package dnsbutler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func ipHandler(configPath string, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")

		qry := r.URL.Query()
		ip := qry.Get("ip")
		if ip == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "failed\n\nquery param 'ip' is missing")
			return
		}

		if !ipRegex.Match([]byte(ip)) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "'%s' is not a valid ip", ip)
			return
		}

		c, err := readConfig(configPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed\n\nInternal Server Error")
			return
		}

		logger.Println("iphandler called - will update now")

		done := make(chan bool)
		defer close(done)
		go updateTargets(c.Targets, ip, logger, done)
		<-done

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	}
}

func Start(configPath string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("DNSButler is starting...")

	c, err := readOrInitConfig(configPath, logger)
	if err != nil {
		return err
	}

	if len(c.Targets) > 0 {
		ip, err := getIPFrom("https://api.ipify.org/", logger)
		if err != nil {
			logger.Printf("Can't retrieve IP from provider '%s'. Retrieved error '%v'\n", c.Provider, err)
		}

		if ip != "" {
			done := make(chan bool)
			go updateTargets(c.Targets, ip, logger, done)
			<-done
			close(done)
		}
	}

	router := http.NewServeMux()
	router.HandleFunc("/", ipHandler(configPath, logger))

	server := &http.Server{
		Addr:              c.ListenAddr,
		Handler:           router,
		ErrorLog:          logger,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	logger.Println("DNSButler is ready to serve at", c.ListenAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", c.ListenAddr, err)
	}

	return nil
}
