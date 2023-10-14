package custom_logic

import "github.com/rs/zerolog"

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
func DoWork(message string, l *zerolog.Logger) error {
	l.Info().Msg("Working with message")
	m, err := transform1(message, l)
	if err != nil {
		return err
	}
	sendToDB(m, l)
	sendToKafka(m, l)
	return nil
}
