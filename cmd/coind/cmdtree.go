package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/iuhjui/autocfg"
	"gitee.com/iuhjui/logger"
	"github.com/spf13/cobra"

	"github.com/coind-io/coind/config"
	"github.com/coind-io/coind/core"
	"github.com/coind-io/coind/version"
)

type CmdTree struct {
	logger *logger.Logger
	config *autocfg.Config
	root   *cobra.Command
}

func NewCmdTree() *CmdTree {
	tree := new(CmdTree)
	tree.makeRoot()
	tree.makeVersion()
	tree.makeWeb()
	return tree
}

func (tree *CmdTree) makeRoot() error {
	root := new(cobra.Command)
	root.Use = "coind"
	root.Short = "Coind is a crypto currency keystone-coin of core node"
	root.PersistentFlags().StringP("config", "c", "", "path to an explicit configuration file")
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// 输出版本信息
		tree.logger.Infof("Coind version:  {{version}} ({{githash}})", map[string]string{
			"version": version.Version,
			"githash": version.GitHash,
		})
		// 初始配置
		cfgname, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		if cfgname == "" {
			cfgname = "./coind.yaml"
		}
		cfg, err := config.LoadConfig(cfgname)
		if err != nil {
			return err
		}
		tree.config = cfg
		tree.logger.Infof("using config file: {{path}}", map[string]string{
			"path": cfgname,
		})
		return nil
	}
	tree.root = root

	root.RunE = func(cmd *cobra.Command, args []string) error {
		// 初始根实例
		deps := core.NewDeps()
		deps.SerLogger(tree.logger)
		deps.SetConfig(tree.config)
		daemon, err := core.NewDaemon(deps)
		if err != nil {
			return err
		}
		// 阻塞终端
		done := make(chan os.Signal, 0)
		signal.Notify(done, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
		<-done
		err = daemon.Close()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (tree *CmdTree) makeWeb() {
	web := new(cobra.Command)
	web.Use = "web"
	web.Short = "Serve an HTTP endpoint on the given host and port."
	web.RunE = func(cmd *cobra.Command, args []string) error {
		// 初始根实例
		deps := core.NewDeps()
		deps.SerLogger(tree.logger)
		deps.SetConfig(tree.config)
		daemon, err := core.NewDaemon(deps)
		if err != nil {
			return err
		}
		// 阻塞终端
		done := make(chan os.Signal, 0)
		signal.Notify(done, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
		<-done
		err = daemon.Close()
		if err != nil {
			return err
		}
		return nil
	}
	// 注册
	tree.root.AddCommand(web)
	return
}

func (tree *CmdTree) makeVersion() {
	ver := new(cobra.Command)
	ver.Use = "version"
	ver.Short = "Prints the version of Coind."
	ver.PersistentPreRun = func(cmd *cobra.Command, args []string) {}
	ver.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Coind version:", version.Version)
		fmt.Println("Git commit hash:", version.GitHash)
		if version.BuildDate != "" {
			fmt.Println("Build date:", version.BuildDate)
		}
	}
	// 注册
	tree.root.AddCommand(ver)
	return
}

func (tree *CmdTree) SetLogger(logger *logger.Logger) {
	tree.logger = logger
	return
}

func (tree *CmdTree) Execute() error {
	return tree.root.Execute()
}
