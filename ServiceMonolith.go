package IrisAPIs

import (
	"context"
	"github.com/pkg/errors"
)

// ServiceMonolith This service monolith contains all service in IrisAPIs
type ServiceMonolith struct {
	ChatBotService          *ChatbotContext
	CurrencyService         CurrencyService
	DatabaseContext         *DatabaseContext
	IpNationService         *IpNationContext
	ApiKeyService           ApiKeyService
	ServiceMgmt             ServiceManagement
	ArticleProcessorService ArticleProcessorService
	PbsTrafficDataService   PbsTrafficDataService
	BuildInfoService        BuildInfoService
	teardownQueue           []TeardownableServices
	cancelFunc              context.CancelFunc
}

func NewServiceMonolith(config *Configuration) *ServiceMonolith {
	return &ServiceMonolith{}
}

func (m *ServiceMonolith) ReInitServices(ctx context.Context, config *Configuration) error {
	logger := GetLogger(ctx)
	db, err := NewDatabaseContext(config.ConnectionString, true, nil)
	if err != nil {
		//If failed to initialize, will stop re-init
		return errors.Wrap(err, "Error initializing database!")
	}

	for _, s := range m.teardownQueue {
		err := s.Teardown()
		if err != nil {
			logger.Warning("%v teardown failed!")
		}
	}

	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	m.CurrencyService = m.registerService(NewCurrencyContextWithConfig(config.FixerIoApiKey,
		config.FixerIoLastFetchSuccessfulPeriod, config.FixerIoLastFetchFailedPeriod, db)).(CurrencyService)
	m.ChatBotService = m.registerService(NewChatbotContext()).(*ChatbotContext)
	m.DatabaseContext = m.registerService(db).(*DatabaseContext)
	m.IpNationService = m.registerService(NewIpNationContext(db)).(*IpNationContext)
	m.ApiKeyService = m.registerService(NewApiKeyService(db)).(ApiKeyService)
	m.ServiceMgmt = m.registerService(func() ServiceManagement {
		service := NewServiceManagement()
		_ = service.RegisterPresetServices(ctx)
		return service
	}()).(ServiceManagement)
	m.ArticleProcessorService = m.registerService(NewArticleProcessorContext()).(ArticleProcessorService)
	m.BuildInfoService = m.registerService(NewBuildInfoService()).(BuildInfoService)
	m.PbsTrafficDataService = m.registerService(NewPbsTrafficDataService(db)).(PbsTrafficDataService)
	return nil
}

func (m *ServiceMonolith) registerService(service interface{}) interface{} {
	if m.teardownQueue == nil {
		m.teardownQueue = make([]TeardownableServices, 0)
	}
	teardownableService, isTeardownable := service.(TeardownableServices)

	if isTeardownable {
		m.teardownQueue = append(m.teardownQueue, teardownableService)
	}

	return service
}
