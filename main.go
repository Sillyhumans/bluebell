package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	mySnowflake "bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	// 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err: %v\n", err)
		return
	}

	fmt.Println(settings.Conf.Name)

	// 加载日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init Logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync() // 用来把缓冲区的日志加载
	zap.L().Debug("logger init success...")

	// 加载mysql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err: %v\n", err)
		return
	}
	defer mysql.Close()

	// 加载redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err: %v\n", err)
		return
	}
	defer redis.Close()

	// 注册路由
	r := routes.Setup()

	if err := mySnowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err: %v\n", err)
		return
	}

	// 初始化gin框架内部用的翻译器(报错翻译)
	err := controller.InitTrans("zh")
	if err != nil {
		fmt.Printf(", err: %v\n", err)
	}

	err = r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
	select {}

	// 启动服务
	//srv := &http.Server{
	//	Addr:    fmt.Sprintf("%d:%d", settings.Conf.Host, settings.Conf.Port),
	//	Handler: r,
	//}
	//go func() {
	//	// 开启一个goroutine启动服务
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		zap.L().Fatal("listen: %s\n", zap.Error(err))
	//	}
	//}()
	//// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	//quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	//// kill 默认会发送 syscall.SIGTERM 信号
	//// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	//// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	//// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	//<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	//zap.L().Info("ShutDown Server....")
	//// 创建一个5秒超时的context
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	//if err := srv.Shutdown(ctx); err != nil {
	//	zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	//}
	//
	//log.Println("Server exiting")
}
