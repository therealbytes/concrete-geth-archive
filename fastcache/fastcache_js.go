//go:build js
// +build js

package fastcache

type fastcache_cache_js struct{}

func (c *fastcache_cache_js) Del(k []byte) {
	return
}

func (c *fastcache_cache_js) Get(dst []byte, k []byte) []byte {
	return nil
}

func (c *fastcache_cache_js) GetBig(dst []byte, k []byte) (r []byte) {
	return nil
}

func (c *fastcache_cache_js) Has(k []byte) bool {
	return false
}

func (c *fastcache_cache_js) HasGet(dst []byte, k []byte) ([]byte, bool) {
	return nil, false
}

func (c *fastcache_cache_js) Reset() {
	return
}

func (c *fastcache_cache_js) SaveToFile(filePath string) error {
	return nil
}

func (c *fastcache_cache_js) SaveToFileConcurrent(filePath string, concurrency int) error {
	return nil
}

func (c *fastcache_cache_js) Set(k []byte, v []byte) {
	return
}

func (c *fastcache_cache_js) SetBig(k []byte, v []byte) {
	return
}

func New(maxBytes int) Cache {
	return &fastcache_cache_js{}
}

func LoadFromFileOrNew(path string, maxBytes int) Cache {
	return &fastcache_cache_js{}
}
