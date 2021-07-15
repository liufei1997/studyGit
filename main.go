package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type User struct {
	Id   int
	Name string `json:"name1"`
	Age  int
}

func (u User) ReflectCallFuncNoArgs() {
	fmt.Println("ReflectCallFuncNoArgs")
}

func (u User) ReflectCallFuncHasArgs(name string, age int) {
	fmt.Println("ReflectCallFuncHasArgs name: ", name, ", age:", age, "and origal User.Name:", u.Name)
}

func (u User) ReflectCallFunc() {
	fmt.Println("Allen.Wu ReflectCallFunc")
}

func DoFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
}

type Enum int

const (
	Zero Enum = 0
)

type Person struct {
	Name string
	age  int
}

func Print1(ch chan int, group *sync.WaitGroup) {
	defer func() {
		waitGroup.Done()
	}()
	fmt.Println(<-ch)
}

var waitGroup sync.WaitGroup

func g1(ch chan int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for i := 1; i < 100; i++ {
		ch <- i
		if i%2 != 0 {
			fmt.Println("g1 => ", i)
		}
	}
}

func g22(ch chan int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for i := 1; i < 100; i++ {
		<-ch
		if i%2 == 0 {
			fmt.Println("g2 => ", i)
		}
	}
}

// Tip:通过cancel主动关闭
func ctxCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())

		case <-time.After(time.Millisecond * 100):
			fmt.Println("Time out")
		}
	}(ctx)
	cancel()
}

// Tip:通过超时，自动触发
func ctxTimeOut() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	// 主动执行 cancel，也会让协程收到消息
	defer cancel()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-time.After(time.Millisecond * 10):
			fmt.Println("Time out")
		}
	}(ctx)

	time.Sleep(time.Second)
}

// Tip:用Key/Value传递参数，可以浅浅封装一层，转化为自己想要的结构体
func ctxValue() {
	ctx := context.WithValue(context.Background(), "User", "liufei")
	go func(ctx context.Context) {
		v, ok := ctx.Value("User").(string)
		if ok {
			fmt.Println("pass user value", v)
		}
	}(ctx)
	time.Sleep(time.Second)
}

func wg() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			if i < 5 {
				fmt.Println(i)
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("finished")
}

// Tip: waitGroup 不要进行copy
func errWg1() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg sync.WaitGroup) {
			fmt.Println(i)
			defer wg.Done()
		}(i, wg)
	}
	wg.Wait()
	fmt.Println("finished")
}

// Tip: waitGroup 的 Add 要在goroutine前执行
func errWg2() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func(i int) {
			wg.Add(1)
			fmt.Println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("finished")
}

// Tip: waitGroup 的 Add 很大会有影响吗？
func errWg3() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(100)
		go func(i int) {
			fmt.Println(i)
			defer wg.Done()
			wg.Add(-100)
		}(i)
	}
	fmt.Println("finished")
}

type Ball struct {
	hits int
}

func passBall() {
	table := make(chan *Ball)
	go player("ping", table)
	go player("pong", table)
	// Tip:核心逻辑：往 channel 里面放数据，作为启动信号；从 channel 里面取出数据，作为关闭信号
	table <- new(Ball)
	time.Sleep(time.Second)
	<-table
}

func player(name string, table chan *Ball) {
	for {
		// 刚进 goroutine 时，阻塞在这里
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(time.Millisecond * 100)
		// Tip：运行到这里时，另一个 goroutine 在接收数据,所以能准确到达
		table <- ball
	}
}

func playerWithClose(name string, table chan *Ball) {
	for {
		ball, ok := <-table
		if !ok {
			break
		}
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(time.Millisecond * 100)
		// Tip：运行到这里时，另一个 goroutine 在接收数据,所以能准确到达
		table <- ball
	}
}

func g(ctx context.Context) {
	fmt.Println(ctx.Value("begin"))
	fmt.Println("你是猪")
	go g2(context.WithValue(ctx, "movie", "正义联盟"))
}

