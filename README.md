static_db
=================
静态库，用kmean聚类加快实现图片的特征搜索，入库，删除，查询等操作
-----------------
 version1版本，run main 即可，产生的随机数组成库，因此kmean迭代不收敛，自行更换即可  

 version2版本，十个nodefirst 0.5起始 间隔1 100个node,每10个node在0-1，1-2.。。之间随机，因此聚类很均匀，10个

 version3版本 wf.go发送http请求，通过直接15行的payload然后执行db_test.go进行基准测试

 version4版本 修改了一些变量名，更”专业“,更改了http请求的方式，用标准json格式，更易懂，然后有两种调用方式，一种是一个url，通过post body确定add 
 search delete，一种是通过url 默认通过url，同时可以使用grpc方式发送message,同时common文件里有解码编码函数，用于和上游对接

未实现：http转static_proto的test函数，test和server_proto通信，最后调用libworker返回       


