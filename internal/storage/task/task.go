package task

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/fatih/color"
	"github.com/lrita/cmap"
)

type TaskData interface {
	GetCreatedAt() time.Time
}

// commands storage in runtime
var Commands = &commandsMapper{
	sorted: &commands{
		commands: make([]*Command, 0),
	},
}

// ResetStorage resets storage of commands for agent
func ResetStorage() {
	Commands = &commandsMapper{
		sorted: &commands{
			commands: make([]*Command, 0),
		},
	}
}

// command information
type Command struct {
	id        int64
	cmd       string
	createdAt time.Time
	closedAt  time.Time
	author    string
	data      *CommandResults
}

// command's results (tasks and messages)
type CommandResults struct {
	messages cmap.Map[int64, *Message]
	tasks    cmap.Map[int64, *Task]
	sorted   []TaskData
}

// AddMessage saves command's message to storage
func (t *CommandResults) AddMessage(m *Message) {
	t.messages.Store(m.GetId(), m)
	t.Fill()
}

// Get returns task data
func (t *CommandResults) Get() []TaskData {
	return t.sorted
}

// GetTaskById returns task by ID
func (t *CommandResults) GetTaskById(id int64) *Task {
	if v, ok := t.tasks.Load(id); ok {
		return v
	}
	return nil
}

// AddTask adds command's task to storage
func (t *CommandResults) AddTask(task *Task) {
	t.tasks.Store(task.GetId(), task)
	t.Fill()
}

// UpdateTask updates task in storage
func (t *CommandResults) UpdateTask(v *Task) {
	t.tasks.Store(v.GetId(), v)
	t.Fill()
}

// Fill fills sorted task data results
func (t *CommandResults) Fill() {
	temp := make([]TaskData, 0)

	t.messages.Range(func(k int64, v *Message) bool {
		temp = append(temp, v)
		return true
	})

	t.tasks.Range(func(k int64, v *Task) bool {
		temp = append(temp, v)
		return true
	})

	sort.SliceStable(temp, func(i, j int) bool {
		return temp[i].GetCreatedAt().Before(temp[j].GetCreatedAt())
	})

	t.sorted = temp
}

// Message implements TaskData interface and stores infromation about command's message
type Message struct {
	TaskData
	id        int64
	kind      shared.TaskMessage
	message   string
	createdAt time.Time
}

func (m *Message) String() string {
	var s strings.Builder
	s.WriteRune('[')
	switch m.kind {
	case shared.NotifyMessage:
		s.WriteString(color.CyanString("*"))
	case shared.InfoMessage:
		s.WriteString(color.GreenString("+"))
	case shared.WarningMessage:
		s.WriteString(color.YellowString("!"))
	case shared.ErrorMessage:
		s.WriteString(color.RedString("-"))
	}
	s.WriteRune(']')
	s.WriteRune(' ')
	s.WriteString(m.message)
	return s.String()
}

// Task implements TaskData interface and stores infromation about command's task
type Task struct {
	TaskData
	id          int64
	isOutputBig bool
	isBinary    bool
	output      []byte
	outputLen   uint64
	status      shared.TaskStatus
	createdAt   time.Time
	capability  shared.Capability
}

func (t *Task) StringStatus() string {
	var s strings.Builder
	s.WriteRune('[')
	switch t.status {
	case shared.StatusNew:
		s.WriteString(color.HiWhiteString("NEW"))
	case shared.StatusInProgress:
		s.WriteString(color.HiCyanString("IN PROGRESS"))
	case shared.StatusCancelled:
		s.WriteString(color.HiYellowString("CANCELLED"))
	case shared.StatusError:
		s.WriteString(color.HiRedString("ERROR"))
	case shared.StatusSuccess:
		s.WriteString(color.HiGreenString("DONE"))
	}
	s.WriteRune(']')
	s.WriteString(fmt.Sprintf(" (%d) received output with length %d bytes", t.id, t.outputLen))
	return s.String()
}

func (t *Task) GetId() int64 {
	return t.id
}

func (t *Task) GetIdHex() string {
	return fmt.Sprintf("%06x", t.id)[:6]
}

func (t *Task) SetId(id int64) {
	t.id = id
}

func (t *Task) GetIsOutputBig() bool {
	return t.isOutputBig
}

func (t *Task) SetIsOutputBig(flag bool) {
	t.isOutputBig = flag
}

func (t *Task) GetOutput() []byte {
	return t.output
}

func (t *Task) GetOutputString() string {
	return string(t.output)
}

func (t *Task) SetOutput(data []byte) {
	if !utils.IsStrPrintable(string(data)) {
		t.SetIsBinary(true)
	} else {
		t.SetIsBinary(false)
	}
	t.output = data
}

func (t *Task) GetOutputLen() uint64 {
	return t.outputLen
}

func (t *Task) SetOutputLen(length uint64) {
	t.outputLen = length
}