func g2(ctx context.Context) {
	fmt.Println(ctx.Value("movie"))
	fmt.Println("电影很好看")
}

func testCtxTimeOut() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("超时")
				return
			case <-time.After(time.Second * 3):
				fmt.Println(666)
			}
		}
	}(ctx)
	time.Sleep(time.Second * 7)
}

func doSomething(ctx context.Context) {
	select {
	case <-time.After(time.Second * 5): // 5 second pass
		fmt.Println("finish do something")
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

type Test struct {
	User string
}

func md5s(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type stPeople struct {
	Gender bool
	Name   string
}

type stStudent struct {
	stPeople
	Class int
}

func f1(arr ...int) {
	fmt.Println(arr)
}

func testDefer() int {
	var i int
	defer func() {
		i++
	}()
	i = 1

	return i
}

func reverse(s []rune) []rune {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type A struct {
	s string
}

// 这是上面提到的 "在方法内把局部变量指针返回" 的情况
func foo(s string) *A {
	a := new(A)
	a.s = s
	return a //返回局部变量a,在C语言中妥妥野指针，但在go则ok，但a会逃逸到堆
}

var ch = make(chan int)

func print5(i int) {
	<-ch
	fmt.Println(i)
}

func g11(ch chan int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for i := 1; i < 10; i++ {
		ch <- i
		if i%2 != 0 {
			fmt.Println("g1 => ", i)
		}
	}
}

func g21(ch chan int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for i := 1; i < 10; i++ {
		<-ch
		if i%2 == 0 {
			fmt.Println("g2 => ", i)
		}
	}
}

var ch1 chan int

func fetch(ch chan int) {

	ch <- 1
}

func set(ch chan int) {

	out := <-ch
	fmt.Println(out)
}

func job(index int) {
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("执行完毕，序号:%d\n", index)
}

func go1(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}
func go2(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 11; i++ {
		out := <-ch
		fmt.Println(out)
	}
}
func Increase() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func update(a []int) {
	a = append(a, 1)
}

type intList []int

func main() {

	var list intList = []int{1,2,3,4}
	for i:= 0 ; i < len(list);i++{
		
	}
	//s1 := make([]int, 0)
	//s2 := []int{1,2,3}
	//update(s1)
	//update(s2)
	//fmt.Println(s1)
	//fmt.Println(s2)

	//bytes := []byte("aaa")
	//fmt.Println(string(bytes))
	//startIdStr := fmt.Sprintf("%d",1706186690)
	//fmt.Println(startIdStr)

	//in := Increase()
	//fmt.Println(in())
	//fmt.Println(in())

	//var wg sync.WaitGroup
	//wg.Add(2)
	//ch := make(chan int, 10)
	//go go1(ch, &wg)
	//go go2(ch, &wg)
	//wg.Wait()
}

func main1() {
	maxNum := 10
	pool := make(chan struct{}, maxNum)
	var wg sync.WaitGroup
	for i := 10; i < 100; i++ {
		pool <- struct{}{}
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			job(index)
		}(i)
	}

	//go set(ch1)
	//go fetch(ch1)
	//
	//time.Sleep(time.Second*5)

	//waitGroup.Add(2)
	//ch := make(chan int)
	//go g11(ch, &waitGroup)
	//go g21(ch, &waitGroup)
	//waitGroup.Wait()

	//for i := 1; i <= 10; i ++ {
	//	go print5(i)
	//	ch <- 1
	//}

	//var count uint32
	//trigger := func(i uint32, fn func()) {
	//	for {
	//		if n := atomic.LoadUint32(&count); n == i {
	//			fn()
	//			atomic.AddUint32(&count, 1)
	//			break
	//		}
	//		time.Sleep(time.Nanosecond)
	//	}
	//}
	//for i := uint32(0); i < 10; i++ {
	//	go func(i uint32) {
	//		fn := func() {
	//			fmt.Println(i)
	//		}
	//		trigger(i, fn)
	//	}(i)
	//}
	//trigger(10, func() {})

	//ch := make(chan int, 1)
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		ch <- i
	//	}()
	//}
	//
	//var wg sync.WaitGroup
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go func(a int,ch chan int) {
	//		//defer func() {
	//		//	wg.Done()
	//		//}()
	//		fmt.Println(<-ch)
	//	}(i,ch)
	//}
	//wg.Wait()
	//

	//c := make(chan int, 1)
	//for i := 0; i < 1; i++ {
	//	go func(i int) {
	//		c <- 8
	//	}(i)
	//	time.Sleep(time.Second * 2)
	//	if i == 0 {
	//		close(c)
	//	}
	//}
	//
	//for i := 0; i < 3; i++ {
	//	out, ok := <-c
	//	fmt.Println(out,ok)
	//}

	//ch := make(chan int)
	//ch <- 1
	//var wg sync.WaitGroup
	//for i:= 0; i < 1; i++ {
	//	wg.Add(1)
	//	go func(chan int, *sync.WaitGroup){
	//		defer func(group *sync.WaitGroup){
	//			wg.Done()
	//		}(&wg)
	//		out , ok :=  <- ch
	//		fmt.Println(out , ok)
	//	}(ch,&wg)
	//	if i==0 {
	//		close(ch)
	//	}
	//}
	//wg.Wait()

	//a := foo("hello")
	//b := a.s + " world"
	//c := b + "!"
	//fmt.Println(c)

	//var s = "abc"
	//runes := reverse([]rune(s))
	//fmt.Println(runes)
	//fmt.Println(string(runes))

	//var s1 []int
	//s2 := make([]int,0)
	//s4 := make([]int,0)
	//s5 := make([]float64,0)
	//fmt.Printf("%+v \n",*(*reflect.SliceHeader)(unsafe.Pointer(&s5)))
	//fmt.Printf("s1 pointer:%+v, s2 pointer:%+v, s4 pointer:%+v, \n", *(*reflect.SliceHeader)(unsafe.Pointer(&s1)),*(*reflect.SliceHeader)(unsafe.Pointer(&s2)),*(*reflect.SliceHeader)(unsafe.Pointer(&s4)))
	//fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s1))).Data==(*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data)
	//fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data==(*(*reflect.SliceHeader)(unsafe.Pointer(&s4))).Data)

	//m1 := make(map[string]int)
	//i,ok := m1["jack"]
	//if ok {
	//	fmt.Println(i)
	//} else {
	//	fmt.Println(i)	// 0
	//}
	//
	//fmt.Println(m1["tom"])

	//str := `sss`
	//fmt.Println([]byte(str))

	//var str string = "aaaa"
	//var data []byte = []byte(str)
	//fmt.Println(data)

	//for i := 0; i < 3; i++ {
	//	defer fmt.Println(i)
	//}

	//for i := 0; i < 3; i++ {
	//	i := i
	//	defer func() {
	//		fmt.Println(i)
	//	}()
	//}

	//fmt.Println(testDefer())

	//s1 := []int{1, 2}
	//fmt.Println(s1)
	//fmt.Println(len(s1),cap(s1))
	//s1 = append(s1, 1, 1,1)
	//fmt.Println(len(s1),cap(s1))
	//fmt.Println(s1)
	//s1[0] = 99
	//fmt.Println(s1)

	//data := []int{1, 2, 3}
	//for i, v := range data {
	//	if i == 0 {
	//		data[i] = 99
	//    }
	//	v *= 10  // data 中原有元素是不会被修改的
	//}
	//fmt.Println("data: ", data) // data:  [1 2 3]

	//student := stStudent{stPeople{true, "222"}, 4}
	//fmt.Println(student)
	//f1(1,23,4)

	//defer func() {
	//	fmt.Println("recovered: ", recover())
	//}()
	//panic("not good")

	//var data = []byte(`{"status": 200}`)
	//var result map[string]interface{}
	//
	//if err := json.Unmarshal(data, &result); err != nil {
	//
	//}
	//sprintf := fmt.Sprintf("%f", result["status"])
	//fmt.Println(sprintf)

	//x := "text"
	//xBytes := []byte(x)
	//xBytes[0] = 'T' // 注意此时的 T 是 rune 类型
	//x = string(xBytes)
	//fmt.Println(x) // Text

	//s := md5s("15056671063")
	//
	//fmt.Print(len(s))
	//fmt.Println()
	//fmt.Print(s)

	// a:= "\tgit.in.codoon.com/backend/common v0.0.0-20210412054509-7c9f3212a286"
	//var slice []int
	//slice[1] = 0

	// 闭包
	//wg := sync.WaitGroup{}
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go func(int, *sync.WaitGroup, ) {
	//		fmt.Printf("i:%d", i)
	//		wg.Done()
	//	}(i, &wg)
	//}
	//wg.Wait()
	//fmt.Println("exit")

	//ch := make(chan struct{User string})
	//go func() {
	//	fmt.Println("start working")
	//	time.Sleep(time.Second * 1)
	//	ch <- struct{User string}{}
	//}()
	//<-ch
	//fmt.Println("finished")

	//// 创建空的 context的两种方法
	//ctx := context.Background() // 返回一个空的context ，不能被 cancel，key为空
	//// todoCtx := context.TODO()	// 和 Background类似，当你不确定的时候使用
	//ctx, cancel := context.WithCancel(ctx)
	//go func() {
	//	time.Sleep(time.Second * 6)
	//	cancel()
	//}()
	//doSomething(ctx)

	//go testCtxTimeOut()
	//time.Sleep(time.Second * 7)

	//a := struct {
	//}{}
	//fmt.Println(a)

	//ctx := context.WithValue(context.Background(), "begin", "从台词看到一句话:")
	//go g(ctx)
	//time.Sleep(time.Second)

	//passBall()
	//errWg2()
	//var mutex sync.Mutex
	//mutex.Lock()
	//fmt.Println("Locked")
	//mutex.Lock()

	//var mutex sync.Mutex
	//fmt.Println("Lock the lock")
	//mutex.Lock()
	//fmt.Println("The lock is locked")
	//channels := make([]chan int, 4)
	//for i := 0; i < 4; i++ {
	//	channels[i] = make(chan int)
	//	go func(i int, c chan int) {
	//		fmt.Println("Not lock: ", i)
	//		mutex.Lock()
	//		fmt.Println("Locked: ", i)
	//		time.Sleep(time.Second)
	//		fmt.Println("Unlock the lock: ", i)
	//		mutex.Unlock()
	//		c <- i
	//	}(i, channels[i])
	//}
	//time.Sleep(time.Second)
	//fmt.Println("Unlock the lock")
	//mutex.Unlock()
	//time.Sleep(time.Second)
	//
	//for _, c := range channels {
	//	<-c
	//}

	//var m map[string]string
	//delete(m,"sdsd")
	//m["result"] = "result"

	//waitGroup.Add(2)
	//ch := make(chan int)
	//go g1(ch, &waitGroup)
	//go g2(ch, &waitGroup)
	//waitGroup.Wait()

	//ch := make(chan int, 2)
	//ch <- 1
	////ch <- 2
	//fmt.Println(<-ch)
	//close(ch)
	//
	//fmt.Println(<-ch)

	//s := "Hello 世界"
	//b := []byte(s)    // 转换为 []byte，数据被自动复制
	//b[5] = ','        // 把空格改为半角逗号
	//b[5] = '中'
	//fmt.Printf("%s\n", s)
	//fmt.Printf("%s\n", b)

	//s1 := make([]int, 0)
	//var s2 []int
	//var s3 []int
	//fmt.Println(&s1,&s3,&s2)

	//fmt.Println(len(s1),cap(s1))	// 0 0
	//fmt.Println(len(s2),cap(s2))	// 0 0
	//s1 = append(s1,[]int{1,2}...)
	//s2 = append(s2,[]int{1,2}...)
	//fmt.Println(len(s1),cap(s1))	// 1 2
	//fmt.Println(len(s2),cap(s2))	// 1 2
	//fmt.Println(s1,s2)

	//fmt.Println(s1 == nil)	// false
	//fmt.Println(s2 == nil)	// true

	//// m := make(map[string]string)
	//var m1 map[string]string
	//m1["name"] = "tom"

	//s1 := [...]int{1,2,3}
	//inits := s1[1:2]
	//fmt.Println(inits)
	//s2 := make([]int, 2)
	//copy(s2,s1[:2])
	//fmt.Println(s2)
	//s2[0] = 99
	//fmt.Println(s1,s2)

	//arr := [2][3]int{{1,2,3},{2,3,4}}
	//for _,r := range arr{
	//	for _,c :=range r{
	//		fmt.Println(c)
	//	}
	//
	//}

	//arr := [2][3]int{{1,2,3},{2,3,4}}
	//fmt.Println(arr)

	//arr := []Person{{
	//	Name: "Bob",
	//	age:  18,
	//},{
	//	Name: "lisa",
	//	age: 19,
	//}}
	//fmt.Println(arr)

	//date := time.Now().AddDate(-2, 0, 0)
	//fmt.Println(date)
	//formatTime := now.Format("2006-01-02 15:04:05")
	//fmt.Println(formatTime)	// 2021-07-02 16:43:25

	//type T struct {
	//	A int
	//	B string
	//}
	//t := T{23, "skidoo"}
	//s := reflect.ValueOf(&t).Elem()
	//fmt.Println(s) // {23 skidoo}
	//typeOfT := s.Type()
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	//}

	//var a = 55
	//valueOf := reflect.ValueOf(&a)
	//value := valueOf.Elem()
	//if valueOf.Elem().CanSet() {
	//	value.SetInt(66)
	//}
	//fmt.Println(a)	// 66

	//fmt.Println(reflect.TypeOf(a))
	//fmt.Println(reflect.TypeOf(a).Name())
	//fmt.Println(reflect.TypeOf(a).Kind())
	//fmt.Println(reflect.ValueOf(a))

	//u1 := User{
	//	Id:   0,
	//	Name: "tom",
	//	Age:  18,
	//}
	//typeOf := reflect.TypeOf(u1)
	//valueOf := reflect.ValueOf(u1)
	//numField := typeOf.NumField()
	//for i:= 0; i<numField; i++ {
	//	field := typeOf.Field(i)
	//	value := valueOf.Field(i)
	//	fmt.Println(field.Name,field.Tag,field.Type,value)
	//}
	//
	//if fieldByName, b := typeOf.FieldByName("Name"); b{
	//	fmt.Println(fieldByName.Tag.Get("json"))
	//}

	//name := typeOf.Name()
	//kind := typeOf.Kind()
	//fmt.Println(name)
	//fmt.Println(kind)

	//s1 := []int{1,2}
	//typeOfS1 := reflect.TypeOf(s1)
	//fmt.Println(len(typeOfS1.Name()))
	//fmt.Println(typeOfS1.Kind())

	//// 声明一个空结构体
	//type cat struct {
	//}
	//// 获取结构体实例的反射类型对象
	//typeOfCat := reflect.TypeOf(cat{})
	//// 显示反射类型对象的名称和种类
	//fmt.Println(typeOfCat.Name(), typeOfCat.Kind())
	//// 获取Zero常量的反射类型对象
	//typeOfA := reflect.TypeOf(Zero)
	//// 显示反射类型对象的名称和种类
	//fmt.Println(typeOfA.Name(), typeOfA.Kind())

	//var a int64 = 5
	//
	//typeOfA := reflect.TypeOf(a)
	//name := typeOfA.Name()
	//kind := typeOfA.Kind()
	//fmt.Println(name,kind)
	//valueOfA := reflect.ValueOf(a)
	//canSet := valueOfA.CanSet()
	//fmt.Println(canSet)

	//var num float64 = 1.2345
	//fmt.Println("old value of pointer:", num)
	//
	//// 通过reflect.ValueOf获取num中的reflect.Value，注意，参数必须是指针才能修改其值
	//pointer := reflect.ValueOf(num)
	//fmt.Printf(pointer.NumField())
	//newValue := pointer.Elem()
	//
	//fmt.Println("type of pointer:", newValue.Type())
	//fmt.Println("settability of pointer:", newValue.CanSet())

	// 重新赋值
	//newValue.SetFloat(77)
	//fmt.Println("new value of pointer:", num)

	////////////////////
	// 如果reflect.ValueOf的参数不是指针，会如何？
	//pointer = reflect.ValueOf(num)
	//newValue = pointer.Elem() // 如果非指针，这里直接panic，“panic: reflect: call of reflect.Value.Elem on float64 Value”

	//user := User{1, "Allen.Wu", 25}
	//
	//DoFiledAndMethod(user)

	//user := User{1, "Allen.Wu", 25}
	//
	//// 1. 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，得到“反射类型对象”后才能做下一步处理
	//getValue := reflect.ValueOf(user)
	//
	//// 一定要指定参数为正确的方法名
	//// 2. 先看看带有参数的调用方法
	//methodValue := getValue.MethodByName("ReflectCallFuncHasArgs")
	//args := []reflect.Value{reflect.ValueOf("wudebao"), reflect.ValueOf(30)}
	//methodValue.Call(args)
	//
	//// 一定要指定参数为正确的方法名
	//// 3. 再看看无参数的调用方法
	//methodValue = getValue.MethodByName("ReflectCallFuncNoArgs")
	//args = make([]reflect.Value, 0)
	//methodValue.Call(args)

	//a := 5
	//type1 := reflect.TypeOf(a)
	//valueOf := reflect.ValueOf(a)
	//
	//fmt.Println(type1, valueOf)

	//s1 := make([]int, 3, 3)
	//s1 = append(s1, 1, 2,3,4)
	//fmt.Println(s1,len(s1), cap(s1))

	//s1 := make([]int, 3, 3)
	//
	//s1 = append(s1, 1)
	//fmt.Println(s1,len(s1), cap(s1))  //  4, 6
	//
	//s1 = append(s1, 2)
	//fmt.Println(s1,len(s1), cap(s1))	// 5 6
	//
	//s1 = append(s1, 3)
	//fmt.Println(s1,len(s1), cap(s1))	// 6,6
	//
	//s1 = append(s1, 4)
	//fmt.Println(s1,len(s1), cap(s1))	// 7 12

	//s1 := make([]int, 2, 2)
	////s2 := []int{1,2}
	//s1 = append(s1, 1,2,3)
	//fmt.Println(len(s1), cap(s1))

	//s1 := make([]int, 2, 2)
	////s2 := []int{1,2}
	//s1 = append(s1, 1)
	//fmt.Println(len(s1), cap(s1))
	//s1 = append(s1, 2)
	//fmt.Println(len(s1), cap(s1))
	//s1 = append(s1, 3)
	//fmt.Println(len(s1), cap(s1))

	// 奇怪
	//s1 := make([]int, 3, 3)
	//fmt.Println(len(s1), cap(s1))
	//s1 = append(s1, 1,2,3,4)
	//fmt.Println(len(s1), cap(s1))			// 7  8

	// 奇怪
	//s2:=[]int{1,2,3}
	//fmt.Println(len(s2), cap(s2))
	//s2 = append(s2, 4,5,6,7,8,9,10,11,12,13)
	//fmt.Println(s2)
	//fmt.Println(len(s2), cap(s2))

	//var slice []int
	//for i := 0; i < 1024; i++ {
	//	slice = append(slice, i)
	//}
	//var slice1 = []int{1024, 1025}
	//fmt.Printf("slice len = %v cap = %v\n", len(slice), cap(slice))
	//fmt.Printf("slice1  len = %v cap = %v\n", len(slice1), cap(slice1))
	//slice = append(slice, slice1...)
	//fmt.Printf("slice len = %v cap = %v\n", len(slice), cap(slice))

	//a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//b := a[1:3]
	//b[0] = 11	// b[0]的改写，即对a[1]的改写
	//fmt.Println(a[1]) // a[1]被写成了11
	//fmt.Println(len(a), cap(a))  // 10  10
	//fmt.Println(len(b), cap(b))	 // 2    9

	//第1条中的例子：
	//var slice = make([]int,2,3)
	//var slice1 = []int{4, 5, 6, 7, 8, 9, 10, 11, 12}
	//fmt.Printf("slice %v len = %v cap = %v\n", slice, len(slice), cap(slice))
	//fmt.Printf("slice1 %v len = %v cap = %v\n", slice1, len(slice1), cap(slice1))
	//slice = append(slice, slice1...)
	//fmt.Printf("slice %v len = %v cap = %v\n", slice, len(slice), cap(slice))

	// 奇怪
	//s1 := make([]int, 3, 6)
	//fmt.Println(len(s1), cap(s1))
	//s1 = append(s1, 1,2,3,4,5,6,7,8,9,10,11,12)
	//fmt.Println(s1)
	//fmt.Println(len(s1), cap(s1))

	//var slice = []int{1, 2, 3, 4, 5, 6, 7}
	//var slice1 = []int{8, 9}
	//fmt.Printf("slice %v len = %v cap = %v\n", slice, len(slice), cap(slice))
	//fmt.Printf("slice1 %v len = %v cap = %v\n", slice1, len(slice1), cap(slice1))
	//slice = append(slice, slice1...)
	//fmt.Printf("slice %v len = %v cap = %v\n", slice, len(slice), cap(slice))

	//var slice = []int{1, 2, 3}
	//var slice1 = []int{4, 5, 6, 7, 8, 9, 10, 11, 12}
	//fmt.Printf("slice %v len = %v cap = %v\n", slice, len(slice), cap(slice))
	//fmt.Printf("slice1 %v len = %v cap = %v\n", slice1, len(slice1), cap(slice1))
	//slice = append(slice, slice1...)
	//fmt.Printf("slice %v len = %v cap = %v\n", slice, len(slice), cap(slice))

	//s1 := make([]int, 3)
	//fmt.Println(len(s1),cap(s1))
	//s1 = append(s1, 3)
	//fmt.Println(len(s1),cap(s1))

	//a := []int{1, 2, 3, 4, 5}
	//shadow := append([]int{}, a[1:3]...)
	//shadow = append(shadow, 100)
	//fmt.Println(shadow, a)
	// [2 3 100] [1 2 3 4 5]

	//var s1 = []int{}
	//fmt.Println(s1 == nil)		// false

	//a := []int{1,2,3}
	//b := a[:]
	//b[1] = 1
	//fmt.Println(a,b)

	//a := []int{1,2,3,4,5}
	//shadow := a[1:3]
	//shadow = append(shadow,100)
	//fmt.Println(shadow,a)		// [2 3 100] [1 2 3 100 5]

	//wg := sync.WaitGroup{}
	//wg.Add(5)
	//for i := 0; i < 5; i++ {
	//	go func(i int) {
	//		defer wg.Done()
	//		fmt.Println(i)
	//	}(i)
	//}
	//wg.Wait()

}
