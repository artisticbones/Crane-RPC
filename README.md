## Crane-RPC

本文主要分三个小章节，首先是 RPC 到底是什么，其次就是到底为什么需要 RPC ，以及本文中最重要的章节，RPC 的具体实现。

##### 什么是 RPC (Remote Procedure Call)

RPC 是指远程过程调用，也就是说两台服务器 A，B，一个应用部署在 A 服务器上，该应用想要调用 B 服务器上应用提供的函数/方法,由于不在一个内存空间，不能直接调用，需要通过网络来表达调用的语义和传达调用的数据。

就比如该笔试题中出现的两个不同的函数，如果要完成调用那么：

* 首先，要解决通信问题，主要是通过在客户端和服务器之间建立 TCP 连接，远程过程调用的所有交换的数据都在这个连接中传输。这个连接可以是按需连接，调用结束后就断掉，也可以是长连接，过个远程过程调用共享同一个连接。
* 第二，要解决寻址的问题，也就是说，A 服务器上的应用怎么告诉底层的 RPC 框架，如何连接到 B 服务器以及特定的端口，方法的名称是什么，怎样才能完成调用。
* 第三，当 A 服务器上的应用发起远程过程调用时，方法的参数需要通过底层的网络协议(在本题中就是 TCP 协议) 传递到 B 服务器，这时就是出现问题，因为网络协议是基于二进制的，内存中的参数的值要序列化成二进制的形式，也就是序列化（Serialize）或编组（marshal），通过寻址和传输将序列化的二进制发送给B 服务器
* 第四，相应的 B 服务器收到请求后，需要对参数进行反序列化，恢复为内存中的表达方式，然后找到对应的方法进行本地调用，然后得到返回值
* 第五，得到的返回值还要发送回服务器 A 上的应用，同样需要经过序列化的方式发送，服务器 A 接到后再反序列化，恢复为内存中的表达方式，交给 A 服务器上的应用。

![preview](https://pic3.zhimg.com/45366c44f775abfd0ac3b43bccc1abc3_r.jpg?source=1940ef5c)

##### 为什么需要 RPC

根据上边的分析可以很简单的得出结论：两个应用无法在一个进程内，甚至一个计算机内通过本地调用的方式完成需求。比如不同的系统间的通讯，甚至不同的组织间的通讯。由于计算能力需要横向扩展，需要在多台机器组成的集群上部署应用。

##### Crane-RPC 实现方案

根据第一小节的分析，目前主要有五大需求：1. 数据包如何定义，以及按照什么方式解包或者加包；2. 客户端怎么连接服务器端，以及对连接的管理； 3. 客户端通过什么技术与服务器端通信；4. 服务器端怎么接收来自客户端的连接，以及对服务的注册的管理等；5. 服务器端采用什么办法从二进制中获取到客户端想要的服务。

1. 数据包的定义

由于只需要关心调用的服务，所以数据包可以简化为三个属性域：serviceName ，args ，error。具体如下：

```go
type Data struct {
	ServiceName string	// 服务名称
	Args []interface{}	// 传递的参数
	Err string			// socket 的错误
}
```

2. 数据的加包与解包

这里使用的包为 `gob` , `gob` 的介绍如下：

gob 包管理 gob 流——在编码器（发送器）和解码器（接受器）之间交换的binary值。一般用于传递远端程序调用（RPC）的参数和结果，如net/rpc包就有提供。

gob 的实现给每一个数据类型都编译生成一个编解码程序，当单个编码器用于传递数据流时，会分期偿还编译的消耗，是效率最高的。

Gob 流是自解码的。流中的所有数据都有前缀（采用一个预定义类型的集合）指明其类型。指针不会传递，而是传递值；也就是说数据是压平了的。递归的类型可以很好的工作，但是递归的值（比如说值内某个成员直接/间接指向该值）会出问题。

要使用 gob，先要创建一个编码器，并向其一共一系列数据：可以是值，也可以是指向实际存在数据的指针。编码器会确保所有必要的类型信息都被发送。在接收端，解码器从编码数据流中恢复数据并将它们填写进本地变量里。

**发送端和接收端的值/类型不需要严格匹配**。对结构体来说，字段（根据字段名识别）如果发送端有而接收端没有，会被忽略；接收端有而发送端没有的字段也会被忽略；发送端和接收端都有的字段其类型必须是可兼容的；发送端和接收端都会在 gob 流和实际 go 类型之间进行必要的指针取址/寻址工作。

结构体、数组和切片都被支持。结构体只编码和解码导出字段。字符串和byte数组/切片有专门的高效表示（参见下述）。当解码切片时，如果当前切片的容量足够会被复用，否则会申请新的底层数组（所以还是用切片地址为好）。此外，生成的切片的长度会修改为解码的成员的个数。

Gob流不支持函数和通道。试图在最顶层编码这些类型的值会导致失败。结构体中包含函数或者通道类型的字段的话，会视作非导出字段（忽略）处理。

Gob可以编码任意实现了 GobEncode r接口或者 encoding.BinaryMarshaler 接口的类型的值（通过调用对应的方法），GobEncoder 接口优先。

Gob 可以解码任意实现了 GobDecoder 接口或者 encoding.BinaryUnmarshaler 接口的类型的值（通过调用对应的方法），同样GobDecoder 接口优先。

所以对于数据包的序列化与反序列化的代码如下：

```go
func encode(data Data) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(b []byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	return data, nil
}
```

3. 客户端的部分

首先，要实现 TCP 连接，那么客户端也需要保存连接，除此之外，暂时不需要其他内容，所以客户端的结构体如下：

```go
type Client struct {
	conn net.Conn
}
```

接下来就是客户端的具体功能函数部分，客户端的主要功能：1. 实例化出一个 client；2. 调用 server 端函数。其中的调用函数是最重要的部分，接下来将讲解这个部分。
