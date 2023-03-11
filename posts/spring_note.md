---
title: "Spring总结"
date: 2023-03-02T11:33:24+08:00
draft: true
tags: []
---

## 基于XML的IOC
要使用spring容器需要在maven中先导入相关依赖.

```xml
<dependency>
  <groupId>org.springframework</groupId>
  <artifactId>spring-context</artifactId>
  <version>5.2.10.RELEASE</version>
</dependency>
```

然后在resources目录下新建一个applicationContext.xml文件, 该文件告诉spring那些对象需要管理.

在spring中对象叫做Bean, 所以在xml中让spring管理的对象就是Bean.

每个bean可以指定id和别名, 并且需要类的全路径来告诉spring从那个类来实例化对象. scope表示这个bean的类型, 默认是单例的. lazy-init属性默认为false, 表示不惰性初始化.

property属性用来给bean中的字段注入依赖, name属性是setter方法的名字去掉set并首字母小写, ref可以是id或者name, value是字面量. 还可以给bean添加初始化函数和销毁函数通过bean属性字段或者实现`InitializingBean`或者`DisposableBean`接口. 同时context需要显示的关闭或者调用`registerShutdownHook`函数.

构造器注入可以用按name匹配,index匹配,type匹配.

autowire属性表示开启自动装配这个bean的属性, 值可以是byType, byName.

可以给bean的字段注入集合, 只需要在property下添加子标签. 可以注入的集合有

1. <array><value></value></array>
2. <list><value></value></list>
3. <set><value></value><set>
4. <map><entry key="" value=""/></map>
5. <props><prop key="" value=""/></props>

```xml
<bean id="aName" name="a,b,c d; e" class="classpath.full.path.name" scope="prototype" init-method="init" destroy-method="destroy">
  <property name="setMethodName" ref="beanId">
  <property name="setMethodName" value="123">
  <constructor-arg name="argName" ref=""/>
  <constructor-arg type="java.lang.String" value="111"/>
  <constructor-arg index="1" ref=""/>
</bean>
```

bean还可以通过工厂来实例化
1. 静态工厂 需要在bean标签上添加`factory-method`属性来指定是类中的那个方法.
2. 实例工厂 需要先声明一个实例工厂的bean, 然后在需要使用该工厂的bean中添加`factory-method`和`factory-bean`属性
3. FactoryBean实例化 需要一个实现了`FactoryBean<T>`的接口的类

### 加载外部数据源

需要在xml中添加context名称空间

```xml
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:context="http://www.springframework.org/schema/context"
       xsi:schemaLocation="
            http://www.springframework.org/schema/beans
            http://www.springframework.org/schema/beans/spring-beans.xsd
            http://www.springframework.org/schema/context
            http://www.springframework.org/schema/context/spring-context.xsd
            "> 

    <context:property-placeholder location="classpath*:*.properties" system-properties-mode="NEVER">
</beans>
```
通过`<context:property-plcaholder>`标签来引用property文件, system-properties-mode属性表示是否引用系统环境变量.

加载完的properties可以在value中通过`${name}`使用.

## 基于注解的IOC容器

要使用注解来初始化spring容器需要先声明一个类作为入口点. 这个类需要有`@Configuration`注解, 同时添加`@ComponentScan`注解告诉spring那些包需要需要扫描bean. `@PropertySource`注解用来加载property属性文件. 使用时通过`@Value("${name}")`注入值.

使用的时候需要实例化`AnnotationConfigApplicationContext`类来初始化容器.

如果要将一个对象交给spring管理需要在类上添加`@Component`注解, 该注解的作用和`<bean>`标签一样. 同时spring为了区分不同的业务添加了`@Service`, `@Repository`,`@Controller`等注解作为`@Component`的别名.

如果要给一个对象注入依赖需要在需要注入的字段上添加`@Autowired`注解, 并且可以添加`@Qualifier`表示自动注入的类名或类型.

### 管理第三方bean

