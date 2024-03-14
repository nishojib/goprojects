package pocketlog

import "io"

// Option defines a functional option to our logger.
type Option func(*Logger)

// WithOutput returns a configuration function that sets the output of logs.
func WithOutput(output io.Writer) Option {
	return func(lgr *Logger) {
		lgr.output = output
	}
}

// WithMaxMessageLength sets the maximum length, in characters, of a message.
// Use 0 for no maximum length.
func WithMaxMessageLength(maxMessageLength uint) Option {
	return func(lgr *Logger) {
		lgr.maxMessageLength = maxMessageLength
	}
}
