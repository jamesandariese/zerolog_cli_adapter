package zerolog_cli_adapter

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func jsonMustMarshal(v interface{}) string {
	m, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(m)
}

func generateTestLogger() (zerolog.Logger, func() string) {
	out := &strings.Builder{}
	writer := zerolog.ConsoleWriter{Out: out}
	writer.NoColor = true
	writer.FormatErrFieldName = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("'ERRFIELDNAME:%s'", i)
	}
	writer.FormatErrFieldValue = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("'ERRFIELDVALUE:%s'", i)
	}
	writer.FormatTimestamp = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return "'TIMESTAMP'"
	}
	writer.FormatLevel = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return strings.ToUpper(fmt.Sprintf("'LEVEL:%s'", i))
	}
	writer.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("'MESSAGE:%s'", i)
	}
	writer.FormatFieldName = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("'FIELDNAME:%s'", i)
	}
	writer.FormatFieldValue = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("'FIELDVALUE:%s'", i)
	}
	writer.FormatCaller = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("'CALLER:TESTVAL'")
	}

	readAndReset := func() string {
		rs := out.String()
		out.Reset()
		return rs
	}
	return zerolog.New(writer).With().Timestamp().Logger(), readAndReset
}

//func TestFlagDebug(t *testing.T) {
//	fs := flag.NewFlagSet("Flaggles", flag.ContinueOnError)
//	baseLogger, readAndReset := generateTestLogger()
//	lg := NewLoggerGenerator(baseLogger)
//	lg.UpdateFlagSet(fs)
//	fs.Parse([]string{"-debug"})
//	if !lg.debug || lg.verbose || lg.trace || lg.quiet || lg.silent {
//		t.Errorf("debug not set in logger generator")
//	}
//	logger := lg.Logger()
//	logger.Info().Msg("hello info")
//	logger.Debug().Msg("hello debug")
//	logger.Trace().Msg("hello trace")
//	expected := `"'TIMESTAMP' 'LEVEL:INFO' 'CALLER:TESTVAL' 'MESSAGE:hello info'\n` +
//		`'TIMESTAMP' 'LEVEL:DEBUG' 'CALLER:TESTVAL' 'MESSAGE:hello debug'\n"`
//	if acc := jsonMustMarshal(readAndReset()); acc != expected {
//		t.Errorf("logged \"%s\", expected \"%s\"", acc, expected)
//	}
//}

func testFlags(t *testing.T, f []string, expected string, testLG func(*LoggerGenerator) (bool, string)) {
	fs := flag.NewFlagSet("Flaggles", flag.ContinueOnError)
	baseLogger, readAndReset := generateTestLogger()
	lg := NewLoggerGenerator(baseLogger)
	lg.UpdateFlagSet(fs)
	fs.Parse(f)
	if ok, msg := testLG(lg); !ok {
		t.Error(msg)
	}
	logger := lg.Logger()
	logger.Trace().Msg("hello trace")
	logger.Debug().Msg("hello debug")
	logger.Info().Msg("hello info")
	logger.Warn().Msg("hello warn")
	logger.Error().Msg("hello error")
	//logger.Fatal().Msg("hello fatal")
	//logger.Panic().Msg("hello panic")
	if acc := jsonMustMarshal(readAndReset()); acc != expected {
		t.Errorf("logged \"%s\", expected \"%s\"", acc, expected)
	}
}

var expectedOutputs map[string]string = make(map[string]string)

func init() {
	expectedOutputs = make(map[string]string)
	for _, level := range []string{"debug", "info", "trace", "error", "warn", "fatal", "panic"} {
		expectedOutputs[level] = fmt.Sprintf(`'TIMESTAMP' 'LEVEL:%s' 'MESSAGE:hello %s'\n`, strings.ToUpper(level), level)
		expectedOutputs[level+"_caller"] = fmt.Sprintf(`'TIMESTAMP' 'LEVEL:%s' 'CALLER:TESTVAL' 'MESSAGE:hello %s'\n`, strings.ToUpper(level), level)
	}
}

func genExpectedOutput(t *testing.T, level string, caller bool) string {
	sb := &strings.Builder{}
	sb.WriteString(`"`)

	c := ""
	if caller {
		c = "_caller"
	}

	levels := []string{"trace", "debug", "info", "warn", "error"} //, "fatal", "panic"}

	on := false
	for _, l := range levels {
		if l == level {
			on = true
		}
		if on {
			t.Logf("Adding %s to expected outputs for %s%s", expectedOutputs[l+c], level, c)
			sb.WriteString(expectedOutputs[l+c])
		}
	}
	sb.WriteString(`"`)
	return sb.String()
}

func testFlagVerboseAndV(t *testing.T, f string) {
	var verboseExpectedOutput string = genExpectedOutput(t, "debug", false) //[]string{"debug", "info", "warn", "error", "fatal", "panic"}, false)
	testFlags(t, []string{f}, verboseExpectedOutput, func(lg *LoggerGenerator) (bool, string) {
		if !lg.verbose || lg.debug || lg.trace || lg.quiet || lg.silent {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
func TestFlagVerbose(t *testing.T) {
	testFlagVerboseAndV(t, "-verbose")
}
func TestFlagV(t *testing.T) {
	testFlagVerboseAndV(t, "-v")
}

func TestFlagDebug(t *testing.T) {
	var expectedOutput string = genExpectedOutput(t, "debug", true) //[]string{"debug", "info", "warn", "error", "fatal", "panic"}, false)
	testFlags(t, []string{"-debug"}, expectedOutput, func(lg *LoggerGenerator) (bool, string) {
		if !lg.debug || lg.verbose || lg.trace || lg.quiet || lg.silent {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
func TestFlagTrace(t *testing.T) {
	var expectedOutput string = genExpectedOutput(t, "trace", true) //[]string{"debug", "info", "warn", "error", "fatal", "panic"}, false)
	testFlags(t, []string{"-trace"}, expectedOutput, func(lg *LoggerGenerator) (bool, string) {
		if !lg.trace || lg.verbose || lg.debug || lg.quiet || lg.silent {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
func TestFlagQuiet(t *testing.T) {
	var expectedOutput string = genExpectedOutput(t, "warn", false) //[]string{"debug", "info", "warn", "error", "fatal", "panic"}, false)
	testFlags(t, []string{"-quiet"}, expectedOutput, func(lg *LoggerGenerator) (bool, string) {
		if !lg.quiet || lg.verbose || lg.trace || lg.silent || lg.debug {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
func TestFlagSilent(t *testing.T) {
	testFlags(t, []string{"-silent"}, `""`, func(lg *LoggerGenerator) (bool, string) {
		if !lg.silent || lg.verbose || lg.trace || lg.quiet || lg.debug {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
func TestDebugSilent(t *testing.T) {
	testFlags(t, []string{"-debug", "-silent"}, `""`, func(lg *LoggerGenerator) (bool, string) {
		if !lg.silent || !lg.debug || lg.trace || lg.quiet || lg.verbose {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
func TestDebugQuiet(t *testing.T) {
	testFlags(t, []string{"-debug", "-quiet"}, genExpectedOutput(t, "warn", true), func(lg *LoggerGenerator) (bool, string) {
		if !lg.quiet || !lg.debug || lg.trace || lg.silent || lg.verbose {
			return false, "verbose not exclusively set when passing verbose flag"
		}
		return true, ""
	})
}
