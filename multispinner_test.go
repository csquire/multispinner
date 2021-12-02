package multispinner

import (
	"testing"
	"time"
)

func TestMultispinner(t *testing.T) {
	frames := []string{"🙈", "🙉", "🙊"}
	spinner := NewMultispinner(frames, time.Millisecond*500)
	spinner.AddOrUpdate(NewSpinner("thing1", "working", RUNNING))
	spinner.AddOrUpdate(NewSpinner("thing2", "working", RUNNING))
	spinner.Start()
	time.Sleep(10 * time.Second)
	spinner.AddOrUpdate(NewSpinner("thing2", "failed", FAILURE))
	spinner.AddOrUpdate(NewSpinner("thing1", "success", SUCCESS))
	spinner.Stop()
}
