package service

import (
	"github.com/go-co-op/gocron"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/dto"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest/mapper"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/storage"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
	"slices"
	"sync/atomic"
	"time"
)

// ComplianceService specifies the function required for the auth complianceService
type ComplianceService interface {
	CreateRequestAndCheckCompliance(checkRequestDto *dto.ComplianceCheckRequestDto) (dto.ComplianceCheckShortResponseDto, error)
	GetAllUnfinishedComplianceRequest() ([]model.ComplianceCheckRequest, error)
	UpdateComplianceRequest(checkRequest *model.ComplianceCheckRequest, checkResponseDto *dto.ComplianceCheckResponseDto) error
}

// complianceService wrapper for the interface implementation
type complianceService struct {
	logger                       logger.Logger
	complianceStorage            storage.ComplianceRepository
	accountOperatorServiceClient AccountOperatorServiceClient
	externalComplianceService    ExternalComplianceService
	updateTaskIndex              int64
}

// InitComplianceService constructs an instance of the Auth ComplianceService
func InitComplianceService(logger logger.Logger, complianceStorage storage.ComplianceRepository,
	externalComplianceService ExternalComplianceService, accountOperatorServiceClient AccountOperatorServiceClient) ComplianceService {
	cs := &complianceService{
		logger:                       logger,
		complianceStorage:            complianceStorage,
		externalComplianceService:    externalComplianceService,
		accountOperatorServiceClient: accountOperatorServiceClient,
		updateTaskIndex:              0,
	}
	// Schedule updateComplianceRequestStatusesTask to run every 10 second
	s := gocron.NewScheduler(time.UTC)

	go func() {
		// Wait for 1 minute before starting the task
		<-time.After(10 * time.Second)
		// Schedule the task to run every minute
		job, err := s.Every(10).Second().Do(cs.updateComplianceRequestStatusesTask)
		if err != nil {
			// handle the error related to setting up the job
			cs.logger.Error("updateComplianceRequestStatusesTask couldn't run. Check code")
			panic(err)
		} else {
			cs.logger.Info("updateComplianceRequestStatusesTask was started successfully, job id: ", job)
		}
	}()

	s.StartAsync()
	return cs
}

//todo run schedule task in async way. Here example
//func parallel() {
//	list := []int{1, 2, 3, 4, 5}
//
//	// Create a channel to receive the results
//	resultCh := make(chan int)
//
//	// Create a wait group to synchronize goroutines
//	var wg sync.WaitGroup
//
//	// Launch goroutines to process each element in parallel
//	for _, value := range list {
//		wg.Add(1)
//		go func(v int) {
//			// Process the element
//			result := process(v)
//
//			// Send the result to the channel
//			resultCh <- result
//
//			wg.Done()
//		}(value)
//	}
//
//	// Start a goroutine to close the result channel once all goroutines are done
//	go func() {
//		wg.Wait()
//		close(resultCh)
//	}()
//
//	// Collect the results from the channel
//	for result := range resultCh {
//		// Do something with the result
//		fmt.Println(result)
//	}
//}
//
//func process(value int) int {
//	// Perform some computation on the value
//	return value * value
//}

func (cs *complianceService) CreateRequestAndCheckCompliance(checkRequestDto *dto.ComplianceCheckRequestDto) (dto.ComplianceCheckShortResponseDto, error) {
	cs.logger.Debugf("Try to send compliance request for pwgId: %w, pwgType: %w to provider: %w, and then save short data to DB",
		checkRequestDto.PwgEntityGuid, checkRequestDto.PwgEntityType, checkRequestDto.Provider)
	requestEntity := mapper.ToComplianceRequestEntity(*checkRequestDto)

	//todo reveal validation and check if data is the same
	////check if requests already existed in DB
	//existedRequestEntity, _ := cs.complianceStorage.FindComplianceCheckRequest(requestEntity.PwgEntityGuid, requestEntity.PwgEntityType)
	//if existedRequestEntity != nil {
	//	response := mapper.ToComplianceResponseDtoFromRequestEntity(*existedRequestEntity)
	//	return response, nil
	//}

	//sending data to external provider
	shortResponseDto, err := cs.externalComplianceService.SendComplianceRequestToProvider(checkRequestDto)
	if err != nil {
		cs.logger.Errorf("Error while sending compliance request: %w, external response: %w, error: %w", checkRequestDto, shortResponseDto, err)
		if shortResponseDto != nil {
			requestEntity.RequestExternalGuid = shortResponseDto.ExternalGuid
		}
		requestEntity.Status = enum.ComplianceStatusFailed
	} else {
		requestEntity.RequestedAt = time.Now()
		requestEntity.RequestExternalGuid = shortResponseDto.ExternalGuid
		if shortResponseDto.Status != "" {
			requestEntity.Status = shortResponseDto.Status
		} else {
			requestEntity.Status = enum.ComplianceStatusCreated
			shortResponseDto.Status = enum.ComplianceStatusCreated
		}
	}
	//save compliance request to DB
	errStorage := cs.complianceStorage.SaveComplianceCheckRequest(&requestEntity)
	if errStorage != nil {
		cs.logger.Panicf("Error occurred while trying to save request entity: %w, raw check request dto: %w", requestEntity, checkRequestDto)
		err = errStorage
	}
	return *shortResponseDto, err
}

