package mail

import "net/smtp"

type noAuth struct{}

func (na *noAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	return "", nil, nil
}

func (na *noAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	return nil, nil
}
