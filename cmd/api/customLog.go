package main

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func levels() {
	// Panic of fatal messages stop the execution flow
	// log.Panic().Msg("This is a panic message")
	// log.Fatal().Msg("This is a fatal message")
	log.Error().Msg("This is an error message")
	log.Warn().Msg("This is a warning message")
	log.Info().Msg("This is an information message")
	log.Debug().Msg("This is a debug message")
	log.Trace().Msg("This is a trace message")
}

// Not every execution requires the same level of logging. For example, you can go from a Debug level to a less verbose level like Info.
func setGlobalLevel() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Debug().Msg("Debug message is displayed")
	log.Info().Msg("Info Message is displayed")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Debug().Msg("Degug message is no longer displayed")
	log.Info().Msg("Info message is displayed")
}

// Logging errors is similar to adding extra fields as shown above.
func logError() {
	err := errors.New("there as an error")

	log.Error().Err(err).Msg("this is the way to log errors")
}

// Subloggers#
//
// We can have several log instances running simultaneously.
// This is particularly useful when we deploy new components in our code
// that require additional information to be displayed.
// In this example, we have a sublogger that always displays the name of the component.
func sublogger() {

	mainLogger := zerolog.New(os.Stderr).With().Logger()
	mainLogger.Info().Msg("This is the output from the main logger")

	subLogger := mainLogger.With().Str("component", "componentA").Logger()
	subLogger.Info().Msg("This is the the extended output from the sublogger")
}

// File Output#
//
// The previous examples use the standard output. To store the output in a file,
// we only have to instantiate a logger with the descriptor of a new file.
func logOutputToFile() {
	// // create a temp file
	// // tempFile, err := ioutil.TempFile(os.TempDir(), "deleteme")
	// tempFile, err := os.CreateTemp(os.TempDir(), "deleteme")
	// if err != nil {
	// 	// Can we log an error before we have our logger? :)
	// 	log.Error().Err(err).Msg("there was an error creating a temporary file four our log")
	// }
	// fileLogger := zerolog.New(tempFile).With().Logger()
	// fileLogger.Info().Msg("This is an entry from my log")
	//
	// fmt.Printf("The log file is allocated at %s\n", tempFile.Name())

	// First, set log level for zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Then, mention the file or path to the file to store logs
	logFile, _ := os.OpenFile(
		"restApiWithPagination.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	// Then, set up the multi logger with multiwriter. Here, one looger for file another for std output(console)
	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
	// Then, include timestamp with the loggers
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	log.Info().Msg("Application is starting...")
}

// Pretty Logging#
//
// Zerolog offers a decorator with a more visual output specially designed for consoles.
// I do recommend this option for command line interfaces (CLI).
func prettyConsole() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Error().Msg("This is an error message")
	log.Warn().Msg("This is a warning message")
	log.Info().Msg("This is an information message")
	log.Debug().Msg("This is a debug message")
}

// Reference:https://jmtirado.net/338-2/
