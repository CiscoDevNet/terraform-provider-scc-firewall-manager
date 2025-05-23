package transactiontype

type Type string

const (
	ONBOARD_ASA                        Type = "ONBOARD_ASA"
	ONBOARD_IOS                        Type = "ONBOARD_IOS"
	ONBOARD_DUO_ADMIN_PANEL            Type = "ONBOARD_DUO_ADMIN_PANEL"
	CREATE_FTD                         Type = "CREATE_FTD"
	REGISTER_FTD                       Type = "REGISTER_FTD"
	DELETE_CDFMC_MANAGED_FTD           Type = "DELETE_CDFMC_MANAGED_FTD"
	RECONNECT_ASA                      Type = "RECONNECT_ASA"
	READ_ASA                           Type = "READ_ASA"
	DEPLOY_ASA_DEVICE_CHANGES          Type = "DEPLOY_ASA_DEVICE_CHANGES"
	MSP_CREATE_TENANT                  Type = "MSP_CREATE_TENANT"
	MSP_ADD_USERS_TO_TENANT            Type = "MSP_ADD_USERS_TO_TENANT"
	MSP_DELETE_USERS_FROM_TENANT       Type = "MSP_DELETE_USERS_FROM_TENANT"
	MSP_ADD_USER_GROUPS_TO_TENANT      Type = "MSP_ADD_USER_GROUPS_TO_TENANT"
	MSP_DELETE_USER_GROUPS_FROM_TENANT Type = "MSP_DELETE_USER_GROUPS_FROM_TENANT"
	UPGRADE_ASA                        Type = "UPGRADE_ASA"
	UPGRADE_FTD                        Type = "UPGRADE_FTD"
)
