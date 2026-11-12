//go:build windows

package server

import (
	"errors"
	"path/filepath"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

func inspectStoragePermission(kind, path string) StoragePermissionItem {
	result := StoragePermissionItem{Kind: kind, State: "unavailable"}
	if path == "" { return result }
	abs, err := filepath.Abs(path)
	if err != nil { return result }
	volume := filepath.VolumeName(abs)
	// Do not connect to UNC or mapped network drives. Ordinary local drive
	// letters are supported, including portable/removable installations.
	if len(volume) != 2 || volume[1] != ':' {
		result.State = "skipped_nonlocal"
		return result
	}
	root, err := windows.UTF16PtrFromString(volume + "\\")
	if err != nil { return result }
	switch windows.GetDriveType(root) {
	case 2, 3, 6: // removable, fixed, RAM disk
	default:
		result.State = "skipped_nonlocal"
		return result
	}
	// Inspect ancestors from the drive root down. Do not intentionally follow
	// junctions/symlinks to other trees. This is a snapshot, not a race-proof
	// filesystem confinement mechanism.
	var parents []string
	for p := filepath.Dir(abs); p != filepath.Dir(p); p = filepath.Dir(p) {
		parents = append(parents, p)
	}
	for i := len(parents)-1; i >= 0; i-- {
		p, err := windows.UTF16PtrFromString(parents[i])
		if err != nil { return result }
		attr, err := windows.GetFileAttributes(p)
		if err != nil {
			if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) || errors.Is(err, windows.ERROR_PATH_NOT_FOUND) { result.State = "missing" }
			return result
		}
		if attr & windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 { result.State = "skipped_reparse"; return result }
	}
	name, err := windows.UTF16PtrFromString(abs)
	if err != nil { return result }
	h, err := windows.CreateFile(name, windows.READ_CONTROL,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) || errors.Is(err, windows.ERROR_PATH_NOT_FOUND) { result.State = "missing" }
		return result
	}
	defer windows.CloseHandle(h)
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(h, &info); err != nil { return result }
	if info.FileAttributes & windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 { result.State = "skipped_reparse"; return result }
	sd, err := windows.GetSecurityInfo(h, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil { return result }
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil { return result }
	return inspectStorageDescriptor(kind, sd, user.User.Sid.String())
}

// Only ordinary allow ACEs are interpreted. Complex ACEs are marked partial;
// deny rules and group membership are not evaluated as effective access.
func inspectStorageDescriptor(kind string, sd *windows.SECURITY_DESCRIPTOR, currentSID string) StoragePermissionItem {
	result := StoragePermissionItem{Kind: kind, State: "unavailable"}
	defer runtime.KeepAlive(sd)
	control, _, err := sd.Control()
	if err != nil { return result }
	result.InheritanceEnabled = control & windows.SE_DACL_PROTECTED == 0
	acl, _, err := sd.DACL()
	if err != nil && !errors.Is(err, windows.ERROR_OBJECT_NOT_FOUND) { return result }
	result.State = "inspected"
	// Null or absent DACL permits unrestricted access, unlike an empty ACL.
	if acl == nil || errors.Is(err, windows.ERROR_OBJECT_NOT_FOUND) {
		result.BroadRead = true
		return result
	}
	for i := uint32(0); i < uint32(acl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(acl, i, &ace); err != nil || ace == nil {
			result.State = "partial"; continue
		}
		if ace.Header.AceFlags & windows.INHERIT_ONLY_ACE != 0 { continue }
		switch ace.Header.AceType {
		case windows.ACCESS_DENIED_ACE_TYPE:
			continue
		case windows.ACCESS_ALLOWED_ACE_TYPE:
		default:
			result.State = "partial"; continue
		}
		if ace.Header.AceSize < 16 { result.State = "partial"; continue }
		if ace.Mask & (windows.FILE_READ_DATA|windows.GENERIC_READ|windows.GENERIC_ALL) == 0 { continue }
		sid := (*windows.SID)(unsafe.Pointer(&ace.SidStart))
		if !sid.IsValid() || sid.Len()+8 > int(ace.Header.AceSize) { result.State = "partial"; continue }
		switch sid.String() {
		case "S-1-1-0", "S-1-5-11", "S-1-5-32-545", "S-1-5-32-546":
			result.BroadRead = true
		case currentSID, "S-1-5-18", "S-1-5-32-544":
		default:
			result.OtherRead = true
		}
	}
	return result
}
