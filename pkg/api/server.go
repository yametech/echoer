package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"unsafe"

	"github.com/gin-gonic/gin"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/yametech/echoer/api"
	"github.com/yametech/echoer/pkg/command"
	"github.com/yametech/echoer/pkg/storage"
	"google.golang.org/grpc"
)

var _ api.EchoServer = &Server{}

type Server struct {
	*Handle
	storage.IStorage
	middlewares []gin.HandlerFunc
	parser      *command.Parser
}

func NewServer(storage storage.IStorage) *Server {
	server := &Server{
		IStorage:    storage,
		middlewares: make([]gin.HandlerFunc, 0),
		parser:      command.NewParser(storage),
	}
	server.Handle = &Handle{Server: server}
	return server
}

func (s *Server) RegistryMiddlewares(h gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, h)
}

func (s *Server) Run(addr string) error {
	router := gin.New()
	router.Use(s.middlewares...)
	router.GET("/", func(g *gin.Context) { g.JSON(http.StatusOK, "echoer") })

	// watch
	router.GET("/watch", s.watch)

	// action
	router.POST("/action", s.actionCreate)
	router.GET("/action", s.actionList)
	router.GET("/action/:name", s.actionGet)
	router.DELETE("/action/:name", s.actionDelete)

	// event
	router.POST("/event", s.eventCreate)
	router.GET("/event", s.eventList)

	//flow
	router.POST("/flow", s.flowCreate)
	router.GET("/flow", s.flowList)
	router.GET("/flow/:name", s.flowGet)
	router.DELETE("/flow/:name", s.flowDelete)

	//flowrun
	router.POST("/flowrun", s.flowRunCreate)
	router.GET("/flowrun", s.flowRunList)
	router.GET("/flowrun/:name", s.flowRunGet)
	router.DELETE("/flowrun/:name", s.flowRunDelete)

	//step
	// POST recv action response state
	router.POST("/step", s.ackStep)

	if err := router.Run(addr); err != nil {
		return err
	}
	return nil
}

func (s *Server) Execute(ctx context.Context, request *api.ExecuteRequest) (*api.ExecuteCommandResponse, error) {
	cmdStr := *(*string)(unsafe.Pointer(&request.Command))
	cmd, args, err := s.parser.Parse(cmdStr)
	if err != nil {
		if err != command.ErrCommandNotFound {
			return nil, fmt.Errorf(`can't not parse (%s)`, cmdStr)
		}
		return &api.ExecuteCommandResponse{
			Reply: api.CommandExecutionReply_ERR,
			Raw:   []byte(err.Error()),
		}, nil
	}
	return s.createResponse(cmd.Execute(args...))
}

func (s *Server) createResponse(reply command.Reply) (resp *api.ExecuteCommandResponse, err error) {
	resp = new(api.ExecuteCommandResponse)
	switch reply.(type) {
	case *command.OkReply:
		resp.Reply = api.CommandExecutionReply_OK
		return
	case *command.RawReply:
		resp.Reply = api.CommandExecutionReply_Raw
		resp.Raw = reply.Value().([]byte)
		return
	case *command.ErrorReply:
		resp.Reply = api.CommandExecutionReply_ERR
		resp.Raw = []byte(fmt.Sprintf("%s", reply.Value()))
		return
	}
	err = fmt.Errorf("unknow reply (%v)", reply)
	return
}

func (s *Server) RpcServer(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("[INFO] listen rpc (%s)\n", addr)
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_validator.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpc_validator.StreamServerInterceptor()),
	)
	api.RegisterEchoServer(srv, s)
	return srv.Serve(listen)
}

type createRawData struct {
	Data string `json:"data"`
}