如果要让spring管理项目需要的第三方bean的话需要先创建一个类, 该类将作为创建该第三方bean的工厂. 在工厂方法上使用`@Bean`注解表示该方法将提供第三方对象. 然后在住配置类上添加`@Import`注解来导入该第三方配置类, 或者直接在该类上添加`@Configuration`注解(不建议)

### junit中使用spring

首先需要设置类运行器`@RunWith(SpringJUnit4ClassRunner.class)`, 然后设置spring环境对应的配置类`@ContextConfiguration(class=SpringConfig.class)`. 之后就可以在测试中使用注解注入spring容器管理的bean了.

## Spring AOP

要使用spring aop需要先导入aspectj的包, 然后在主配置类上添加`@EnableAspectJAutoProxy`注解.
```xml
<dependency>
  <groupId>org.aspectj</groupId>
  <artifactId>aspectjweaver</artifactId>
  <version>1.9.4</version>
</dependency>
```

### 使用AOP

首先需要创建一个类并交给spring管理只要在类上添加`@Component`和`@Aspect`注解.

### 概念

* @Pointcut("execution()") 切入点表达式
  * @Pointcut("execution(void com.xxx.dao.update())")
  * @Pointcut("execution(void com.xxx.dao.impl.yyy.zzz())")
  * @Pointcut("execution(* com.xxx.dao.impl.yyy.zzz(*))")
  * @Pointcut("execution(void com.*.*.*.zzz())")
  * @Pointcut("execution(* *..*(..))")
  * @Pointcut("execution(* *..*e(..))")
  * @Pointcut("execution(void com..*())")
  * @Pointcut("execution(* com.xxx.*.*Service.zzz*(..))")
  * @Pointcut("execution(* com.xxx.*.*Service.save(..))")

* Advice类型
  * @Before("")  JoinPoint
  * @After("")   JoinPoint
  * @Around("")  ProceedingJoinPoint
  * @AfterReturing(value="", returning="") JoinPoint returnType  returning的值要和returnType参数的名字相同
  * @AfterThrowing(vaule="", throwing="")


## Spring事务 Mybatis Druid

### Mybatis

开发步骤

1. 创建数据表
2. 创建实体类
3. 编写映射文件ClassNameMapper.xml
4. 编写SQLMapCofnig.xml

#### 创建实体类和sql表之间的关系

mybatis会根据xml中描述的内容将实体类和sql操作进行映射.

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE mapper
  PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
  "https://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="org.mybatis.example.BlogMapper">
  <select id="selectBlog" resultType="Blog">
    select * from Blog where id = #{id}
  </select>
</mapper>
```

sql中的变量使用`#{val}`, 返回值用resultType, 参数用parameterType

#### 配置核心约束文件

mybatis的核心配置是一个`SqlSessionFactory`类, 该类将通过这个xml来实例化. 之后操作数据库只用打开一个session就行了(别忘了关闭). **需要注意mybatis默认开启事务并且不自动提交**.

```java
SqlSession sess = sqlSessionFactory.openSession();
sess.selectList("ns.id");

sess.commit();
sess.close();
```

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE configuration
  PUBLIC "-//mybatis.org//DTD Config 3.0//EN"
  "https://mybatis.org/dtd/mybatis-3-config.dtd">
<configuration>
  <properties resource="jdbc.properties"/>

  <typeAliases>
    <typeAliase type="fullpathname" alias="name"/>
  </typeAliases>

  <environments default="development">
    <environment id="development">
      <transactionManager type="JDBC"/>
      <dataSource type="POOLED">
        <property name="driver" value="${driver}"/>
        <property name="url" value="${url}"/>
        <property name="username" value="${username}"/>
        <property name="password" value="${password}"/>
      </dataSource>
    </environment>
  </environments>
  <mappers>
    <mapper resource="org/mybatis/example/BlogMapper.xml"/>
  </mappers>
