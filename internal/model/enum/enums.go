package enum

type ComplianceStatusType string

const (
	ComplianceStatusCreated        ComplianceStatusType = "created"
	ComplianceStatusRejected       ComplianceStatusType = "rejected"
	ComplianceStatusActionRequired ComplianceStatusType = "actionRequired"
	ComplianceStatusFailed         ComplianceStatusType = "failed"
	//below entries used for ComplianceCheck
	ComplianceStatusPending ComplianceStatusType = "pending"
	ComplianceStatusPassed  ComplianceStatusType = "passed"
)

var ComplianceUnfinishedStatuses = []ComplianceStatusType{ComplianceStatusCreated, ComplianceStatusFailed, ComplianceStatusActionRequired, ComplianceStatusPending}

var ComplianceFinishedStatuses = []ComplianceStatusType{ComplianceStatusRejected, ComplianceStatusPassed}

type KYCProviderType string

const (
	KYCProviderFinclusive KYCProviderType = "finclusive"
	KYCProviderPwgMock    KYCProviderType = "pwgMock"
)

type PwgEntityType string

const (
	PwgEntityUsersAccount PwgEntityType = "usersAccount"
	PwgEntityInstitution  PwgEntityType = "institution"
)
