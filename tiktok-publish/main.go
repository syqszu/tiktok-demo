/*
 * @Author: shanshan
 * @Date: 2023-8-12
 * @Description: Publish RPC Server 端初始化
 具体：
    # Publish RPC Server 端配置初始化
    1. 导入所需的包。
    2. 初始化配置信息。
    3. 初始化JWT对象。
       - 从配置中获取JWT签名密钥，并使用它创建JWT对象。
    # Publish RPC Server 端运行
    1. 初始化日志。
    2. 设置日志记录器。
    3. 创建Etcd注册中心实例，并指定Etcd地址。
    4. 解析服务地址。
       - 使用服务器地址和端口创建TCP地址。
    5. 创建OpenTelemetry提供者实例。
       - 使用指定的服务名称和导出端点创建OpenTelemetry提供者。
       - 使用不安全模式创建提供者。
    6. 初始化。
       - 执行初始化操作，可能是一些额外的初始化工作。
    7. 创建RPC服务器实例。
       - 创建一个新的Publish RPC服务器实例。
       - 指定服务实现对象。
       - 设置服务器地址。
       - 添加常用中间件。
       - 添加服务器中间件。
       - 指定注册中心。
       - 设置限制选项，包括最大连接数和最大QPS。
       - 使用Multiplex传输。
       - 使用OpenTelemetry提供的服务器套件。
       - 设置服务器的基本信息，包括服务名称。
    8. 运行RPC服务器。
       - 启动RPC服务器并等待请求。
 */
package main

import (
	"context"
	"fmt"
	"net"
	etcd 	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

var (
	Config      = ttviper.ConfigInit("TIKTOK_PUBLISH", "publishConfig")// 初始化配置信息
	ServiceName = Config.Viper.GetString("Server.Name")// 获取服务名称
	ServiceAddr = fmt.Sprintf("%s:%d", Config.Viper.GetString("Server.Address"), Config.Viper.GetInt("Server.Port"))// 获取服务地址
	EtcdAddress = fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))// 获取Etcd地址
	Jwt         *jwt.JWT   //JWT对象
)

// Publish RPC Server 端配置初始化
func Init() {
	dal.Init()
	Jwt = jwt.NewJWT([]byte(Config.Viper.GetString("JWT.signingKey")))  // 使用配置中的JWT签名密钥创建JWT对象
}

// Publish RPC Server 端运行
func main() {
	var logger = dlog.InitLog(3) // 初始化日志记录器
	defer logger.Sync()  // 在main函数结束时关闭日志记录器

	klog.SetLogger(logger) // 设置日志记录器

	r, err := etcd.NewEtcdRegistry([]string{EtcdAddress}) // 创建Etcd注册中心实例
	if err != nil {
		klog.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", ServiceAddr) // 解析服务地址
	if err != nil {
		klog.Fatal(err)
	}

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(ServiceName),  // 设置服务名称
		provider.WithExportEndpoint("localhost:4317"),  // 设置导出端点
		provider.WithInsecure(), // 使用不安全模式
	)
	defer p.Shutdown(context.Background())   // 在main函数结束时关闭OpenTelemetry提供者

	Init()

	svr := publish.NewServer(
		new(PublishSrvImpl),
		server.WithServiceAddr(addr),                                       // address
		server.WithMiddleware(middleware.CommonMiddleware),                 // middleware
		server.WithMiddleware(middleware.ServerMiddleware),                 // middleware
		server.WithRegistry(r),                                             // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(tracing.NewServerSuite()),                         // trace
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: ServiceName}),
	)

	if err := svr.Run(); err != nil {
		klog.Fatalf("%s stopped with error:", ServiceName, err)
	}
}