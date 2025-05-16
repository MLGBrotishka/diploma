package entity

// Перечисление типов прав доступа
type Permission int32

const (
	PermissionNone          Permission = 0
	PermissionCreate        Permission = 1 // Право на создание сущностей.
	PermissionApply         Permission = 2 // Право на применение изменений.
	PermissionRollback      Permission = 3 // Право на откат изменений.
	PermissionList          Permission = 4 // Право на просмотр списков.
	PermissionGet           Permission = 5 // Право на получение конкретной сущности.
	PermissionApplyOther    Permission = 6 // Право на применение изменений, созданных другими.
	PermissionRollbackOther Permission = 7 // Право на откат изменений, созданных другими.
)

func (p Permission) String() string {
	return Permission_name[p]
}

var Permission_name = map[Permission]string{
	PermissionNone:          "PERMISSION_NONE",
	PermissionCreate:        "PERMISSION_CREATE",
	PermissionApply:         "PERMISSION_APPLY",
	PermissionRollback:      "PERMISSION_ROLLBACK",
	PermissionList:          "PERMISSION_LIST",
	PermissionGet:           "PERMISSION_GET",
	PermissionApplyOther:    "PERMISSION_APPLY_OTHER",
	PermissionRollbackOther: "PERMISSION_ROLLBACK_OTHER",
}
