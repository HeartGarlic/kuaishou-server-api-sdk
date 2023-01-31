package cache

import (
	"sync"
	"time"
)

// Cache interface
type Cache interface {
	Get(key string) interface{}
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) bool
	Delete(key string) error
}

// data 存储数据用的
type data struct {
	Data    interface{}
	Expired time.Time
}

// Memory 实现一个内存缓存
type Memory struct {
	sync.Mutex // 读写锁
	data       map[string]*data
}

// NewMemory 实例化一个内存缓存器
func NewMemory() Cache {
	return &Memory{
		data: map[string]*data{},
	}
}

// Get 获取缓存的值
func (mem *Memory) Get(key string) interface{} {
	if val, ok := mem.data[key]; ok {
		// 判断缓存是否过期
		if val.Expired.Before(time.Now()) {
			// 删除这个key
			mem.deleteKey(key)
			return nil
		}
		return val.Data
	}
	return nil
}

// Set 设置一个值
func (mem *Memory) Set(key string, val interface{}, timeout time.Duration) error {
	mem.Lock()
	defer mem.Unlock()
	mem.data[key] = &data{
		Data:    val,
		Expired: time.Now().Add(timeout),
	}
	return nil
}

// IsExist 判断值是否存在
func (mem *Memory) IsExist(key string) bool {
	if val, ok := mem.data[key]; ok {
		if val.Expired.Before(time.Now()) {
			return false
		}
		return true
	}
	return false
}

// Delete 删除一个值
func (mem *Memory) Delete(key string) error {
	mem.deleteKey(key)
	return nil
}

// deleteKey 删除一个缓存key
func (mem *Memory) deleteKey(key string) {
	mem.Lock()
	defer mem.Unlock()
	delete(mem.data, key)
	return
}
