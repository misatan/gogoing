# GO-BASE



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

### 基本概念

> 指针地址

程序运行时，变量都会在内存中拥有一个地址，这个地址代表变量在内存中的位置，指针地址就是变量的内存地址



> 指针类型

Go语言中的值类型`int、float、bool、string、array、struct`都有对应的指针类型，：如：`*int,*int64,*string`，指针类型给出解引用（引用指向变量的值）最多能够操作多少字节的信息



> 指针取值

指针取值就是解引用（引用指向变量的值）



> 演示示例

指针地址、指针类型、指针取值示例：

```go
func main() {
	a := 10
	b := &a                         
	fmt.Printf("type of b:%T\n", b) 
	c := *b
	fmt.Printf("type of c:%T\n", c)
	fmt.Printf("value of c:%v\n", c)
}
/**
打印结果：
    a:10,addr:0xc00001a088
    b:0xc00001a088,addr:0xc00000a028
    type of b:*int
    type of c:int
    value of c:10
*/
```

以上示例说明：

- a 是基本类型变量，通过取地址操作（&），将a变量地址赋值给了b变量，此时b变量为指针变量，指针类型为`*int`
- 通过指针取值操作（*），将指针变量b指向的变量赋值给了c，此时c变量为基本类型变量，变量类型为`int`

小结：

- `&`为取地址操作符，`*`为取值操作符，是一对互补操作



指针传值示例:

```go
func modify1(x int) {
    x = 100
}

func modify2(x *int) {
    *x = 100
}

func main() {
    a := 10
    modify1(a)
    fmt.Println(a) // 10
    modify2(&a)
    fmt.Println(a) // 100
}
```



### 空指针

指针被定义未分配变量，则它的初始值为nil

```go
func main() {
    var p *string
    fmt.Println(p)
    fmt.Printf("p的值是%v\n", p)
    if p != nil {
        fmt.Println("非空")
    } else {
        fmt.Println("空值")
    }
}
/**
打印结果:
    <nil>
    p的值是<nil>
    空值      
*/
```



### 内置函数（new、make）

引入示例：

```go
func main() {
    var a *int
    *a = 100
    fmt.Println(*a)

    var b map[string]int
    b["测试"] = 100
    fmt.Println(b)
}
```

以上代码执行有panic报错；原因是：GO语言在使用引用变量时，不仅要声明它，还需要为它分配地址。而对于值类型则不需要，因为值类型变量在声明时就已经默认分配好了内存。

关于内存分配，GO语言主要是使用`new`，`make` 两个内置函数进行内存分配

> new

```go
func new(Type) *Type
```

- `Type`表示数据类型
- `*Type`表示类型指针

new函数 较少使用，new函数主要用于获取一个类型的指针，且指针对应的值为该类型的默认值

示例：

```go
func main() {
    a := new(int)
    b := new(bool)
    fmt.Printf("%T\n", a) // *int
    fmt.Printf("%T\n", b) // *bool
    fmt.Println(*a)       // 0
    fmt.Println(*b)       // false
}    
```



> make

`make`也是用于分配内存。区别于`new`，它**只用于slice、map以及channel**的内存创建，并且返回类型是这三个类型本身，因为这三个类型都是引用类型，所以没有必要像`new`函数，返回类型指针。

```go
func make(t Type, size ...IntegerType) Type
```



> 两者区别

1. 都是用于内存分配
2. `make`只用于slice、map、channel的初始化，返回是这三个引用类型本身
3. `new`用于类型的内存分配，内存对应值为类型默认值，返回是类型指针





## Map

map是一种无序的基于key-value的数据结构，Go语言中的map是引用类型，必须初始化才能使用

### 定义

```go
map[KeyType]ValueType
```

- KeyType：表示键的类型
- ValueType：表示键对应值的类型



map需要使用make来进行初始化

```go
make(map[KeyType]ValueType，[cap])
```

`cap`表示map容量，非必须参数，但最好创建时设置合适值



### 基本使用

