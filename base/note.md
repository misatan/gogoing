# GO-Base



## Slice

slice是由动态数组的概念而来的，本身不是动态数据或者数组指针，可以为开发者提供更加方便增减数据的数据结构，同时slice也具有索引，可迭代的性质

### 数据结构

切片内部通过指针引用底层数组，设定相关属性将数据读写操作限定在指定区域内。本身是一个**只读对象**，工作机制类似数组指针封装

**切片是对数组一个连续片段的引用**，所以切片是一个引用类型。这个片段可以是整个数组，或者是数组连续子集。切片提供了一个与指向数组的动态窗口

和数组最大的区别是，切片的长度可以在运行时修改，最小为0，最大为数组长度，即切片时一个长度可变的数组

数据结构定义如下：

```go
type slice struct{
    array unsafe.Pointer
    len int
    cap int
}
```

![image-20230424112531298](.\note.assets\image-20230424112531298.png)

Pointer 是指向数组的指针，len代表当前切片的长度，cap是当前切片的容量。cap >= len

![image-20230424112816261](.\note.assets\image-20230424112816261.png)



### 切片

`make(type,len,cap)`函数运行在运行期动态指定数组长度，绕开数组类型必须使用编译常量的限制

创建切片有两种方式：**make创建切片**，**空切片**

#### make和切片字面量

以下为go20.3 源码：

```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
    //mem:切片最大容量,overflow:mem是否内存溢出
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	//切片容量是否在，[0,maxAllow]中
    if overflow || mem > maxAlloc || len < 0 || len > cap {
        //比较切片的长度是否在[0,maxAlloc]中
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panic(errorString("makeslice: len out of range"))
		}
		panic(errorString("makeslice: cap out of range"))
	}
	//申请内存地址,返回切片首地址
	return mallocgc(mem, et, true)
}
```

int64版本

```go
func makeslice64(et *_type, len64, cap64 int64) unsafe.Pointer {
	len := int(len64)
	if int64(len) != len64 {
		panicmakeslicelen()
	}

	cap := int(cap64)
	if int64(cap) != cap64 {
		panicmakeslicecap()
	}

	return makeslice(et, len, cap)
}
```

原理一致，只是多一个转int64的操作

![image-20230424133550205](.\note.assets\image-20230424133550205.png)

上图用make函数创建了一个len=4，cap=6的切片。内存空间申请了6个int类型的内存大小。由于len = 4，所以之后2个暂时访问不到，但容量还是在的。数组指针指向数组中内容为初始化默认值0。

除了make函数可以创建切片外，字面量也可以创建切片。

![image-20230424133848067](.\note.assets\image-20230424133848067.png)

上图用字面量创建了一个len = 6，cap = 6的切片，数组中的每个元素都初始化过了。**需要注意的是**：`[]`中不要写具体容量，因为写了就是数组了，要做区分

![image-20230424134244629](.\note.assets\image-20230424134244629.png)

还有一种字面量创建切片的方法（Array[low:hign,max]，len = hign-low，cap = max-low）。如上图所示，array[2:5:5]：SliceA 创建出了一个 len=3，cap=3的切片；从Array的第二位元素（0是第一位）开始，一直到第四位元素（不包括第五位）。同理array[1:3:5]：SliceB创建出了一个len=3，cap=4的切片；从Array的第一位元素开始，一直到第二位元素。



#### nil切片和空切片

> nil切片

```go
var slice []int
```

![image-20230424135111654](.\note.assets\image-20230424135111654.png)

nil切片用于很多标准库和内置函数中，当描述不存在的切片时，就会使用。比如函数发生异常时，返回的切片就是nil切片。nil切片的指针指向nil。

> 空切片

```go
slice := make([]int,0)
slice := []int{}
```

![image-20230424135320358](.\note.assets\image-20230424135320358.png)

空切片一般会用于表示一个空的集合。比如数据库查询，一条结果也没查到，那么就可以返回一个空切片



**nil切片 和 空切片 区别：**

- 空切片指向不是nil，指向的是一个内存地址，但是它没有分配任何内存空间，即底层元素包含0个元素





### 扩容

