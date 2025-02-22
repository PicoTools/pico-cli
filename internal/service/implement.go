package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/PicoTools/pico-cli/internal/version"
	operatorv1 "github.com/PicoTools/pico/pkg/proto/operator/v1"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/fatih/color"
	"github.com/go-faster/errors"
	"google.golang.org/grpc"
)

// HelloInit connects to hello topic
func HelloInit(ctx context.Context) (grpc.ServerStreamingClient[operatorv1.HelloResponse], error) {
	return getSvc().Hello(ctx, &operatorv1.HelloRequest{
		Version: version.Version(),
	})
}

// HelloHandshake processes hadnshake from hello topic
func HelloHandshake(ctx context.Context) error {
	msg, err := conn.ss.controlStream.Recv()
	if err != nil {
		return err
	}
	if msg.GetHandshake() == nil {
		return fmt.Errorf("unexpected hello response (no handshake data)")
	}
	conn.metadata.username = msg.GetHandshake().GetUsername()
	conn.metadata.cookie = msg.GetHandshake().GetCookie().GetValue()
	conn.metadata.delta = time.Since(msg.GetHandshake().GetTime().AsTime())
	return nil
}

// HelloMonitor maintained control session
func HelloMonitor(ctx context.Context) error {
	for {
		if _, err := conn.ss.controlStream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}
	return nil
}

// SubscribeChat subscribes operator on gathering chat events
func SubscribeChat(ctx context.Context) error {
	stream, err := getSvc().SubscribeChat(ctx, &operatorv1.SubscribeChatRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "open chat subscription stream")
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		// get message
		if msg.GetMessage() != nil {
			v := msg.GetMessage()
			if v.GetIsServer() {
				// server message in chat
				notificator.Printf("[%s] %s", color.GreenString("chat"), v.GetMessage())
				continue
			}
			if strings.Compare(v.GetFrom().GetValue(), GetUsername()) == 0 {
				// do not print message from operator itself
				continue
			}
			notificator.Printf("[%s] %s: %s", color.GreenString("chat"), color.RedString(v.GetFrom().GetValue()), v.GetMessage())
		}
	}
	return nil
}

// SubscribeAgents subscribes operator on gathering agents events
func SubscribeAgents(ctx context.Context) error {
	stream, err := getSvc().SubscribeAgents(ctx, &operatorv1.SubscribeAgentsRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "open agent subscription stream")
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		// get list of agents
		if msg.GetAgents() != nil {
			for _, v := range msg.GetAgents().GetAgents() {
				b := &agent.Agent{}
				b.SetId(v.GetId())
				b.SetListenerId(v.GetLid())
				b.SetExtIp(v.GetExtIp().GetValue())
				b.SetIntIp(v.GetIntIp().GetValue())
				b.SetOs(shared.AgentOs(v.GetOs()))
				b.SetOsMeta(v.GetOsMeta().GetValue())
				b.SetHostname(v.GetHostname().GetValue())
				b.SetUsername(v.GetUsername().GetValue())
				b.SetDomain(v.GetDomain().GetValue())
				b.SetIsPrivileged(v.GetPrivileged().GetValue())
				b.SetProcessName(v.GetProcName().GetValue())
				b.SetPid(v.GetPid().GetValue())
				b.SetArch(shared.AgentArch(v.GetArch()))
				b.SetSleep(v.GetSleep())
				b.SetJitter(uint8(v.GetJitter()))
				b.SetCaps(v.GetCaps())
				b.SetColor(v.GetColor().GetValue())
				b.SetNote(v.GetNote().GetValue())
				b.SetFirst(v.GetFirst().AsTime().Add(conn.metadata.delta))
				b.SetLast(v.GetLast().AsTime().Add(conn.metadata.delta))
				// add agent to storage
				agent.Agents.Add(b)
			}
			continue
		}
		// get agent
		if msg.GetAgent() != nil {
			b := &agent.Agent{}
			v := msg.GetAgent()
			b.SetId(v.GetId())
			b.SetListenerId(v.GetLid())
			b.SetExtIp(v.GetExtIp().GetValue())
			b.SetIntIp(v.GetIntIp().GetValue())
			b.SetOs(shared.AgentOs(v.GetOs()))
			b.SetOsMeta(v.GetOsMeta().GetValue())
			b.SetHostname(v.GetHostname().GetValue())
			b.SetUsername(v.GetUsername().GetValue())
			b.SetDomain(v.GetDomain().GetValue())
			b.SetIsPrivileged(v.GetPrivileged().GetValue())
			b.SetProcessName(v.GetProcName().GetValue())
			b.SetPid(v.GetPid().GetValue())
			b.SetArch(shared.AgentArch(v.GetArch()))
			b.SetSleep(v.GetSleep())
			b.SetJitter(uint8(v.GetJitter()))
			b.SetCaps(v.GetCaps())
			b.SetColor(v.GetColor().GetValue())
			b.SetNote(v.GetNote().GetValue())
			b.SetFirst(v.GetFirst().AsTime().Add(conn.metadata.delta))
			b.SetLast(v.GetLast().AsTime().Add(conn.metadata.delta))
			// add agent to storage
			agent.Agents.Add(b)
			continue
		}
		// get note
		if msg.GetNote() != nil {
			v := msg.GetNote()
			if b := agent.Agents.GetById(v.GetId()); b != nil {
				b.SetNote(v.GetNote().GetValue())
			}
			continue
		}
		// get color
		if msg.GetColor() != nil {
			v := msg.GetColor()
			if b := agent.Agents.GetById(v.GetId()); b != nil {
				b.SetColor(v.GetColor().GetValue())
			}
			continue
		}
		// get last checkout timestamp
		if msg.GetLast() != nil {
			v := msg.GetLast()
			if b := agent.Agents.GetById(v.GetId()); b != nil {
				b.SetLast(v.GetLast().AsTime().Add(conn.metadata.delta))
			}
			continue
		}
		// get sleep value
		if msg.GetSleep() != nil {
			v := msg.GetSleep()
			if b := agent.Agents.GetById(v.GetId()); b != nil {
				b.SetSleep(v.GetSleep())
				b.SetJitter(uint8(v.GetJitter()))
			}
			continue
		}
	}
	return nil
}

