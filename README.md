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

Task必须继承自ATask，并实现ITask的Config、Work两个函数

Config 为通过配置文件初始化Task

Work 为执行任务，表示Task已经获得Module的调度权




