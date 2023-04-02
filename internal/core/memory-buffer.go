package core

import "fmt"

type MemoryBuffer struct {
	capacity    int
	buffer      map[string]Data
	diskManager DiskManager
}

func CreateNewMemoryBuffer(capacity int, diskManager *DiskManager) (*MemoryBuffer, error) {
	if capacity < 1 {
		return nil, fmt.Errorf("capacity should be greater than or equals to 1")
	}
	if diskManager == nil {
		return nil, fmt.Errorf("diskManger should not be nil")
	}

	return &MemoryBuffer{capacity: capacity, buffer: make(map[string]Data), diskManager: *diskManager}, nil
}

func (memoryBuffer *MemoryBuffer) Find(key string) (*Data, error) {
	data, found := memoryBuffer.buffer[key]
	if !found {
		return nil, fmt.Errorf("Can't find the key in memory.")
	}
	if data.value == "" {
		return nil, fmt.Errorf("The key is removed. Tombstone found.")
	}
	return &data, nil
}

func (memoryBuffer *MemoryBuffer) Insert(key string, value string) error {
	return memoryBuffer.createEntityInBuffer(CreateData(key, value))
}

func (memoryBuffer *MemoryBuffer) Delete(key string) error {
	return memoryBuffer.createEntityInBuffer(CreateData(key, ""))
}

func (memoryBuffer *MemoryBuffer) createEntityInBuffer(data Data) error {
	if len(memoryBuffer.buffer) > memoryBuffer.capacity {
		// TODO: lock the memory buffer
		memoryBuffer.flushToDisk()
		memoryBuffer.clearBuffer()
	}
	memoryBuffer.buffer[data.key] = data
	return nil
}

func (memoryBuffer *MemoryBuffer) flushToDisk() error {
	memoryBuffer.diskManager.CreateSSTable(memoryBuffer.buffer)
	return nil
}

func (memoryBuffer *MemoryBuffer) clearBuffer() {
	memoryBuffer.buffer = make(map[string]Data)
}
