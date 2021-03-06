// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import "github.com/spf13/hugo/source"

type Handler interface {
	// Read the Files in and register
	Read(*source.File, *Site, HandleResults)

	// Convert Pages to prepare for templatizing
	// Convert Files to their final destination
	Convert(interface{}, *Site, HandleResults)

	// Extensions to register the handle for
	Extensions() []string
}

type HandledResult struct {
	page *Page
	file *source.File
	err  error
}

type HandleResults chan<- HandledResult

type ReadFunc func(*source.File, *Site, HandleResults)
type PageConvertFunc func(*Page, *Site, HandleResults)
type FileConvertFunc ReadFunc

type Handle struct {
	extensions  []string
	read        ReadFunc
	pageConvert PageConvertFunc
	fileConvert FileConvertFunc
}

var handlers []Handler

func (h Handle) Extensions() []string {
	return h.extensions
}

func (h Handle) Read(f *source.File, s *Site, results HandleResults) {
	h.read(f, s, results)
}

func (h Handle) Convert(i interface{}, s *Site, results HandleResults) {
	if h.pageConvert != nil {
		h.pageConvert(i.(*Page), s, results)
	} else {
		h.fileConvert(i.(*source.File), s, results)
	}
}

func RegisterHandler(h Handler) {
	handlers = append(handlers, h)
}

func Handlers() []Handler {
	return handlers
}

func FindHandler(ext string) Handler {
	for _, h := range Handlers() {
		if HandlerMatch(h, ext) {
			return h
		}
	}
	return nil
}

func HandlerMatch(h Handler, ext string) bool {
	for _, x := range h.Extensions() {
		if ext == x {
			return true
		}
	}
	return false
}
