package binr

import (
	"fmt"
)

type sizeOpts struct {
	bump  uint64
	units []string
}

var (
	SizeOptBinary *sizeOpts = &sizeOpts{
		bump:  1024,
		units: []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"},
	}
	SizeOptDecimal *sizeOpts = &sizeOpts{
		bump:  1000,
		units: []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"},
	}
)

func defaultSizeOpts() *sizeOpts {
	return SizeOptBinary
}

func HumanSize(s uint64, o *sizeOpts) string {
	if o == nil {
		o = defaultSizeOpts()
	}

	i := 0
	if s < o.bump {
		return fmt.Sprintf("%d%s", s, o.units[i])
	}

	size := float64(s)
	for size >= float64(o.bump) {
		size = size / float64(o.bump)
		i++
	}

	return fmt.Sprintf("%.2f%s", size, o.units[i])
}