</configuration>
```

#### 代理开发方式

使用这种方式开发首先需要定义数据操作的接口, 并且在mapper中的namespace必须是这个接口的全限定名.
mapper中操作的id必须和类中的方法名相同, resultType, parameterType也要和方法的参数相同.
获取代理对象通过`sess.getMapper(interface.class)`然后通过代理对象访问数据库.

#### 动态sql
```xml
<sql id="aaa"> common sql </sql>

<select id="findByCond" resultType="user" parameterType="user"> 
  <include refid="aaa"/>

  select * from user
  <where>
    <if test="id!=0">
        and id=#{id}
    </if>
    <foreach collection="collection type" open="id in(" close=")" item="id" separator=",">
         #{id}
    </foreach>
  </where>
</select>
```

#### typeHandlers处理sql和java之间的类型转换
继承`BaseTypeHandler<T>`并且在核心文件中注册.
```xml
<typeHandlers>
  <typeHandler handler="xxx">
</typeHandlers>
<plugins>
  <plugin interceptor="分页标签插件">
</plugins>
```

#### 多表操作关系映射
需要手动指定实体如何对应.
```xml
<resultMap id="" type="">
  <id column="" property=""/>
  <result column="" property="">
  <association property="" javaType="">
    <id = column="" property="">
  </association>
</resultMap>

<select resultMap="id">
</select>
```

### 注解开发
常用注解
1. @Insert
2. @Update
3. @Select
4. @Result 
5. @Results
6. @One
7. @Many

使用@Mapper来标注接口是一个mapper类.

这些注解直接加到dao接口中的方法上. 然后核心配置中的mappers中指定要扫描的包`<package name="package  full name"/>`

对于复杂类型使用`@Result`和`@Results`

### 整合spring

```java
public class MybaticConfig {
  @Bean
  public SqlSessionFactoryBean sqlSessionFactory(DataSource dataSource) {
    SqlSessionFactoryBean b = new SqlSessionFactoryBean();
    b.setTypeALiasesPackage("xxx");
    b.setDataSource(dataSource);
    return b;
  }

  @Bean
  public MapperScannerConfigure mapperScannerConfigure() {
    MapperScannerConfigure msc = new MapperScannerConfigure();
    msc.setBasePackage("xx");
    return msc;
  }
}
```


@EnableTransactionManagement   PlatformTransactionManager  @Transactional


## SpringMVC

引入依赖
```xml
<dependency>
  <groupId>javax.servlet</groupId>
  <artifactId>javax.servlet-api</artifactId>
  <version>3.1.0</version>
  <scope>provided</scope>
</dependency>
<dependency>
  <groupId>org.springframework</groupId>
  <artifactId>spring-webmvc</artifactId>
  <version>5.2.10.RELEASE</version>
</dependency>
 <dependency>
  <groupId>com.fasterxml.jackson.core</groupId>
  <artifactId>jackson-databind</artifactId>
  <version>2.9.0</version>
</dependency>

<plugin>
  <groupId>org.apache.tomcat.maven</groupId>
  <artifactId>tomcat7-maven-plugin</artifactId>
  <version>2.1</version>
  <configuration>
    <port>80</port>
    <path>/</path>
  </configuration>
</plugin>
```

springmvc的主配置类和spring主配置类的注解基本相同,springmvc需要开启`@EnableWebMvc`注解, 但是springmvc需要有一个继承了`AbstractDispatcherServletInitializer`的类.
并在这个类中分别初始化spring和springmvc的ioc容器. 也可以继承`AbstractAnnotationConfigDispatcherServletInitializer`类简化.
```java
public class ServletInitConfig extends AbstractDispatcherServletInitializer {
  protected WebApplicationContext createServletApplicationContext() {
    AnnotationConfigWebApplicationContext ctx = new AnnotationConfigWebApplicationContext(SpringMVCConfig.class);
    return ctx;
  }
  protected String[] getServletMappings() {
    return new String[]{"/"};
  }
  protected WebApplicationContext createRootApplicationContext() {
    AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(SpringConfig.class);
    return ctx;
  }

