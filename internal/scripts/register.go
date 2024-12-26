package scripts

import (
	"embed"
	"fmt"

	acat "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_cat"
	acd "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_cd"
	acp "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_cp"
	adestruct "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_destruct"
	adownload "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_download"
	aexec "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_exec"
	aexecassembly "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_exec_assembly"
	aexecdetach "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_exec_detach"
	aexit "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_exit"
	ajobkill "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_jobkill"
	ajobs "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_jobs"
	akill "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_kill"
	als "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_ls"
	amkdir "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_mkdir"
	amv "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_mv"
	apause "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_pause"
	appid "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_ppid"
	aps "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_ps"
	apwd "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_pwd"
	ashell "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_shell"
	asleep "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_sleep"
	aupload "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_upload"
	awhoami "github.com/PicoTools/pico-cli/internal/scripts/aliases/a_whoami"
	"github.com/PicoTools/pico-cli/internal/scripts/aliases/alias"
	isarm32 "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_arm32_arch"
	isarm64 "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_arm64_arch"
	islinux "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_linux_os"
	ismacos "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_macos_os"
	iswindows "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_windows_os"
	isx64 "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_x64_arch"
	isx86 "github.com/PicoTools/pico-cli/internal/scripts/aliases/is_x86_arch"
	merror "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_error"
	minfo "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_info"
	mnotify "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_notify"
	mwarning "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_warning"
	tcancel "github.com/PicoTools/pico-cli/internal/scripts/aliases/t_cancel"
	"github.com/PicoTools/plan/pkg/engine/object"
	"github.com/PicoTools/plan/pkg/engine/storage"
	"github.com/PicoTools/plan/pkg/engine/types"
	mlanUtils "github.com/PicoTools/plan/pkg/engine/utils"
	"github.com/PicoTools/plan/pkg/engine/visitor"
	"github.com/go-faster/errors"
)

