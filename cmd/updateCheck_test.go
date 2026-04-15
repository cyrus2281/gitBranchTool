package cmd

import (
	"runtime"
	"testing"
)

// --- compareVersions tests ---

func TestCompareVersions_UpToDate(t *testing.T) {
	status, err := compareVersions("1.0.0", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionUpToDate {
		t.Errorf("expected VersionUpToDate, got %d", status)
	}
}

func TestCompareVersions_PatchUpdate(t *testing.T) {
	status, err := compareVersions("1.0.0", "1.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionMinorUpdate {
		t.Errorf("expected VersionMinorUpdate for patch, got %d", status)
	}
}

func TestCompareVersions_MinorUpdate(t *testing.T) {
	status, err := compareVersions("1.0.0", "1.1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionMinorUpdate {
		t.Errorf("expected VersionMinorUpdate, got %d", status)
	}
}

func TestCompareVersions_MajorUpdate(t *testing.T) {
	status, err := compareVersions("1.0.0", "2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionMajorUpdate {
		t.Errorf("expected VersionMajorUpdate, got %d", status)
	}
}

func TestCompareVersions_CurrentNewerPatch(t *testing.T) {
	status, err := compareVersions("1.0.1", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionUnofficial {
		t.Errorf("expected VersionUnofficial, got %d", status)
	}
}

func TestCompareVersions_CurrentNewerMinor(t *testing.T) {
	status, err := compareVersions("1.1.0", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionUnofficial {
		t.Errorf("expected VersionUnofficial, got %d", status)
	}
}

func TestCompareVersions_CurrentNewerMajor(t *testing.T) {
	status, err := compareVersions("2.0.0", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionUnofficial {
		t.Errorf("expected VersionUnofficial, got %d", status)
	}
}

func TestCompareVersions_ZeroVersions(t *testing.T) {
	status, err := compareVersions("0.0.0", "0.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionUpToDate {
		t.Errorf("expected VersionUpToDate, got %d", status)
	}
}

func TestCompareVersions_RealVersionPatch(t *testing.T) {
	status, err := compareVersions("3.2.5", "3.2.6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionMinorUpdate {
		t.Errorf("expected VersionMinorUpdate for patch bump, got %d", status)
	}
}

func TestCompareVersions_RealVersionMajor(t *testing.T) {
	status, err := compareVersions("3.2.5", "4.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionMajorUpdate {
		t.Errorf("expected VersionMajorUpdate, got %d", status)
	}
}

func TestCompareVersions_LargeNumbers(t *testing.T) {
	status, err := compareVersions("100.200.300", "100.200.301")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != VersionMinorUpdate {
		t.Errorf("expected VersionMinorUpdate, got %d", status)
	}
}

func TestCompareVersions_ErrorTooFewParts(t *testing.T) {
	_, err := compareVersions("1.0", "1.0.0")
	if err == nil {
		t.Error("expected error for version with 2 parts")
	}
}

func TestCompareVersions_ErrorTooManyParts(t *testing.T) {
	_, err := compareVersions("1.0.0.1", "1.0.0")
	if err == nil {
		t.Error("expected error for version with 4 parts")
	}
}

func TestCompareVersions_ErrorNonNumericCurrent(t *testing.T) {
	_, err := compareVersions("abc.0.0", "1.0.0")
	if err == nil {
		t.Error("expected error for non-numeric current version")
	}
}

func TestCompareVersions_ErrorNonNumericLatest(t *testing.T) {
	_, err := compareVersions("1.0.0", "1.abc.0")
	if err == nil {
		t.Error("expected error for non-numeric latest version")
	}
}

func TestCompareVersions_ErrorEmptyString(t *testing.T) {
	_, err := compareVersions("", "1.0.0")
	if err == nil {
		t.Error("expected error for empty version string")
	}
}

// --- getAssetPrefix tests ---

func TestGetAssetPrefix_CurrentOS(t *testing.T) {
	prefix := getAssetPrefix("3.2.5")

	switch runtime.GOOS {
	case "linux":
		expected := "g-linux-v3.2.5"
		if prefix != expected {
			t.Errorf("expected %q, got %q", expected, prefix)
		}
	case "darwin":
		expected := "g-macos-v3.2.5"
		if prefix != expected {
			t.Errorf("expected %q, got %q", expected, prefix)
		}
	case "windows":
		expected := "g-win-v3.2.5.exe"
		if prefix != expected {
			t.Errorf("expected %q, got %q", expected, prefix)
		}
	default:
		if prefix != "" {
			t.Errorf("expected empty prefix for unknown OS, got %q", prefix)
		}
	}
}

func TestGetAssetPrefix_StripsVPrefix(t *testing.T) {
	prefixWithV := getAssetPrefix("v3.2.5")
	prefixWithoutV := getAssetPrefix("3.2.5")

	// Both should produce the same result since getAssetPrefix strips "v"
	if prefixWithV != prefixWithoutV {
		t.Errorf("v prefix should be stripped: with=%q, without=%q", prefixWithV, prefixWithoutV)
	}
}
