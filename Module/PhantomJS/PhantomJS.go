package PhantomJS

import (
	"context"
	"os/exec"
	"time"
)

func GetPageHtml(url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx,"./rsrc/phantomjs.exe",
		"./rsrc/page.js", url)
	outPut, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(outPut)
}
