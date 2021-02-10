package cmd

import (
	"fmt"
	"os"

	"github.com/linuxing3/vpsman/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dbPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vpsman",
	Short: "Simple vps manager",
	Long: `A simple vps manager that can control
system services backend by sqlite.`,
	Run: func(cmd *cobra.Command, args []string) {
		mainMenu(dbPath)
	},
}

func mainMenu(dbPath string) {
exit:
	for {
		fmt.Println()
		fmt.Println(util.Cyan("欢迎使用xray管理程序"))
		fmt.Println()
		menuList := []string{"用户管理", "Xray管理", "Nginx管理", "Trojan管理", "web管理"}
		switch util.LoopInput("请选择: ", menuList, false) {
		case 1:
			userMenu(dbPath)
		case 2:
			xrayMenu()
		case 3:
			nginxMenu()
		case 4:
			trojanMenu()
		case 5:
			webMenu()
		default:
			break exit
		}
	}
}

// Execute 执行rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vpsman.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&dbPath, "db", "", "./vpsman.db", "数据库目录地址")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".vpsman")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
