package grpc

import (
	"fmt"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"time"
)

func GetNewPushC2ClientConnection() *grpc.ClientConn {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt)))
	connectionString := fmt.Sprintf("%s:%d", utils.MythicConfig.MythicServerHost, utils.MythicConfig.MythicServerGRPCPort)
	for {
		logging.LogDebug("Attempting to connect to grpc...")
		if conn, err := grpc.Dial(connectionString, opts...); err != nil {
			logging.LogError(err, "Failed to connect to GRPC port for Mythic, trying again...", "connection", connectionString)
		} else {
			return conn
		}
		time.Sleep(grpcReconnectDelay * time.Second)
	}

}
