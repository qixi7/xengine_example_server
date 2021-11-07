# 示例服务器

## 说明
由于只是示例，没有添加太多复杂的内容，比如性能收集，rcp。后续时间空闲了补充示例。

## 让服务器跑起来
### step_1：修改配置文件
【config/devconfig.json】用于设置服务器需要配置的字段，所有传到git or svn上的配置文件都只是一个模块文件(xxx.json.template)，用这个devconfig.json文件去生成具体的服务器配置【serverconfig.json】

### step_2：生成配置文件

```shell
sh run.sh genconfig
```

即可生成配置文件，观察example下生成了serverconfig.json，原理见脚本。

### step_3：编译

```shell
sh run.sh build
```

### step_4：运行（编译和运行也能直接在IDE【goland】上进行，这里针对linux）

```shell
# 方式一（守护进程）
sh run.sh startall
# 方式二（非守护进程）
cd example
./example
```

### step_5：关闭

```shell
sh run.sh stopall
```



## 剩下的后续有时间补充~