// SubscribeTasks subscribes operator on gathering commands events
func SubscribeTasks(ctx context.Context) error {
	stream, err := getSvc().SubscribeTasks(ctx)
	if err != nil {
		return errors.Wrap(err, "open tasks subscription stream")
	}
	// operator's authorization message
	if err = stream.Send(&operatorv1.SubscribeTasksRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Type: &operatorv1.SubscribeTasksRequest_Hello{
			Hello: &operatorv1.SubscribeTasksHelloRequest{},
		},
	}); err != nil {
		return errors.Wrap(err, "send hello message to tasks topic")
	}
	// save stream
	conn.ss.tasksStream = stream
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		// get command
		if msg.GetCommand() != nil {
			command := &task.Command{}
			v := msg.GetCommand()
			command.SetId(v.GetId())
			command.SetCmd(v.GetCmd())
			command.SetCreatedAt(v.GetCreated().AsTime().Add(conn.metadata.delta))
			command.SetAuthor(v.GetAuthor())
			// add command to storage
			task.Commands.Add(command)
			continue
		}
		// get command's message
		if msg.GetMessage() != nil {
			m := &task.Message{}
			v := msg.GetMessage()
			command := task.Commands.GetById(v.GetId())
			if command == nil {
				continue
			}
			m.SetId(v.GetMid())
			m.SetKind(shared.TaskMessage(v.Type))
			m.SetMessage(v.GetMessage())
			m.SetCreatedAt(v.Created.AsTime().Add(conn.metadata.delta))
			command.AddMessage(m)
			continue
		}
		if msg.GetTask() != nil {
			t := &task.Task{}
			v := msg.GetTask()
			command := task.Commands.GetById(v.GetId())
			if command == nil {
				continue
			}
			t.SetId(v.GetTid())
			t.SetIsOutputBig(v.GetOutputBig())
			t.SetCreatedAt(v.GetCreated().AsTime().Add(conn.metadata.delta))
			t.SetOutput(v.GetOutput().GetValue())
			t.SetOutputLen(v.GetOutputLen())
			t.SetStatus(shared.TaskStatus(v.GetStatus()))
			t.SetCapability(shared.Capability(v.GetCap()))
			command.AddTask(t)
			continue
		}
		// get command's task status
		if msg.GetTaskStatus() != nil {
			v := msg.GetTaskStatus()
			command := task.Commands.GetById(v.GetId())
			if command == nil {
				continue
			}
			t := command.GetTaskById(v.GetTid())
			if t == nil {
				continue
			}
			t.SetStatus(shared.TaskStatus(v.GetStatus()))
			command.UpdateTask(t)
			continue
		}
		// get command's task results
		if msg.GetTaskDone() != nil {
			v := msg.GetTaskDone()
			command := task.Commands.GetById(v.GetId())
			if command == nil {
				continue
			}
			t := command.GetTaskById(v.GetTid())
			if t == nil {
				continue
			}
			t.SetStatus(shared.TaskStatus(v.GetStatus()))
			t.SetIsOutputBig(v.GetOutputBig())
			t.SetOutput(v.GetOutput().GetValue())
			t.SetOutputLen(v.GetOutputLen())
			command.UpdateTask(t)
			continue
		}
	}
	return nil
}

