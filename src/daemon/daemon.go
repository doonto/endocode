package daemon

import (
	"context"
	myerror "endocode/src/error"
	"endocode/src/httpserver"
	mylog "endocode/src/log"
	"github.com/sasbury/mini"
	"os"
	"os/signal"
	"syscall"
)

// Daemon repesent top level daemon
type Daemon struct {
	ctx    context.Context    // корневой контекст
	cancel context.CancelFunc // функция закрытия корневого контекста
	cfg    *Config            // конфигурация демона

	// Сервисы демона
	httpServer      *httpserver.Server // HTTP сервер
	httpServerErrCh chan error         // канал ошибок для HTTP сервера
}

// Config repesent daemon options
type Config struct {
	ConfigFileName string // основной файл конфигурации
	ListenSpec     string // строка HTTP листенера
	JwtKey         []byte // JWT secret key
	HTTPUserID     string // пользователь для HTTP Basic Authentication
	HTTPUserPwd    string // пароль для HTTP Basic Authentication

	// Конфигурация вложенных сервисов
	httpServerCfg httpserver.Config // конфигурация HTTP сервера
}

// New create Daemon
func New(ctx context.Context, cfg *Config) (*Daemon, error) {
	var err error
	var config *mini.Config

	mylog.PrintfInfoMsg("Create new daemon")

	{ // входные проверки
		if cfg == nil {
			return nil, myerror.New("6030", "Empty daemon config").PrintfInfo()
		}
		if cfg.ConfigFileName == "" {
			return nil, myerror.New("6030", "Empty config file name").PrintfInfo()
		}
	} // входные проверки

	// Создаем новый демон
	daemon := &Daemon{
		cfg:             cfg,
		httpServerErrCh: make(chan error, 1), // канал ошибок HTTP сервера
	}

	// создаем корневой контекст с отменой
	if ctx == nil {
		daemon.ctx, daemon.cancel = context.WithCancel(context.Background())
	} else {
		daemon.ctx, daemon.cancel = context.WithCancel(ctx)
	}

	// Загружаем конфигурационный файл
	if config, err = loadConfigFile(daemon.cfg.ConfigFileName); err != nil {
		return nil, err
	}

	{ // создаем HTTP server
		// Настраиваем конфигурацию HTTP server
		if err = loadHTTPServerConfig(config, &daemon.cfg.httpServerCfg); err == nil {
			daemon.cfg.httpServerCfg.ListenSpec = daemon.cfg.ListenSpec
		} else {
			return nil, err
		}

		{ // Настраиваем конфигурацию HTTP service
			// Параметры из командной строки
			daemon.cfg.httpServerCfg.ServiceCfg.HTTPUserID = daemon.cfg.HTTPUserID
			daemon.cfg.httpServerCfg.ServiceCfg.HTTPUserPwd = daemon.cfg.HTTPUserPwd
			daemon.cfg.httpServerCfg.ServiceCfg.JwtKey = daemon.cfg.JwtKey
			// Параметры уровня HTTP сервера
			daemon.cfg.httpServerCfg.ServiceCfg.UseTLS = daemon.cfg.httpServerCfg.UseTLS
			daemon.cfg.httpServerCfg.ServiceCfg.UseHSTS = daemon.cfg.httpServerCfg.UseHSTS
			daemon.cfg.httpServerCfg.ServiceCfg.MaxBodyBytes = daemon.cfg.httpServerCfg.MaxBodyBytes

			if err = loadHTTPServiceConfig(config, &daemon.cfg.httpServerCfg.ServiceCfg); err == nil {

				// Настраиваем конфигурацию HTTP Logger
				if err = loadHTTPLoggerConfig(config, &daemon.cfg.httpServerCfg.ServiceCfg.LogCfg); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} // Настраиваем конфигурацию HTTP service

		// Создаем HTTP server
		if daemon.httpServer, err = httpserver.New(daemon.ctx, daemon.httpServerErrCh, &daemon.cfg.httpServerCfg); err != nil {
			return nil, err
		}
	} // создаем HTTP server

	mylog.PrintfInfoMsg("New daemon is created")

	return daemon, nil
}

// Run daemon and wait for system signal or error in error chanel
func (d *Daemon) Run() error {
	mylog.PrintfInfoMsg("Starting daemon")

	// запускаем в фоне HTTP сервер, возврат в канал ошибок
	go func() { d.httpServerErrCh <- d.httpServer.Run() }()

	mylog.PrintfInfoMsg("Daemon is running. For exit <CTRL-c>")

	// подписываемся на системные прирывания
	syscalCh := make(chan os.Signal, 1) // канал системных прирываний
	signal.Notify(syscalCh, syscall.SIGINT, syscall.SIGTERM)

	// ожидаем прерывания или возврат в канал ошибок
	select {
	case s := <-syscalCh: //  ожидаем системное прирывание
		mylog.PrintfInfoMsg("Exiting, got signal", s)
		d.Shutdown() // останавливаем daemon
		return nil
	case err := <-d.httpServerErrCh: // возврат от HTTP сервера в канал ошибок
		mylog.PrintfInfoMsg("Exiting, got error")
		mylog.PrintfErrorInfo(err) // логируем ошибку
		d.Shutdown()               // останавливаем daemon
		return err
	}
}

// Shutdown daemon
func (d *Daemon) Shutdown() {
	mylog.PrintfInfoMsg("Shutting down daemon")

	// Закрываем корневой контекст
	defer d.cancel()

	// Останавливаем HTTP сервер, ожидаем завершения активных подключений
	if myerr := d.httpServer.Shutdown(); myerr != nil {
		mylog.PrintfErrorInfo(myerr) // дополнительно логируем результат остановки
	}

	// ... Останавливаем остальные сервисы

	mylog.PrintfInfoMsg("Daemon is shutdown")
}