```go
func main() {
	//使用案例1
	test1()
	//使用案例2 [支持在声明时填充元素]
	test2()
}

func test1() {
	scoreMap := make(map[string]int, 8)
	scoreMap["zhangsan"] = 90
	scoreMap["xiaoming"] = 100
	fmt.Println(scoreMap)
	fmt.Println(scoreMap["xiaoming"])
	fmt.Printf("type of a:%T\n", scoreMap)
}

func test2() {
	userInfo := map[string]string{
		"username": "misatan",
		"password": "12356",
	}
	fmt.Println(userInfo)
}

/**
打印结果：
	map[xiaoming:100 zhangsan:90]
    100
    type of a:map[string]int
    map[password:12356 username:misatan]
*/
```



### 判断是否存在

判断键是否存在可以利用`map[key]`的返回值来判断

```go
value,ok := map[key]
```

示例:

```go
func main() {
	scoreMap := make(map[string]int)
	scoreMap["zhangsan"] = 90
	scoreMap["lisi"] = 100
	//如果key存在ok为true，v为对应的值，不存在ok为false，v为值类型的默认值
	v, ok := scoreMap["lisi"]
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("查无此人")
	}
}
/**
打印结果：
	100
*/
```



### 遍历

使用`for range`遍历map

```go
func main() {
	fmt.Println("---test1---")
	test1()
	fmt.Println("---test2---")
	test2()
	fmt.Println("---test3---")
	test3()
}

func test1() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["李四"] = 60
	for k, v := range scoreMap {
		fmt.Println(k, v)
	}
}

//只想遍历key
func test2() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["李四"] = 60
	for k := range scoreMap {
		fmt.Println(k)
	}
}

//只想遍历value
func test3() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["李四"] = 60
	for _, v := range scoreMap {
		fmt.Println(v)
	}
}
```

注意：遍历map的元素顺序与添加键值对的顺序无关



### 删除

删除`map`键值对使用内置函数`delete()`

````go
delete(map,key)
````

解释：

- map：表示要删除键值对的map
- key：表示要删除的键值对的键

示例代码：

```go
func main() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["王五"] = 60
	delete(scoreMap, "小明") //将小明:100从map中删除
	for k, v := range scoreMap {
		fmt.Println(k, v)
	}
}
/**
	王五 60
	张三 90
*/
```



### Map类型的切片

```go
func main() {
	var mapSlice = make([]map[string]string, 3)
	for index, value := range mapSlice {
		fmt.Printf("index:%d,value:%v\n", index, value)
	}
	fmt.Println("before init")
	//对切片中的map元素初始化
	mapSlice[0] = make(map[string]string, 10)
	mapSlice[0]["name"] = "王五"
	mapSlice[0]["password"] = "123456"
	mapSlice[0]["addr"] = "沈阳大街"
	for index, arr := range mapSlice {
		fmt.Printf("index:%d,value:%v\n", index, arr)
	}
}
/**
打印结果：
    index:0,value:map[]
    index:1,value:map[]
    index:2,value:map[]
    before init
    index:0,value:map[addr:沈阳大街 name:王五 password:123456]
    index:1,value:map[]
    index:2,value:map[]
*/
```

访问切片具体Map：`slice[index][key]`





### Value为切片的Map

```go
func main() {
	sliceMap := make(map[string][]string, 3)
	fmt.Println(sliceMap)
	fmt.Println("before init")
	key := "china"
	value, ok := sliceMap[key]
	if !ok {
		value = make([]string, 0, 2)
	}
	value = append(value, "北京", "上海")
	sliceMap[key] = value
	fmt.Println(sliceMap)
}
/**
打印结果：
    map[]
    before init
    map[china:[北京 上海]]
*/
```





## Struct

Go没有类的概念，也不支持类继承等面向对象的概念，Go语言中通过结构体的内嵌配合接口，程序具有更改的扩展性和灵活性

### 前置知识

#### 自定义类型

Go语言可以使用`type`关键字来自定义类型，可以基于基本类型如`int`，`string`等来定义，也可以通过`struct`定义

```go
//将MyInt定义为Int类型
type MyInt int
//定义struct类型
type Student struct{
    Name:string
    age:int8
}
```

通过`type`定义的类型是全新的类型，`MyInt`就是定义的全新类型，具有int特性



#### 类型别名

类型别名是Go1.9添加新功能，也是使用`type`关键字定义类型的别名

