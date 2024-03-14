package merkledag

import (
	"encoding/hex"
	"hash"
)

func Add(store KVStore, node Node, h hash.Hash) []byte {
	switch node.Type() {
	case FILE:
		// 类型断言，确认node确实是File类型
		file, ok := node.(File)
		if !ok {
			return nil
		}
		// 获取文件内容并计算哈希值
		content := file.Bytes()
		h.Reset()
		h.Write(content)
		hashedContent := h.Sum(nil)

		// 将哈希值（作为key）和内容（作为value）保存到KVStore
		store.Put(hashedContent, content)
		return hashedContent
	case DIR:
		// 类型断言，确认node确实是Dir类型
		dir, ok := node.(Dir)
		if !ok {
			return nil
		}

		it := dir.It()
		var hashes [][]byte
		for it.Next() {
			// 递归调用Add处理子节点
			childHash := Add(store, it.Node(), h)
			if childHash != nil {
				hashes = append(hashes, childHash)
			}
		}

		// 对所有子节点的哈希值进行排序（如果需要）并计算总哈希
		h.Reset()
		for _, hash := range hashes {
			h.Write(hash)
		}
		dirHashed := h.Sum(nil)

		// 将目录的哈希值（作为key）和所有子哈希值（作为value）保存到KVStore
		// 注意：这里简化处理，实际应用中可能需要以某种形式存储子哈希值，例如拼接或使用结构化格式
		store.Put(dirHashed, bytes.Join(hashes, []byte{}))
		return dirHashed
	default:
		return nil
	}
}
 