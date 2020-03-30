package utils

import (
	"fmt"
	"time"
)

// thank you stack overflow user: https://stackoverflow.com/a/23695774
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}
