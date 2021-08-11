package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"math"
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

//puclic static void mergeSort(int[] arr, int l, int r) {
//	if(l >= r) {
//		return
//	}
//	int mid = (l + r)>>1
//	mergerSort(arr, l, mid)
//	mergeSort(arr, mid + 1, r)
//	int i = l, j = mid, k = 0;
//	int[] temp = new int[r - l + 1];
//	while(i <= mid && j <= r){
//		if(arr[i] < arr[j]){
//			temp[k++] = i++;
//		}else{
//			temp[k++] = j++;
//		}
//	}
//	while(i <= mid){
//		temp[k++] = arr[i++];
//	}
//	while(j <= l){
//		temp[k++] = arr[j++];
//	}
//	for(int i = l, j= 0; i < r; i++, j++){
//		arr[i] = temp[j];
//	}
//}

var Exists = struct{}{}

// 构造一个set，首先定义set的类型
// set类型
type Set struct {
	m map[interface{}]struct{}
}

// 初始化 在声明的同时可以选择传入或者不传入进去
func New(items ...interface{}) *Set {
	// 获取Set的地址
	s := &Set{}
	// 声明map类型的数据结构
	s.m = make(map[interface{}]struct{})
	s.Add(items...)
	return s
}

// 添加 简化操作可以添加不定个数的元素进入到Set中，用变长参数的特性来实现这个需求即可，因为Map不允许Key值相同，所以不必有排重操作。同时将Value数值指定为空结构体类型。
func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		s.m[item] = Exists
	}
	return nil
}

// 包含 Contains操作其实就是查询操作，看看有没有对应的Item存在，可以利用Map的特性来实现，但是由于不需要Value的数值，所以可以用 _,ok来达到目的
func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

// 获取Set长度很简单，只需要获取底层实现的Map的长度即可：
func (s *Set) Size() int {
	return len(s.m)
}

// 清除操作的话，可以通过重新初始化Set来实现
func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

// 判断两个Set是否相等，可以通过循环遍历来实现，即将A中的每一个元素，查询在B中是否存在，只要有一个不存在，A和B就不相等
func (s *Set) Equal(other *Set) bool {
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}

	// 迭代查询遍历
	for key := range s.m {
		// 只要有一个不存在就返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 判断A是不是B的子集，也是循环遍历的过程
func (s *Set) IsSubset(other *Set) bool {
	// s的size长于other，不用说了
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

type Car struct {
	weight int
	name   string
}

func (p *Car) Run() {
	fmt.Println("running")
}

type Bike struct {
	Car
	lunzi int
}
type Train struct {
	Car
}

func (p *Train) String() string {
	str := fmt.Sprintf("name=[%s] weight=[%d]", p.name, p.weight)
	return str
}

type Student struct {
	Name  string
	Age   int
	Score int
}

//将Pupil和Graduate共有的方法绑定到 *Student
func (stu *Student) ShowInfo() {
	fmt.Printf("学生名：%v 年龄：%v 成绩：%v \n", stu.Name, stu.Age, stu.Score)
}
func (stu *Student) SetScore(score int) {
	//业务判断
	stu.Score = score
}

//小学生
type Pupil struct {
	Student //嵌套了Student匿名结构体
}

func (p *Pupil) testing() {
	fmt.Println("小学生正在考试中..")
}

//大学生
type Graduate struct {
	Student //嵌套了Student匿名结构体
}

func (p *Graduate) testing() {
	fmt.Println("大学生正在考试中...")
}

// 实现多态
//声明一个接口
type Usb interface {
	//声明两个没有实现的方法
	//Start()
	//Stop()
}

type Phone struct {
	name string
}

//让Phone实现Usb接口的方法
func (p Phone) Start() {
	fmt.Println("手机开始工作...")
}
func (p Phone) Stop() {
	fmt.Println("手机停止工作...")
}

type Camera struct {
	name string
}

func (c Camera) Start() {
	fmt.Println("相机开始工作...")
}
func (c Camera) Stop() {
	fmt.Println("相机停止工作...")
}

type BaseNum struct {
	num1 int
	num2 int
} // BaseNum 即为父类型名称

type Add struct {
	BaseNum
} //加法子类, 定义加法子类的主要目的, 是为了定义对应子类的方法

type Sub struct {
	BaseNum
} //减法子类

func (a *Add) Opt() (value int) {
	return a.num1 + a.num2
} //加法的方法实现

func (s *Sub) Opt() (value int) {
	return s.num1 + s.num2
} //减法的方法实现

type Opter interface { //接口定义
	Opt() int //封装, 归纳子类方法, 注意此处需要加上返回值, 不然没有办法输出返回值(因为方法中使用了返回值)
}

/**
 * 定义几何接口
 */
type geometry interface {
	area() float64
	perim() float64
}

/**
 * 定义矩形结构
 */
type rect struct {
	width, height float64
}

/**
 * 定义圆形结构
 */
type circle struct {
	radius float64
}

/**
 * 实现矩形面积方法
 */
func (r rect) area() float64 {
	return r.width * r.height
}

/**
 * 实现矩形周长方法
 */
func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

/**
 * 实现圆形面积方法
 */
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

/**
 * 实现圆形周长方法
 */
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

/**
 * 接口做参数实现计算方法
 * @param  g geometry      接口参数
 * @return
 */
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

type Animal interface {
	Sleep()
	Age() int
	Type() string
}
type Cat struct {
	MaxAge int
}

func (this *Cat) Sleep() {
	fmt.Println("Cat need sleep")
}
func (this *Cat) Age() int {
	return this.MaxAge
}
func (this *Cat) Type() string {
	return "Cat"
}

type Dog struct {
	MaxAge int
}

func (this *Dog) Sleep() {
	fmt.Println("Dog need sleep")
}
func (this *Dog) Age() int {
	return this.MaxAge
}
func (this *Dog) Type() string {
	return "Dog"
}

func Factory(name string) Animal {
	switch name {
	case "dog":
		return &Dog{MaxAge: 20}
	case "cat":
		return &Cat{MaxAge: 10}
	default:
		panic("No such animal")
	}
}

// Definition for a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderTraversal(root *TreeNode) (res []int) {
	stack := []*TreeNode{}
	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = append(res, root.Val)
		root = root.Right
	}
	return
}

