package model

type SkippedScanner struct{}

func (SkippedScanner) Scan(interface{}) error {
	return nil
}
