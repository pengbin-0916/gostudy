package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	//分析它应该有哪些字段
	Conn net.Conn
	Buf  [8096]byte //这是传输时，使用缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据...")
	//conn.Read在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了conn就不会阻塞.
	_, err = this.Conn.Read(this.Buf[:4])
	//客户端关闭后，返回的err是EOF
	//fmt.Println(err)
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[0:4]转成一个unit32的类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}
	//把pkgLen反序列化成message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes) //注意是&mes，这里是个坑
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}
func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32 //uint32足够大了，能传4G左右大小的数据
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}
