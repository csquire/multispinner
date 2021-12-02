// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multispinner

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

type Multispinner struct {
	frames      []string
	interval    time.Duration
	newSpinners chan Spinner
	running     bool
	stop        chan struct{}
	writer      io.Writer
}

type Spinner struct {
	Name    string
	Message string
	Status  status
}

type status int

const (
	RUNNING status = iota
	SUCCESS
	FAILURE
)

func NewMultispinner(frames []string, interval time.Duration) *Multispinner {
	return &Multispinner{
		frames:      frames,
		interval:    interval,
		newSpinners: make(chan Spinner, 100),
		running:     false,
		stop:        make(chan struct{}, 1),
		writer:      color.Output,
	}
}

func (m *Multispinner) AddOrUpdate(spinner Spinner) {
	m.newSpinners <- spinner
}

func (m *Multispinner) Start() {
	if m.running {
		return
	}
	m.running = true

	go func() {
		var spinnerList []string
		spinnerMap := make(map[string]Spinner)

		displayedLines := len(spinnerList)
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		fmt.Printf("\033[%dE", displayedLines)
		for {
			for _, frame := range m.frames {
				select {
				case <-m.stop:
					_ = m.draw(frame, displayedLines, spinnerList, spinnerMap)
					return
				case newSpinner := <-m.newSpinners:
					if !contains(spinnerList, newSpinner.Name) {
						spinnerList = append(spinnerList, newSpinner.Name)
					}
					spinnerMap[newSpinner.Name] = newSpinner
				case <-ticker.C:
					displayedLines = m.draw(frame, displayedLines, spinnerList, spinnerMap)
				}
			}
		}
	}()
}

func (m *Multispinner) draw(frame string, displayedLines int, spinnerList []string, spinnerMap map[string]Spinner) int {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\033[%dF", displayedLines)
	for _, name := range spinnerList {
		s := spinnerMap[name]
		currentFrame := frame
		lineColor := yellow
		switch s.Status {
		case SUCCESS:
			lineColor = green
			currentFrame = green("✔")
		case FAILURE:
			lineColor = red
			currentFrame = red("✘")
		}
		framedName := fmt.Sprintf("[%s]", name)
		line := fmt.Sprintf("\033[2K%s %s %s\n", currentFrame, cyan(framedName), lineColor(s.Message))
		fmt.Fprint(m.writer, line)
	}
	return len(spinnerList)
}

func (m *Multispinner) Stop() {
	if m.running {
		m.running = false
		m.stop <- struct{}{}
	}
}

func NewSpinner(name string, message string, status status) Spinner {
	return Spinner{
		Name:    name,
		Message: message,
		Status:  status,
	}
}

func contains(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
