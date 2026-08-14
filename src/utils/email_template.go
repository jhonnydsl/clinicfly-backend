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

func BuildPasswordResetEmailBody(token string) string {
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)

	return fmt.Sprintf(`
	<h2>Recuperação de senha</h2>

	<p>Recebemos uma solicitação para redefinir sua senha.</p>

	<p>
		Clique no botão abaixo para criar uma nova senha:
	</p>

	<p>
		<a href="%s"
			style="
				display:inline-block;
				padding:10px 20px;
				background-color:#007bff;
				color:#ffffff;
				text-decoration:none;
				border-radius:5px;
			">
			Redefinir senha
		</a>
	</p>

	<p>Este link é válido por <strong>30 minutos</strong>.</p>

	<p>
		Se você não solicitou a recuperação de senha,
		ignore esre e-mail.
	</p>
	`, resetLink)
}