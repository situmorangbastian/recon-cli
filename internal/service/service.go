package service

import "github.com/situmorangbastian/recon-cli/internal/reader"

type Service struct {
	reader *reader.Reader
}

func NewService(reader *reader.Reader) *Service {
	return &Service{
		reader: reader,
	}
}
