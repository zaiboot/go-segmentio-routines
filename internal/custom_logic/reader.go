package custom_logic

import (
	"context"
	"zaiboot/segmentIO.tests/internal/configs"

	"github.com/rs/zerolog"
)

func transform1(message string, l *zerolog.Logger) (string, error) {
	l.Info().Msg("transforming1 message")
	return message, nil
}

func sendToKafka(message string, l *zerolog.Logger) {
	l.Info().Msg("sending message to kafka")
}
func sendToDB(message string, l *zerolog.Logger) {
	l.Info().Msg("sending message to db")
}

func DoWork(ctx context.Context, message string, l *zerolog.Logger, d Dedupe_storage, config *configs.Config) error {
	defer d.close()
	d.connect()
	subLogger := l.With().Str("message", message).Logger()
	k, err := getDuplicateKey(message)

	if err != nil {
		// this is just a message we can ignore since the key for dedupe cannot be found
		// some of the reasons might be:
		// * malformed message
		// * empty string
		subLogger.Err(err).Msg("Ignoring message")
		return nil
	}

	seen, err := d.seen(ctx, k, config)
	if err != nil {
		return err
	}

	if seen {
		subLogger.Warn().Msg("Already seen, ignoring")
		return nil
	}

	m, err := transform1(message, l)
	if err != nil {
		return err
	}
	sendToDB(m, l)
	sendToKafka(m, l)
	return nil
}
