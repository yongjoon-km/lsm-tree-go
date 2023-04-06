package memory

import (
	"fmt"
	"lsm-tree-go/internal/core/common"
	"lsm-tree-go/internal/core/disk"
)

type MemoryBuffer struct {
	capacity    int
	buffer      map[string]common.Data
	diskManager disk.DiskManager
}

func CreateNewMemoryBuffer(capacity int, diskManager *disk.DiskManager) (*MemoryBuffer, error) {
	if capacity < 1 {
		return nil, fmt.Errorf("capacity should be greater than or equals to 1")
	}
	if diskManager == nil {
		return nil, fmt.Errorf("diskManger should not be nil")
	}

	return &MemoryBuffer{capacity: capacity, buffer: make(map[string]common.Data), diskManager: *diskManager}, nil
}

func (memoryBuffer *MemoryBuffer) Find(key string) (*common.Data, error) {
	data, found := memoryBuffer.buffer[key]
	if !found {
		return nil, fmt.Errorf("Can't find the key in memory.")
	}
	if data.Value == "" {
		return nil, fmt.Errorf("The key is removed. Tombstone found.")
	}
	return &data, nil
}

func (memoryBuffer *MemoryBuffer) Insert(key string, value string) error {
	return memoryBuffer.createEntityInBuffer(common.CreateData(key, value))
}

func (memoryBuffer *MemoryBuffer) Delete(key string) error {
	return memoryBuffer.createEntityInBuffer(common.CreateData(key, ""))
}

func (memoryBuffer *MemoryBuffer) createEntityInBuffer(data common.Data) error {
	if len(memoryBuffer.buffer) > memoryBuffer.capacity {
		// TODO: lock the memory buffer
		memoryBuffer.flushToDisk()
		memoryBuffer.clearBuffer()
	}
	memoryBuffer.buffer[data.Key] = data
	return nil
}

func (memoryBuffer *MemoryBuffer) flushToDisk() error {
	memoryBuffer.diskManager.CreateSSTable(memoryBuffer.buffer)
	return nil
}

func (memoryBuffer *MemoryBuffer) clearBuffer() {
	memoryBuffer.buffer = make(map[string]common.Data)
}
