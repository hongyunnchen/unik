package os

import (
	"github.com/emc-advanced-dev/pkg/errors"
	"fmt"
	"os"
)

type DiskSize interface {
	ToPartedFormat() string
	ToBytes() Bytes
}

type Bytes int64

func (s Bytes) ToPartedFormat() string {
	return fmt.Sprintf("%dB", uint64(s))
}

func (s Bytes) ToBytes() Bytes {
	return s
}

type MegaBytes int64

func (s MegaBytes) ToPartedFormat() string {
	return fmt.Sprintf("%dMiB", uint64(s))
}

func (s MegaBytes) ToBytes() Bytes {
	return Bytes(s << 20)
}

type GigaBytes int64

func (s GigaBytes) ToPartedFormat() string {
	return fmt.Sprintf("%dGiB", uint64(s))
}

func (s GigaBytes) ToBytes() Bytes {
	return Bytes(s << 30)
}

type Sectors int64

const SectorSize = 512

func (s Sectors) ToPartedFormat() string {
	return fmt.Sprintf("%ds", uint64(s))
}

func (s Sectors) ToBytes() Bytes {
	return Bytes(s * SectorSize)
}

func ToSectors(b DiskSize) (Sectors, error) {
	inBytes := b.ToBytes()
	if inBytes%SectorSize != 0 {
		return 0, errors.New("can't convert to sectors", nil)
	}
	return Sectors(inBytes / SectorSize), nil
}

type BlockDevice string

func (b BlockDevice) Name() string {
	return string(b)
}

type Partitioner interface {
	MakeTable() error
	MakePart(partType string, start, size DiskSize) error
}

type Resource interface {
	Acquire() (BlockDevice, error)
	Release() error
}

type Part interface {
	Resource

	Size() DiskSize
	Offset() DiskSize

	Get() BlockDevice
}

func IsExists(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}
