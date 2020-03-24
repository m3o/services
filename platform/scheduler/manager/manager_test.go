package manager

import (
	"testing"
)

type tcase struct {
	files  []fileToStatus
	expect map[string]serviceStatus
}

func TestFilesToServiceStatus(t *testing.T) {
	cases := []tcase{
		{
			files: []fileToStatus{
				{
					fileName: "asim/main.go",
					status:   "created",
				},
				{
					fileName: "asim/handler/something.go",
					status:   "changed",
				},
			},
			expect: map[string]serviceStatus{
				"asim":         serviceStatusCreated,
				"asim/handler": serviceStatusUpdated,
			},
		},
		{
			files: []fileToStatus{
				{
					fileName: "asim/scheduler/main.go",
					status:   "removed",
				},
				{
					fileName: "asim/service/handler/somehandler.go",
					status:   "changed",
				},
			},
			expect: map[string]serviceStatus{
				"asim/scheduler": serviceStatusDeleted,
				"asim/service":   serviceStatusUpdated,
				"asim":           serviceStatusUpdated,
			},
		},
	}
	for i, c := range cases {
		ss := folderStatuses(c.files)
		for folder, status := range ss {
			if c.expect[folder] != status {
				t.Errorf("case %v: Expected %v for folder %v, got: %v", i, c.expect[folder], folder, status)
			}
		}
	}
}
