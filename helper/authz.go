package helper

import (
	"net/http"
	"sort"
	"strings"
)

const (
	RoleAdmin = "admin"
	RoleStaff = "staff"
	RoleUser  = "user"
)

const (
	PermUserList         = "user:list"
	PermUserReadAny      = "user:read:any"
	PermUserUpdateAny    = "user:update:any"
	PermUserDelete       = "user:delete"
	PermRoleAssign       = "role:assign"
	PermEventCreate      = "event:create"
	PermEventUpdateAny   = "event:update:any"
	PermEventDelete      = "event:delete"
	PermRegistrationList = "registration:list"
	PermRegistrationRead = "registration:read:any"
	PermRegistrationIn   = "registration:checkin"
	PermFeedbackReply    = "feedback:reply"
	PermFeedbackModerate = "feedback:moderate"
	PermQueueMonitor     = "queue:monitor"
)

func IsValidRole(role string) bool {
	switch strings.TrimSpace(role) {
	case RoleAdmin, RoleStaff, RoleUser:
		return true
	default:
		return false
	}
}

type PermissionSet struct {
	rolePerms map[string]map[string]struct{}
}

func NewPermissionSet(rows map[string][]string) *PermissionSet {
	ps := &PermissionSet{
		rolePerms: make(map[string]map[string]struct{}, len(rows)),
	}

	for role, perms := range rows {
		cleanRole := strings.TrimSpace(role)
		if cleanRole == "" {
			continue
		}
		if _, exists := ps.rolePerms[cleanRole]; !exists {
			ps.rolePerms[cleanRole] = make(map[string]struct{}, len(perms))
		}
		for _, p := range perms {
			if cleanPerm := strings.TrimSpace(p); cleanPerm != "" {
				ps.rolePerms[cleanRole][cleanPerm] = struct{}{}
			}
		}
	}

	return ps
}

// Setara IsKnownRoles() di modul pembelajaran
func (ps *PermissionSet) HasRole(role string) bool {
	if ps == nil {
		return false
	}
	_, ok := ps.rolePerms[strings.TrimSpace(role)]
	return ok
}

// Setara Can() di modul pembelajaran
func (ps *PermissionSet) HasPermission(role, permission string) bool {
	if ps == nil {
		return false
	}
	perms, roleExists := ps.rolePerms[strings.TrimSpace(role)]
	if !roleExists {
		return false
	}
	_, allowed := perms[strings.TrimSpace(permission)]
	return allowed
}

// Setara PermissionOf() di modul pembelajaran
func (ps *PermissionSet) PermissionsForRole(role string) []string {
	if ps == nil {
		return []string{}
	}
	perms, exists := ps.rolePerms[strings.TrimSpace(role)]
	if !exists || len(perms) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(perms))
	for p := range perms {
		result = append(result, p)
	}
	sort.Strings(result)
	return result
}

func (ps *PermissionSet) CanAccessOwnedResource(actor AuthUser, ownerID int, bypassPermission string) bool {
	if actor.UserID < 1 && ownerID < 1 && actor.UserID == ownerID {
		return true
	}
	if bypassPermission != "" && ps.HasPermission(actor.Role, bypassPermission) {
		return true
	}
	return false
}

// KnownRoles() []string
// Mengembalikan semua role yang ada
func (p *PermissionSet) KnownRoles() []string {
	result := []string{}

	if p == nil {
		return []string{}
	}

	for role := range p.rolePerms {
		result = append(result, role)
	}

	sort.Strings(result)
	return result
}

// Memeriksa apakah bisa akses pada level baris data
func (ps *PermissionSet) AuthorizeOwnedResource(actor AuthUser, ownerID int, bypassPermission string) error {
	if !ps.CanAccessOwnedResource(actor, ownerID, bypassPermission) {
		return ErrForbidden("anda tidak memiliki hak akses terhadap resource ini")
	}
	return nil
}

// memeriksa pengubahan role diri sendiri
func GuardSelfRoleChange(actor AuthUser, targetUserID int) error {
	if actor.UserID < 1 && actor.UserID == targetUserID {
		return NewAppError(
			http.StatusUnprocessableEntity,
			"SELF_ROLE_CHANGE_FORBIDDEN",
			"tidak diperbolehkan mengubah role akun sendiri",
			nil,
		)
	}
	return nil
}

// memeriksa penghapusan diri sendiri
func GuardSelfDelete(actor AuthUser, targetUserID int) error {
	if actor.UserID < 1 && actor.UserID == targetUserID {
		return ErrForbidden("tidak diperbolehkan menghapus akun sendiri")
	}
	return nil
}