```go
type MyInt = int
```

以下数据类型在`go`语言中都是以类型别名实现的

```go
type byte = uint8
type rune = int32
```



#### 两者区别

- 自定义类型是定义了一个新的数据类型
- 类型别名并没有定义新的数据类型

以下代码解释：

```go
//类型定义
type NewInt int

//类型别名
type MyInt = int

func main() {
    var a NewInt
    var b MyInt

    fmt.Printf("type of a:%T\n", a) //type of a:main.NewInt
    fmt.Printf("type of b:%T\n", b) //type of b:int
} 
```



### 定义

Go语言提供一种自定义的数据类型，封装多个基本数据类型，这种数据类型叫做结构体，GO语言使用结构体来实现面向对象

使用`type`和`struct`关键字定义结构体，格式如下：

```go
type 类型名 struct{
    字段名 字段类型
    字段名 字段类型
    ...
}
/**
注意点:
	1.类型名：标识自定义结构体的名称，同一个包内不能重复
	2.字段名：在结构体中必须唯一
	3.字段类型：具体数据类型，可以是结构体类型
*/
```

举个例子，定义一个Person的结构体：

```go
type Person1 struct{
    Name string
    city string
    age int8
}
type Person2 struct{
    Name,city string
    age int8
}
```

结构体同样的数据类型可以定义在同一行



### 实例化

#### 基本实例化

```go
type Person struct {
	name string
	age  int8
}

// 可以像声明内置类型一样声明结构体类型
var person Person

func main() {
	var p Person
	p.name = "sw"
	p.age = 23
	fmt.Printf("p=%v\n", p)//p={sw 23}
	fmt.Printf("p=%#v\n", p)//p=main.Person{name:"sw", age:23}

	p0 := Person{}
	fmt.Printf("p=%v\n", p0)//p={ 0}
	fmt.Printf("p=%#v\n", p0)//p=main.Person{name:"", age:0}   
}

```

通过`实例名.属性名`的方式访问结构体中的成员变量



#### 匿名结构体

定义一些临时数据等场景下可以使用匿名结构体这一特性：

```go
/**
	var 实例名 struct{
		字段名 字段类型;
		字段名 字段类型;
		...
	}
*/
func main() {
    var user struct{Name string; Age int}
    user.Name = "pprof.cn"
    user.Age = 18
    fmt.Printf("%#v\n", user)
} 
```



#### 指针类型结构体

可通过`new`函数，对结构体实例化，得到结构体指针

```go
type Person struct{
    name:string
    age:int8
}
func main(){
	var p = new(Person)
	p.name = "somebody"
	p.age = 100
	fmt.Printf("p type of %T\n", p) //p2 type of *main.Person
	fmt.Printf("p = %#v\n", p)	//&main.Person{name:"somebody", age:100}
}
```

go语言 指针类型结构体 支持 使用 `p.property`的方式 直接访问结构体字段 【go语言的语法糖;实际底层是 *p.property】



#### 取结构体的地址实例化

使用取地址符`&`，对结构体进行实例化

```go
type Person struct{
    name:string
    age:int8
}
func main(){
    var p = &Person{}
    fmt.Printf("%T\n", p)     //*main.Person
    fmt.Printf("p=%#v\n", p) //p=&main.Person{name:"", age:0}
    p.name = "博客"
    p.age = 30
    fmt.Printf("p=%#v\n", p) //p=&main.Person{name:"博客", age:30} 
}
```



### 初始化

结构体是值类型；实例化后，不对字段做初始化，则为相关类型初始化默认值

```go
type person struct {
    name string
    city string
    age  int8
}

func main() {
    var p4 person
    fmt.Printf("p4=%#v\n", p4) //p4=main.person{name:"", city:"", age:0}
} 
```



#### 键值对初始化

示例引用结构体定义：

```go
type person struct {
    name string
    city string
    age  int8
}
```



```go
func main(){
    p := person{
        name: "pprof.cn",
        city: "北京",
        age:  18,
	}
	fmt.Printf("p=%#v\n", p) //p=main.person{name:"pprof.cn", city:"北京", age:18}
}
```

也可以对指针结构体进行键值对初始化

