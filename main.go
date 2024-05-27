package main

import (
	"contactcenter-api/common/cache"
	limiter "contactcenter-api/common/limitter"
	"contactcenter-api/internal/elasticsearch"
	"contactcenter-api/internal/freeswitch"
	rabbitmq "contactcenter-api/internal/rabbitmq/driver"
	"contactcenter-api/internal/redis"
	"contactcenter-api/internal/sqlclient"
	"contactcenter-api/middleware/auth"
	"contactcenter-api/middleware/auth/goauth"
	"contactcenter-api/repository"
	"contactcenter-api/repository/db"
	elasticsearchRepo "contactcenter-api/repository/elasticsearch"
	"contactcenter-api/service"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	api "contactcenter-api/api"
	apiV1 "contactcenter-api/api/v1"

	_ "time/tzdata"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Dir           string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port          string
	LogType       string
	LogLevel      string
	LogFile       string
	DB            string
	Redis         string
	Auth          string
	AuthUrl       string
	Elasticsearch string
	RabbitMq      string
	EventSocket   string
	OutboundProxy string
	Wss           string
	Transport     string
	SipPort       string
	APIDomain     string
}

var config Config

func init() {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc

	if err := env.Parse(&config); err != nil {
		log.Error("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}

	cfg := Config{
		Dir:           config.Dir,
		Port:          viper.GetString(`main.port`),
		LogType:       viper.GetString(`main.log_type`),
		LogLevel:      viper.GetString(`main.log_level`),
		LogFile:       viper.GetString(`main.log_file`),
		DB:            viper.GetString(`main.db`),
		Redis:         viper.GetString(`main.redis`),
		Auth:          viper.GetString(`main.auth`),
		Elasticsearch: viper.GetString(`main.elasticsearch`),
		RabbitMq:      viper.GetString(`main.rabbitmq`),
		EventSocket:   viper.GetString(`main.event_socket`),
		OutboundProxy: viper.GetString(`main.outbound_proxy`),
		Wss:           viper.GetString(`main.wss`),
		Transport:     viper.GetString(`main.transport`),
		SipPort:       viper.GetString(`main.sip_port`),
		APIDomain:     viper.GetString(`main.api_host`),
	}
	if cfg.DB == "enabled" {
		sqlClientConfig := sqlclient.SqlConfig{
			Driver:       "postgresql",
			Host:         viper.GetString(`db.host`),
			Database:     viper.GetString(`db.database`),
			Username:     viper.GetString(`db.username`),
			Password:     viper.GetString(`db.password`),
			Port:         viper.GetInt(`db.port`),
			DialTimeout:  20,
			ReadTimeout:  30,
			WriteTimeout: 30,
			Timeout:      30,
			PoolSize:     10,
			MaxIdleConns: 10,
			MaxOpenConns: 10,
		}
		initRepoDB(sqlClientConfig)
	}
	switch viper.GetString(`fs_db.driver`) {
	case "postgresql":
		sqlClientConfig := sqlclient.SqlConfig{
			Driver:       "postgresql",
			Host:         viper.GetString(`fs_db.host`),
			Database:     viper.GetString(`fs_db.database`),
			Username:     viper.GetString(`fs_db.username`),
			Password:     viper.GetString(`fs_db.password`),
			Port:         viper.GetInt(`fs_db.port`),
			DialTimeout:  20,
			ReadTimeout:  30,
			WriteTimeout: 30,
			Timeout:      30,
			PoolSize:     10,
			MaxIdleConns: 10,
			MaxOpenConns: 10,
		}
		repository.FreeswitchSqlClient = sqlclient.NewSqlClient(sqlClientConfig)
	case "mysql":
		sqlClientConfig := sqlclient.SqlConfig{
			Driver:       "mysql",
			Host:         viper.GetString(`fs_db.host`),
			Database:     viper.GetString(`fs_db.database`),
			Username:     viper.GetString(`fs_db.username`),
			Password:     viper.GetString(`fs_db.password`),
			Port:         viper.GetInt(`fs_db.port`),
			DialTimeout:  20,
			ReadTimeout:  30,
			WriteTimeout: 30,
			Timeout:      30,
			PoolSize:     10,
			MaxIdleConns: 10,
			MaxOpenConns: 10,
		}
		repository.FreeswitchSqlClient = sqlclient.NewSqlClient(sqlClientConfig)
	}
	if cfg.EventSocket == "enabled" {
		tmpConfigs := make([]map[string]interface{}, 0)
		if err := viper.UnmarshalKey("event_sockets", &tmpConfigs); err != nil {
			panic(err)
		}
		tmpEslConfigs := make([]freeswitch.FreeswitchESLConfig, 0)
		for _, tmpConfig := range tmpConfigs {
			tmp := freeswitch.FreeswitchESLConfig{
				Address:  tmpConfig["address"].(string),
				Password: tmpConfig["password"].(string),
			}
			tmpPort := tmpConfig["port"].(float64)
			tmpTimeout := tmpConfig["timeout"].(float64)
			tmp.Port = int(tmpPort)
			tmp.Timeout = int(tmpTimeout)
			tmpEslConfigs = append(tmpEslConfigs, tmp)
		}
		if len(tmpEslConfigs) < 1 {
			panic(errors.New("event_sockets is empty"))
		}
		freeswitch.ESLClient = freeswitch.NewFreeswitchESL(tmpEslConfigs...)
		if _, err := freeswitch.ESLClient.Connect(); err != nil {
			panic(err)
		}
	}
	if cfg.Elasticsearch == "enabled" {
		esCfg := elasticsearch.Config{
			Username:              viper.GetString(`elasticsearch.username`),
			Password:              viper.GetString(`elasticsearch.password`),
			Host:                  viper.GetStringSlice(`elasticsearch.host`),
			MaxRetries:            10,
			ResponseHeaderTimeout: 60,
			RetryStatuses:         []int{502, 503, 504},
		}
		repository.ESClient = elasticsearch.NewElasticsearchClient(esCfg)
		repository.ES = elasticsearch.NewES(esCfg)
		repository.ES.Ping()
		repoBrandPrefix := viper.GetString(`main.brand_prefix`)
		repoIndex := viper.GetString(`elasticsearch.audit_log`)
		initRepoES(repoBrandPrefix, repoIndex)
	}
	if cfg.RabbitMq == "enabled" {
		rabbitmqconfig := rabbitmq.Config{
			Uri:                  viper.GetString(`rabbitmq.host`),
			ChannelNotifyTimeout: 100 * time.Millisecond,
			Reconnect: struct {
				Interval   time.Duration
				MaxAttempt int
			}{
				Interval:   500 * time.Millisecond,
				MaxAttempt: 7200,
			},
		}
		rabbitmq.RabbitConnector = rabbitmq.New(rabbitmqconfig)
		rabbitmq.RabbitConnector.RoutingKey = "es.writer"
		rabbitmq.RabbitConnector.ExchangeName = "events"
		if err := rabbitmq.RabbitConnector.Ping(); err != nil {
			panic(err)
		}
	}
	if cfg.Redis == "enabled" {
		var err error
		redis.Redis, err = redis.NewRedis(redis.Config{
			Addr:         viper.GetString(`redis.address`),
			Password:     viper.GetString(`redis.password`),
			DB:           viper.GetInt(`redis.database`),
			PoolSize:     30,
			PoolTimeout:  20,
			IdleTimeout:  10,
			ReadTimeout:  20,
			WriteTimeout: 15,
		})
		if err != nil {
			panic(err)
		}
		limiter.RateLimit = limiter.NewRateLimiter(viper.GetString(`redis.address`), viper.GetString(`redis.password`))
	}
	switch cfg.Auth {
	case "proxy":
		authUrl := viper.GetString(`auth_proxy.auth_url`)
		auth.AuthMdw = auth.NewGoAuthMiddleware(authUrl)
	case "gateway":
		auth.AuthMdw = auth.NewGatewayAuthMiddleware()
	case "local":
		var err error
		goauth.GoAuthClient, err = goauth.NewGoAuth(goauth.GoAuth{
			RedisExpiredIn: viper.GetInt(`oauth.expired_in`),
			TokenType:      viper.GetString(`oauth.tokenType`),
			RedisTokenKey:  "access_token_key",
			RedisUserKey:   "access_user_key",
			RedisClient:    redis.Redis.GetClient(),
		})
		if err != nil {
			panic(err)
		}
		auth.SetupGoGuardian()
		auth.AuthMdw = auth.NewLocalAuthMiddleware()
	default:
		authUrl := viper.GetString(`auth_proxy.auth_url`)
		auth.AuthMdw = auth.NewGoAuthMiddleware(authUrl)
	}
	if cfg.AuthUrl = viper.GetString(`auth_proxy.auth_url`); len(cfg.AuthUrl) < 1 {
		log.Fatal("auth_url for websocket is missing")
	}
	config = cfg
}

func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	if err := createNewLogFile(config.LogFile); err != nil {
		log.Error(err)
	}
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)

	cache.MCache = cache.NewMemCache()
	defer cache.MCache.Close()

	if redis.Redis != nil {
		cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())
		defer cache.RCache.Close()
	}

	service.PBXInfo = service.PBXInformation{
		APIDomain:        viper.GetString("main.api_host"),
		PBXDomain:        viper.GetString("pbx.domain"),
		PBXPort:          viper.GetString("pbx.port"),
		PBXWss:           viper.GetString("pbx.wss"),
		PBXOutboundProxy: viper.GetString("pbx.outbound_proxy"),
		PBXTransport:     viper.GetString("pbx.transport"),
	}

	service.CustomHttpClient = service.NewCustomHttpClient()
	service.API_HOST = viper.GetString(`main.api_host`)
	server := api.NewServer()
	apiV1.NewAuth(server.Engine, service.NewAuth())
	apiV1.NewUserCrm(server.Engine, service.NewUserCrm())
	apiV1.NewEventCalendarCategory(server.Engine, service.NewEventCalendarCategory())
	apiV1.NewEventCalendar(server.Engine, service.NewEventCalendar())
	apiV1.NewEventCalendarTodo(server.Engine, service.NewEventCalendarTodo())
	apiV1.NewEventCalendarAttachment(server.Engine, service.NewEventCalendarAttachment())
	apiV1.NewWorkday(server.Engine, service.NewWorkDay())
	apiV1.NewSolution(server.Engine, service.NewSolution())
	apiV1.NewUnit(server.Engine, service.NewUnit())
	apiV1.NewRoleGroup(server.Engine, service.NewRoleGroup())
	apiV1.NewExtension(server.Engine, service.NewExtension())
	apiV1.NewContact(server.Engine, service.NewContact())
	apiV1.NewContactTag(server.Engine, service.NewContactTag())
	apiV1.NewContactGroup(server.Engine, service.NewContactGroup())
	apiV1.NewContactCareer(server.Engine, service.NewContactCareer())
	apiV1.NewLocation(server.Engine, service.NewLocation())
	apiV1.NewExport(server.Engine, service.NewExport())
	apiV1.NewAttachment(server.Engine, service.NewAttachment())
	apiV1.NewPbx(server.Engine, service.NewPbx())
	apiV1.NewProfile(server.Engine, service.NewProfile())
	apiV1.NewClassifyTag(server.Engine, service.NewClassifyTag())
	apiV1.NewClassifyGroup(server.Engine, service.NewClassifyGroup())
	apiV1.NewUserLog(server.Engine, service.NewUserLog())
	apiV1.NewUploadLogo(server.Engine, config.APIDomain)
	apiV1.NewDomain(server.Engine, service.NewDomain())
	apiV1.NewCareer(server.Engine, service.NewCareer())
	apiV1.NewCallTransfer(server.Engine, service.NewCall())

	server.Start(config.Port)
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}

