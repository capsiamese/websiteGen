---
title: "Pathfinding"
date: 2022-11-06T14:24:37+08:00
draft: true
tags: ["cg", "algorithm"]
---

## 1.通用算法
这个部分对于整个系列是最重要的, 这个部分是你理解路径查找必须掌握的, 其他的(包括A*)只是细节. 这一小节会给你一些启示.

这部分非常的简单.

让我们试着将上面的例子模式化成一些伪代码.

我们需要继续追踪那些我们已经知道如何从开始到达的节点. 如果是刚刚开始, 这里只有一个起始节点, 但是当我们开始遍历网格时, 我们会指出如何到达这些节点, 我们叫它可到达列表(reachable).

`reachable = [start_node]`

同时我们还需要一个追踪我们已经到达过的节点的列表(explored).

`explored = []`

这就是我们的核心算法: 在每一步的搜索中, 我们选择一个可达的并且没有到达过的节点, 然后查看这个节点可以到达其他的节点. 如果我们找到了我们要找的节点, 结束否则重复上面的步骤.

我们一直查找直到到达目标节点(这种情况, 我们已经找到了一条通向目标节点的路径), 或者直到遍历完所有节点(另一种情况, 没有一条通往目标节点的路径).

`where reachable is not empty`

我们选择一个已知并且没有到达过的节点.

`node = choose_node(reachable)`

如果我们找出如何到达目标点时, 搜索完成, 然后我们只要根据之前的节点构建出一条从开始到结束的路径.

```pseudocode
if node == goal_node
  path = []
  while node != None
    path.add(node)
    node = node.previous
  return path
```

这里并没有指出一个节点是否访问了多次, 所以我们继续下面的操作.

`reachable.remove(node)`
`explored.add(node)`

我们要弄清楚那个节点是从上个节点到达的, 我们从这个节点的临近节点列表中移除那些我们已经浏览过的.

`new_reachable = get_adjacent_nodes(node) - explored`

然后遍历这些新节点

`for adjacent in now_reachable`

如果我们已经知道如何到达了就无视他, 否则我们添加该节点到`reachable`列表中继续追踪.

```pseudocode
if adjecent not in reachable
  adjacent.previous = node
  reachable.add(adjecent)
```

退出循环的一个方式是找到目标, 另一种是`reachable`列表为空, 如果我们检查了所有的节点, 并且我们没有找到目标节点, 这意味这没有路径到达目标节点.

`return None`

下面是完整代码

```
function find_path(start_node, end_node):
  reachable = [start_node] # 可到达的点集合
  explored = [] # 已浏览过的点的集合

  while reachable is not empty:
    node = choose_node(reachable) # 从可到达集合取出一个

    if node == goal_node:
      return build_path(goal_node)
    
    reachable.remove(node) # 将该点从可到达的集合移除
    explored.add(node) # 加入到已浏览过的集合中

    new_reachable = get_adjacent_node(node) - explored # 找到该点的所有没有到达过的离近点
    for adjacent in new_reachable:
      if adjecent not in reachable: # 临近点没有在可到达集合
        adjecent.previous = node # 将临近点的前驱设为该点
        reachable.add(adjecent) # 将临近点加入可到达集合
    
  return None

function build_path(to_node):
  path = []
  while to_node != None:
    path.add(to_node)
    to_node = to_node.previous
  return path
```

## 2.搜索策略

### 秘方
解释搜索算法行为的关键在于`node = choose_node(reachable)`

## 路径长度很重要
在深入研究`choose_node`的不同行为之前, 我们需要修复一处小小的错误.

当我们选择其中一个当前节点的临近节点时, 我们忽略那些已经到达过的节点.

```
if adjacent not in reachable:
  adjacent.previous = node 
  reachable.add(adjacent)
```

这里的错误是: 如果我们发现了一个更好的到达目的地的节点, 在这个例子中, 我们应该调整节点的前置节点反应出这是个更短的路径.

要实现这种方式, 我们需要知道从起始点到所有可达节点的长度, 将他叫做`cost` of path, 现在让我们假设从一个节点移动到他相邻的节点的花费时1.

在开始路径查找之前我们将所有节点的花费都设置成+inf, 这将使任何路径都比这个值小, 同时将`start_node`的花费设置成0.