// PollAgentTasks starts polling of tasks for agent
func PollAgentTasks(agent *agent.Agent) error {
	if err := conn.ss.tasksStream.Send(&operatorv1.SubscribeTasksRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Type: &operatorv1.SubscribeTasksRequest_Start{
			Start: &operatorv1.StartPollAgentRequest{
				Id: agent.GetId(),
			},
		},
	}); err != nil {
		return errors.Wrapf(err, "poll tasks for agent %s", agent.GetIdHex())
	}
	return nil
}

// UnpollAgentTasks stop polling of tasks for agent
func UnpollAgentTasks(agent *agent.Agent) error {
	if err := conn.ss.tasksStream.Send(&operatorv1.SubscribeTasksRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Type: &operatorv1.SubscribeTasksRequest_Stop{
			Stop: &operatorv1.StopPollAgentRequest{
				Id: agent.GetId(),
			},
		},
	}); err != nil {
		return errors.Wrapf(err, "unpoll tasks for agent %s", agent.GetIdHex())
	}
	return nil
}

// NewCommand creates new command
func NewCommand(id uint32, cmd string, visible bool) error {
	stream, err := getSvc().NewCommand(context.Background())
	if err != nil {
		return errors.Wrap(err, "open command submition stream")
	}
	if err = stream.Send(&operatorv1.NewCommandRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Type: &operatorv1.NewCommandRequest_Command{
			Command: &operatorv1.CreateCommandRequest{
				Id:      id,
				Cmd:     cmd,
				Visible: visible,
			},
		},
	}); err != nil {
		return errors.Wrap(err, "open command")
	}
	// save stream
	conn.ss.commandStreams.Store(id, stream)
	return nil
}

// CloseCommand closes opened command
func CloseCommand(id uint32) error {
	stream, ok := conn.ss.commandStreams.Load(id)
	if !ok {
		return fmt.Errorf("unable load stream for agent %d", id)
	}
	defer func() {
		// remove command from storage
		conn.ss.commandStreams.Delete(id)
	}()
	if _, err := stream.CloseAndRecv(); err != nil {
		if !errors.Is(err, io.EOF) {
			return errors.Wrap(err, "close stream")
		}
	}
	return nil
}

// NewCommandMessage saves message for command
func NewCommandMessage(id uint32, tm shared.TaskMessage, message string) error {
	stream, ok := conn.ss.commandStreams.Load(id)
	if !ok {
		return fmt.Errorf("unable load stream for agent %d", id)
	}
	return stream.Send(&operatorv1.NewCommandRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Type: &operatorv1.NewCommandRequest_Message{
			Message: &operatorv1.CreateMessageRequest{
				Type: uint32(tm),
				Msg:  message,
			},
		},
	})
}

// NewTask creates new task in command
func NewTask(id uint32, v *operatorv1.CreateTaskRequest) error {
	stream, ok := conn.ss.commandStreams.Load(id)
	if !ok {
		return fmt.Errorf("unable load stream for agent %d", id)
	}
	return stream.Send(&operatorv1.NewCommandRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Type: &operatorv1.NewCommandRequest_Task{
			Task: v,
		},
	})
}

// CancelTasks cancels all tasks from operator which in status "NEW"
func CancelTasks(id uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := getSvc().CancelTasks(ctx, &operatorv1.CancelTasksRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Id: id,
	})
	return err
}

// GetTaskOutput returns full task's output
func GetTaskOutput(id int64) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rep, err := getSvc().GetTaskOutput(ctx, &operatorv1.GetTaskOutputRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return rep.GetOutput().GetValue(), nil
}

// SendChatMessage saves message from operator in chat
func SendChatMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := getSvc().NewChatMessage(ctx, &operatorv1.NewChatMessageRequest{
		Cookie: &operatorv1.SessionCookie{
			Value: conn.getMetadata().GetCookie(),
		},
		Message: message,
	})
	return err
}
