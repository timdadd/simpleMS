package common

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//
//func (c *AppConfig) ConnGRPC(serviceName string,conn **grpc.ClientConn) {
//	var err error
//	c.KeyPrefix(serviceName)
//	addr := c.ServiceAddress()
//	*conn, err = grpc.DialContext(c.Ctx, addr,
//		grpc.WithInsecure(),
//		grpc.WithTimeout(time.Second*3),
//		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
//	if err != nil {
//		panic(fmt.Errorf("grpc: failed to connect %s : %w", addr, err))
//	}
//	c.SvcConn[serviceName] = *conn
//}

func (c *AppConfig) ConnGRPC(serviceName string) {
	var opts []grpc.DialOption
	c.KeyPrefix(serviceName)
	if c.TLS() {
		cred, err := credentials.NewClientTLSFromFile(c.CertFile(), c.HostOverride())
		if err != nil {
			c.Log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())
	if c.ServiceAddress() == "" {
		c.Log.Fatalf("No service address %v", ErrNoConfigSettings)
	}
	conn, err := grpc.Dial(c.ServiceAddress(), opts...)
	if err != nil {
		c.Log.Fatalf("fail to dial %s: %v", c.ServiceAddress(), err)
	}
	c.SvcConn[serviceName] = conn
	c.Log.Infof("Established GRPC onnection to %s", serviceName)

	//defer c.SvcConn[serviceName].Close()
}
