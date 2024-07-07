## 介绍
gopool 是一个高性能的 goroutine 池，旨在重复使用 goroutine 并限制 goroutine 的数量。

## 能力
* 复用goroutine，高性能
* 限制goroutine num

## 使用
```
// 方式1：自定义pool
p := NewPool("test", 100, NewConfig(2))
p.Go(func() {
    // your logic
}}

// 方式2：使用默认的pool
gopool.Go(func() {
    // your logic
}}
```