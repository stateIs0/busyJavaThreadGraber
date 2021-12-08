# BusyJavaThreadGraber
抓取最繁忙的 N 个 Java 线程。

# 原理
基于 gopsutil 监控目标进程的 cpu 使用率（效果同 top），如果触发阈值，则抓取此时此刻，该进程下的所有进程 CPU 使用率，并进行排序，dump 最繁忙的 10 个线程的堆栈信息。
