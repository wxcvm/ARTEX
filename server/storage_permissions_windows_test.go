//go:build windows

package server

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

func TestInspectStorageDescriptor(t *testing.T) {
	const current = "S-1-5-21-123-456-789-1001"
	for _, tc := range []struct {
		name, sddl string
		broad, other, inheritance bool
	}{
		{"broad", "D:P(A;;FR;;;WD)(A;;FR;;;AU)", true, false, false},
		{"trusted", "D:P(A;;FA;;;SY)(A;;FA;;;BA)(A;;FR;;;"+current+")", false, false, false},
		{"other", "D:P(A;;FR;;;LS)", false, true, false},
		{"deny", "D:P(D;;FR;;;WD)", false, false, false},
		{"inherit only", "D:P(A;OIIO;FR;;;WD)", false, false, false},
		{"generic", "D:P(A;;GR;;;WD)", true, false, false},
		{"unprotected", "D:(A;;FR;;;WD)", true, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sd, err := windows.SecurityDescriptorFromString(tc.sddl)
			if err != nil { t.Fatal(err) }
			got := inspectStorageDescriptor("fixture", sd, current)
			if got.State != "inspected" || got.BroadRead != tc.broad || got.OtherRead != tc.other || got.InheritanceEnabled != tc.inheritance {
				t.Fatalf("unexpected ACL metadata: %+v", got)
			}
		})
	}
}

func TestInspectStoragePermissionFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture.txt")
	const content = "fixture only"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil { t.Fatal(err) }
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil { t.Fatal(err) }
	sd, err := windows.SecurityDescriptorFromString("D:P(A;;FA;;;"+user.User.Sid.String()+")(A;;FR;;;WD)")
	if err != nil { t.Fatal(err) }
	dacl, _, err := sd.DACL()
	if err != nil { t.Fatal(err) }
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, dacl, nil); err != nil { t.Fatal(err) }
	before, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil { t.Fatal(err) }
	got := inspectStoragePermission("fixture", path)
	if got.State != "inspected" || !got.BroadRead || got.InheritanceEnabled {
		t.Fatalf("unexpected actual file ACL: %+v", got)
	}
	after, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil { t.Fatal(err) }
	if before.String() != after.String() { t.Fatal("inspection modified ACL") }
	data, err := os.ReadFile(path)
	if err != nil || string(data) != content { t.Fatal("inspection modified fixture") }
	if got := inspectStoragePermission("fixture", filepath.Join(dir, "missing.txt")); got.State != "missing" {
		t.Fatalf("missing path: %+v", got)
	}
	if got := inspectStoragePermission("fixture", "\\\\server\\share\\fixture.txt"); got.State != "skipped_nonlocal" {
		t.Fatalf("UNC path: %+v", got)
	}
}
