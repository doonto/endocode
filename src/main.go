package main

import (
	"context"
	"endocode/src/daemon"
	myerror "endocode/src/error"
	mylog "endocode/src/log"
	"fmt"
	"github.com/urfave/cli/v2" // All security-sensitive server parameters are passed via the command line.
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Build parameters
var (
	version   = "0.0.3"
	commit    = "unset"
	buildTime = "unset"
)

// Globab variables
var (
	httpConfigFileFlag string
	logFileFlag        string
	listenStringFlag   string
	debugFlag          string
	httpUserIDFlag     string
	httpUserPwdFlag    string
	jwtKeyFlag         string
)

// flags
var flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "httpconfig, httpcfg",
		Usage:       "HTTP Config file name",
		Required:    true,
		Destination: &httpConfigFileFlag,
	},
	&cli.StringFlag{
		Name:        "listenstring, l",
		Usage:       "Listen string in format <host>:<port>",
		Required:    true,
		Destination: &listenStringFlag,
		Value:       "localhost:3000",
	},
	&cli.StringFlag{
		Name:        "httpuser, httpu",
		Usage:       "User name for access to HTTP server",
		Required:    true,
		Destination: &httpUserIDFlag,
	},
	&cli.StringFlag{
		Name:        "httppassword, httppwd",
		Usage:       "User password for access to HTTP server",
		Required:    true,
		Destination: &httpUserPwdFlag,
	},
	&cli.StringFlag{
		Name:        "jwtkey, jwtk",
		Usage:       "JSON web token secret key",
		Required:    false,
		Destination: &jwtKeyFlag,
	},
	&cli.StringFlag{
		Name:        "debug, d",
		Usage:       "Debug mode: DEBUG, INFO, ERROR",
		Required:    false,
		Destination: &debugFlag,
		Value:       "INFO",
	},
	&cli.StringFlag{
		Name:        "logfile, log",
		Usage:       "Log file name",
		Required:    false,
		Destination: &logFileFlag,
	},
}

//main function
func main() {
	// Create new Application
	app := cli.NewApp()
	app.Name = "HTTP Server"
	app.Version = fmt.Sprintf("%s, commit '%s', build time '%s'", version, commit, buildTime)
	app.Flags = flags // присваиваем ранее определенные флаги
	app.Writer = os.Stderr

	// Определяем действие - запуск демона
	app.Action = func(ctx *cli.Context) (myerr error) {

		// настраиваем параллельное логирование в файл
		if logFileFlag != "" {
			// добавляем в имя лог файла дату и время
			logFileFlag = strings.Replace(logFileFlag, "%s", time.Now().Format("2006_01_02_150405"), 1)

			// открываем лог файл на запись
			logFile, err := os.OpenFile(logFileFlag, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				myerr = myerror.WithCause("6020", "Error open log file: Filename", err, logFileFlag)
				mylog.PrintfErrorMsg(fmt.Sprintf("%+v", myerr))
				return
			}

			// закрываем лог файл по выходу
			defer func() {
				if logFile != nil {

					defer logFile.Close() // ошибку закрытия игнорируем

					// flushing write buffers out to disks
					err := logFile.Sync()

					if err != nil {
						// ошибку через закрытие передаем на уровень выше
						myerr = myerror.WithCause("6020", "Error sync log file before closing", err).PrintfInfo()
					}
				}
			}()

			// Параллельно пишем в os.Stderr и файл
			wrt := io.MultiWriter(os.Stderr, logFile)

			// Переопределяем глобальный логер на кастомный
			mylog.InitLogger(wrt)
		} else {
			mylog.InitLogger(os.Stderr)
		}

		mylog.PrintfInfoMsg("Server is starting up: Version, Logfile", app.Version, logFileFlag)

		// Установим фильтр логирования
		if debugFlag != "" {
			mylog.PrintfInfoMsg("Set log level", debugFlag)
			switch debugFlag {
			case "DEBUG", "ERROR", "INFO":
				mylog.SetFilter(debugFlag)
			default:
				myerr = myerror.New("9001", "Incorrect debugFlag. Only avaliable: DEBUG, INFO, ERROR.", debugFlag)
				mylog.PrintfErrorMsg(fmt.Sprintf("%+v", myerr))
				return
			}
		}

		// Создаем конфигурацию демона
		daemonCfg := &daemon.Config{
			ConfigFileName: httpConfigFileFlag,
			ListenSpec:     listenStringFlag,
			JwtKey:         []byte(jwtKeyFlag),
			HTTPUserID:     httpUserIDFlag,
			HTTPUserPwd:    httpUserPwdFlag,
		}

		// Создаем демон
		daemon, myerr := daemon.New(context.Background(), daemonCfg)
		if myerr != nil {
			mylog.PrintfErrorMsg(fmt.Sprintf("%+v", myerr)) // верхний уровень логирования с трассировкой
			return
		}

		// Стартуем демон и ожидаем завершения
		if myerr = daemon.Run(); myerr != nil {
			mylog.PrintfErrorMsg(fmt.Sprintf("%+v", myerr)) // верхний уровень логирования с трассировкой
			return
		}

		mylog.PrintfInfoMsg("Server is shutdown")
		return
	}

	// Запускаем приложение
	if err := app.Run(os.Args); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

//
//func main() {
//	http.HandleFunc("/helloworld", helloHandler)
//	http.HandleFunc("/versionz", versionz)
//
//	fmt.Printf("Starting server at port 8080\n")
//	if err := http.ListenAndServe(":8080", nil); err != nil {
//		log.Fatal(err)
//	}
//}
//
//func helloHandler(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/helloworld" {
//		http.Error(w, "404 not found.", http.StatusNotFound)
//		return
//	}
//
//	if r.Method != "GET" {
//		http.Error(w, "Method is not supported.", http.StatusNotFound)
//		return
//	}
//
//	fmt.Fprintf(w, "Hello Stranger")
//}
//
//func versionz(w http.ResponseWriter, r *http.Request) {
//	name := chi.URLParam(r, "name")
//	response := HelloResponse{
//		Message: fmt.Sprintf("Hello %s!", name),
//	}
//	jsonResponse(w, response, http.StatusOK)
//}
