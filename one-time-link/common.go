package one_time_link

import "errors"

type SecretContent struct {
	ReadCount            int
	OwnerReadCount       int
	DeniedOwnerReadCount int
	Name                 string
	TextContent          string
	Id                   string
	OwnerAccessKey              string
	ContentType          string
	IsActive             bool
}

func (sc *SecretContent) GuestReadable() bool {
	if sc.ReadCount == 0 && sc.IsActive {
		return true
	} else {
		return false
	}
}

func (sc *SecretContent) ValidOwner(ownerKey string) bool {
	return ownerKey == sc.OwnerAccessKey
}

type storage interface {
	getSecretContent(id string) (*SecretContent, error)
	insertSecretContent(textContent string)
	incrementReadCount(id string)
	removeContent(id string)
	incrementOwnerReadCount(id string)
	incrementDeniedOwnerReadCount(id string)
	deactivateLink(id string)
}

type Service struct {
	storage storage
}

func (s Service) GetSecretContent(id string) (*SecretContent, error) {
	sc, err := s.storage.getSecretContent(id)
	if sc != nil {
		s.storage.incrementReadCount(sc.Id)
		if sc.GuestReadable() {
			s.storage.removeContent(id)
			return sc, nil
		} else {
			return nil, errors.New("sc already read")
		}
	}
	return nil, err
}

func (s Service) GetSecretContentAsOwner(id string,ownerKey string) (*SecretContent, error) {
	sc, err := s.storage.getSecretContent(id)
	if sc != nil {
		s.storage.incrementOwnerReadCount(sc.Id)
		if !sc.ValidOwner(ownerKey) {
			s.storage.incrementDeniedOwnerReadCount(id)
			return nil, errors.New("access denied")
		}
		s.storage.incrementOwnerReadCount(id)
		sc.TextContent = "hidden"
		return sc, nil

	}
	return nil, err
}


func (s Service) DeactivateLink(id string, ownerKey string) error {
	sc, err := s.storage.getSecretContent(id)
	if err != nil {
		return err
	}  
	if sc.ValidOwner(ownerKey) {
		s.storage.removeContent(id)
		s.storage.deactivateLink(id)
		return nil
	} else {
		return errors.New("access denied")
	}
}

func New() Service {
	return Service{}
}