package custom_logic

func transform1(message string) (string, error) {
	return message, nil
}

func sendToKafka(message string) {

}
func sendToDB(message string) {

}
func DoWork(message string) error {
	m, err := transform1(message)
	if err != nil {
		return err
	}
	sendToDB(m)
	sendToKafka(m)
	return nil
}