```go
func main(){
    p := &person{
        name: "pprof.cn",
        city: "北京",
        age:  18,
    }
    fmt.Printf("p=%#v\n", p) //p=&main.person{name:"pprof.cn", city:"北京", age:18} 
}
```

键值初始化，可以选择性的初始化字段，并非必须全部指定初始化值

```go
p := &person{
    city: "北京",
}
fmt.Printf("p=%#v\n", p) //p=&main.person{name:"", city:"北京", age:0}
```



#### 列表初始化

初始化不写键，直接写值列表进行初始化

```go
p := &person{
    "pprof.cn",
    "北京",
    18,
}
fmt.Printf("p=%#v\n", p) //p=&main.person{name:"pprof.cn", city:"北京", age:18}
```

需要注意：

1. 必须初始化结构体所有字段
2. 值的顺序必须喝结构体字段定义顺序一致
3. 不能喝键值初始化方式混用



### 构造函数

go的结构体没有构造函数，但是可以自己实现。避免结构体复杂，值拷贝的性能开销，一般构造函数的返回类型为指针类型

```go
func newPerson(name, city string, age int8) *person {
	return &person{
		name: name,
		city: city,
		age:  age,
	}
}
```





### 方法和接受者

Go语言函数功能作用于某一个变量，可以将函数的接受者定义为这个变量；接受者类型其他语言中的`this`,`self`；需要自己定义

```go
func (接受者变量 接受者类型) 方法名(形参)(返回参数){
    函数体
}
```

说明:

1. 接受者变量：命名官网建议，以类型的第一个小写字母命名，不建议使用`this` or`self`
2. 接受者类型：可以是指针类型 or 非指针类型

函数接受者这一特性可以为结构体定义方法

```go
func main() {
	p1 := newInstance("sw", 25)
	fmt.Printf("p1 age = %v\n", p1.age) //p1 age = 25
	p1.Dream()
	fmt.Printf("p1 age = %v\n", p1.age) //p1 age = 25
	p1.Sleep()
	fmt.Printf("p1 age = %v\n", p1.age) //p1 age = 18 
}

type Person struct {
	name string
	age  int8
}

func newInstance(name string, age int8) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

// Dream 值类型接受者	修改不影响原值
func (p Person) Dream() {
	fmt.Printf("%s的梦想是学好go语言\n", p.name)
	p.age = 20
}

// Sleep 指针类型接受者	修改影响原值
func (p *Person) Sleep() {
	fmt.Printf("%s要睡觉了\n", p.name)
	p.age = 18
}
```

#### 什么时候使用指针类型接受者

1. 需要修改接受者中的值
2. 接受者的拷贝代价比较大的大对象
3. 保证一致性，如果有某个方法使用了指针接受者，那么其他的方法也应该使用指针接受者【要用都用】



### 匿名

#### 匿名字段

结构体声明，允许成员字段只声明类型，不命名。这种字段就叫匿名字段

```go
func main() {
	p := Anonymous{
		"yes",
		18,
	}
	fmt.Printf("%#v\n", p)	//main.Anonymous{string:"yes", int:18}
	fmt.Println(p.string, p.int) //yes 18
}

// Anonymous 匿名字段：类型名作为字段名；结构体要求字段名必须唯一，所以每种字段类型最多只有一个匿名字段
type Anonymous struct {
	string
	int
}
```

匿名字段默认采用类型名作为字段名，结构体要求字段名称唯一，所以同种类型只能有一个匿名字段



#### 嵌套结构体

结构体A中有结构体B`类型`或者`指针`的成员字段



#### 嵌套匿名结构体

```go
//Address 地址结构体
type Address struct {
    Province string
    City     string
}

//User 用户结构体
type User struct {
    Name    string
    Gender  string
    Address //匿名结构体
}
func main() {
    var user2 User
    user2.Name = "pprof"
    user2.Gender = "女"
    user2.Address.Province = "黑龙江"    //通过匿名结构体.字段名访问
    user2.City = "哈尔滨"                //直接访问匿名结构体的字段名
    fmt.Printf("user2=%#v\n", user2) //user2=main.User{Name:"pprof", Gender:"女", Address:main.Address{Province:"黑龙江", City:"哈尔滨"}}
} 
```