func spiralOrder(matrix [][]int) []int {
	m := len(matrix)
	n := len(matrix[0])
	// 已知结果集大小，初始化时便可以设置大小
	res := make([]int, m*n)
	up, l := 0, 0
	down, r := m-1, n-1
	// 索引，不要用append添加，避免多次扩容
	var index int
	for {
		for i := l; i <= r; i++ {
			res[index] = matrix[up][i] //向右移动，直到最右,该行已经遍历完，接下来会向下遍历，，所以需要更新未遍历的上边界
			index++
		}
		up++ //更新上边界，上边界是向下移动一行，对应坐标值则是加1
		if up > down {
			break
		}

		for i := up; i <= down; i++ {
			res[index] = matrix[i][r] //向下移动，直到最下,该列已经遍历完，接下来会向左遍历，所以需要更新未遍历的右边界
			index++
		}
		r-- //更新右边界，右边界是向左移动一列，对应坐标值则是减1
		if r < l {
			break
		}

		for i := r; i >= l; i-- {
			res[index] = matrix[down][i] //向左移动，直到最左,该行已经遍历完，接下来会向上遍历，所以需要更新未遍历的下边界
			index++
		}
		down-- //更新右边界，右边界是向左移动一列，对应坐标值则是减1
		if down < up {
			break
		}

		for i := down; i >= up; i-- {
			res[index] = matrix[i][l] //向上移动，直到最上,该列已经遍历完，接下来会向右遍历，所以需要更新未遍历的左边界
			index++
		}
		l++ //更新左边界，左边界是向右移动一列，对应坐标值则是加1
		if l > r {
			break
		}
	}
	return res
}

type CDCodoonShareBase struct {
	SportsType int32  `protobuf:"varint,1,opt,name=sports_type,json=sportsType" json:"sports_type,omitempty"`
	StartTime  string `protobuf:"bytes,2,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	Location   string `protobuf:"bytes,3,opt,name=location" json:"location,omitempty"`
	ProductId  string `protobuf:"bytes,4,opt,name=product_id,json=productId" json:"product_id,omitempty"`
	Version    string `protobuf:"bytes,5,opt,name=version" json:"version,omitempty"`
	// 是否分享位置，0指不分享，1指分享
	SharePosition int32 `protobuf:"varint,6,opt,name=share_position,json=sharePosition" json:"share_position,omitempty"`
	// 里程
	Distance float64 `protobuf:"fixed64,7,opt,name=distance" json:"distance,omitempty"`
	// 时长 00:00:00
	Timing float64 `protobuf:"fixed64,8,opt,name=timing" json:"timing,omitempty"`
	// 时速
	Speed float64 `protobuf:"fixed64,9,opt,name=speed" json:"speed,omitempty"`
	// 配速 00'00''
	AveragePace float64 `protobuf:"fixed64,10,opt,name=average_pace,json=averagePace" json:"average_pace,omitempty"`
	// 卡路里
	Calorie float64 `protobuf:"fixed64,11,opt,name=calorie" json:"calorie,omitempty"`
	// 1表示运动中，2表示结束
	State int32 `protobuf:"varint,12,opt,name=state" json:"state,omitempty"`
	// 装备名称
	EquipName string `protobuf:"bytes,13,opt,name=equip_name,json=equipName" json:"equip_name,omitempty"`
	// 步数
	Steps float64 `protobuf:"fixed64,14,opt,name=steps" json:"steps,omitempty"`
	// 海拔
	Altitude float64 `protobuf:"fixed64,15,opt,name=altitude" json:"altitude,omitempty"`
	// 平均时速
	AverageSpeed float64 `protobuf:"fixed64,16,opt,name=average_speed,json=averageSpeed" json:"average_speed,omitempty"`
}

// 活动实况第一次进入时单个用户信息
type CDCodoonActivityOneUserMsg struct {
	// 概要信息
	Base *CDCodoonShareBase `protobuf:"bytes,1,opt,name=base" json:"base,omitempty"`
	// 坐标信息
	Position *CDCodoonPositionMsg `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
	// 用户信息
	User *CDUserInfo `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
}

type CDUserInfo struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	Nick   string `protobuf:"bytes,2,opt,name=nick" json:"nick,omitempty"`
	Avatar string `protobuf:"bytes,3,opt,name=avatar" json:"avatar,omitempty"`
	// 性别 0女，1男，2未知
	Gender string `protobuf:"bytes,4,opt,name=gender" json:"gender,omitempty"`
}
type CDCodoonPositionMsg struct {
	// 类型：1指暂停，2指恢复
	Type int32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	// 经度
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude" json:"longitude,omitempty"`
	// 纬度
	Latitude float64 `protobuf:"fixed64,3,opt,name=latitude" json:"latitude,omitempty"`
	// 海拔
	Elevation float64 `protobuf:"fixed64,4,opt,name=elevation" json:"elevation,omitempty"`
	// 耗时
	TostartcostTime float64 `protobuf:"fixed64,5,opt,name=tostartcostTime" json:"tostartcostTime,omitempty"`
	// 长度
	Tostartdistance float64 `protobuf:"fixed64,6,opt,name=tostartdistance" json:"tostartdistance,omitempty"`
	// 速度
	Topreviousspeed float64 `protobuf:"fixed64,7,opt,name=topreviousspeed" json:"topreviousspeed,omitempty"`
}

