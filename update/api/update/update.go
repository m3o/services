package update

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/micro/go-micro/v2/metadata"
)

var (
	lock sync.RWMutex
	// the latest update
	update = new(Update)
)

type Update struct {
	Commit  string `json:"commit"`
	Image   string `json:"image"`
	Release string `json:"release"`
}

// get the latest commit
func getLatestCommit() (string, error) {
	rsp, err := http.Get("https://api.github.com/repos/micro/micro/commits")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	// unmarshal commits
	var commits []map[string]interface{}
	err = json.Unmarshal(b, &commits)
	if err != nil {
		return "", err
	}
	// get the commits
	if len(commits) == 0 {
		return "", err
	}
	// the latest commit
	commit := commits[0]["sha"].(string)
	return commit, nil
}

func getLatestRelease() (string, error) {
	rsp, err := http.Get("https://api.github.com/repos/micro/micro/releases")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	// unmarshal commits
	var releases []map[string]interface{}
	err = json.Unmarshal(b, &releases)
	if err != nil {
		return "", err
	}
	// get the commits
	if len(releases) == 0 {
		return "", err
	}
	// the latest commit
	release := releases[0]["tag_name"].(string)
	return release, nil
}

func getLatestImage() (string, error) {
	rsp, err := http.Get("https://hub.docker.com/v2/repositories/micro/micro/tags/latest")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	// unmarshal commits
	var images map[string]interface{}
	err = json.Unmarshal(b, &images)
	if err != nil {
		return "", err
	}
	// get the commits
	updated := images["last_updated"].(string)
	return updated, nil
}

// set new update
func Init() {
	commit, _ := getLatestCommit()
	release, _ := getLatestRelease()
	image, _ := getLatestImage()

	// update commit and release
	lock.Lock()
	defer lock.Unlock()

	update.Commit = commit
	update.Release = release
	update.Image = image
}

func Get() *Update {
	lock.RLock()
	defer lock.RUnlock()
	return update
}

func Event(ctx context.Context, data []byte) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil
	}

	sig, _ := md.Get("X-Hub-Signature")

	// if signature is blank assume its the docker webhook
	if len(sig) == 0 {
		// check if its the docker webhook
		var request map[string]interface{}
		if err := json.Unmarshal(data, &request); err != nil {
			return nil
		}
		// check a callback url exists
		v := request["callback_url"]
		if v == nil {
			return nil
		}
		callback, ok := v.(string)
		if ok && strings.HasPrefix(callback, "https://registry.hub.docker.com") {
			image, err := getLatestImage()
			if err != nil {
				return nil
			}
			lock.Lock()
			update.Image = image
			lock.Unlock()
		}
	}

	// assume github push
	parts := strings.Split(sig, "=")
	if len(parts) < 2 {
		log.Print("not enough parts in X-Hub-Signature")
		return nil
	}

	sha, _ := hex.DecodeString(parts[1])
	b, _ := ioutil.ReadAll(bytes.NewReader(data))
	mac := hmac.New(sha1.New, []byte("5634780310"))
	mac.Write(b)
	expect := mac.Sum(nil)
	equals := hmac.Equal(sha, expect)

	if !equals {
		log.Print("hmac not equal expected")
		return nil
	}

	ev, _ := md.Get("X-Github-Event")

	// update the latest values based on what type of event was received
	switch ev {
	case "push":
		commit, err := getLatestCommit()
		if err != nil {
			return err
		}
		lock.Lock()
		update.Commit = commit
		lock.Unlock()
	case "release":
		release, err := getLatestRelease()
		if err != nil {
			return err
		}
		lock.Lock()
		update.Release = release
		lock.Unlock()
	default:
		log.Print("received unknown git event", ev)
	}

	return nil
}
