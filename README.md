static_db
=================
静态库，用kmean聚类加快实现图片的特征搜索，入库，删除，查询等操作
-----------------
 version1版本，run main 即可，产生的随机数组成库，因此kmean迭代不收敛，自行更换

 version2版本，十个nodefirst 0.5起始 间隔1 100个node,每10个node在0-1，1-2...之间随机，因此聚类很均匀，10个

 version3版本 wf.go发送http请求，通过直接15行的payload然后执行db_test.go进行基准测试

 version4版本 修改了一些变量名,更改了http请求的方式，用标准json格式，共有两种调用方式，一种是一个url，通过post body确定add 
 search delete，一种是通过url 默认通过url，同时可以使用grpc方式发送message,common文件里有解码编码函数，用于和上游对接
未实现：http转static_proto的test函数，test和server_proto通信，最后调用libworker返回

version5版本 （gateway分支） 使用gateway让它代替人工进行http请求的反序列化，减少代码冗余，一开始空库，通过api_ips的ips_add添加然后搜索，
调用ips接口从图片获取特征然后存至DB，ips_search进行搜索，之后可以将ips替换
       


