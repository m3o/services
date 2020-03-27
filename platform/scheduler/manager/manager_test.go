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
					status:   githubFileChangeStatusCreated,
				},
				{
					fileName: "asim/handler/something.go",
					status:   githubFileChangeStatusChanged,
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
					status:   githubFileChangeStatusRemoved,
				},
				{
					fileName: "asim/service/handler/somehandler.go",
					status:   githubFileChangeStatusChanged,
				},
			},
			expect: map[string]serviceStatus{
				"asim/scheduler": serviceStatusDeleted,
				"asim/service":   serviceStatusUpdated,
				"asim":           serviceStatusUpdated,
			},
		},
		{
			files: []fileToStatus{
				{
					fileName: "build.sh",
					status:   githubFileChangeStatusChanged,
				},
				{
					fileName: "asim/scheduler/something.go",
					status:   githubFileChangeStatusChanged,
				},
				{
					fileName: "asim/scheduler/hander/something.go",
					status:   githubFileChangeStatusRemoved,
				},
			},
			expect: map[string]serviceStatus{
				"asim/scheduler": serviceStatusUpdated,
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
