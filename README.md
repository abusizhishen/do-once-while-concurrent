# golang-do-something-once-while-cocurrent
解决并发时对同一资源的重复请求、缓存穿透、数据重复计算，如第三方接口重复请求、redis、数据库重复请求，避免对第三方依赖造成过大压力，以至于服务崩溃，影响到其他服务，同时配合分布式锁，可以解决分布式幂等问题

根据请求标识，对单机同一时刻同一资源的重复请求做限制，只有一个请求等获得锁，其他请求均处于阻塞状态，直到获得执行权限的线程执行结束，释放信号。此时所有请求均能得到数据。

![img](https://github.com/abusizhishen/justOnceWhileCocurrent/blob/master/example.jpg?raw=true)
