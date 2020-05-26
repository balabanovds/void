package service

func (s *Service) debugLog(err error) {
	s.log.Debug().Caller(1).Err(err).Send()
}
