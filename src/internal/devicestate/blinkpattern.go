/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package devicestate

type blinkPatterns int

const (
	off blinkPatterns = iota
	blink
	on
	exit
)

type blinkPattern struct {
	pattern blinkPatterns
}

func (b *blinkPattern) changePattern(newPattern blinkPatterns) {
	b.pattern = newPattern
}

func (b blinkPattern) getPattern() blinkPatterns {
	return b.pattern
}
