
![visitors](https://visitor-badge.glitch.me/badge?page_id=stateIs0.busyJavaThreadGraber&left_color=green&right_color=red)

# BusyJavaThreadGraber
抓取最繁忙的 N 个 Java 线程，快速排查 java 程序的 cpu 问题，看看是哪个线程在作妖。



# 使用
1. 安装 go
2. git clone {ssh/http}, cd main, go build
3. ./main -pid {java pid} -tick 1 -threshold 50 & 
4. 如触发阈值(例如 50），程序会自动 dump 文件，文件格式：pid + 时间 + .txt

参数介绍：

-pid java 进程 id

-tick 抓取间隔

-threshold cpu 阈值， 触发阈值则抓取堆栈

-threadNum 要抓取的线程数

-user java 进程所属用户名

# 效果
![](stack.png)


# 原理
基于 gopsutil 监控目标进程的 cpu 使用率（效果同 top）。
如果触发阈值，则使用 top -Hp 抓取此时此刻，该进程下的所有线程，并获取 CPU 使用率，再进行排序。然后拿出最忙的 10 个线程。

同时执行 jstack，获取 Java 堆栈，拿之前获取的线程 id 在堆栈中查找，快速获取线程堆栈。

最后的效果就是 dump 最繁忙的 10 个线程的堆栈信息。

# 注意
仅支持 linux
