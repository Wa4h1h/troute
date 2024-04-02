package cli

import (
	"strconv"
)

type (
	uint8Flag  uint8
	uint16Flag uint16
)

func (u *uint8Flag) String() string {
	return strconv.FormatUint(uint64(*u), 10)
}

func (u *uint8Flag) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 8)
	*u = uint8Flag(v)

	return err
}

func (u *uint16Flag) String() string {
	return strconv.FormatUint(uint64(*u), 10)
}

func (u *uint16Flag) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 16)
	*u = uint16Flag(v)

	return err
}
