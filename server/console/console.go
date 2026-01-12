package console

import (
	"context"
	"github.com/longzekun/specter/client/command"
	"github.com/longzekun/specter/client/console"
	"github.com/longzekun/specter/client/constants"
	client_transport "github.com/longzekun/specter/client/transport"
	"github.com/longzekun/specter/protobuf/rpcpb"
	"github.com/longzekun/specter/server/transport"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"net"
)

func Start() {
	_, ln, _ := transport.LocalListener()
	ctxDialer := grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		return ln.Dial()
	})

	options := []grpc.DialOption{
		ctxDialer,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(client_transport.ClientMaxReceiveMessageSize)),
	}

	conn, err := grpc.Dial(ln.Addr().String(), options...)
	if err != nil {
		zap.S().Warnf("fail to dial: %v", err)
		return
	}
	defer conn.Close()

	localRPC := rpcpb.NewSpecterRPCClient(conn)
	con := console.NewConsole(false)

	console.StartClient(con, localRPC, command.ServerCommands(con, serverOnlyCommands), nil)
}

func serverOnlyCommands() (commands []*cobra.Command) {
	startMultiplayer := &cobra.Command{
		Use:     constants.MultiplayerModeStr,
		Short:   "Enable multiplayer mode",
		Long:    ``,
		Run:     startMultiplayerModeCmd,
		GroupID: constants.MultiplayerHelpGroup,
	}
	command.Bind("multiplayer", false, startMultiplayer, func(f *pflag.FlagSet) {
		f.StringP("lhost", "l", "", "interface to bind server to")
		f.Uint16P("lport", "p", 31337, "tcp listen port")
	})

	commands = append(commands, startMultiplayer)

	return commands
}