当访问结构体成员时会先在结构体中查找该字段，找不到再去匿名结构体中查找。



#### 嵌套结构体字段冲突问题

假如外部struct的字段名和内部的字段名相同

有以下两个名称冲突的规则：

1. 外部struct覆盖内部struct的同名字段、同名方法
2. 同级别的struct出现同名字段、同名方法将报错

第一个规则使得Go Struct能够实现面向对象中的重写（override），而且可以重写字段、重写方法

第二个规则使得同名属性不会出现歧义

示例：

```go
type A struct{
    a int
    b int
}

type B struct{
    b float32
    c string
    d string
}

type C struct{
    A
    B
    a string
    c string
}

var c C
```

按照规则1：属于C的`a，c`会分别覆盖`A.a，B.c`。可以直接使用`c.a，c.c`分别访问属于C的`a，c`字段，使用`c.d`或`c.B.d`都能访问属于嵌套的`B.d`字段。如果想要访问内部struct中被覆盖的属性，可以使用`c.A.a`的方式访问

按照规则2：A和B在C中都是同级嵌套结构体，所以`A.b`和`B.b`是冲突的，将会报错，因为调用`c.b`不知道是要调用`A.b`还是`B.b`



### 继承

在结构体中通过嵌套匿名结构体实现继承`*Father`

```go
//Animal 动物
type Animal struct {
    name string
}

func (a *Animal) move() {
    fmt.Printf("%s会动！\n", a.name)
}

//Dog 狗
type Dog struct {
    Feet    int8
    *Animal //通过嵌套匿名结构体实现继承
}

func (d *Dog) wang() {
    fmt.Printf("%s会汪汪汪~\n", d.name)
}

func main() {
    d1 := &Dog{
        Feet: 4,
        Animal: &Animal{ //注意嵌套的是结构体指针
            name: "乐乐",
        },
    }
    d1.wang() //乐乐会汪汪汪~
    d1.move() //乐乐会动！
}
```



### 可见性

- 结构体名，结构体字段，首字母大写；则包外可见（公开的）；否则仅在包内进行访问（私有的）



### 标签（Tag）

`Tag`是结构体的元信息，可以在运行时，通过反射机制读取出来

具体格式如下：

```go
`key1`:"value1" `key2`:"value2"
```

示例演示：

```go
type P struct {
	Name string `json:"name" bson:"Name"`
	Age  int8   `json:"age" bson:"Age"`
}

func main() {
	p := P{
		"sw",
		18,
	}

	pType := reflect.TypeOf(p)
	fieldName, isOk := pType.FieldByName("Name")
	if isOk {
		jsonTag := fieldName.Tag.Get("json")
		bsonTag := fieldName.Tag.Get("bson")
		fmt.Println("Name json Tag=", jsonTag, ",bson Tag=", bsonTag)
	} else {
		fmt.Println("no Name field")
	}

	fieldName, isOk = pType.FieldByName("Age")
	if isOk {
		jsonTag := fieldName.Tag.Get("json")
		bsonTag := fieldName.Tag.Get("bson")
		fmt.Println("Age json Tag=", jsonTag, ",bson Tag=", bsonTag)
	} else {
		fmt.Println("no Age field")
	}
}
```

反射之后学习



## Flow_Control

#### IF

```go
if 布尔表达式 {
    执行体
}else if 布尔表达式{

    执行体
}else{

    执行体
}
```

示例：

```go
func main() {
	ifTest(100)
}

func ifTest(num int) {
	if num < 0 {
		fmt.Println("num < 0,num:", num)
	} else if num < 100 {
		fmt.Println("num < 100,num:", num)
	} else {
		fmt.Println("num >= 100,num:", num)
	}
}
```

Go语言不支持三目运算符，避免影响GO语言简洁性和代码的可读性



#### Switch

Go switch 分支表达式**可以是不同的类型，不限于常量**。可以省略break，默认自动终止

复习一下Java分支表达式：**byte、short、char、int、Byte、Short、Character、Integer、Enum、String**

```go
switch <expression>{
	case var1:
   		...
    case var2:
    	...
    case var3:
    	...
    default:
        ...
}
```

示例：

```go
```