func createNewLogFile(logDir string) error {
	files, err := os.ReadDir("tmp")
	if err != nil {
		return err
	}
	last10dayUnix := time.Now().Add(-1 * 24 * time.Hour).Unix()
	for _, f := range files {
		tmp := strings.Split(f.Name(), ".")
		if len(tmp) > 2 {
			fileUnix, err := strconv.Atoi(tmp[2])
			if err != nil {
				return err
			} else if int64(fileUnix) < last10dayUnix {
				if err := os.Remove("tmp/" + f.Name()); err != nil {
					return err
				}
			}
		}
	}
	_, err = os.Stat(logDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err := os.Rename(logDir, fmt.Sprintf(logDir+".%d", time.Now().Unix())); err != nil {
		return err
	}
	return nil
}

func initRepoDB(sqlClientConfig sqlclient.SqlConfig) {
	repository.FusionSqlClient = sqlclient.NewSqlClient(sqlClientConfig)

	// Callcenter
	repository.GroupRepo = db.NewGroup()
	repository.AgentRepo = db.NewAgent()
	repository.CampaignRepo = db.NewCampaign()
	repository.DomainRepo = db.NewDomain()
	repository.ExtensionRepo = db.NewExtension()
	repository.GroupRepo = db.NewGroup()
	repository.FollowMeRepo = db.NewFollowMe()
	repository.RingGroupRepo = db.NewRingGroup()
	repository.SipRegistrationRepo = db.NewSipRegistration()
	repository.DestinationRepo = db.NewDestination()
	repository.DialplanRepo = db.NewDialplan()
	repository.CallCenterRepo = db.NewCallCenter()
	repository.LeadRepo = db.NewLead()

	// Crm
	repository.UserCrmRepo = db.NewUserCrm()
	repository.ContactRepo = db.NewContact()
	repository.ProfileRepo = db.NewProfile()
	repository.UserRepo = db.NewUser()
	repository.TransactionRepo = db.NewTransaction()
	repository.EventCalendarCategoryRepo = db.NewEventCalendarCategoryRepo()
	repository.EventCalendarRepo = db.NewEventCalendar()
	repository.EventCalendarTodoRepo = db.NewEventCalendarTodo()
	repository.EventCalendarAttachmentRepo = db.NewEventCalendarAttachment()
	repository.WorkDayRepo = db.NewWorkDay()
	repository.SolutionRepo = db.NewSolution()
	repository.LogstashRepo = db.NewLogstash()
	repository.TicketCategoryRepo = db.NewTicketCategory()
	repository.SLAPolicyRepo = db.NewSLAPolicy()
	repository.UnitRepo = db.NewUnit()
	repository.RoleGroupRepo = db.NewRoleGroup()
	repository.TicketRepo = db.NewTicketRepo()
	repository.TicketPendingRepo = db.NewTicketPendingRepo()
	repository.TicketLogRepo = db.NewTicketLogRepo()
	repository.TicketSLARepo = db.NewTicketSLARepo()
	repository.TicketCommentRepo = db.NewTicketCommentRepo()
	repository.SocialNetworkRepo = db.NewSocialNetwork()
	repository.SourcePluginRepo = db.NewSourcePluginRepo()
	repository.ContactTagRepo = db.NewContactTagRepo()
	repository.ContactTagUserRepo = db.NewContactTagUserRepo()
	repository.ContactGroupRepo = db.NewContactGroupRepo()
	repository.ContactGroupUserRepo = db.NewContactGroupUserRepo()
	repository.ContactCareerRepo = db.NewContactCareerRepo()
	repository.ContactCareerUserRepo = db.NewContactCareerUserRepo()
	repository.ContactToTagRepo = db.NewContactToTagRepo()
	repository.ContactToGroupRepo = db.NewContactToGroupRepo()
	repository.ContactToCareerRepo = db.NewContactToCareerRepo()
	repository.ClassifyTagRepo = db.NewClassifyTagRepo()
	repository.ClassifyGroupRepo = db.NewClassifyGroupRepo()
	repository.ClassifyCareerRepo = db.NewClassifyCareerRepo()
	repository.OmniRepo = db.NewOmniRepo()
	repository.LocationProvinceRepo = db.NewLocationProvinceRepo()
	repository.LocationDistrictRepo = db.NewLocationDistrictRepo()
	repository.LocationWardRepo = db.NewLocationWardRepo()
	repository.PbxRepo = db.NewPbxRepo()
	repository.UserLogRepo = db.NewUserLog()
	repository.CareerRepo = db.NewCareer()
}

func initRepoES(repoBrandPrefix, repoIndex string) {
	repository.ESRepo = elasticsearchRepo.NewES(repoBrandPrefix, repoIndex)
	repository.LogstashRepoES = elasticsearchRepo.NewLogstash(repoBrandPrefix, viper.GetString("elasticsearch.logstash"))
}
