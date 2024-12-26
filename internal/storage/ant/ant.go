package ant

import (
	"fmt"
	"sort"
	"time"

	"github.com/PicoTools/pico-shared/shared"
	"github.com/lrita/cmap"
)

// active ant for polling
var ActiveAnt *Ant

// ants storage in runtime
var Ants = &antsMapper{
	sorted: &ants{
		ants: make([]*Ant, 0),
	},
}

// ant information
type Ant struct {
	id           uint32
	listenerId   int64
	extIp        string
	intIp        string
	os           shared.AntOs
	osMeta       string
	hostname     string
	username     string
	domain       string
	isPrivileged bool
	processName  string
	pid          uint64
	arch         shared.AntArch
	sleep        uint32
	jitter       uint8
	caps         uint32
	color        uint32
	note         string
	first        time.Time
	last         time.Time
}

type ants struct {
	ants []*Ant
}

type antsMapper struct {
	ants   cmap.Map[uint32, *Ant]
	sorted *ants
}

func (b *Ant) GetId() uint32 {
	return b.id
}

func (b *Ant) GetIdHex() string {
	return fmt.Sprintf("%08x", b.id)
}

func (b *Ant) SetId(id uint32) {
	b.id = id
}

func (b *Ant) GetListenerId() int64 {
	return b.listenerId
}

func (b *Ant) SetListenerId(id int64) {
	b.listenerId = id
}

func (b *Ant) GetExtIp() string {
	return b.extIp
}

func (b *Ant) SetExtIp(data string) {
	b.extIp = data
}

func (b *Ant) GetIntIp() string {
	return b.intIp
}

func (b *Ant) SetIntIp(data string) {
	b.intIp = data
}

func (b *Ant) GetOs() shared.AntOs {
	return b.os
}

func (b *Ant) SetOs(os shared.AntOs) {
	b.os = os
}

func (b *Ant) GetOsMeta() string {
	return b.osMeta
}

func (b *Ant) SetOsMeta(data string) {
	b.osMeta = data
}

func (b *Ant) GetHostname() string {
	return b.hostname
}

func (b *Ant) SetHostname(data string) {
	b.hostname = data
}

func (b *Ant) GetUsername() string {
	return b.username
}

func (b *Ant) SetUsername(data string) {
	b.username = data
}

func (b *Ant) GetDomain() string {
	return b.domain
}

func (b *Ant) SetDomain(data string) {
	b.domain = data
}

func (b *Ant) GetIsPrivileged() bool {
	return b.isPrivileged
}

func (b *Ant) SetIsPrivileged(flag bool) {
	b.isPrivileged = flag
}

func (b *Ant) GetProcessName() string {
	return b.processName
}

func (b *Ant) SetProcessName(data string) {
	b.processName = data
}

func (b *Ant) GetPid() uint64 {
	return b.pid
}

func (b *Ant) SetPid(pid uint64) {
	b.pid = pid
}

func (b *Ant) GetArch() shared.AntArch {
	return b.arch
}

func (b *Ant) SetArch(arch shared.AntArch) {
	b.arch = arch
}

func (b *Ant) GetSleep() uint32 {
	return b.sleep
}

func (b *Ant) SetSleep(sleep uint32) {
	b.sleep = sleep
}

func (b *Ant) GetJitter() uint8 {
	return b.jitter
}

func (b *Ant) SetJitter(jitter uint8) {
	b.jitter = jitter
}

func (b *Ant) GetCaps() uint32 {
	return b.caps
}

func (b *Ant) SetCaps(caps uint32) {
	b.caps = caps
}

func (b *Ant) GetColor() uint32 {
	return b.color
}

func (b *Ant) SetColor(color uint32) {
	b.color = color
}

func (b *Ant) GetNote() string {
	return b.note
}

func (b *Ant) SetNote(data string) {
	b.note = data
}

func (b *Ant) GetFirst() time.Time {
	return b.first
}

func (b *Ant) SetFirst(t time.Time) {
	b.first = t
}

func (b *Ant) GetLast() time.Time {
	return b.last
}

func (b *Ant) SetLast(t time.Time) {
	b.last = t
}

// IsDelay returns true if ant delayed on sleep + sleep * jitter
func (b *Ant) IsDelay(delta time.Duration) bool {
	sleep := int(b.sleep * 1000)
	jitter := int(b.jitter / 100)
	return time.Now().After(b.GetLast().Add(time.Duration(sleep+sleep*jitter) * time.Millisecond))
}

// IsDead returns true if ant delayed on 3 * (sleep + sleep * jitter)
func (b *Ant) IsDead(delta time.Duration) bool {
	sleep := int(b.sleep * 1000)
	jitter := int(b.jitter / 100)
	return time.Now().After(b.GetLast().Add(time.Duration(3*(sleep+sleep*jitter)) * time.Millisecond))
}

// Sort sorts ant list by last checkout timestamp
func (b *ants) Sort() {
	sort.SliceStable(b.ants, func(i, j int) bool {
		return b.ants[i].GetLast().Before(b.ants[j].GetLast())
	})
}

// Add adds ant to storage
func (b *antsMapper) Add(v *Ant) {
	b.ants.Store(v.GetId(), v)
	b.Fill()
}

// Get returns list of ants
func (b *antsMapper) Get() []*Ant {
	return b.sorted.ants
}

// GetById returns ant specified by ID
func (b *antsMapper) GetById(id uint32) *Ant {
	if v, ok := b.ants.Load(id); ok {
		return v
	}
	return nil
}

// Count returns count of ants in storage
func (b *antsMapper) Count() int {
	return b.ants.Count()
}

// Fill fills sorted array of ants
func (b *antsMapper) Fill() {
	temp := &ants{
		ants: make([]*Ant, 0),
	}

	b.ants.Range(func(k uint32, v *Ant) bool {
		temp.ants = append(temp.ants, v)
		return true
	})

	temp.Sort()

	b.sorted = temp
}