```
if adjacent not in reachable:
  reachable.add(adjacent)

# 如果临近节点是一个新的路径 或更短的路径
if node.cost + 1 < adjacent.cost:
  adjacent.previous = node
  adjacent.cost = node.cost + 1

```
### 统一按cost搜索
现在让我们将注意力放到`choose_node`上, 随机选择一个节点在查找最短路径时显然不是一个好的办法.

一个好的主意是我们选择从开始节点到临近节点花费最少的节点, 这样通常会优先选择最短的路径, 并且这不代表较长的路径不会被选择到. 这个算法从开始选择一条有效路径时应该尽可能短.

```
function choose_node(reachable):
  best_node = None
  for node in reachable:
    if best_node == None or best_node.cost > node.cost:
      best_node = node

  return best_node
```

## 3.A*揭秘

### The Magical Algorithm
想想一下我们运行这个搜索算法在一个可以做一些黑箱操作的芯片的计算机上, 有了这个芯片我们可以表达`choose_node`用一个非常简单的方式保证选择出的最短路径不会将时间浪费在一些永远不会到达的地方.

```
function choose_node (reachable):
  return magic(reachable, "whatever node is next in the shortest path")
```

但是这个奇妙的芯片仍然需要一些底层的代码, 这可能是个好的近似值.
```
function choose_node(reachable):
  min_cost = infinity
  best_node = None

  for node in reachable:
    cost_start_to_node = node.cost
    cost_node_to_goal = magic(node, "shortest path to the goal")
    total_cost = cost_start_to_node + cost_node_to_goal

    if min_cost > total_cost:
      min_cost = total_cost
      best_node = node
  
  return best_node
```
这是一个非常好的方法用来选择下一个节点, 我们选择一个产生一条从起点到目标节点的最短的节点.

### The Non-Magical but Pretty Awesome A*
不幸的是这个芯片和我们的硬件不兼容, 大多数的代码都是好的除了这一行
`cost_node_to_goal = magic(node, "shortest path to the goal")`

因此我们不能使用magic知道我们还没有浏览过的节点的花费, 好了, 让我们猜一下, 我们是乐观的, 所以我们将假设没有任何东西在我们的起点和终点之间, 我们就直接走过去.
`cost_node_to_goal = distance(node, goal)`

注意这里最短路径和最小距离是不一样的, 最小距离是假设他们之间绝对没有任何障碍.

这个估值会有一点简单, 在这个网格样例中, 两点之间的距离叫做Manhattan distance(曼哈顿距离), (abs(Ax - Bx) + abs(Ay - By)), 如果能对角线移动则是sqrt((Ax-Bx)^2 + (Ay-By)^2)等等, 最重要的是永远不要过高的估算成本.

```
function choose_node(reachable):
  min_cost = infinity
  best_node = None

  for node in reachable:
    cost_start_to_node = node.cost
    cost_node_to_goal = estimate_distance(node, goal_node)
    total_cost = cost_start_to_node + cost_node_to_goal

    if min_cost > total_cost:
      min_cost = total_cost
      best_node = node
  return best_node
```

#### 完整代码
```
function find_path(start_node, goal_node):
  reachable = [start_node]
  explored = []

  while reachable is not empty:
    node = choose_node(reachable)
    if node == goal_node:
      return build_path(goal_node)

    reachable.remove(node)
    explored.add(node)
    
    new_reachable = get_adjacent_nodes(node) - explored
    for adjacent in new_reachable:
      if adjacent not in reachable:
        reachable.add(adjacent)
      
      if node.cost + 1 < adjacent.cost:
        adjacent.previous = node
        adjacent.cost = node.cost + 1
  
  return None

function build_path(goal_node):
  path = []
  while to_node != None:
    path.add(to_node)
    to_node = tode.previous
  return path

function choose_node(reachable):
  min_cost = infinity
  best_node = None

  for node in reachable:
    cost_start_to_node = node.cost
    cost_node_to_goal = estimate_distance(node, goal_node)
    total_cost = cost_start_to_node + cost_node_to_goal

    if min_cost > total_cost:
      min_cost = total_cost
      best_node = node

  return best_node
```

## 4.实用A*


# Ref
[Pathfinding Demystified](https://gabrielgambetta.com/generic-search.html)
[RBG](https://www.redblobgames.com/)