type CDCodoonActivityStartMsg struct {
	// 活动ID
	ActivityId int64 `protobuf:"varint,1,opt,name=activity_id,json=activityId" json:"activity_id,omitempty"`
	// 参加人数
	JoinNum int32 `protobuf:"varint,2,opt,name=join_num,json=joinNum" json:"join_num,omitempty"`
	// 跑团名称
	GroupName string `protobuf:"bytes,3,opt,name=group_name,json=groupName" json:"group_name,omitempty"`
	// 每个用户的信息
	UserList []*CDCodoonActivityOneUserMsg `protobuf:"bytes,4,rep,name=user_list,json=userList" json:"user_list,omitempty"`
	// 当前用户信息
	CurUser *CDCodoonActivityOneUserMsg `protobuf:"bytes,5,opt,name=cur_user,json=curUser" json:"cur_user,omitempty"`
	// 当前用户排名
	CurRank int32 `protobuf:"varint,6,opt,name=cur_rank,json=curRank" json:"cur_rank,omitempty"`
	// 当前用户是否参加了该活动 1参加
	IsJoin int32 `protobuf:"varint,7,opt,name=is_join,json=isJoin" json:"is_join,omitempty"`
	// 活动位置
	// Location *CDCodoonActivityGPS `protobuf:"bytes,8,opt,name=location" json:"location,omitempty"`
	// 活动名称
	ActivityName string `protobuf:"bytes,9,opt,name=activity_name,json=activityName" json:"activity_name,omitempty"`
}

func writeMap() {

	for i := 0; i < 5; i++ {
		if _, ok := M[i]; ok {
			fmt.Println("已存在，返回")
			break
		}
	}
}

//M1 := make(map[int]bool)

var M map[int]bool

func add11(a int) (e error) {
	fmt.Println(111)
	return
}

public class 快排 {
public static void main(String[] args) {
int[] arr = new int[]{5, 2, 3, 1};
quickSort(arr, 0, arr.length -1);
for (int i = 0; i < arr.length; i++) {
System.out.print(arr[i]);
}
}

public static void quickSort(int[] arr, int l, int r){
if (l >= r) return; //递归结束条件
// 因为循环的内容是不管三七二十一先把两个指针向中间移动一位再进行判断，
// 只有 把 i, j  取为边界的两侧，向中间移动时才会指到真正的边界
int i = l - 1; //左边的指针
int j = r + 1; //右边的指针
int mid = arr[l + ((r - l) >> 1)];    //选取数组中间的那个数作为分界点，把数组分成两部分
while(i < j){
do{
i++;
}while(arr[i] < mid);             //在左边大于等于分界点的数停下

do{
j--;
}while(arr[j] > mid);             //在右边小于等于分界点的数停下
if (i < j){                      //交换两个数
int temp = arr[i];
arr[i] = arr[j];
arr[j] = temp;
}
}
quickSort(arr, l, j);                 //递归处理左边
quickSort(arr, j+1, r); //递归处理右边
}
}