// registerApi registers API for PLAN integration
func registerApi() {
	// alias: register new alias
	storage.UserFunctions[alias.GetApiName()] = object.NewNativeFunc(alias.GetApiName(), alias.FrontendAlias)
	// m_notify: command's message with NOTIFY type
	storage.UserFunctions[mnotify.GetApiName()] = object.NewNativeFunc(mnotify.GetApiName(), mnotify.FrontendMessageNotify)
	// m_info: command's message with INFO type
	storage.UserFunctions[minfo.GetApiName()] = object.NewNativeFunc(minfo.GetApiName(), minfo.FrontendMessageInfo)
	// m_warning: command's message with WARNING type
	storage.UserFunctions[mwarning.GetApiName()] = object.NewNativeFunc(mwarning.GetApiName(), mwarning.FrontendMessageWarning)
	// m_error: command's message with ERROR type
	storage.UserFunctions[merror.GetApiName()] = object.NewNativeFunc(merror.GetApiName(), merror.FrontendMessageError)
	// a_sleep: change sleep/jitter ant's parameters
	storage.UserFunctions[asleep.GetApiName()] = object.NewNativeFunc(asleep.GetApiName(), asleep.FrontendAntSleep)
	// a_ls: get directory listing
	storage.UserFunctions[als.GetApiName()] = object.NewNativeFunc(als.GetApiName(), als.FrontendAntLs)
	// a_pwd: get process working directory
	storage.UserFunctions[apwd.GetApiName()] = object.NewNativeFunc(apwd.GetApiName(), apwd.FrontendAntPwd)
	// a_cd: change process working directory
	storage.UserFunctions[acd.GetApiName()] = object.NewNativeFunc(acd.GetApiName(), acd.FrontendAntCd)
	// a_whoami: get current user and its grants
	storage.UserFunctions[awhoami.GetApiName()] = object.NewNativeFunc(awhoami.GetApiName(), awhoami.FrontendAntWhoami)
	// a_ps: get processes listing
	storage.UserFunctions[aps.GetApiName()] = object.NewNativeFunc(aps.GetApiName(), aps.FrontendAntPs)
	// a_cat: print content of file
	storage.UserFunctions[acat.GetApiName()] = object.NewNativeFunc(acat.GetApiName(), acat.FrontendAntCat)
	// a_exec: execute binary with arguments
	storage.UserFunctions[aexec.GetApiName()] = object.NewNativeFunc(aexec.GetApiName(), aexec.FrontendAntExec)
	// a_cp: copy files/directories
	storage.UserFunctions[acp.GetApiName()] = object.NewNativeFunc(acp.GetApiName(), acp.FrontendAntCp)
	// a_jobs: get active jobs on ant
	storage.UserFunctions[ajobs.GetApiName()] = object.NewNativeFunc(ajobs.GetApiName(), ajobs.FrontendAntJobs)
	// a_jobkill: kill active job on ant
	storage.UserFunctions[ajobkill.GetApiName()] = object.NewNativeFunc(ajobkill.GetApiName(), ajobkill.FrontendAntJobkill)
	// a_kill: kill process on target OS
	storage.UserFunctions[akill.GetApiName()] = object.NewNativeFunc(akill.GetApiName(), akill.FrontendAntKill)
	// a_mv: move files/directories
	storage.UserFunctions[amv.GetApiName()] = object.NewNativeFunc(amv.GetApiName(), amv.FrontendAntMv)
	// a_mkdir: create directory
	storage.UserFunctions[amkdir.GetApiName()] = object.NewNativeFunc(amkdir.GetApiName(), amkdir.FrontendAntMkdir)
	// a_exec_assembly: execute .NET in CLR
	storage.UserFunctions[aexecassembly.GetApiName()] = object.NewNativeFunc(aexecassembly.GetApiName(), aexecassembly.FrontendAntExecuteAssembly)
	// a_download: download file from target FS
	storage.UserFunctions[adownload.GetApiName()] = object.NewNativeFunc(adownload.GetApiName(), adownload.FrontendAntDownload)
	// a_upload: upload file to target FS
	storage.UserFunctions[aupload.GetApiName()] = object.NewNativeFunc(aupload.GetApiName(), aupload.FrontendAntUpload)
	// a_pause: one-time pause of ant's execution (one-time sleep)
	storage.UserFunctions[apause.GetApiName()] = object.NewNativeFunc(apause.GetApiName(), apause.FrontendAntPause)
	// a_destruct: ant's self-destruction
	storage.UserFunctions[adestruct.GetApiName()] = object.NewNativeFunc(adestruct.GetApiName(), adestruct.FrontendAntDestruct)
	// a_exec_detach: execute binary with arguments with detaching
	storage.UserFunctions[aexecdetach.GetApiName()] = object.NewNativeFunc(aexec.GetApiName(), aexecdetach.FrontendAntExecDetach)
	// a_shell: execute shell command
	storage.UserFunctions[ashell.GetApiName()] = object.NewNativeFunc(ashell.GetApiName(), ashell.FrontendAntShell)
	// a_ppid: spoof PPID
	storage.UserFunctions[appid.GetApiName()] = object.NewNativeFunc(appid.GetApiName(), appid.FrontendAntPpid)
	// a_exit: stop ant's execution
	storage.UserFunctions[aexit.GetApiName()] = object.NewNativeFunc(aexit.GetApiName(), aexit.FrontendAntExit)
	// t_cancel: cancel all operator's tasks with status NEW
	storage.UserFunctions[tcancel.GetApiName()] = object.NewNativeFunc(tcancel.GetApiName(), tcancel.FrontendTasksCancel)
	// is_windows: is ant running on windows
	storage.UserFunctions[iswindows.GetApiName()] = object.NewNativeFunc(iswindows.GetApiName(), iswindows.FrontendIsWindows)
	// is_linux: is ant running on linux
	storage.UserFunctions[islinux.GetApiName()] = object.NewNativeFunc(islinux.GetApiName(), islinux.FrontendIsLinux)
	// is_macos: is ant running on macos
	storage.UserFunctions[ismacos.GetApiName()] = object.NewNativeFunc(ismacos.GetApiName(), ismacos.FrontendIsMacos)
	// is_x64: is arch x64
	storage.UserFunctions[isx64.GetApiName()] = object.NewNativeFunc(isx64.GetApiName(), isx64.FrontendIsX64)
	// is_x86: is arch x86
	storage.UserFunctions[isx86.GetApiName()] = object.NewNativeFunc(isx86.GetApiName(), isx86.FrontendIsX86)
	// is_arm64: is arch arm64
	storage.UserFunctions[isarm64.GetApiName()] = object.NewNativeFunc(isarm64.GetApiName(), isarm64.FrontendIsArm64)
	// is_arm32: is arch arm32
	storage.UserFunctions[isarm32.GetApiName()] = object.NewNativeFunc(isarm32.GetApiName(), isarm32.FrontendIsArm32)
}

var (
	//go:embed builtin/*.pico
	builtinScriptsFS embed.FS
)

// registerBuiltin registers builtin scripts
func registerBuiltin() error {
	// list of scripts
	e, err := builtinScriptsFS.ReadDir("builtin")
	if err != nil {
		return err
	}
	for _, v := range e {
		// read data from script
		data, err := builtinScriptsFS.ReadFile(fmt.Sprintf("builtin/%s", v.Name()))
		if err != nil {
			return errors.Wrapf(err, "read %s", v.Name())
		}
		// create AST
		tree, err := mlanUtils.CreateAST(string(data))
		if err != nil {
			return errors.Wrap(err, v.Name())
		}
		// visit AST
		visitor := visitor.NewVisitor()
		if res := visitor.Visit(tree); res != types.Success {
			return errors.Wrapf(visitor.GetError(), "evaluation %s", v.Name())
		}
	}
	return nil
}
