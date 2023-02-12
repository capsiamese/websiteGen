---
title: "Learncss"
date: 2022-08-06T15:46:06+08:00
draft: true
tags: ["js"]
---

## display
css中最重要的用于控制布局的属性, 每个元素都有默认的display值`inline`或`block`.

`block`通常被叫做块级元素,`inline`叫做行内元素.

### block
`div`是一个标准的块级元素, 块级元素会新开始一行并且尽可能的撑满放置这个元素的容器, 其他常用块级元素`p`, `form`, `header`, `footer`, `section`...

### inline
`span`是一个标准的行内元素, 一个行内元素可以在段落中包裹一些文字而不打乱布局, 其他常用行内元素`a`...


### none
`none`属性一般用来不删除元素的情况下不显示元素,  并且不占据空间.

`visiablity`属性设置成`hidden`也不会显示元素, 但是会占据空间.

[display属性所有值](https://developer.mozilla.org/en-US/docs/Web/CSS/display)

## width
设置块级元素的`width`可以防止他从左到右撑满整个容器, 可以将左右边距设置成auto来将块级元素剧中, 但是如果当浏览器窗口小于元素`width`时会会出现滚动条.

这时可以设置`max-width`来替代`width`, 如果浏览器窗口大于`max-width`块级元素的最大长度就是`max-width`, 如果浏览器窗口小于`max-width`则块级元素会相应缩小, 而不是出现滚动条.

## 盒模型

当设置了元素的宽度, 但是实际展现的元素却超出了设置, 这是因为元素的边框和内边框会撑开元素, 内边距和外边距会增加他们的宽度.

```css
.small {
    width: 500px;
    margin: 20px auto;
}
.big {
    width: 500px;
    margin: 20px auto;
    padding: 50px;
    border-width: 10px;
}
```

此时虽然两个css的宽度相同, 但是实际大小却是`.big`会更大, 会在`width`的基础上向外扩展`padding * 2`的大小.

如果新增一个`box-sizing: border-box`属性, 此时内外边距都不会再向外增加边距了, 而是向内挤压.

## position
`static`是position的默认值, 任意`postion:static`的元素不会被特殊定位, 表示不会被`positioned`, 其他的属性会被`positioned`.

`relative`表现的和`static`一样如果不加额外属性的话. 如果在一个相对定位的元素上设置`top`, `right`, `bottom`, `left`属性的话, 会使其偏离正确位置.其他元素不会受该属性的元素影响而偏离原来的位置.

`fixed`固定定位的元素不会随着视窗的滚动而滚动.

`absolute`与`fixed`类似, 但是不是相对于视窗而是相对于最近的`positioned`祖先的元素. 如果没有`positioned`祖先则相对于整个文档的body元素, 并且会随着页面滚动.

1. float 可实现文字环绕图片.
2. clear属性被用于控制浮动. 可以使用和float浮动方向相同的值来将元素移动到浮动元素之后
3. overflow: auto 如果浮动子元素超过该元素的大小, 该属性会扩展该元素直到将浮动元素包住.
4. width: 50% 百分比是一种相对于包含块的计量单位, 

## @media
媒体查询可以实现响应式设计
```css
@media (min-width:600px) {
    nav {
        float: left;
        width: 25%;
    }
    section {
        margin-left: 25%;
    }
}
@media (max-width:599px) {
    nav li {
        display: inline;
    }
}
```
```html
<nav>
    <ul>
        <li></li>
    </ul>
</nav>
<section></section>
<section></section>
```
当宽度大于600px时nav和section呈现出左右排列, 宽度小于600px是nav和section恢复元素默认值, 列表变成inline.

[media queries](https://developer.mozilla.org/en-US/docs/Web/CSS/Media_Queries/Using_media_queries)
[view port](https://dev.opera.com/articles/an-introduction-to-meta-viewport-and-viewport/)

## inline-block
```css
.box {
    float: left;
    width:200px;
    height:100px;
    margin: 1em;
}
.after-box{
    clear:left;
}
```
```html
<div class="box"></div>
<div class="box"></div>
<div class="box"></div>
<div class="box"></div>
<div class="box"></div>
<div class="after-box"></div>
```

该布局会呈现出很多方格来铺满浏览器, 对after-box使用clear来将元素放到一堆盒子下方.

```css
.box2 {
    display: inline-block;
    width:200px;
    height:100px;
    margin: 1em;
}
```
可以使用inline-block来达到同样的效果, 并且after-box不再需要clear.

使用inline-block需要注意
* vertical-align属性会影响inline-block元素
* 需要设置每一列的宽度
* 如果html源代码中元素之间有空格, 那么列于列之间会产生空隙.

## column
可以轻松实现文字的多列布局.

# flex

# reference

[学习css布局](https://zh.learnlayout.com/)