func (cs *complianceService) GetAllUnfinishedComplianceRequest() ([]model.ComplianceCheckRequest, error) {
	return cs.complianceStorage.GetUnfinishedComplianceCheckRequest()
}

func (cs *complianceService) updateComplianceRequestStatusesTask() {
	atomic.AddInt64(&cs.updateTaskIndex, 1) // Increment the counter by 1
	cs.logger.Debugf("task-%d, Start running updateComplianceRequestStatusesTask...", cs.updateTaskIndex)
	unfinishedComplianceRequests, err := cs.GetAllUnfinishedComplianceRequest()
	if err != nil {
		cs.logger.Panicf("task-%d, Error occurred while trying to load all unfinished compliance requests. Check error: %w", cs.updateTaskIndex, err)
		return
	}
	if len(unfinishedComplianceRequests) == 0 {
		cs.logger.Debugf("task-%d, Nothing to do in updateComplianceRequestStatusesTask...", cs.updateTaskIndex)
		return
	}

	for _, checkRequest := range unfinishedComplianceRequests {
		cs.logger.Debugf("task-%d, Try to send update compliance status for request externalGuid: %v, for pwgId: %v,"+
			" pwgType: %v to provider: %v, and then save short data to DB", cs.updateTaskIndex, checkRequest.RequestExternalGuid,
			checkRequest.PwgEntityGuid, checkRequest.PwgEntityType, checkRequest.Provider)
		//get new data from external provider
		responseDto, err := cs.externalComplianceService.GetComplianceStatusFromProvider(checkRequest)
		if err != nil || responseDto == nil {
			cs.logger.Panicf("task-%d, Error occurred while trying to load all unfinished compliance requests. Check error: %v",
				cs.updateTaskIndex, err)
			continue
		}
		//check Compliance status
		if responseDto.Status == checkRequest.Status {
			//there is nothing to do. Compliance request is still unchanged in Provider
			continue
		}
		//send updated status to AO service
		err = cs.accountOperatorServiceClient.SendComplianceResponse(*responseDto)
		if err != nil {
			cs.logger.Errorf("task-%d, Error while sending compliance request: %v, external response: %v, error: %v",
				cs.updateTaskIndex, checkRequest, responseDto, err)
			continue
		}
		//update status in DB entries
		_ = cs.UpdateComplianceRequest(&checkRequest, responseDto)
		cs.logger.Debugf("task-%d, updateComplianceRequestStatusesTask finished", cs.updateTaskIndex)
	}
}

func (cs *complianceService) UpdateComplianceRequest(checkRequest *model.ComplianceCheckRequest, checkResponseDto *dto.ComplianceCheckResponseDto) error {
	mapper.ApplyChangesToComplianceRequestEntity(checkRequest, checkResponseDto)
	var err error
	if slices.Contains(enum.ComplianceFinishedStatuses, checkResponseDto.Status) { //check if we need to create finished compliance entry for current check request
		//save finished compliance check to DB as new entry
		checkEntity := mapper.ToComplianceCheckEntity(*checkResponseDto, *checkRequest)
		err = cs.createOrUpdateComplianceCheck(&checkEntity)
		if err != nil {
			cs.logger.Panicf("Error occurred while trying to save compliance entity: %w, raw check response dto: %w", checkEntity, checkResponseDto)
			return err
		}
	}
	//update existed compliance request in DB
	err = cs.complianceStorage.UpdateComplianceCheckRequest(checkRequest)
	if err != nil {
		cs.logger.Panicf("Error occurred while trying to save compliance check request entity: %w, raw check response dto: %w", checkRequest, checkResponseDto)
	}
	return err
}

func (cs *complianceService) createOrUpdateComplianceCheck(newCheckEntity *model.ComplianceCheck) error {
	existedCheck, err := cs.complianceStorage.FindComplianceCheckByExternalGuid(newCheckEntity.ExternalGuid)
	if existedCheck != nil {
		mapper.ApplyChangesToComplianceCheckEntity(existedCheck, newCheckEntity)
		err = cs.complianceStorage.UpdateComplianceCheck(existedCheck)
	} else {
		err = cs.complianceStorage.SaveComplianceCheck(newCheckEntity)
	}
	return err
}
