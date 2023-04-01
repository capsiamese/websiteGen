---
title: "rabbitMQ入门"
date: 2023-03-11T14:39:19+08:00
draft: true
tags: []
---

mq的功能

1. 流量削峰
2. 应用解耦
3. 异步处理

rabbitmq: 接受, 存储, 转发消息

四大核心概念:

1. 生产者
2. 交换机
3. 队列
4. 消费者

核心部分

1. simple
2. work queues
3. pub/sub
4. routing
5. topics
6. publisher confirms

```
               Connection         Broker (RabbitMQ)             Connection
                                          /->  queue
                                 exchange-
producer ->      channel                  \-> queue                channel         -> consumer
                 channel    ->            /-> queue      ->        channel
producer ->      channel         exchange-                         channel         -> consumer
                                          \-> queue
```

每个生产者对应一个channel

Broker: 接受和分发消息

virtual host: 出于多租户和安全因素, 把amqp基本组件划分到一个虚拟的组钟, 类似(namespace), 多个不同的用户使用同一个rabbitmq时可以划分出多个

vhost, 每个用户在自己的vhost创建 exchange/queue等.

connection: 生产/消费者与服务器之间的tcp连接

channel: 类似goroutine

exchange: message到达broker的第一站, 根据分发规则匹配查询表中的routing key分发消息到queue中
常用类型direct(p2p) topic(pub/sub) fanout(multicast)

插件:rabbitmq-plugins enable xxx

1. rabbitmq_management web15672管理
2. rabbitmq-delayed-message-exchange 延迟消息exchange

消息丢失:

1. 自动应答

消息发送后立即被确认传送成功, 如果消息在接受到之前, 消费者出现连接或者关闭channel那么消息就丢失了

或者消费者积压消息来不及处理, 而进程又结束了, 消息也会丢失

2. 手动

发送确认

1. 单消息确认
2. 多消息确认
3. 异步确认

交换机(exchange), 发布订阅
交换机决定如何处理消息

1. 无名exchange
2. direct
3. topic
4. headers
5. fanout 广播到所有队列中

binding: exchange与queue之间的绑定关系

routing key: exchange与queue之间的绑定规则

## 死信队列

producer将消息发送到broker而consumer从queue中取出进行消费, 但是由于某些原因导致queue中的消息不能被消费

应用场景: 为了保证定安业务的消息数据不丢失, 需要使用死信队列机制, 当消息消费发生异常时, 将消息投入到死信队列中.

来源:

1. 消息ttl过期
2. 队列达到最大长度
3. 消息被拒绝

---

## 入门

### 队列

队列(Queue)是rabbitmq中的内部对象, 用来存储消息, 多个消费者可以订阅同一个队列, 这时rabbitmq会通过round-robin来分配消息到消费者,
消费者可以通过qos机制控制获取消息的数量, 来达到负载均衡.

rabbitmq不支持队列层面的广播, 可以使用redis pub/sub来实现.

### 交换器

交换器(exchange)是rabbitmq中的dispatcher用来将生产者的消息分配到不同的队列中去.

rabbitmq中的交换器共有四种类型, 如果要加不同的交换器类型可以通过安装插件添加.

1. Topic 约定路由键和绑定键使用'.'分隔的字符串, 绑定键中可以使用'*'(匹配多个)和'#'(匹配一个)做模糊匹配
2. Direct 路由键和绑定键完全匹配
3. Fanout
4. Headers 一般不用

### 路由键

路由键(routing-key), 生产者将消息发送给broker时会携带路由键, exchange会通过路由键和匹配规则将消息放入队列中.[相当于ip]

### 绑定键

绑定键(binding-key), 用来将exchange和队列绑定, 如果exchange是fanout将无视binding-key.[相当于mask]

## 函数

### Connection和Channel通知的最佳实践

链接和管道通知 Channel.NotifyClose, Connection.NotifyClose 当链接或者管道被关闭时, 如果是graceful close将不会被通知,
如果是网络异常, 在管理界面强制关闭等会通知错误.

`notifyConnCloseCh := conn.NotifyClose(make(chan *amqp091.Error, 1))`, 为了避免死锁, 管道里的内容必须消费.

```go
package main

func a() {
	go func() {
		for notifyConnClose != nil || notifyChanClose != nil {
			select {
			case err, ok := <-notifyConnClose:
				if !ok {
					notifyConnClose = nil
				} else {
					print(err)
				}
			case err, ok := <-notifyChanClose:
				if !ok {
					notifyChanClose = nil
				} else {
					print(err)
				}
			}
		}
	}()
}
```

### NotifyPublish通知的最佳实践

使用`Channel.NotifyPublish`允许调用者被通知,

### Channel

channel中的常用函数

#### Consume

定义: `Channel.Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error)`

Consume必须消费所有接收到的消息, 如果不接收连接上的所有方法.

amqp所有的deliveries都必须被确认, 需要在处理完成后调用Delivery.Ack, 如果消费者被取消了或者管道,链接关闭,
任何未被确认的消息都会重新入队.

所有的消费者使用一个字符串作为唯一标识, 如果是空的话, 库将会生成一个标识, 消费者标识将包含在每一个Delivery中的ConsumerTag字段.

autoAck如果是true的话, 服务器将会在deliveries被写入网络之前被确认, 所以这样会导致丢失消息, 并且不能再调用Delivery.Ack.

exclusive如果是true的话, 服务器将保证这条链接是队列中唯一的消费者.

noLocal 没用 ***建议对Publish和Consumer使用单独的Connection***

noWait如果是true的话, 不等待服务器确认请求并且立即开始deliveries. 如果无法消费, 通道将会异常关闭.

当channel或者connection被关闭时所有被缓存的或者正在路上的消息将会被丢弃.

#### ExchangeBind  ExchangeUnbind ExchangeDelete

#### ExchangeDeclare

#### ExchangeDeclarePassive

#### Flow Confirm

#### Get

#### Nack

#### NotifyCancel NotifyClose NotifyConfirm NotifyFlow NotifyPublish NotifyReturn

#### PublishWithContext

#### PublishWithDeferredConfirmWithContext

#### Qos

#### QueueBind

#### QueueDeclare

#### QueueDeclarePassive

#### QueueDelete QueuePurge QueueUnbind

#### Reject

#### Tx TxCommit TxRollback

### Delivery

#### Ack

#### Nack

#### Reject

### Publishing


## 死信队列

当消息在一个队列中编程死信之后, 就可以被重新发送到另一个交换器(死信交换器DLX)中, 绑定死信交换器的队列叫死信队列.

消息编程私信的几种情况
1. 消息被拒绝 Reject/Nack并设置requeue为false
2. 消息过期
3. 队列到达最大长度

通过在queueDeclare中设置`x-dead-letter-exchange`将队列添加到DLX中.

通过死信队列可以实现延迟队列.

## 常用mq结构

1. 简单队列 单生产者单消费者
2. 工作队列 单生产者多消费者
3. 发布订阅 单生产者多队列
4. 路由模式 根据routing-key选择队列
5. 主题模式 根据模式匹配选择队列
6. rpc 通过队列进行rpc
7. 发布确认模式 可靠的发布消息和发送者确认


## References

[Concepts](https://www.rabbitmq.com/tutorials/amqp-concepts.html)
[AMQP091 GoDocs](https://pkg.go.dev/github.com/rabbitmq/amqp091-go#Publishing)
[RabbitMQ实战指南](https://google.com)