func (t *Task) GetIsBinary() bool {
	return t.isBinary
}

func (t *Task) SetIsBinary(flag bool) {
	t.isBinary = flag
}

func (t *Task) GetStatus() shared.TaskStatus {
	return t.status
}

func (t *Task) SetStatus(status shared.TaskStatus) {
	t.status = status
}

func (t *Task) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) SetCreatedAt(createdAt time.Time) {
	t.createdAt = createdAt
}

func (t *Task) GetCapability() shared.Capability {
	return t.capability
}

func (t *Task) SetCapability(cap shared.Capability) {
	t.capability = cap
}

func (m *Message) GetId() int64 {
	return m.id
}

func (m *Message) SetId(id int64) {
	m.id = id
}

func (m *Message) GetKind() shared.TaskMessage {
	return m.kind
}

func (m *Message) SetKind(kind shared.TaskMessage) {
	m.kind = kind
}

func (m *Message) GetMessage() string {
	return m.message
}

func (m *Message) SetMessage(message string) {
	m.message = message
}

func (m *Message) GetCreatedAt() time.Time {
	return m.createdAt
}

func (m *Message) SetCreatedAt(t time.Time) {
	m.createdAt = t
}

func (t *Command) GetTaskById(id int64) *Task {
	return t.data.GetTaskById(id)
}

func (t *Command) UpdateTask(task *Task) {
	t.data.UpdateTask(task)
}

func (t *Command) AddMessage(m *Message) {
	t.data.AddMessage(m)
}

func (t *Command) AddTask(task *Task) {
	t.data.AddTask(task)
}

func (t *Command) GetId() int64 {
	return t.id
}

func (t *Command) GetIdHex() string {
	return fmt.Sprintf("%06x", t.id)[:6]
}

func (t *Command) SetId(id int64) {
	t.id = id
}

func (t *Command) GetCmd() string {
	return t.cmd
}

func (t *Command) SetCmd(cmd string) {
	t.cmd = cmd
}

func (t *Command) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Command) SetCreatedAt(createdAt time.Time) {
	t.createdAt = createdAt
}

func (t *Command) GetClosedAt() time.Time {
	return t.closedAt
}

func (t *Command) SetClosedAt(closedAt time.Time) {
	t.closedAt = closedAt
}

func (t *Command) GetAuthor() string {
	return t.author
}

func (t *Command) SetAuthor(author string) {
	t.author = author
}

func (t *Command) GetData() *CommandResults {
	return t.data
}

type commandsMapper struct {
	commands cmap.Map[int64, *Command]
	sorted   *commands
}

type commands struct {
	commands []*Command
}

// GetLast returns last command for agent
func (t *commandsMapper) GetLast() *Command {
	data := t.Get()
	if len(data) == 0 {
		return nil
	}
	return data[len(data)-1]
}

// Add adds command to storage
func (t *commandsMapper) Add(v *Command) {
	v.data = &CommandResults{
		sorted: make([]TaskData, 0),
	}
	t.commands.Store(v.GetId(), v)
	t.Fill()
}

// Get returns sorted list of commands
func (t *commandsMapper) Get() []*Command {
	return t.sorted.commands
}

// GetTasks returns all tasks in all commands
func (t *commandsMapper) GetTasks() []*Task {
	temp := make([]*Task, 0)
	t.commands.Range(func(k int64, v *Command) bool {
		v.data.tasks.Range(func(key int64, value *Task) bool {
			temp = append(temp, value)
			return true
		})
		return true
	})
	sort.Slice(temp, func(i, j int) bool {
		return temp[i].GetCreatedAt().Before(temp[j].GetCreatedAt())
	})
	return temp
}

// GetById returns command by ID
func (t *commandsMapper) GetById(id int64) *Command {
	if v, ok := t.commands.Load(id); ok {
		return v
	}
	return nil
}

func (t *commandsMapper) GetTaskById(id int64) *Task {
	var task *Task
	t.commands.Range(func(k int64, v *Command) bool {
		v.data.tasks.Range(func(key int64, value *Task) bool {
			if key == id {
				task = value
				return false
			}
			return true
		})
		return task == nil
	})
	return task
}

// Count returns number of commands in storage
func (t *commandsMapper) Count() int {
	return t.commands.Count()
}

// Sort sorts commands by create timestamp
func (t *commands) Sort() {
	sort.SliceStable(t.commands, func(i, j int) bool {
		return t.commands[i].GetCreatedAt().Before(t.commands[j].GetCreatedAt())
	})
}

// Fill fills sorted list with commands
func (t *commandsMapper) Fill() {
	temp := &commands{
		commands: make([]*Command, 0),
	}

	t.commands.Range(func(k int64, v *Command) bool {
		temp.commands = append(temp.commands, v)
		return true
	})

	temp.Sort()

	t.sorted = temp
}
