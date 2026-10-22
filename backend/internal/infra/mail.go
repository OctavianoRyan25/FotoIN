package infra

import (
	"fmt"
	"net/smtp"
	"strings"

	"fotoin/backend/internal/config"
)

type MailClient struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewMailClient(cfg config.SMTPConfig) *MailClient {
	mailClient := &MailClient{
		host: cfg.Host,
		port: cfg.Port,
		user: cfg.Username,
		pass: cfg.Password,
		from: cfg.From,
	}

	return mailClient
}

type DownloadLinkData struct {
	BuyerEmail   string
	GalleryTitle string
	OrderID      string
	Items        []DownloadItem
	ExpiresAt    string
}

type DownloadItem struct {
	Filename    string
	DownloadURL string
}

// func (m *MailClient) SendDownloadLinks(data DownloadLinkData) error {
// 	itemsHTML := ""
// 	for i, item := range data.Items {
// 		itemsHTML += fmt.Sprintf(`
// 			<tr>
// 				<td style="padding:8px 12px;border-bottom:1px solid #eee;">%d.</td>
// 				<td style="padding:8px 12px;border-bottom:1px solid #eee;">%s</td>
// 				<td style="padding:8px 12px;border-bottom:1px solid #eee;">
// 					<a href="%s" style="color:#6366f1;text-decoration:none;font-weight:600;">⬇ Download</a>
// 				</td>
// 			</tr>`, i+1, item.Filename, item.DownloadURL)
// 	}

// 	body := fmt.Sprintf(`<!DOCTYPE html>
// <html>
// <head><meta charset="UTF-8"></head>
// <body style="font-family:'Inter',Arial,sans-serif;background:#f9fafb;margin:0;padding:0;">
//   <div style="max-width:600px;margin:40px auto;background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,0.08);">
//     <div style="background:linear-gradient(135deg,#6366f1,#8b5cf6);padding:32px;text-align:center;">
//       <h1 style="color:#fff;margin:0;font-size:28px;letter-spacing:-0.5px;">📸 FotoIN</h1>
//       <p style="color:rgba(255,255,255,0.85);margin:8px 0 0;">Foto Anda Siap Diunduh!</p>
//     </div>
//     <div style="padding:32px;">
//       <p style="color:#374151;font-size:16px;">Halo,</p>
//       <p style="color:#374151;font-size:15px;">Terima kasih telah membeli foto dari galeri <strong>%s</strong>. Berikut adalah link unduhan untuk foto-foto Anda:</p>

//       <div style="background:#f3f4f6;border-radius:8px;padding:4px;margin:24px 0;">
//         <table style="width:100%%;border-collapse:collapse;">
//           <thead>
//             <tr style="background:#e5e7eb;">
//               <th style="padding:10px 12px;text-align:left;color:#6b7280;font-size:12px;font-weight:600;text-transform:uppercase;">#</th>
//               <th style="padding:10px 12px;text-align:left;color:#6b7280;font-size:12px;font-weight:600;text-transform:uppercase;">File</th>
//               <th style="padding:10px 12px;text-align:left;color:#6b7280;font-size:12px;font-weight:600;text-transform:uppercase;">Aksi</th>
//             </tr>
//           </thead>
//           <tbody>%s</tbody>
//         </table>
//       </div>

//       <div style="background:#fef3c7;border:1px solid #f59e0b;border-radius:8px;padding:16px;margin:16px 0;">
//         <p style="color:#92400e;margin:0;font-size:14px;">⏰ <strong>Perhatian:</strong> Link unduhan ini akan kedaluwarsa pada <strong>%s</strong>. Segera unduh foto Anda!</p>
//       </div>

//       <p style="color:#6b7280;font-size:13px;">ID Order: <code style="background:#f3f4f6;padding:2px 6px;border-radius:4px;">%s</code></p>
//     </div>
//     <div style="background:#f9fafb;padding:20px 32px;text-align:center;border-top:1px solid #e5e7eb;">
//       <p style="color:#9ca3af;font-size:13px;margin:0;">© 2024 FotoIN · Platform Fotografi Event Indonesia</p>
//       <p style="color:#9ca3af;font-size:12px;margin:4px 0 0;">Jika ada pertanyaan, balas email ini atau hubungi support kami.</p>
//     </div>
//   </div>
// </body>
// </html>`, data.GalleryTitle, itemsHTML, data.ExpiresAt, data.OrderID)

// 	return smtp.SendMail(data.BuyerEmail, fmt.Sprintf("📸 Foto Anda dari %s Siap Diunduh - FotoIN", data.GalleryTitle), body)
// }

func (m *MailClient) SendOTP(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", m.from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	return smtp.SendMail(addr, auth, m.user, []string{to}, []byte(msg))
}
