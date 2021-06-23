
<h1 id=""> </h1>







## 1 

 
### 1.1  

* 请求方式：POST
* 请求路径：/trpc.playground.emptyTest.Hello/SayHello


#### 请求参数
 





##### request body参数描述
|  名称   | 类型  |是否必须| 描述  |
|  ----  | ----  | ---- | ---- |
|  age  |    [google.protobuf.UInt64Value](#schema_google.protobuf.UInt64Value)    |  否  |   无   |
|  name  |    string     |  否  |   无   |
##### 请求json示例: 
```json
{
    "age": {
        "value": 0
    },
    "name": "string"
}
```


 




#### http 状态码为：200,描述：
##### response body描述
|  名称   | 类型  |是否必须| 描述  |
|  ----  | ----  | ---- | ---- |
|  code  |     number      |  否  |   无    |
|  msg  |     string      |  否  |   无    |
##### 响应json示例:

```json
{
    "code": 0,
    "msg": "string"
}
```






## 数据定义


<h3 id='schema_HelloReq'>HelloReq</h3>

```json
{
    "age": {
        "value": 0
    },
    "name": "string"
}
```
|  名称   | 类型  |是否必须| 描述  |
|  ----  | ----  | ---- | ---- |
|  age  |     [google.protobuf.UInt64Value](#schema_google.protobuf.UInt64Value)   |  否  |   无    |
|  name  |     string    |  否  |   无    |



<h3 id='schema_HelloRsp'>HelloRsp</h3>

```json
{
    "code": 0,
    "msg": "string"
}
```
|  名称   | 类型  |是否必须| 描述  |
|  ----  | ----  | ---- | ---- |
|  code  |     number    |  否  |   无    |
|  msg  |     string    |  否  |   无    |



<h3 id='schema_google.protobuf.UInt64Value'>google.protobuf.UInt64Value</h3>

```json
{
    "value": 0
}
```
|  名称   | 类型  |是否必须| 描述  |
|  ----  | ----  | ---- | ---- |
|  value  |     number    |  否  |   无    |



## 枚举类型定义

