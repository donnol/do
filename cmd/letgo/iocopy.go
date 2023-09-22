package main

import (
	"errors"
	"fmt"
	"io"
)

// Copy from io.Copy

// errInvalidWrite means that a write returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	// if wt, ok := src.(io.WriterTo); ok {
	// 	println("write to")
	// 	return wt.WriteTo(dst)
	// }
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	// if rt, ok := dst.(io.ReaderFrom); ok { // NOTE: TCP conn 是实现了ReaderFrom接口的，里面会优先尝试使用splice和sendfile系统调用来读取数据（减少副本）；这里我们为了能把读到的数据打印出来，就先注释掉了
	// 	println("read from")
	// 	return rt.ReadFrom(src)
	// }
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			data := buf[0:nr]

			fmt.Printf("=== read data:\n%s\n===\n", data)

			nw, ew := dst.Write(data)
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
