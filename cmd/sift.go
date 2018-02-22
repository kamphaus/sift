package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/kamphaus/sift"
	"github.com/svent/go-flags"
)

var (
	errorLogger = log.New(os.Stderr, "Error: ", 0)
)

func main() {
	var targets []string
	var args []string
	var err error
	var options = sift.Options{}

	parser := flags.NewNamedParser("sift", flags.HelpFlag|flags.PassDoubleDash)
	parser.AddGroup("Options", "Options", &options)
	parser.Name = "sift"
	parser.Usage = "[OPTIONS] PATTERN [FILE|PATH|tcp://HOST:PORT]...\n" +
		"  sift [OPTIONS] [-e PATTERN | -f FILE] [FILE|PATH|tcp://HOST:PORT]...\n" +
		"  sift [OPTIONS] --targets [FILE|PATH]..."

	// temporarily parse options to see if the --no-conf/--conf options were used and
	// then discard the result
	options.LoadDefaults()
	args, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			fmt.Println(e.Error())
			os.Exit(0)
		} else {
			errorLogger.Println(err)
			os.Exit(2)
		}
	}
	noConf := options.NoConfig
	configFile := options.ConfigFile
	options = sift.Options{}

	s := sift.NewSearch(&options, nil, os.Stdout, errorLogger)

	// perform full option parsing respecting the --no-conf/--conf options
	s.GetOptions().LoadDefaults()
	s.LoadConfigs(noConf, configFile)
	args, err = parser.Parse()
	if err != nil {
		errorLogger.Println(err)
		os.Exit(2)
	}

	targets = s.ProcessArgs(args)

	if err := s.Apply(s.GetMatchPatterns(), targets); err != nil {
		errorLogger.Fatalf("cannot process options: %s\n", err)
	}
	runtime.GOMAXPROCS(s.GetOptions().Cores)

	retVal, err := s.ExecuteSearch(targets)
	if err != nil {
		errorLogger.Println(err)
	}
	os.Exit(retVal)
}
