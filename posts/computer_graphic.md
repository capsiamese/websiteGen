---
title: "Computer_graphic"
date: 2022-10-20T17:28:33+08:00
draft: false
tags: ["cg"]
---

# introductory concepts

raytracer 
rasterizer

## coordinate systems
一个canvas具有width和height, 以像素为单位我们叫做$C_w$和$C_h$, 对于大多数屏幕origin在左上角
x从左向右增长, y从上向下增长, 这种坐标系对电脑的内存来说是非常自然的, 但是对人脑不是很友好,
然而3d图形程序员倾向使用字面书写的坐标系, 原点在中间x轴向右增加向左减少, y轴向上增加向下减少.

使用原点在中间的坐标系时x的范围是
$[-\frac{C_w}{2}, \frac{C_w}{2})$, 
y的范围是$[-\frac{C_h}{2}, \frac{C_h}{2})$

从第一个坐标系转到另一个坐标系  
$S_x = \frac{C_w}{2} + C_x$  
$S_y = \frac{C_h}{2} - C_y$

## color models
1. subtractive color model CMYK
2. additive color model RGB

## the scene

y is up and x and z are horizontal, and all theree axes are perpendicular to each other.

# basic raytracing

## canvas to viewport

canvas坐标轴$C_x$, $C_y$

从canvas坐标转换到空间坐标只需要改变比例 
$V_x = C_x * \frac{V_w}{C_w}$
$V_y = C_y * \frac{V_h}{C_h}$ 
$V_z = d$

对于每个在canvas上的像素$(C_x, C_y)$, 我们能决定相应的点在viewport上$(V_x, V_y, V_z)$

## tracing rays
## the ray equation
光线通过O并且光线是直的, 因此我们可以表达任意一个光线的点P
$P = O + t(V - O)$
t可以是任意实数, $(V - O)$是光线的方向$\vec{D}$
之后方程就变成了$P = O + t\vec{D}$

## the sphere equation
圆是一个点的集合, 这些点到一个固定点的距离相同, 这个距离叫做半径Radius
这个固定点叫做圆心Center. $distance(P, C) = r$

P和C之间的距离是向量P到C的长度$|P - C| = r$

一个向量($|\vec{V}|$)的长度它自己点积($\langle\vec{V},\vec{V}\rangle$)的平方根
$ \sqrt{\langle P-C, P-C \rangle} = r$, $\langle P-C, P-C \rangle = r^2$

## ray meets sphere




# ref
[ComputerGraphics](https://gabrielgambetta.com/computer-graphics-from-scratch/)