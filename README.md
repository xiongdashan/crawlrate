定时抓取汇率
============

##设定任务：

从xe.com上抓取指定的汇率，并保存在数据库中，同时把抓取的所有币种的汇率发送到指定邮箱

## 提供Api

查询汇率，返回最新的抓取记录


```
GET 

http://localhost/currency/:convert

convert: 当前币种-目标币种，如：

http://127.0.0.1:1000/currency/CNY-HKD

```

返回


```json

{
    "Convert": "CNY-HKD", //转换的币种
    "Rate": 1.17302,//汇率
    "Date": "2018-07-15" //汇率对应日期
}

```