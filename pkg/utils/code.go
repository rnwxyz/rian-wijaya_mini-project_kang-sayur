package utils

import (
	"math/rand"
	"time"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
)

type code struct {
	length int
}

func NewCode() *code {
	return &code{
		length: constants.Item_code_length,
	}
}

func (c *code) GenerateCode() string {
	rand.Seed(time.Now().Unix())

	ran_str := make([]byte, c.length)

	// Generating Random string
	for i := 0; i < c.length; i++ {
		ran_str[i] = byte(48 + rand.Intn(10))
	}
	str := string(ran_str)
	return str
}
