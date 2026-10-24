package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "example.invalid/denisenkom/go-mssqldb" // driver
	"example.invalid/gin-contrib/cors"
	"example.invalid/gin-gonic/gin"
	"example.invalid/spf13/cobra"
	"example.invalid/spf13/viper"

	"example.invalid/ServerTeam/Servers/cmd/slot88/auth/service"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

var port int

var svc service.Service

var rootCmd = &cobra.Command{
	Use:   "slot88-auth",
	Short: "slot88-auth provide token stuff, generator, exchange user info",
	RunE: func(cmd *cobra.Command, args []string) error {

		go func() {
			for {
				service.CleanExpireToken(db)
				time.Sleep(viper.GetDuration("token.clean_during"))
			}
		}()

		r := gin.Default()
		r.LoadHTMLGlob("views/**/*")

		r.GET("", func(c *gin.Context) {
			c.String(200, "HelloWorld")
		})

		adm := r.Group("adm")
		{
			adm.Use(gin.BasicAuth(viper.GetStringMapString("admin.users")))
			adm.GET("", func(c *gin.Context) {
				c.HTML(http.StatusOK, "token/gen.html", gin.H{})
			})
		}

		token := r.Group("token")

		cfg := cors.DefaultConfig()
		cfg.AllowOrigins = []string{"*"}
		token.Use(cors.New(cfg))
		{
			token.GET("/", func(c *gin.Context) {

				token := c.Query("mD5")

				info, err := service.GetToken(db, token)
				if err != nil {
					c.IndentedJSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": fmt.Sprintf("Token Error, %v", err),
					})
					return
				}

				c.IndentedJSON(http.StatusOK, gin.H{
					"success": true,
					"message": "OK",
					"data":    info,
				})
				return
			})

			token.GET(":token", func(c *gin.Context) {

				token := c.Param("token")

				info, err := service.GetToken(db, token)
				if err != nil {
					c.IndentedJSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": fmt.Sprintf("Token Error, %v", err),
					})
					return
				}

				c.IndentedJSON(http.StatusOK, gin.H{
					"success": true,
					"message": "OK",
					"info":    info,
				})
				return
			})

			token.POST("gen", func(c *gin.Context) {

				info := new(example.invalid)

				if err := c.BindJSON(info); err != nil {
					c.IndentedJSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": fmt.Sprintf("Argument Error %v", err),
					})
					return
				}

				fmt.Println(info.PlayerInfo)
				fmt.Println(info.ServerInfo)
				now := time.Now()
				exp := now.Add(viper.GetDuration("token.expire"))

				info.Expire = exp.Format(time.RFC3339)

				_i, err := svc.GetInfo(db, info.PlayerInfo.PlayerID, info.ServerInfo.GameID, info.ServerInfo.ServerID)
				if err != nil {
					mlog.Error("Get Info Error, %v", err)
					c.IndentedJSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": "Find Info Erorr, Please Check arguments is right",
					})
					return
				}

				info.PlayerInfo = _i.PlayerInfo
				info.ServerInfo = _i.ServerInfo

				token, err := service.GenToken(*info)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": fmt.Sprintf("Generate Token Error, %v", err),
					})
					return
				}

				code, err := service.SaveToken(db, token)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": fmt.Sprintf("Save Token Error, %v", err),
					})
					return
				}

				c.IndentedJSON(http.StatusOK, gin.H{
					"success": true,
					"message": "OK",
					"token":   code,
				})
				return
			})
		}

		serv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		}

		mlog.System("Start on %s", serv.Addr)
		if err := serv.ListenAndServe(); err != nil {
			return err
		}

		return nil
	},
}

var db *sql.DB

func init() {

	conf.CLogLevel = mlog.DebugLevel
	conf.ClearLog = true

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		mlog.Fatal("Read Config File Error, %v\n", err)
	}

	example.invalid("Using config file: %s", viper.ConfigFileUsed())

	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 9001, "Listen port")

	initDB()

	initSvc()
}

func initDB() {
	var err error
	db, err = sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=%s;",
		viper.GetString("db.host"),
		viper.GetString("db.database"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetInt("db.port"),
		viper.GetString("db.encrypt"),
	))
	if err != nil {
		mlog.Fatal("Connect DB Error, %v", err)
	}
}

func initSvc() {
	t := strings.ToLower(viper.GetString("type"))

	switch t {
	case "dd":
		svc = new(service.DDService)
	case "target":
		svc = new(service.PGSService)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
