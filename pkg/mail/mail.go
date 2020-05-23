package mail

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type Request struct {
	auth    *noAuth
	config  *Config
	to      []string
	subject string
	body    string
}

func NewRequest(config *Config, to []string, subject string) *Request {
	return &Request{
		auth:    &noAuth{},
		config:  config,
		to:      to,
		subject: subject,
	}
}
