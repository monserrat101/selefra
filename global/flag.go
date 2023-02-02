package global

import (
	"os"
	"sync"

	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/spf13/cobra"
)

// Variable store some global variable
type Variable struct {
	readOnlyVariable

	mux sync.RWMutex

	// token is not empty when user is login
	token   string
	orgName string
	stage   string
}

// readOnlyVariable will only be set when programmer started
type readOnlyVariable struct {
	once sync.Once

	// workspace store where selefra worked
	workspace string

	// cmd store what command is running
	cmd string

	logLevel string
}

var g = Variable{
	readOnlyVariable: readOnlyVariable{
		once: sync.Once{},
	},
	mux: sync.RWMutex{},
}

func Init(cmd, workspace, level string) {
	g.once.Do(func() {
		g.cmd = cmd
		if _, ok := levelMap[level]; ok {
			g.logLevel = level
		} else {
			g.logLevel = defaultLogLevel
		}

		if workspace != "" {
			g.workspace = workspace
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}

		g.workspace = cwd
	})
}

// WrappedInit wrapper the Init function to a cobra func
func WrappedInit(workspace, level string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		Init(cmd.Name(), workspace, level)
	}
}

// DefaultWrappedInit is a cobra func that will use default value to init Variable
func DefaultWrappedInit() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		Init(cmd.Name(), "", "")
	}
}

func SetToken(token string) {
	g.mux.Lock()
	defer g.mux.Unlock()

	g.token = token
}

func SetStage(stage string) {
	g.mux.Lock()
	defer g.mux.Unlock()

	g.stage = stage
}

func WorkSpace() string {
	return g.workspace
}

func Token() string {
	g.mux.RLock()
	defer g.mux.RUnlock()

	return g.token
}

func OrgName() string {
	g.mux.RLock()
	defer g.mux.RUnlock()

	return g.orgName
}

func Cmd() string {
	return g.cmd
}

func Stage() string {
	g.mux.RLock()
	defer g.mux.RUnlock()

	return g.stage
}

func LogLevel() string {
	return g.logLevel
}

var levelMap = map[string]bool{
	"trace":   true,
	"debug":   true,
	"info":    true,
	"warning": true,
	"error":   true,
	"fatal":   true,
}

var defaultLogLevel = "error"

// TODO: will be deprecated
var (
	WORKSPACE  = pointer.ToStringPointer(".")
	LOGINTOKEN = ""
	ORGNAME    = ""
	CMD        = ""
	STAG       = ""
	LOGLEVEL   = "error"
)

var o sync.Once

func ChangeLevel(level string) {
	if levelMap[level] {
		o.Do(func() {
			LOGLEVEL = level
		})
	}
}

const PkgBasePath = "ghcr.io/selefra/postgre_"
const PkgTag = ":latest"

var SERVER = "main-api.selefra.io"