```go
//1.17源码
func growslice_17(et *_type, old slice, cap int) slice {
    //...
    
    newcap := old.cap
    // 两倍扩容
	doublecap := newcap + newcap
    // 新切片需要的容量大于两倍扩容的容量，则直接按照新切片需要的容量扩容
	if cap > doublecap {
		newcap = cap
	} else {
        // 原 slice 容量小于 1024 的时候，新 slice 容量按2倍扩容
		if old.cap < 1024 {
			newcap = doublecap
		} else { // 原 slice 容量超过 1024，进入一个循环，新 slice 容量每次增加原来的1.25倍。
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			if newcap <= 0 {
				newcap = cap
			}
		}
    }
    
    //...
}

//1.18源码
func growslice_18(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
    //...

	newcap := oldCap
	//两倍扩容
	doublecap := newcap + newcap
    // 新切片需要的容量大于两倍扩容的容量，则直接按照新切片需要的容量扩容
	if newLen > doublecap {
		newcap = newLen
	} else {
		const threshold = 256
        // 原 slice 容量小于 256 的时候，新 slice 容量按2倍扩容
		if oldCap < threshold {
			newcap = doublecap
        } else { // 原slice容量超过 256,进入一个循环，新slice容量每次增加 (oldCap + 3*threshold) / 4
			for 0 < newcap && newcap < newLen {
				newcap += (newcap + 3*threshold) / 4
			}
			if newcap <= 0 {
				newcap = newLen
			}
		}
	}
    
    //...
}
```



总结：

