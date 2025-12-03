package utils

import "fmt"

func BuildAppointmentEmailBody(date, startTime, endTime string) string {
	return fmt.Sprintf(`
	<h2>Confirmação de Agendamento<h2>
	<p>Seu atendimento foi agendado com sucesso!</p>
	<p><strong>Data:</strong> %s</p>
	<p><strong>Início:</strong> %s</p>
	<p><strong>Término:</strong> %s</p>
	`, date, startTime, endTime)
}