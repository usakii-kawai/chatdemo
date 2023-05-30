package server

type ServerStartOptions struct {
	Id     string
	Listen string
}

func RunServerStart(opts *ServerStartOptions, version string) error {
	server := NewServer(opts.Id, opts.Listen)
	defer server.Shutdown()

	return server.Start()
}

// func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
// 	opts := &ServerStartOptions{}
// 	cmd := &cobra.Command{
// 		Use:   "chat",
// 		Short: "Starts a chat server",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return RunServerStart(ctx, opts, version)
// 		},
// 	}

// 	cmd.PersistentFlags().StringVarP(&opts.id, "serverid", "i", "demo", "server id")
// 	cmd.PersistentFlags().StringVarP(&opts.listen, "listen", "l", ":8080", "listen address")

// 	return cmd
// }
