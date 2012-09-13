package mimemail

import "testing"

func TestSender(t *testing.T) {
	addr := "foo@example.com"
	message := Mail{
		From: Address{"Jogi Bär", addr}}
	if from := message.Sender(); from != addr {
		t.Errorf("message.From() = %v, want %v", from, addr)
	}
}

func TestRecipients(t *testing.T) {
	addrs := []Address{
		Address{"John Foo", "john@example.com"},
		Address{"Jane Bar", "jane@example.com"},
		Address{"Jack Crow", "jack@example.com"},
		Address{"Jö Shy", "joe@example.com"}}
	message := Mail{
		To:  addrs[:2],
		Cc:  addrs[2:3],
		Bcc: addrs[3:4]}
	recipients := message.Recipients()
	if n := len(recipients); n != 4 {
		t.Fatalf("len(message.Recipients()) = %v, want %v", n, 4)
	}
	for i, _ := range addrs {
		if addrs[i].Email != recipients[i] {
			t.Errorf("message.Recipients()[%v] = %v, want %v", i,
				recipients[i], addrs[i].Email)
		}
	}
}

func TestMessage(t *testing.T) {
	mail := Mail{
		From: Address{"The Dot", "dot@www.example.com"},
		To: []Address{
			Address{"Jogi Bär", "jogi@example.com"},
			Address{"John Doe", "doe@example.com"}},
		Subject: "The Ä is an umlaut",
		Body: []byte("Gotta go, it's to läit.")}
	valid := `Content-Type: text/plain; charset="utf-8"
MIME-Version: 1.0
Content-Transfer-Encoding: base64
From: "=?utf-8?b?VGhlIERvdA==?=" <dot@www.example.com>
To: "=?utf-8?b?Sm9naSBCw6Ry?=" <jogi@example.com>, "=?utf-8?b?Sm9obiBEb2U=?=" <doe@example.com>
Subject: =?utf-8?b?VGhlIMOEIGlzIGFuIHVtbGF1dA==?=

R290dGEgZ28sIGl0J3MgdG8gbMOkaXQu
`
	if msg := string(mail.Message()); msg != valid {
		t.Errorf("message.Message() = %v, want %v", msg, valid)
	}
}
