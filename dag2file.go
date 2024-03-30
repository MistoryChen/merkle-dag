package merkledag

// Hash2File 从 KVStore 中读取 hash 对应的数据，并根据 path 返回对应的文件内容
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 从 KVStore 中获取 hash 对应的数据
	data := store.Get(hash)
	if data == nil {
		return nil
	}

	// 根据 path 获取对应的文件内容
	return getFileContent(data, path, hp)
}

// getFileContent 根据路径从文件内容中获取对应文件的内容
func getFileContent(data []byte, path string, hp HashPool) []byte {
	// 假设 path 形如 "/dir1/dir2/file.txt"
	// 首先将 path 拆分成各级目录和文件名
	components := strings.Split(path, "/")

	// 从根节点开始逐级遍历路径
	for _, component := range components {
		// 如果当前节点是文件，则返回文件内容
		file, ok := hp.NewNode(data).(File)
		if ok {
			return file.Bytes()
		}

		// 如果当前节点是文件夹，则继续向下遍历
		dir, ok := hp.NewNode(data).(Dir)
		if ok {
			// 获取当前文件夹的迭代器
			it := dir.It()
			// 遍历文件夹中的每一个文件/文件夹
			for it.Next() {
				node := it.Node()
				// 如果节点的名称与当前路径组件匹配，则继续向下遍历
				if node.Name() == component {
					data = node.Bytes() // 更新当前节点的数据
					break
				}
			}
		} else {
			// 如果当前节点既不是文件也不是文件夹，则返回空
			return nil
		}
	}

	// 如果路径遍历完毕仍未找到对应文件，则返回空
	return nil
}
