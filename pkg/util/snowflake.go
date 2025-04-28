package util

import (
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

func InitSonyflake() *sonyflake.Sonyflake {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Date(2021, 7, 28, 0, 0, 0, 0, time.UTC),
	})

	if sf == nil {
		panic("sonyflake not created")
	}

	return sf
}

func GetUniqueId(s *sonyflake.Sonyflake) (string, error) {
	r, err := s.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(r, 10), nil
}
