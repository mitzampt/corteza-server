package expr

import (
	"fmt"
	"github.com/spf13/cast"
	"time"
)

func SetIDWithPath(dst *uint64, val interface{}, pp ...string) error {
	if len(pp) > 0 {
		return fmt.Errorf("can not set ID with path")
	}

	if tmp, err := cast.ToUint64E(val); err != nil {
		return err
	} else {
		*dst = tmp
	}

	return nil
}

func SetTimeWithPath(dst *time.Time, val interface{}, pp ...string) error {
	if len(pp) > 0 {
		return fmt.Errorf("can not set time with path")
	}

	if tmp, err := cast.ToTimeE(val); err != nil {
		return err
	} else {
		*dst = tmp
	}

	return nil
}

func SetKVWithPath(dst *map[string]string, val interface{}, pp ...string) error {
	switch len(pp) {
	case 0:
		tmp, err := cast.ToStringMapStringE(val)
		if err != nil {
			return err
		}

		*dst = tmp
	case 1:
		tmp, err := cast.ToStringE(val)
		if err != nil {
			return err
		}

		(*dst)[pp[0]] = tmp
	default:
		return fmt.Errorf("can not set KV with path deeper than 1 level")
	}

	return nil
}
