package command

import (
	"authentik-go/config"
	"authentik-go/core"
	"authentik-go/database"
	"authentik-go/i18n"
	"authentik-go/router"
	"authentik-go/router/middleware"
	"authentik-go/utils/common"
	"authentik-go/web"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	ginI18n "github.com/gin-contrib/i18n"
)

var rootCommand = &cobra.Command{
	Use:   "hios",
	Short: "启动服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.CONF.System.Host == "" {
			config.CONF.System.Host = "0.0.0.0"
		}
		if config.CONF.System.Port == "" {
			config.CONF.System.Port = "3376"
		}
		if config.CONF.System.Cache == "" {
			config.CONF.System.Cache = common.RunDir("/.cache")
		}
		if config.CONF.System.Dsn == "" {
			config.CONF.System.Dsn = fmt.Sprintf("sqlite3://%s/%s", config.CONF.System.Cache, "database.db")
		}
		config.CONF.System.Start = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
		//
		err := common.WriteFile(config.CONF.System.Cache+"/config.json", common.StructToJson(config.CONF.System))
		if err != nil {
			common.PrintError(fmt.Sprintf("配置文件写入失败: %s", err.Error()))
			os.Exit(1)
		}
		// 初始化db
		err = core.InitDB()
		if err != nil {
			common.PrintError(fmt.Sprintf("数据库加载失败: %s", err.Error()))
			os.Exit(1)
		}
		// 初始化数据库
		err = database.Init()
		if err != nil {
			common.PrintError(fmt.Sprintf("数据库初始化失败: %s", err.Error()))
			os.Exit(1)
		}
		// 初始化工作目录
		common.Mkdir("work", 0777)
		common.Mkdir("work/logs", 0777)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 启动服务端
		t, err := template.New("index").Parse(string(web.IndexByte))
		if err != nil {
			common.PrintError(fmt.Sprintf("模板解析失败: %s", err.Error()))
			os.Exit(1)
		}
		if config.CONF.System.Mode == "debug" {
			gin.SetMode(gin.DebugMode)
		} else if config.CONF.System.Mode == "test" {
			gin.SetMode(gin.TestMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		// 设置日志输出到文件
		file, _ := os.OpenFile("./work/logs/request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer file.Close()
		gin.DefaultWriter = file
		gin.DefaultErrorWriter = file

		// 设置路由
		routers := gin.Default()
		routers.Use(middleware.OperationLog())
		routers.Use(gzip.Gzip(gzip.DefaultCompression))
		routers.SetHTMLTemplate(t)
		routers.Use(i18n.GinI18nLocalize())
		routers.SetFuncMap(template.FuncMap{
			"Localize": ginI18n.GetMessage,
		})
		routers.Any("/*path", func(context *gin.Context) {
			router.Init(context)
		})
		//
		common.PrintSuccess("启动成功: http://localhost:" + config.CONF.System.Port)
		//
		routers.Run(fmt.Sprintf("%s:%s", config.CONF.System.Host, config.CONF.System.Port))
		//
	},
}

func Execute() {
	godotenv.Load(".env")
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	rootCommand.Flags().StringVar(&config.CONF.System.Host, "host", os.Getenv("HOST"), "主机名，默认：0.0.0.0")
	rootCommand.Flags().StringVar(&config.CONF.System.Port, "port", os.Getenv("PORT"), "端口号，默认：3376")
	rootCommand.Flags().StringVar(&config.CONF.System.Mode, "mode", os.Getenv("MODE"), "运行模式，可选：debug|test|release")
	rootCommand.Flags().StringVar(&config.CONF.System.Cache, "cache", "", "数据缓存目录，默认：{RunDir}/.cache")
	rootCommand.Flags().StringVar(&config.CONF.System.WssUrl, "wss", "", "服务端生成的url")
	rootCommand.Flags().StringVar(&config.CONF.System.Dsn, "dsn", os.Getenv("DB_DSN"), "数据来源名称，如：sqlite://{CacheDir}/database.db")
	rootCommand.Flags().StringVar(&config.CONF.Jwt.SecretKey, "secret_key", "base64:ONdadQs1W4pY3h3dzr1jUSPrqLdsJQ9tCBZnb7HIDtk=", "jwt密钥")
	rootCommand.Flags().StringVar(&config.CONF.Redis.RedisUrl, "redis_url", "redis://localhost:56379", "RedisUrl")
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
