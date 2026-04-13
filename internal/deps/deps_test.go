package deps

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestDeps(t *testing.T) {
	fmt.Println(os.UserCacheDir())
	fmt.Println(os.UserConfigDir())
	// fmt.Println(os.User())
	fmt.Println(url.QueryEscape("git@github.com:pbm-org/pbm.git"))
}
