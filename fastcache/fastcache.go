package fastcache

type Cache interface {
	Del(k []byte)
	Get(dst []byte, k []byte) []byte
	GetBig(dst []byte, k []byte) (r []byte)
	Has(k []byte) bool
	HasGet(dst []byte, k []byte) ([]byte, bool)
	Reset()
	SaveToFile(filePath string) error
	SaveToFileConcurrent(filePath string, concurrency int) error
	Set(k []byte, v []byte)
	SetBig(k []byte, v []byte)
	// UpdateStats(s *Stats)
}