  @Override
  // 设置乱码处理
  protected Filter[] getServletFilters() {
    CharacterEncodingFilter filter = new CharacterEncodingFilter();
    filter.setEncoding("UTF-8");
    return new Filter[]{filter};
  }
}
```

使用两个ioc容器时需要注意排除目录可以使用`@ComponentScan(excludeFilters = @ComponentScan.Filter(type=FilterType.ANNOTATION, class=Controller.class))`排除.

springmvc中的路由类需要使用`@Controller`注解. 其中路由路径通过`@RequestMapping("")`声明, 可以用在类上表示路由组, 用在方法上表示具体路由.
`@ResponseBody`表示方法的返回值放到响应体中.

如果要将请求的查询参数放到方法的参数上需要使用`@RequestParam`注解, 默认情况下如果查询参数key的名字和参数名相同会自动匹配, 不用加注解.
使用`@RequestBody`在方法参数前可以将请求体绑定到对象上. `@DataTimeFormat`用来格式化时间.

方法的返回值如果是字符串并且没有`@ResponseBody`首先匹配jsp页面, 有该注解会自动解析返回类型.

### Rest风格注解

对于开发Rest风格的路由, springmvc将`@Controller`和`@ResponseBody`整合成了`@RestController`简化.
同时将`@RequestMapping(vaule="", method="")`简化成了`@GetMapping(value="")`(Http方法+Mapping).

对于路径参数使用`@PathVariable`注解放到方法参数前. 路径参数变量使用`@GetMapping("/{var}")`其中var要和方法参数名一致.

### 异常处理

使用`@RestControllerAdvice`标识当前类为异常处理器, `@ExceptionHandler(exception.class)`标识该方法处理那种异常.

### 静态资源

将webapp中的静态资源交给springmvc处理, 需要增加一个配置类该类需要继承`WebMvcConfigurationSupport`类,并覆盖`addResourceHandlers`方法.
在该方法中使用`reg.addResourceHandler("").addResourceLocations("")`添加资源映射.

### 拦截器

要增加拦截器需要在实现了`WebMvcConfigurationSupport`或者`WebMvcConfigure`类中覆盖`addInterceptors`方法然后调用`reg.addInterceptor(interceptor).addPathPatterns("", "")`增加拦截器. 增加的拦截器需要通过DI注入.

这里需要注意使用`WebMvcConfigure`具有一定的侵入性.

定义拦截器实现`HandlerInterceptor`接口, 并且交给ioc容器.

前置拦截器的调用顺序为增加顺序. 后置拦截器为反向. 如果前置拦截器中断了调用链, 则后置拦截器还会执行该拦截器之前注册的拦截器的后置方法.

## SpringBoot

springboot使用`@ConfigurationProperties`来定义yaml中要有那些类容.

yaml中使用`---`来分割块使得不同块中可以有同名字段.

在maven中安装插件, 来动态使用yaml配置源
```xml
<build>
    <plugins>
        <plugin>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-maven-plugin</artifactId>
        </plugin>
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-resources-plugin</artifactId>
            <version>3.2.0</version>
            <configuration>
                <encoding>UTF-8</encoding>
                <useDefaultDelimiters>true</useDefaultDelimiters>
            </configuration>
        </plugin>
    </plugins>
</build>

<profiles>
    <!--开发环境-->
    <profile>
        <id>dev</id>
        <properties>
            <profile.active>dev</profile.active>
        </properties>
    </profile>
    <!--生产环境-->
    <profile>
        <id>pro</id>
        <properties>
            <profile.active>pro</profile.active>
        </properties>
        <activation>
            <activeByDefault>true</activeByDefault>
        </activation>
    </profile>
    <!--测试环境-->
    <profile>
        <id>test</id>
        <properties>
            <profile.active>test</profile.active>
        </properties>
    </profile>
</profiles>
```

要在测试中使用springboot直接用`@SpringBootTest(class=xxx.class)`注解.
