package fs

// inode提前生成，占好位值，虚拟文件路径按照inode大小计算hash值，可以快速路由到指定inode位置，inode指向一个文件起始位置
// hash函数

type FileHead struct {
}

type FileBlock struct {
	length int64
}
