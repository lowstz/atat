atat
===================

TODO
-------------------

*  SQL语句优化
*  搜索中文分词支持
*  /book/top 热门搜索 API
*  单书 目录
*  单书 标签
*  /user/:userid API 支持用户登陆，暂时没必要

### IN PROCESS
*  field images 有bug待修
*  优化 /book/search API的排序
*  安全性检测，比如SQL注入等
*  增加单元测试
*  上Redis，增加缓存功能

### DONE
* 各API增加fields参数
* 处理高并发下导致mysql连接池出现 Too many connections 的错误
* 增加 HTTP Basic API 验证 
