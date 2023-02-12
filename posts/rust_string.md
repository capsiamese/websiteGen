---
title: "Rust_string"
date: 2022-12-29T11:07:44+08:00
draft: true
tags: ['rust']
---

# String类型

字符串字面值, 即被硬编码进程序中的字符串值.

String类型的字符串管理分配在堆上的数据, 所以能够存储编译时未知大小的数据.

可以使用from函数基于字符串字面值创建String`let s = String::from("Hello")`

String的实现由三部分组成, 这部分存储在栈上.
1. ptr指向底层数组
2. len底层数组已使用长度
3. capacity底层数组容量

当我们将一个String赋值给另一个String时, 这意味着我们从栈上拷贝了他的指针, 长度, 容量, 并没有赋值指针指向的堆上的数据.

```rust
fn main() {
    let s1 = String::from("Hello");
    let s2 = s1;

    println!("{}", s1);
}
```
这段代码不能编译成功, 因为rust在拷贝指针, 长度和容量时使第一个变量无效, 这个操作被称为移动, 上面的例子可以被解读为s1被移动到了s2.

如果我们需要深度赋值String中堆上的数据, 可以使用`clone`通用函数
```rust
fn main() {
    let s1  = String::from("World");
    let s2 = s1.clone();
    println!("{}, {}", s1, s2);
}
```

# 字符串Slice

字符串slice是String中一部分值的引用
```rust
fn main() {
    let s = String::from("hello world");
    let hello = &s[0...5];
    let world = &s[6..11];
}
```
可以使用一个由中括号的[starting_index..ending_index]指定的range创建一个slice.
对于rust的`..range`语法, 如果想要从0开始可以不写starting_index, 如果要包含String的最后一个字节可以不写ending_index. 也可以同时舍弃两个索引.

*** 字符串slice ***的类型声明写作*** &str ***, 并且字符串字面量就是&str类型.

