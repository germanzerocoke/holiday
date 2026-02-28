package service

func (s *Service) SaveNewMemberId(idRaw []byte) error {
	err := s.repository.SaveNewMemberId(idRaw)
	if err != nil {
		return err
	}
	return nil
}
