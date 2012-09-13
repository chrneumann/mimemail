/* Construct text/plain MIME messages for use with net/smtp.
 * 
 * Base64 is used as transfer encoding and utf-8 as charset.
 */
package mimemail

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// A recipient or sender address.
type Address struct {
	Name  string
	Email string
}

// Format address for use in headers.
func (a Address) Format() string {
	return fmt.Sprintf("\"%v\" <%v>", hencode(a.Name), a.Email)
}

type Mail struct {
	From    Address
	To      []Address
	Cc      []Address
	Bcc     []Address
	Subject string
	Body    []byte
}

// Get sender address without name.
func (m Mail) Sender() string {
	return m.From.Email
}

// Get recipient addresses without names.
func (m Mail) Recipients() (rcpts []string) {
	for _, addr := range append(m.To, append(m.Cc, m.Bcc...)...) {
		rcpts = append(rcpts, addr.Email)
	}
	return
}

func encode(x string) string {
	return base64.StdEncoding.EncodeToString([]byte(x))
}

func hencode(x string) string {
	return fmt.Sprintf("=?utf-8?b?%v?=", encode(x))
}

// Get formatted MIME message (headers + body).
func (m Mail) Message() (msg []byte) {
	msg = []byte(`Content-Type: text/plain; charset="utf-8"
MIME-Version: 1.0
Content-Transfer-Encoding: base64
`)
	msg = append(msg, []byte(fmt.Sprintf("From: %v\n", m.From.Format()))...)
	var tos, ccs, bccs []string
	for _, addr := range m.To {
		tos = append(tos, addr.Format())
	}
	for _, addr := range m.Cc {
		ccs = append(ccs, addr.Format())
	}
	for _, addr := range m.Bcc {
		bccs = append(bccs, addr.Format())
	}
	msg = append(msg, []byte(fmt.Sprintf("To: %v\n",
		strings.Join(tos, ", ")))...)
	if len(ccs) > 1 {
		msg = append(msg, []byte(fmt.Sprintf("Cc: %v\n",
			strings.Join(ccs, ", ")))...)
	}
	if len(bccs) > 1 {
		msg = append(msg, []byte(fmt.Sprintf("Bcc: %v\n",
			strings.Join(bccs, ",")))...)
	}
	msg = append(msg, []byte(fmt.Sprintf("Subject: %v\n\n%v\n",
		hencode(m.Subject), encode(string(m.Body))))...)
	return
}
