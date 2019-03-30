# IOTBox

## 1、概念说明：
IOTBox中包含三个概念

Kernel 内核，负责调度和执行整个系统，包含定时器定时执行特定的任务。

Module  系统组成模块，包括通讯类模块及存储类模块等，模块负责执行任务。

Task    任务，执行某种特定任务的功能。


## 2、使用方式

### 2.1、开发Module

Module必须继承自common.AModule，并实现IModule的Config、Send两个函数

Config 为通过配置文件初始化模块

Send 表示Task发送数据包给模块，并期待模块的返回

### 2.2 开发Task

Task必须继承自ATask，由于go语言没有构造函数，因此必须定义CreateXXX函数(该函数为动态加载做准备)
该函数必须设置该对象的SetOtherConfig函数。
<code>
func (d *Demo) XXXConfig(kernel common.IKernel, config map[string]interface{}) error {
</code>

Task为状态机模式，需要配置当前执行的Work函数

<code>
func XXXWork(module common.IModule) (common.WorkState, error) 
</code>

Config 为通过配置文件初始化Task

Work 为执行任务，表示Task已经获得Module的调度权

## 3、Example说明

### 3.1 定时任务

本例子用于模仿定时执行的多帧任务，由模块core启动该任务，任务执行后，将结果发送给handler配置的
处理事件(Sqlite及上报)。

>> 注意：由于每次调用任务对象为固定对象，因此需考虑避免并发问题（如被某个模块调用时，不可被其他模块所调用）

#### 3.1.1 涉及代码

/module/core/core.go        定时启动的任务组

/module/sqlite/sqlite.go    Sqlite执行任务

/task/demo/demo.go          Demo多帧任务

/task/sql/sql.go            将结果转换为SQL语句，并发送给Sqlite执行

### 3.2 