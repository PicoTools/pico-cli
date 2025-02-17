package agent

import (
	"fmt"
	"sort"
	"time"

	"github.com/PicoTools/pico/pkg/shared"
	"github.com/lrita/cmap"
)

// active agent for polling
var ActiveAgent *Agent

// agents storage in runtime
var Agents = &agentsMapper{
	sorted: &agents{
		agents: make([]*Agent, 0),
	},
}

// agent information
type Agent struct {
	id           uint32
	listenerId   int64
	extIp        string
	intIp        string
	os           shared.AgentOs
	osMeta       string
	hostname     string
	username     string
	domain       string
	isPrivileged bool
	processName  string
	pid          uint64
	arch         shared.AgentArch
	sleep        uint32
	jitter       uint8
	caps         uint32
	color        uint32
	note         string
	first        time.Time
	last         time.Time
}

type agents struct {
	agents []*Agent
}

type agentsMapper struct {
	agents cmap.Map[uint32, *Agent]
	sorted *agents
}

func (b *Agent) GetId() uint32 {
	return b.id
}

func (b *Agent) GetIdHex() string {
	return fmt.Sprintf("%08x", b.id)
}

func (b *Agent) SetId(id uint32) {
	b.id = id
}

func (b *Agent) GetListenerId() int64 {
	return b.listenerId
}

func (b *Agent) SetListenerId(id int64) {
	b.listenerId = id
}

func (b *Agent) GetExtIp() string {
	return b.extIp
}

func (b *Agent) SetExtIp(data string) {
	b.extIp = data
}

func (b *Agent) GetIntIp() string {
	return b.intIp
}

func (b *Agent) SetIntIp(data string) {
	b.intIp = data
}

func (b *Agent) GetOs() shared.AgentOs {
	return b.os
}

func (b *Agent) SetOs(os shared.AgentOs) {
	b.os = os
}

func (b *Agent) GetOsMeta() string {
	return b.osMeta
}

func (b *Agent) SetOsMeta(data string) {
	b.osMeta = data
}

func (b *Agent) GetHostname() string {
	return b.hostname
}

func (b *Agent) SetHostname(data string) {
	b.hostname = data
}

func (b *Agent) GetUsername() string {
	return b.username
}

func (b *Agent) SetUsername(data string) {
	b.username = data
}

func (b *Agent) GetDomain() string {
	return b.domain
}

func (b *Agent) SetDomain(data string) {
	b.domain = data
}

func (b *Agent) GetIsPrivileged() bool {
	return b.isPrivileged
}

func (b *Agent) SetIsPrivileged(flag bool) {
	b.isPrivileged = flag
}

func (b *Agent) GetProcessName() string {
	return b.processName
}

func (b *Agent) SetProcessName(data string) {
	b.processName = data
}

func (b *Agent) GetPid() uint64 {
	return b.pid
}

func (b *Agent) SetPid(pid uint64) {
	b.pid = pid
}

func (b *Agent) GetArch() shared.AgentArch {
	return b.arch
}

func (b *Agent) SetArch(arch shared.AgentArch) {
	b.arch = arch
}

func (b *Agent) GetSleep() uint32 {
	return b.sleep
}

func (b *Agent) SetSleep(sleep uint32) {
	b.sleep = sleep
}

func (b *Agent) GetJitter() uint8 {
	return b.jitter
}

func (b *Agent) SetJitter(jitter uint8) {
	b.jitter = jitter
}

func (b *Agent) GetCaps() uint32 {
	return b.caps
}

func (b *Agent) SetCaps(caps uint32) {
	b.caps = caps
}

func (b *Agent) GetColor() uint32 {
	return b.color
}

func (b *Agent) SetColor(color uint32) {
	b.color = color
}

func (b *Agent) GetNote() string {
	return b.note
}

func (b *Agent) SetNote(data string) {
	b.note = data
}

func (b *Agent) GetFirst() time.Time {
	return b.first
}

func (b *Agent) SetFirst(t time.Time) {
	b.first = t
}

func (b *Agent) GetLast() time.Time {
	return b.last
}

func (b *Agent) SetLast(t time.Time) {
	b.last = t
}

// IsDelay returns true if agent delayed on sleep + sleep * jitter
func (b *Agent) IsDelay(delta time.Duration) bool {
	sleep := int(b.sleep * 1000)
	jitter := int(b.jitter / 100)
	return time.Now().After(b.GetLast().Add(time.Duration(sleep+sleep*jitter) * time.Millisecond))
}

// IsDead returns true if agent delayed on 3 * (sleep + sleep * jitter)
func (b *Agent) IsDead(delta time.Duration) bool {
	sleep := int(b.sleep * 1000)
	jitter := int(b.jitter / 100)
	return time.Now().After(b.GetLast().Add(time.Duration(3*(sleep+sleep*jitter)) * time.Millisecond))
}

// Sort sorts agent list by last checkout timestamp
func (b *agents) Sort() {
	sort.SliceStable(b.agents, func(i, j int) bool {
		return b.agents[i].GetLast().Before(b.agents[j].GetLast())
	})
}

// Add adds agent to storage
func (b *agentsMapper) Add(v *Agent) {
	b.agents.Store(v.GetId(), v)
	b.Fill()
}

// Get returns list of agents
func (b *agentsMapper) Get() []*Agent {
	return b.sorted.agents
}

// GetById returns agent specified by ID
func (b *agentsMapper) GetById(id uint32) *Agent {
	if v, ok := b.agents.Load(id); ok {
		return v
	}
	return nil
}

// Count returns count of agents in storage
func (b *agentsMapper) Count() int {
	return b.agents.Count()
}

// Fill fills sorted array of agents
func (b *agentsMapper) Fill() {
	temp := &agents{
		agents: make([]*Agent, 0),
	}

	b.agents.Range(func(k uint32, v *Agent) bool {
		temp.agents = append(temp.agents, v)
		return true
	})

	temp.Sort()

	b.sorted = temp
}
