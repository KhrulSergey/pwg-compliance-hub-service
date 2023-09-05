package storage

import (
	"github.com/google/uuid"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/model/enum"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
	"gorm.io/gorm"
)

// complianceRepository is a type that handles
type complianceRepository struct {
	dbConnector *gorm.DB
	logger      logger.Logger
}

// ComplianceRepository defines Compliance object actions.
type ComplianceRepository interface {
	FindComplianceCheckByExternalGuid(guid string) (*model.ComplianceCheck, error)
	SaveComplianceCheck(c *model.ComplianceCheck) error
	UpdateComplianceCheck(c *model.ComplianceCheck) error

	FindComplianceCheckRequest(pwgEntityGuid string, pwgEntityType enum.PwgEntityType) (*model.ComplianceCheckRequest, error)
	GetUnfinishedComplianceCheckRequest() ([]model.ComplianceCheckRequest, error)
	SaveComplianceCheckRequest(c *model.ComplianceCheckRequest) error
	UpdateComplianceCheckRequest(c *model.ComplianceCheckRequest) error
}

const (
	complianceCheckTable        = "compliance_check"
	complianceCheckRequestTable = "compliance_check_request"
)

// InitComplianceRepository creates repository layer for Compliance entities.
func InitComplianceRepository(db *gorm.DB, logger logger.Logger) ComplianceRepository {
	return &complianceRepository{dbConnector: db, logger: logger}
}

func (cr *complianceRepository) FindComplianceCheckByExternalGuid(guid string) (*model.ComplianceCheck, error) {
	var checkEntity model.ComplianceCheck
	tx := cr.dbConnector.Raw("SELECT * FROM compliance_checks WHERE external_guid = ?", guid).First(&checkEntity)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &checkEntity, nil
}
func (cr *complianceRepository) SaveComplianceCheck(check *model.ComplianceCheck) error {
	check.Id = uuid.New()
	result := cr.dbConnector.Create(check)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (cr *complianceRepository) UpdateComplianceCheck(check *model.ComplianceCheck) error {
	result := cr.dbConnector.Save(check)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cr *complianceRepository) FindComplianceCheckRequest(pwgEntityGuid string, pwgEntityType enum.PwgEntityType) (*model.ComplianceCheckRequest, error) {
	var checkRequest model.ComplianceCheckRequest
	tx := cr.dbConnector.Raw("SELECT * FROM compliance_check_requests WHERE pwg_entity_guid = ? and pwg_entity_type = ?", pwgEntityGuid, pwgEntityType).First(&checkRequest)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &checkRequest, nil
}

func (cr *complianceRepository) GetUnfinishedComplianceCheckRequest() ([]model.ComplianceCheckRequest, error) {
	unfinishedStatuses := enum.ComplianceUnfinishedStatuses
	var checkRequests []model.ComplianceCheckRequest
	checkRequestList := cr.dbConnector.Raw("SELECT * FROM compliance_check_requests WHERE status in ? ORDER BY created_at", unfinishedStatuses).Find(&checkRequests)
	if checkRequestList.Error != nil {
		return nil, checkRequestList.Error
	}
	return checkRequests, nil
}

func (cr *complianceRepository) SaveComplianceCheckRequest(checkRequest *model.ComplianceCheckRequest) error {
	result := cr.dbConnector.Create(checkRequest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cr *complianceRepository) UpdateComplianceCheckRequest(checkRequest *model.ComplianceCheckRequest) error {
	result := cr.dbConnector.Save(checkRequest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (cr *complianceRepository) SaveCompliance(c *model.Compliance) error {
	result := cr.dbConnector.Create(c)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
