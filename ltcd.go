// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"

	"github.com/ltcsuite/ltcd/blockchain/indexers"
	"github.com/ltcsuite/ltcd/limits"
)

var (
	cfg *config
)

// winServiceMain is only invoked on Windows.  It detects when ltcd is running
// as a service and reacts accordingly.
var winServiceMain func() (bool, error)

// ltcdMain is the real main function for ltcd.  It is necessary to work around
// the fact that deferred functions do not run when os.Exit() is called.  The
// optional serverChan parameter is mainly used by the service code to be
// notified with the server once it is setup so it can gracefully stop it when
// requested from the service control manager.
func ltcdMain(serverChan chan<- *server) error {
	// Load configuration and parse command line.  This function also
	// initializes logging and configures it accordingly.
	tcfg, _, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = tcfg
	defer backendLog.Flush()

	// Get a channel that will be closed when a shutdown signal has been
	// triggered either from an OS signal such as SIGINT (Ctrl+C) or from
	// another subsystem such as the RPC server.
	interruptedChan := interruptListener()
	defer ltcdLog.Info("Shutdown complete")

	// Show version at startup.
	ltcdLog.Infof("Version %s", version())

	// Enable http profiling server if requested.
	if cfg.Profile != "" {
		go func() {
			listenAddr := net.JoinHostPort("", cfg.Profile)
			ltcdLog.Infof("Profile server listening on %s", listenAddr)
			profileRedirect := http.RedirectHandler("/debug/pprof",
				http.StatusSeeOther)
			http.Handle("/", profileRedirect)
			ltcdLog.Errorf("%v", http.ListenAndServe(listenAddr, nil))
		}()
	}

	// Write cpu profile if requested.
	if cfg.CPUProfile != "" {
		f, err := os.Create(cfg.CPUProfile)
		if err != nil {
			ltcdLog.Errorf("Unable to create cpu profile: %v", err)
			return err
		}
		pprof.StartCPUProfile(f)
		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	// Perform upgrades to ltcd as new versions require it.
	if err := doUpgrades(); err != nil {
		ltcdLog.Errorf("%v", err)
		return err
	}

	// Return now if an interrupt signal was triggered.
	if interruptRequested(interruptedChan) {
		return nil
	}

	// Load the block database.
	db, err := loadBlockDB()
	if err != nil {
		ltcdLog.Errorf("%v", err)
		return err
	}
	defer func() {
		// Ensure the database is sync'd and closed on shutdown.
		ltcdLog.Infof("Gracefully shutting down the database...")
		db.Close()
	}()

	// Return now if an interrupt signal was triggered.
	if interruptRequested(interruptedChan) {
		return nil
	}

	// Drop indexes and exit if requested.
	//
	// NOTE: The order is important here because dropping the tx index also
	// drops the address index since it relies on it.
	if cfg.DropAddrIndex {
		if err := indexers.DropAddrIndex(db); err != nil {
			ltcdLog.Errorf("%v", err)
			return err
		}

		return nil
	}
	if cfg.DropTxIndex {
		if err := indexers.DropTxIndex(db); err != nil {
			ltcdLog.Errorf("%v", err)
			return err
		}

		return nil
	}

	// Create server and start it.
	server, err := newServer(cfg.Listeners, db, activeNetParams.Params)
	if err != nil {
		// TODO: this logging could do with some beautifying.
		ltcdLog.Errorf("Unable to start server on %v: %v",
			cfg.Listeners, err)
		return err
	}
	defer func() {
		ltcdLog.Infof("Gracefully shutting down the server...")
		server.Stop()
		server.WaitForShutdown()
		srvrLog.Infof("Server shutdown complete")
	}()
	server.Start()
	if serverChan != nil {
		serverChan <- server
	}

	// Wait until the interrupt signal is received from an OS signal or
	// shutdown is requested through one of the subsystems such as the RPC
	// server.
	<-interruptedChan
	return nil
}

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Block and transaction processing can cause bursty allocations.  This
	// limits the garbage collector from excessively overallocating during
	// bursts.  This value was arrived at with the help of profiling live
	// usage.
	debug.SetGCPercent(10)

	// Up some limits.
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}

	// Call serviceMain on Windows to handle running as a service.  When
	// the return isService flag is true, exit now since we ran as a
	// service.  Otherwise, just fall through to normal operation.
	if runtime.GOOS == "windows" {
		isService, err := winServiceMain()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if isService {
			os.Exit(0)
		}
	}

	// Work around defer not working after os.Exit()
	if err := ltcdMain(nil); err != nil {
		os.Exit(1)
	}
}