func main() {

	s := make([]int, 10)

	s = append(s, 1, 2, 3)

	fmt.Println(s)

	//add11(1)
	//
	//sn1 := struct {
	//	age  int
	//	name string
	//}{age: 11, name: "qq"}
	//
	//sn2 := struct {
	//	age  int
	//	name string
	//}{age: 11, name: "qq"}
	//
	//if sn1 == sn2 {
	//	fmt.Println("sn1 == sn2")
	//}

	//var intervals = [][]int{{3,3},{2,2}}
	//
	//sort.Slice(intervals, func(i, j int) bool {
	//	return intervals[i][0]<intervals[j][0]
	//})
	//fmt.Println(intervals)

	//m := make(map[int]bool)
	//M = m
	//msg := &live_proto.CDCodoonReviewRankMsg{}
	//b := []byte("CM294AIQUxoV6ZW/5a6J5b6u6ams6ZSm5Y6m6ZifIrEBCiRiMGFlNGQ1My1mZjMzLTRlNzctOGZjMi0zOGU2YjI1MzgzOWQSDeS4veS4vTEwMjEyMTYabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2IwYWU0ZDUzLWZmMzMtNGU3Ny04ZmMyLTM4ZTZiMjUzODM5ZC8yMDIxLTAzLTIwVDEyLjQ3LjMwL3UzUkVZdGUyb3NCdmVEQnJIUS5qcGVnIPcfKVg5tMh2uMZAIrABCiQ0NzA4NDQ3Mi05ZGEyLTRjYzUtODY0MC1mZjg5ZjI2ZWM2OTYSDDc4MjQyMTg1ZXVsaRpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvNDcwODQ0NzItOWRhMi00Y2M1LTg2NDAtZmY4OWYyNmVjNjk2LzIwMjAtMDktMjJUMDkuNTcuMjMvMkdqTEFMMHlGaXdPa0wyT2x5LmpwZWcgyR0pPQrXoyCRxEAitgEKJGU5OTg3ZmYwLTdlMGUtNDJjMi1hNWNmLTFlODViOWUwN2RjZRIS6K645aSp57+UODc3ODQ4NjQ4Gm5odHRwczovL3MxLmNvZG9vbmNkbi5jb20vdXNlci9lOTk4N2ZmMC03ZTBlLTQyYzItYTVjZi0xZTg1YjllMDdkY2UvMjAyMC0wOC0wNlQxMS4xMC41Ni9JV2FzOFR1NHE2amRMNjFJaE0uanBlZyD7Gilfl9wNDInDQCKoAQokZTc4Zjc2MjItNWZmMy00MWExLWI4Y2EtNGVjMDVlM2ZjOTJhEgnolpvkuprliKkaaWh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2U3OGY3NjIyLTVmZjMtNDFhMS1iOGNhLTRlYzA1ZTNmYzkyYS8yMDE5LTAzLTE2VDA3LjE3LjA5L2x6Uk1xdklxRGFRWkNYa1AxMCDRGSnufD81/ljBQCKkAQokZDg4OWYzMWUtMWM4OC00YjhmLWI5NmQtZjRhYjI0YjhiNmUyEgV3eeiLsRppaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvZDg4OWYzMWUtMWM4OC00YjhmLWI5NmQtZjRhYjI0YjhiNmUyLzIwMTktMDgtMjJUMDYuNTQuMDYvSXRFUWJEUTU4WFRrbWkyR2w3INojKQAAAAAAWMFAIpwBCiQ3NjM3NDJkOS04ZjBiLTQ3MzgtODI2MS0yMzI1YWFkY2RhODkSCuaXtuWFiTA5NDcaXGh0dHBzOi8vaW1nMy5jb2Rvb24uY29tL3BvcnRyYWl0L2luaXQvNzYzNzQyZDktOGYwYi00NzM4LTgyNjEtMjMyNWFhZGNkYTg5LzE1MTU1NDU3NDQ3MTEuanBnIOwYKTMzMzMTNcFAIqwBCiRmYjAyZDQ0NC00YjkxLTQ2Y2UtOWM2Yy05M2JlZGQ4OTkwNWMSEOW6t+WNjuWdhzEwMjEwNTMaZmh0dHBzOi8vaW1nMy5jb2Rvb24uY29tL3BvcnRyYWl0L2ZiMDJkNDQ0LTRiOTEtNDZjZS05YzZjLTkzYmVkZDg5OTA1Yy9jNjVlZGFhMWE3ZTU0NGI5OWIxYzNlZjJhNTlhMTJiYyDdGSnb+X5qnDDBQCKrAQokZTNhMDRlNzAtYWVkMi00NGFkLTg0YjQtNDIyYzNlYmJhMGY2EgzmuIXmuIUxMDE2MTEaaWh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2UzYTA0ZTcwLWFlZDItNDRhZC04NGI0LTQyMmMzZWJiYTBmNi8yMDE4LTEwLTI1VDIxLjIxLjAyL1lNQ3hFWWNkU3p0dU1KQlR2eCCPGimM22gA75C/QCKtAQokMmI0YjkxNjgtNzdlMS00N2MzLTk3YzUtMDcxZjIxNGM5OWJlEgnmm77npaXok4kabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyLzJiNGI5MTY4LTc3ZTEtNDdjMy05N2M1LTA3MWYyMTRjOTliZS8yMDIwLTA0LTI5VDEwLjA4LjU2LzBPTmx1SEhLTkJnN25hSm1EUC5qcGVnIOorKX1M4kjco7xAIrQBCiQyMjRlYTM5Ni04Yzk1LTRiMTUtYjNhYy03YTEyMTBlNGE1ZWISEOW8oOWKsuWKmzEwMjA5ODgabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyLzIyNGVhMzk2LThjOTUtNGIxNS1iM2FjLTdhMTIxMGU0YTVlYi8yMDIwLTExLTA5VDEyLjA3LjQyL3NlZHcxZ3drelNCaDZsdFFLUi5qcGVnIIEQKcNkqmDUlbVAIrQBCiRjYzQ2NThlOS1lYjU0LTQ3YWYtOTg0ZC0wYzBkM2U2YTQ2ZWMSEOaItOa1t+i+iTEwMTkwNzkabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2NjNDY1OGU5LWViNTQtNDdhZi05ODRkLTBjMGQzZTZhNDZlYy8yMDIwLTAxLTA5VDE5LjMxLjQ3L3VCb3dKRkZ0ZHBFUFVhV1FBOS5qcGVnIK4TKaabxCDgFLVAIrEBCiRkY2I0NWQzYy01MWQ1LTQ3MmYtYTU3NC1hMTljOWVhODAwNmISDeWonOWonDEwMjE4OTAabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2RjYjQ1ZDNjLTUxZDUtNDcyZi1hNTc0LWExOWM5ZWE4MDA2Yi8yMDIwLTA5LTE3VDA3LjUyLjI2L1NzZ0p1dFRkRmxVYjFubTAyUC5qcGVnIL0RKY/C9Sjc2rNAIrUBCiQ3MWQ4MzFkYS02NDI5LTQzMDAtOTUzNi1jYjhiNTNiNGRkNmMSEee+hemHkeiMuS0xMDE4ODUyGm5odHRwczovL3MxLmNvZG9vbmNkbi5jb20vdXNlci83MWQ4MzFkYS02NDI5LTQzMDAtOTUzNi1jYjhiNTNiNGRkNmMvMjAxOS0xMi0yN1QxOC41NS4wNC9MaVJzU2hPaDh5THg2aFZsUU0uanBlZyC9EilMJD78ntKzQCKwAQokNDVlZWMyM2QtNjVjOS00MTYxLWJhYzAtYzZiZGUwZTBiYmMyEgzogIHllJDlm73lk6Uabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyLzQ1ZWVjMjNkLTY1YzktNDE2MS1iYWMwLWM2YmRlMGUwYmJjMi8yMDE5LTExLTIxVDEyLjI3LjU4L0VQZzVua0pxNDFhUVRGancybS5qcGVnIPARKbTIdr4vtLNAIq0BCiQ4YTc0NjUxMi0xNjc5LTQ2MDYtODA2Mi04YmNhNWQ4ZmZjZmESCeWImOWFg+mcnhpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvOGE3NDY1MTItMTY3OS00NjA2LTgwNjItOGJjYTVkOGZmY2ZhLzIwMjAtMDgtMDNUMjEuMTEuNDAvODZkNVRHeFZ1bVJqUWt5QzdILmpwZWcgkA0pnMQgsDJqr0AitAEKJGI0OTE3ZWM5LWIwMTItNDE5Ni1hNDhjLTVlNTBjZmVhZDdkYhIQ6ZmI5YWI5pWPMTAyMTA3MhpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvYjQ5MTdlYzktYjAxMi00MTk2LWE0OGMtNWU1MGNmZWFkN2RiLzIwMjEtMDYtMDVUMTAuMzEuMDUvVWRpN3dLbTFNbUM0cjlqdzhqLmpwZWcglQkpQBNhw1MArEAirQEKJGYwN2Q3MTMzLTE1NTEtNDUzNi05ZGIzLTFhM2FmZmM3YTUwNBIR6JKL6IyC5YWJNjQ3MDA3ODgaZmh0dHBzOi8vaW1nMy5jb2Rvb24uY29tL3BvcnRyYWl0L2YwN2Q3MTMzLTE1NTEtNDUzNi05ZGIzLTFhM2FmZmM3YTUwNC8xNjc1MDczYzM1YjM0NmNkOGYwNDAwNDU1YzgzODY3YyClPilIv30d+N2rQCK0AQokM2FiOTVjNWItYmI5NS00MmIzLWEwMDMtNTg4MGQwNjRmZGYzEhDlvKDkvJrorqkxMDE2NTg4Gm5odHRwczovL3MxLmNvZG9vbmNkbi5jb20vdXNlci8zYWI5NWM1Yi1iYjk1LTQyYjMtYTAwMy01ODgwZDA2NGZkZjMvMjAyMC0wOC0wOVQxOS40My40MS85THMyVzBaaU43RTg0d1BRZ2YuanBlZyCoBykldQKaaKeoQCow4oCc5YWo5rCR5YGl6Lqr5pel4oCd6ZSm5Y6m6Zif6IyF5rSy5rKz5aSn6IGa6LeR")
	//
	//// b := []byte("CKe44AIQExoMUkZT5a2d5b635rmWKrIBGq8BCiQ2ZjVkMzlhZi0wZTliLTQ3MzYtYjcxNi1lYWRjMzZmNTczY2QSFeS4gOWvuOays+WxseS4gOWvuOihgBpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvNmY1ZDM5YWYtMGU5Yi00NzM2LWI3MTYtZWFkYzM2ZjU3M2NkLzIwMjEtMDgtMDdUMTUuMTQuMTYvZzlSbkV5UFFwbmVrNzdHRmpKLmpwZWciADAAOAFCEgkAAADAsUNcQBEAAACgVhE3QEoVUkZT5a2d5b635rmW6IGa6LeR5pel")
	//proto.Unmarshal(b, msg)
	//
	//fmt.Println(msg)

	//[]byte	'CM294AIQUxoV6ZW/5a6J5b6u6ams6ZSm5Y6m6ZifIrEBCiRiMGFlNGQ1My1mZjMzLTRlNzctOGZjMi0zOGU2YjI1MzgzOWQSDeS4veS4vTEwMjEyMTYabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2IwYWU0ZDUzLWZmMzMtNGU3Ny04ZmMyLTM4ZTZiMjUzODM5ZC8yMDIxLTAzLTIwVDEyLjQ3LjMwL3UzUkVZdGUyb3NCdmVEQnJIUS5qcGVnIPcfKVg5tMh2uMZAIrABCiQ0NzA4NDQ3Mi05ZGEyLTRjYzUtODY0MC1mZjg5ZjI2ZWM2OTYSDDc4MjQyMTg1ZXVsaRpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvNDcwODQ0NzItOWRhMi00Y2M1LTg2NDAtZmY4OWYyNmVjNjk2LzIwMjAtMDktMjJUMDkuNTcuMjMvMkdqTEFMMHlGaXdPa0wyT2x5LmpwZWcgyR0pPQrXoyCRxEAitgEKJGU5OTg3ZmYwLTdlMGUtNDJjMi1hNWNmLTFlODViOWUwN2RjZRIS6K645aSp57+UODc3ODQ4NjQ4Gm5odHRwczovL3MxLmNvZG9vbmNkbi5jb20vdXNlci9lOTk4N2ZmMC03ZTBlLTQyYzItYTVjZi0xZTg1YjllMDdkY2UvMjAyMC0wOC0wNlQxMS4xMC41Ni9JV2FzOFR1NHE2amRMNjFJaE0uanBlZyD7Gilfl9wNDInDQCKoAQokZTc4Zjc2MjItNWZmMy00MWExLWI4Y2EtNGVjMDVlM2ZjOTJhEgnolpvkuprliKkaaWh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2U3OGY3NjIyLTVmZjMtNDFhMS1iOGNhLTRlYzA1ZTNmYzkyYS8yMDE5LTAzLTE2VDA3LjE3LjA5L2x6Uk1xdklxRGFRWkNYa1AxMCDRGSnufD81/ljBQCKkAQokZDg4OWYzMWUtMWM4OC00YjhmLWI5NmQtZjRhYjI0YjhiNmUyEgV3eeiLsRppaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvZDg4OWYzMWUtMWM4OC00YjhmLWI5NmQtZjRhYjI0YjhiNmUyLzIwMTktMDgtMjJUMDYuNTQuMDYvSXRFUWJEUTU4WFRrbWkyR2w3INojKQAAAAAAWMFAIpwBCiQ3NjM3NDJkOS04ZjBiLTQ3MzgtODI2MS0yMzI1YWFkY2RhODkSCuaXtuWFiTA5NDcaXGh0dHBzOi8vaW1nMy5jb2Rvb24uY29tL3BvcnRyYWl0L2luaXQvNzYzNzQyZDktOGYwYi00NzM4LTgyNjEtMjMyNWFhZGNkYTg5LzE1MTU1NDU3NDQ3MTEuanBnIOwYKTMzMzMTNcFAIqwBCiRmYjAyZDQ0NC00YjkxLTQ2Y2UtOWM2Yy05M2JlZGQ4OTkwNWMSEOW6t+WNjuWdhzEwMjEwNTMaZmh0dHBzOi8vaW1nMy5jb2Rvb24uY29tL3BvcnRyYWl0L2ZiMDJkNDQ0LTRiOTEtNDZjZS05YzZjLTkzYmVkZDg5OTA1Yy9jNjVlZGFhMWE3ZTU0NGI5OWIxYzNlZjJhNTlhMTJiYyDdGSnb+X5qnDDBQCKrAQokZTNhMDRlNzAtYWVkMi00NGFkLTg0YjQtNDIyYzNlYmJhMGY2EgzmuIXmuIUxMDE2MTEaaWh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2UzYTA0ZTcwLWFlZDItNDRhZC04NGI0LTQyMmMzZWJiYTBmNi8yMDE4LTEwLTI1VDIxLjIxLjAyL1lNQ3hFWWNkU3p0dU1KQlR2eCCPGimM22gA75C/QCKtAQokMmI0YjkxNjgtNzdlMS00N2MzLTk3YzUtMDcxZjIxNGM5OWJlEgnmm77npaXok4kabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyLzJiNGI5MTY4LTc3ZTEtNDdjMy05N2M1LTA3MWYyMTRjOTliZS8yMDIwLTA0LTI5VDEwLjA4LjU2LzBPTmx1SEhLTkJnN25hSm1EUC5qcGVnIOorKX1M4kjco7xAIrQBCiQyMjRlYTM5Ni04Yzk1LTRiMTUtYjNhYy03YTEyMTBlNGE1ZWISEOW8oOWKsuWKmzEwMjA5ODgabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyLzIyNGVhMzk2LThjOTUtNGIxNS1iM2FjLTdhMTIxMGU0YTVlYi8yMDIwLTExLTA5VDEyLjA3LjQyL3NlZHcxZ3drelNCaDZsdFFLUi5qcGVnIIEQKcNkqmDUlbVAIrQBCiRjYzQ2NThlOS1lYjU0LTQ3YWYtOTg0ZC0wYzBkM2U2YTQ2ZWMSEOaItOa1t+i+iTEwMTkwNzkabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2NjNDY1OGU5LWViNTQtNDdhZi05ODRkLTBjMGQzZTZhNDZlYy8yMDIwLTAxLTA5VDE5LjMxLjQ3L3VCb3dKRkZ0ZHBFUFVhV1FBOS5qcGVnIK4TKaabxCDgFLVAIrEBCiRkY2I0NWQzYy01MWQ1LTQ3MmYtYTU3NC1hMTljOWVhODAwNmISDeWonOWonDEwMjE4OTAabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyL2RjYjQ1ZDNjLTUxZDUtNDcyZi1hNTc0LWExOWM5ZWE4MDA2Yi8yMDIwLTA5LTE3VDA3LjUyLjI2L1NzZ0p1dFRkRmxVYjFubTAyUC5qcGVnIL0RKY/C9Sjc2rNAIrUBCiQ3MWQ4MzFkYS02NDI5LTQzMDAtOTUzNi1jYjhiNTNiNGRkNmMSEee+hemHkeiMuS0xMDE4ODUyGm5odHRwczovL3MxLmNvZG9vbmNkbi5jb20vdXNlci83MWQ4MzFkYS02NDI5LTQzMDAtOTUzNi1jYjhiNTNiNGRkNmMvMjAxOS0xMi0yN1QxOC41NS4wNC9MaVJzU2hPaDh5THg2aFZsUU0uanBlZyC9EilMJD78ntKzQCKwAQokNDVlZWMyM2QtNjVjOS00MTYxLWJhYzAtYzZiZGUwZTBiYmMyEgzogIHllJDlm73lk6Uabmh0dHBzOi8vczEuY29kb29uY2RuLmNvbS91c2VyLzQ1ZWVjMjNkLTY1YzktNDE2MS1iYWMwLWM2YmRlMGUwYmJjMi8yMDE5LTExLTIxVDEyLjI3LjU4L0VQZzVua0pxNDFhUVRGancybS5qcGVnIPARKbTIdr4vtLNAIq0BCiQ4YTc0NjUxMi0xNjc5LTQ2MDYtODA2Mi04YmNhNWQ4ZmZjZmESCeWImOWFg+mcnhpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvOGE3NDY1MTItMTY3OS00NjA2LTgwNjItOGJjYTVkOGZmY2ZhLzIwMjAtMDgtMDNUMjEuMTEuNDAvODZkNVRHeFZ1bVJqUWt5QzdILmpwZWcgkA0pnMQgsDJqr0AitAEKJGI0OTE3ZWM5LWIwMTItNDE5Ni1hNDhjLTVlNTBjZmVhZDdkYhIQ6ZmI5YWI5pWPMTAyMTA3MhpuaHR0cHM6Ly9zMS5jb2Rvb25jZG4uY29tL3VzZXIvYjQ5MTdlYzktYjAxMi00MTk2LWE0OGMtNWU1MGNmZWFkN2RiLzIwMjEtMDYtMDVUMTAuMzEuMDUvVWRpN3dLbTFNbUM0cjlqdzhqLmpwZWcglQkpQBNhw1MArEAirQEKJGYwN2Q3MTMzLTE1NTEtNDUzNi05ZGIzLTFhM2FmZmM3YTUwNBIR6JKL6IyC5YWJNjQ3MDA3ODgaZmh0dHBzOi8vaW1nMy5jb2Rvb24uY29tL3BvcnRyYWl0L2YwN2Q3MTMzLTE1NTEtNDUzNi05ZGIzLTFhM2FmZmM3YTUwNC8xNjc1MDczYzM1YjM0NmNkOGYwNDAwNDU1YzgzODY3YyClPilIv30d+N2rQCK0AQokM2FiOTVjNWItYmI5NS00MmIzLWEwMDMtNTg4MGQwNjRmZGYzEhDlvKDkvJrorqkxMDE2NTg4Gm5odHRwczovL3MxLmNvZG9vbmNkbi5jb20vdXNlci8zYWI5NWM1Yi1iYjk1LTQyYjMtYTAwMy01ODgwZDA2NGZkZjMvMjAyMC0wOC0wOVQxOS40My40MS85THMyVzBaaU43RTg0d1BRZ2YuanBlZyCoBykldQKaaKeoQCow4oCc5YWo5rCR5YGl6Lqr5pel4oCd6ZSm5Y6m6Zif6IyF5rSy5rKz5aSn6IGa6LeR'
	//		fmt.Println(string(b

	//var index int
	//for i := 0; i < 3; i++ {
	//	[index] = true
	//	index++
	//}
	//fmt.Println("11")
	//writeMap()
	//go func(a int ){
	//	if _, ok := M[a]; ok{
	//		fmt.Println("重复")
	//	}
	//}(index)

	//var index int
	//for i := 0; i < 3; i++ {
	//	if index = index + 1; index < 3{
	//		fmt.Println()
	//	}
	//}

	//s := [][]int{{1,2,3},{4,5,6},{7,8,9}}
	////s := [][]int{{1},{2}}
	//order := spiralOrder(s)
	//fmt.Println(order)

	//	i := -1
	//	fmt.Println(i)
	//Notice:
	//	i = i + 1
	//	fmt.Println(i)
	//
	//
	//	if i < 1 {
	//		goto Notice
	//	}

	//totalLength, _ := strconv.ParseFloat("555", 64)
	//fmt.Println(totalLength)

	//startDate := time.Now().Local()
	//fmt.Println(startDate)
	//s := []string{"111","222","333"}
	//join1 := strings.Join(s, ";")
	//fmt.Println(join1)
	//fmt.Println(s)

	// 2006-01-02T15:04:05Z07:00
	//fmt.Println(time.Now())

	//animal := Factory("dog")
	//animal.Sleep()
	//fmt.Printf("%s max age is: %d", animal.Type(), animal.Age())

	//r := rect{width: 3, height: 4}
	//c := circle{radius: 5}
	//
	///**
	// * 多态调用
	// */
	//measure(r)
	//measure(c)

	//// 继承
	//data:= BaseNum{2,3}
	//var a Add = Add{data}
	//var b Sub= Sub{data}
	//
	////使用接口
	//var i Opter
	//i = &a
	//value := i.Opt()
	//i = &b
	//value := i.Opt()
	//
	////使用多态
	//value := MultiState(&a)
	//value := MultiState(&b)
	////输出测试
	//fmt.Println(value)
	//
	////定义一个Usb接口数组，可以存放Phone和Canera的结构体变量
	////这里就体现出多态数组
	//var  usbArr [3]Usb
	//usbArr[0] = Phone{"小米"}
	//usbArr[1] = Phone{"vivo"}
	//usbArr[2] = Camera{"诺基亚"}
	//fmt.Println(usbArr)
	//
	//test := Phone{}
	//test.name = "小米"
	//fmt.Println(test.name)
	//test.Start()
	//test.Stop()
	//
	//test1 := Camera{}
	//test1.name = "索尼相机"
	//fmt.Println(test1.name)
	//test1.Start()
	//test1.Stop()

	//当我们对结构体嵌入了匿名结构体，使用方法会发生变化
	//pupil := &Pupil{}
	//pupil.Student.Name = "tom"
	//pupil.Student.Age = 8
	//pupil.testing()
	//pupil.SetScore(80)
	//pupil.Student.ShowInfo()
	//
	//graduate := &Graduate{}
	//graduate.Student.Name = "mary"
	//graduate.Student.Age = 24
	//graduate.testing()
	//graduate.SetScore(95)
	//graduate.Student.ShowInfo()

	//
	//var a Bike
	//a.weight = 100
	//a.name = "bike"
	//a.lunzi = 2
	//fmt.Println(a)
	//a.Run()
	//
	//var b Train
	//b.weight = 100
	//b.name = "train"
	//b.Run()
	//fmt.Printf("%s", &b)

	//m := map[string]string{
	//	"1": "one",
	//	"2": "two",
	//	"3": "three",
	//}
	//m["1"] = "eee"
	//fmt.Println(m)

	//var map1=make(map[int]string) //这种是直接通过make函数创建。
	//map1[5]="xiaxia"
	//map1[1]="kdkdd"
	//map1[8]="abc"
	//map1[4]="china"
	//for k,v:= range map1 {
	//	fmt.Printf("%v-%v\n", k, v)
	//}

	//s1 := []string{}
	//m1 := map[string]int{}
	//m1["a"] = 1
	//m1["b"] = 2
	//m1["c"] = 3
	//m1["d"] = 4
	//for v := range m1 {
	//
	//	s1 = append(s1, v)
	//}
	//fmt.Println(s1)
	//sort.Strings(s1)
	//fmt.Println(s1)

	//urltest := "dnwYEgksgFE"
	//fmt.Println(urltest)
	//encodeurl:= url.QueryEscape(urltest)
	//fmt.Println(encodeurl)

	//strings.ReplaceAll()

	//var environment string = os.Getenv("GOENV")

	// GOENV="/Users/codoon/Library/Application Support/go/env"

	//fmt.Println("InitConfig env: ", environment)

	//s := fmt.Sprintf("%.2f", float64(1043)/100)
	//s1 := fmt.Sprintf("%.2f", float64(4632)/100)
	//s2 := fmt.Sprintf("%.2f", float64(10)/100)
	//fmt.Println(s)
	//fmt.Println(s1)
	//fmt.Println(s2)

	//fmt.Println(time.Now().Local())

	//yesterdayEndTime, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Now().Location())
	//yesterdayStartTime := yesterdayEndTime.Add(-24 * time.Hour)
	//// 例如 当前时间（time.Now()）是 2021-07-22 10:14:06  则 theDayBeforeYesterdayEndTime 为 2021-07-21 00:00:00 +0800 CST
	//// 如果该用户的  continue_end_time < theDayBeforeYesterdayEndTime 则说明该用户断签，
	//// 则将 此时的 group_sports_groupmember表记录迁移到 user_sign_in_history 表
	//theDayBeforeYesterdayEndTime := yesterdayStartTime
	//if time.Now().Before(theDayBeforeYesterdayEndTime) {
	//	fmt.Println(1)
	//} else {
	//	fmt.Println(2)
	//}
	//
	//before := time.Now().Add(-1 * time.Millisecond).Before(time.Now())
	//fmt.Println(before)

	//now := time.Now()
	//location, err := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), time.Local)
	//
	//if err != nil{
	//
	//}
	//fmt.Println(location)
	//
	//
	//
	//fmt.Println(time.Now().Local(),555)
	//fmt.Println(!(time.Now().Local().IsZero()))

	// startDate, _ := time.ParseInLocation("2006-01-02", time.Now().Local().Format("2006-01-02"), time.Local)
	//
	//fmt.Println(time.Now())
	//yesterdayEndTime, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Now().Location())
	//yesterdayStartTime := yesterdayEndTime.Add(-24 * time.Hour)
	//fmt.Println(yesterdayEndTime)
	//fmt.Println(yesterdayStartTime)

	//h := md5.New()
	//h.Write([]byte("15390427825"))
	//md5PhoneMun := fmt.Sprintf("%x", h.Sum(nil))
	//fmt.Println(md5PhoneMun)

	// (1<<TrainingPlatformBitPosEnd - 1) & ^(1 << TrainingPlatformBitPosMagicCustomer)

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
