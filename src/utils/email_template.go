package utils

import "fmt"

func BuildAppointmentEmailBody(date, startTime, endTime string) string {
	return fmt.Sprintf(`
	<h2>Confirmação de Agendamento</h2>
	<p>Seu atendimento foi agendado com sucesso!</p>
	<p><strong>Data:</strong> %s</p>
	<p><strong>Início:</strong> %s</p>
	<p><strong>Término:</strong> %s</p>
	`, date, startTime, endTime)
}

func BuildAppointmentCancellationEmailBody(date, startTime, endTime string) string {
	return fmt.Sprintf(`
	<h2>Cancelamento de Agendamento</h2>
	<p>Seu atendimento foi cancelado.</p>
	<p><strong>Data:</strong> %s</p>
	<p><strong>Início:</strong> %s</p>
	<p><strong>Término:</strong> %s</p>
	`, date, startTime, endTime)
}

func BuildAppointmentUpdateEmailBody(oldDate, oldStartTime, oldEndTime, newDate, newStartTime, newEndTime string) string {
	return fmt.Sprintf(`
	<h2>Alteração de Agendamento</h2>

	<p>Seu atendimento foi remarcado.</p>

	<h3>Horário anterior</h3>
	<p><strong>Data:</strong> %s</p>
	<p><strong>Início:</strong> %s</p>
	<p><strong>Término:</strong> %s</p>

	<h3>Novo horário</h3>
	<p><strong>Data:</strong> %s</p>
	<p><strong>Início:</strong> %s</p>
	<p><strong>Término:</strong> %s</p>
	`, oldDate, oldStartTime, oldEndTime, newDate, newStartTime, newEndTime)
}