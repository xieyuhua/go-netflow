(WIP) go-netflow
-------
一个简单的Linux下端口流量采集工具

```bash
$ make

$ go-netflow -ports 8080,443
```

目前采集的端口流量汇总主要以日志的形式输出；

本工具目前主要是我用于测试开发环境的端口流量监控，不建议用于生产环境

TODO

- 支持windows环境的端口流量采集
- 流量输出形式支持可扩展（目前打日志）
- 

```
[netflow port: 54126, in_bytes: 0, out_bytes:52, date:2023-05-13 10:38, Timestamp:1683945493]
[netflow port: 8888, in_bytes: 0, out_bytes:350, date:2023-05-13 10:38, Timestamp:1683945497]
```