go 1.18之前：
$$
threshold = 1024 \\
newCap = \left\{\begin{array} {l}
cap  \qquad \qquad \quad if \space cap > oldcap*2 \\
oldcap*2 \qquad \quad if \space oldcap < threshold \\
oldcap + \frac{oldcap}{4} \quad if \space oldcap > threshold
\end{array} \right.
$$


go 1.18及之后：
$$
threshold = 256 \\
newCap = \left\{\begin{array}{l}
cap  \qquad \qquad \quad if \space cap > oldcap*2 \\
oldcap*2 \qquad \quad if \space oldcap < threshold \\
oldcap + \frac{oldcap + 3*threshold}{4} \quad if \space oldcap > threshold
\end{array} \right.
$$
在1.18后go语言优化了扩容策略，使得底层数组大小增长更加平滑：通过减小与之并固定增加一个常数（`(3*threshold)/4`）,使得扩容系数不会出现从2到1.25的突变。不同`oldcap`的扩容系数

| oldcap | 扩容系数 |
| ------ | -------- |
| 256    | 2.0      |
| 512    | 1.63     |
| 1024   | 1.44     |
| 2048   | 1.35     |
| 4096   | 1.30     |

随着容量增大，扩容系数越来越小，可以更好节省内存。并且取极限，扩容系数会变为1.25

扩容策略代码示例：

```go
package main

import "fmt"

func main() {
	s := []int{}
	for i := 0; i < 4098; i++ {
		var oldCap = cap(s)
		temp := s
		s = append(s, i)
		var newCap = cap(s)
		if oldCap != newCap {
			fmt.Printf("oldPointer:%p,newPoint:%p,oldCap:%d,newCap:%d\n", &temp, &s, oldCap, newCap)
		}
	}
}
/**
打印结果：
	oldPointer:0xc000008090,newPoint:0xc000008078,oldCap:0,newCap:1
    oldPointer:0xc0000080a8,newPoint:0xc000008078,oldCap:1,newCap:2      
    oldPointer:0xc0000080c0,newPoint:0xc000008078,oldCap:2,newCap:4      
    oldPointer:0xc0000080f0,newPoint:0xc000008078,oldCap:4,newCap:8      
    oldPointer:0xc000008150,newPoint:0xc000008078,oldCap:8,newCap:16     
    oldPointer:0xc000008210,newPoint:0xc000008078,oldCap:16,newCap:32    
    oldPointer:0xc000008390,newPoint:0xc000008078,oldCap:32,newCap:64    
    oldPointer:0xc000008690,newPoint:0xc000008078,oldCap:64,newCap:128   
    oldPointer:0xc000008c90,newPoint:0xc000008078,oldCap:128,newCap:256  
    oldPointer:0xc000009890,newPoint:0xc000008078,oldCap:256,newCap:512  
    oldPointer:0xc000131098,newPoint:0xc000008078,oldCap:512,newCap:848  
    oldPointer:0xc00013d020,newPoint:0xc000008078,oldCap:848,newCap:1280 
    oldPointer:0xc0001498a8,newPoint:0xc000008078,oldCap:1280,newCap:1792
    oldPointer:0xc00015a8b8,newPoint:0xc000008078,oldCap:1792,newCap:2560
    oldPointer:0xc0001690c8,newPoint:0xc000008078,oldCap:2560,newCap:3408
    oldPointer:0xc000182060,newPoint:0xc000008078,oldCap:3408,newCap:5120
*/
```



注意扩容时机： 期望len > cap

```go
package main

import "fmt"

func main() {
	var array = []int{10, 20, 30, 40, 50}
	var s1 = array[0:2]
	var s2 = s1
	fmt.Printf("Before s1 = %v,Pointer = %p,len = %d,cap = %d\n", s1, &s1, len(s1), cap(s1))
	fmt.Printf("Before s2 = %v,Pointer = %p,len = %d,cap = %d\n", s2, &s2, len(s2), cap(s2))
	s2 = append(s2, 50)
	s2[1] += 10
	fmt.Printf("After s1 = %v,Pointer = %p,len = %d,cap = %d\n", s1, &s1, len(s1), cap(s1))
	fmt.Printf("After s2 = %v,Pointer = %p,len = %d,cap = %d\n", s2, &s2, len(s2), cap(s2))
	fmt.Printf("After array = %v,Pointer = %p,len = %d,cap = %d\n", array, &array, len(array), cap(array))
}
/**
	打印结果：
		Before s1 = [10 20],Pointer = 0xc000008090,len = 2,cap = 5
        Before s2 = [10 20],Pointer = 0xc0000080a8,len = 2,cap = 5
        After s1 = [10 30],Pointer = 0xc000008090,len = 2,cap = 5
        After s2 = [10 30 50],Pointer = 0xc0000080a8,len = 3,cap = 5
        After array = [10 30 50 40 50],Pointer = 0xc000008078,len = 5,cap = 5
*/
```

以上示例有个**严重的问题**：当`len < cap`新切片(s2) 进行append操作时 或者 赋值操作时 会影响老切片(s1)和原数组(array)的数据；



### 拷贝

```go
//将无指针元素从字符串或切片复制到切片中
func slicecopy(toPtr unsafe.Pointer, toLen int, fromPtr unsafe.Pointer, fromLen int, width uintptr) int {
	if fromLen == 0 || toLen == 0 {
		return 0
	}

    //n 记录原切片 和 目标切片最短的len
	n := fromLen
	if toLen < n {
		n = toLen
	}

	if width == 0 {
		return n
	}

    
	size := uintptr(n) * width
    ...

    //只有一个元素,直接转换指针即可
	if size == 1 {
		*(*byte)(toPtr) = *(*byte)(fromPtr) 
    } else { //多个元素，则把size个 bytes 从fromPtr地址开始，copy到toPtr地址之后
		memmove(toPtr, fromPtr, size)
	}
	return n
}
```

拷贝就是把源切片值复制到目标切片中，并返回被复制的元素个数，source切片和target切片必须类型一致，当较短的切片复制完成，复制操作完成

![image-20230424154121234](.\note.assets\image-20230424154121234.png)



拷贝操作代码示例：

```go
package main

import "fmt"

func main() {
    test1() // 拷贝切片
    test2() // 拷贝字符串
}

func test1(){
    array := []int{10, 20, 30, 40}
	slice := make([]int, 6)
	n := copy(slice, array)
	fmt.Println(n, slice)
}

func test2(){
	slice := make([]int, 3)
	n := copy(slice, "abcdef")
	fmt.Println(n, slice)
}
/**
打印结果：
	4 [10 20 30 40 0 0]
	3 [97,98,99]
*/
```



**注意示例**：`for range 遍历 slice`

```go
func main() {
	slice := []int{10, 20, 30, 40}
	for index, value := range slice {
		fmt.Printf("value = %d , value-addr = %x , slice-addr = %x\n", value, &value, &slice[index])
	}
}
/**
	打印结果：
		value = 10 , value-addr = c4200aedf8 , slice-addr = c00014e020
        value = 20 , value-addr = c4200aedf8 , slice-addr = c00014e028
        value = 30 , value-addr = c4200aedf8 , slice-addr = c00014e030
        value = 40 , value-addr = c4200aedf8 , slice-addr = c00014e038
*/
```

`for range 遍历 slice` 其中的value只是slice具体元素的值拷贝，所以value地址不变

![image-20230424164708240](.\note.assets\image-20230424164708240.png)



### Slice与Array对比

GO中，数组是值类型数据结构

```go
func main() {
    arrayA := [2]int{100, 200}
    var arrayB [2]int

    arrayB = arrayA

    fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA)
    fmt.Printf("arrayB : %p , %v\n", &arrayB, arrayB)

    testArray(arrayA)
}

func testArray(x [2]int) {
    fmt.Printf("func Array : %p , %v\n", &x, x)
}

/**
打印结果:
    arrayA : 0xc4200bebf0 , [100 200]
    arrayB : 0xc4200bec00 , [100 200]
    func Array : 0xc4200bec30 , [100 200]
*/
```

由此三个内存地址不同可验证，**赋值**和**函数传参**操作都会**复制整个数组数据**

因此产生的问题：传参使用数组，如果数组大小有100w，在64位机器大约会消耗8M内存。这样会消耗大量内存存储重复的数据。

为此如何解决？

1. 函数传参使用数组指针

   ```go
   func main() {
       arrayA := []int{100, 200}
   	testArrayPoint(&arrayA) //1.传数组指针
       fmt.Printf("arrayA: %p,%v\n", &arrayA, arrayA)
   }
   
   func testArrayPoint(x *[]int) {
   	fmt.Printf("func Array: %p,%v\n", x, *x)
   	(*x)[1] += 100
   }
   /**
   打印结果:
       func Array: 0xc00010e060,[100 200]
   	arrayA: 0xc00010e060,[100 300]
   ```

   此方式，就算传入十亿的数组，只需在栈上分配8字节（64位指针变量占用）内存给指针即可。高效的你用内存，性能也大大提升

   不过此方式也有一个**共享内存的弊端**：假如arrayA指针更改，则会影响函数中的指针，使用切片传参不会有这个影响

2. 函数传参使用切片

   ```go
   func main() {
   	arrayA := []int{100, 200}
   	arrayB := arrayA[:]
   	testArrayPoint(&arrayB) //2.传切片
   	fmt.Printf("arrayA: %p,%v\n", &arrayA, arrayA)
   }
   
   func testArrayPoint(x *[]int) {
   	fmt.Printf("func Array: %p,%v\n", x, *x)
   	(*x)[1] += 100
   }
   /**
   打印结果:
       func Array: 0xc000092078,[100 200]
       arrayA: 0xc000092060,[100 300]
   */
   ```

   以上示例表现出，切片既可以达到**节省内存**的目的，还可以**处理好共享内存**的问题

   但并不是用切片替换数组就能够提升程序效率，示例如下：

   ```go
   package main
   
   import "testing"
   
   func array() [1024]int {
   	var x [1024]int
   	for i := 0; i < len(x); i++ {
   		x[i] = i
   	}
   	return x
   }
   
   func slice() []int {
   	x := make([]int, 1024)
   	for i := 0; i < len(x); i++ {
   		x[i] = i
   	}
   	return x
   }
   
   func BenchmarkArray(b *testing.B) {
   	for i := 0; i < b.N; i++ {
   		array()
   	}
   }
   
   func BenchmarkSlice(b *testing.B) {
   	for i := 0; i < b.N; i++ {
   		slice()
   	}
   }
   ```

   使用性能测试，禁用内联和优化，观察切片堆上内存分配情况

   ```bash
   go test -bench . -benchmem -gcflags "-N -l"
   ```

   输出结果：

   ```base
   goos: windows
   goarch: amd64
   cpu: Intel(R) Core(TM) i7-10700 CPU @ 2.90GHz
   BenchmarkArray-16         598572              2019 ns/op               0 B/op          0 allocs/op
   BenchmarkSlice-16         445610              2428 ns/op            8192 B/op          1 allocs/op
   ```

   输出解释：

   - （Array）16核CPU，循环：598572， 平均执行时间：2019ns，每次执行栈上分配内存总量：0，每次执行内存分配次数：0
   - （Slice）16核CPU，循环：445610，平均执行事件：2428ns，每次执行栈上分配内存总理：8192，每次执行内存分配次数：1

   如此看来，并非所有时候都适合用切片代替数组，因为**切片底层数组可能会在堆上分配内存，而且小数组在栈上拷贝消耗未必比make消耗大**





## Pointer

区别于C/C++中的指针，Go语言指针不能进行偏移和运算，是安全指针。

学习指针，需要弄清三个概念：指针地址，指针类型，指针取值



> 指针地址

程序运行时，变量都会在内存中拥有一个地址，这个地址代表变量在内存中的位置，指针地址就是变量的内存地址

> 指针类型

Go语言中的值类型`int、float、bool、string、array、struct`都有对应的指针类型，：如：`*int,*int64,*string`，指针类型给出解引用（引用指向变量的值）最多能够操作多少字节的信息

> 指针取值

指针取值就是解引用（引用指向变量的值）

> 演示示例

```go
```

