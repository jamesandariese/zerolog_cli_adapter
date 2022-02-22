package zerolog_cli_adapter

import (
	"flag"

	"github.com/rs/zerolog"
)

type LoggerGenerator struct {
	baseLogger zerolog.Logger
	verbose    bool
	debug      bool
	trace      bool
	quiet      bool
	silent     bool
}

// Modifies the passed-in FlagSet with the options used by this library:
//     -v, -verbose, -debug, -trace, -quiet, and -silent
func (lg *LoggerGenerator) UpdateFlagSet(flagSet *flag.FlagSet) {
	flagSet.BoolVar(&lg.verbose, "verbose", false, "Verbose logging (info level logging)")
	flagSet.BoolVar(&lg.verbose, "v", false, "Alias for -verbose")
	flagSet.BoolVar(&lg.debug, "debug", false, "Enable call site logging with debug logging by default")
	flagSet.BoolVar(&lg.trace, "trace", false, "Enable call site logging with trace logging by default")
	flagSet.BoolVar(&lg.quiet, "quiet", false, "Log warnings and errors only")
	flagSet.BoolVar(&lg.silent, "silent", false, "No logs and only show cubbyhole token at end (overrides all other logging flags)")
}

// Returns the configured logger
// flags must be parsed before calling this.
func (lg *LoggerGenerator) Logger() zerolog.Logger {
	logger := lg.baseLogger
	if lg.trace {
		logger = logger.Level(zerolog.TraceLevel).With().Caller().Logger()
	}
	if lg.debug {
		logger = logger.Level(zerolog.DebugLevel).With().Caller().Logger()
	}
	if lg.verbose {
		logger = logger.Level(zerolog.DebugLevel)
	}
	if lg.quiet {
		logger = logger.Level(zerolog.WarnLevel)
	}
	if lg.silent {
		logger = zerolog.Nop()
	}
	return logger
}

// Create a new *LoggerGenerator with the passed base logger
//
// Configure it with lg.UpdateFlagSet and flagset.Parse.
//
// Retrieve configured logger with lg.Logger.
func NewLoggerGenerator(l zerolog.Logger) *LoggerGenerator {
	return &LoggerGenerator{
		baseLogger: l,
	}
}
