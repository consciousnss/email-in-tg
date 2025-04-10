package tg

const emailTemplate = `✉️ <b>Новое письмо</b> ✉️
<b>От:</b> {{.MailFrom}}
<b>Кому:</b> {{.MailTo}}
<b>Тема:</b> {{.Subject}}
<b>Дата:</b> {{.Date}}

{{.Text}}`
