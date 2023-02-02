package configure

import (
	"rbi-nontrading-operation/internal/config"
	"rbi-nontrading-operation/pkg/controller"
	"rbi-nontrading-operation/pkg/logic/backup_brokerage_agreements"
	"rbi-nontrading-operation/pkg/logic/main_flow"
	"rbi-nontrading-operation/pkg/storage/handlers/calling_the_CF_procedure"
	"rbi-nontrading-operation/pkg/storage/handlers/get_last_open_date"
	"rbi-nontrading-operation/pkg/storage/handlers/get_last_update_date_time"
	"rbi-nontrading-operation/pkg/storage/handlers/get_new_data_above_update_date_time"
	"rbi-nontrading-operation/pkg/storage/handlers/get_userId_by_crmId"
	"rbi-nontrading-operation/pkg/storage/handlers/insert_new_backup_data"
	"rbi-nontrading-operation/pkg/storage/handlers/insert_new_data"
	"rbi-nontrading-operation/pkg/storage/handlers/new_backup_brokerage_agreements"

	"github.com/jmoiron/sqlx"
)

func UsrInfoController(
	ctrl *controller.ControllerImpl,
	investDb *sqlx.DB,
	diaDb *sqlx.DB,
	env config.Environment,
) (e error) {
	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	// db invest common handlers
	propogateErr(ctrl.RegisterHandler(get_last_update_date_time.NewGetLastUpdateDateTimeHandler(investDb, env.InvDBTableName)))
	propogateErr(ctrl.RegisterHandler(insert_new_data.NewInsertNewDataHandler(investDb, env.InvDBTableName)))
	// db dia common handlers
	propogateErr(ctrl.RegisterHandler(get_new_data_above_update_date_time.NewGetNewDataAboveUpdateDateTimeHandler(diaDb, env.DiaDBTableName)))
	// flow
	propogateErr(ctrl.RegisterHandler(main_flow.NewMainFlowHandler(ctrl)))
	// уточнить
	propogateErr(ctrl.RegisterHandler(backup_brokerage_agreements.NewBackupBrokerageAgreementsHandler(ctrl)))

	propogateErr(ctrl.RegisterHandler(get_last_open_date.NewGetLastUpdateDateTimeHandler(investDb)))
	//
	propogateErr(ctrl.RegisterHandler(new_backup_brokerage_agreements.NewGetLastUpdateDateTimeHandler(diaDb)))

	propogateErr(ctrl.RegisterHandler(calling_the_CF_procedure.NewcallingTheCFProcedureHandler(investDb)))

	propogateErr(ctrl.RegisterHandler(get_userId_by_crmId.NewGetUserIdByCrmIdHandler(investDb)))

	propogateErr(ctrl.RegisterHandler(insert_new_backup_data.NewInsertNewBackupBrokerageQueryHandler(investDb)))

	return
}
