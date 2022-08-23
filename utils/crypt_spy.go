package utils

type GenerateHashSpy struct {
	GenerateResp []byte
	GenerateErr  error
}

func (c *GenerateHashSpy) GenerateHash(password string) ([]byte, error) {
	return c.GenerateResp, c.GenerateErr
}
