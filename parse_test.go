package fir

import (
	"bytes"
	"strings"

	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func Test_query(t *testing.T) {
	type args struct {
		fi fileInfo
	}
	tests := []struct {
		name string
		args args
		want fileInfo
	}{

		{
			name: "no filters in event string and key is not present",
			args: args{
				fi: fileInfo{
					name: "test.html",
					content: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>`),
				},
			},
			want: fileInfo{
				name: "test.html",
				content: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>`),
				eventTemplates: eventTemplates{
					"event:ok": eventTemplate{
						"tmpl1": struct{}{},
						"-":     struct{}{},
					},
				}},
		},
		{
			name: "same event multiple elements: no filters in event string and key is not present",
			args: args{
				fi: fileInfo{
					name: "test.html",
					content: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					`),
				},
			},
			want: fileInfo{
				name: "test.html",
				content: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					`),
				eventTemplates: eventTemplates{
					"event:ok": eventTemplate{
						"tmpl1": struct{}{},
						"-":     struct{}{},
					},
				}},
		},
		{
			name: "multiple events, multiple templates, multiple elements: no filters in event string and key is not present",
			args: args{
				fi: fileInfo{
					name: "test.html",
					content: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event1:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					<div
						 @fir:event2:ok::tmpl2=""
						 @fir:event:ok=""> 
					</div>
					`),
				},
			},
			want: fileInfo{
				name: "test.html",
				content: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event1:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					<div
						 @fir:event2:ok::tmpl2=""
						 @fir:event:ok=""> 
					</div>
					`),
				eventTemplates: eventTemplates{
					"event:ok": eventTemplate{
						"-": struct{}{},
					},
					"event1:ok": eventTemplate{
						"tmpl1": struct{}{},
					},
					"event2:ok": eventTemplate{
						"tmpl2": struct{}{},
					},
				}},
		},
		{
			name: "filters in event string and key is present",
			args: args{
				fi: fileInfo{
					name: "test.html",
					content: []byte(`<!DOCTYPE html> 
					<div key="1"
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok::tmpl2=""  
						 @fir:event:ok::tmpl2="" 
						 @fir:event:ok=""
						 @fir:[event1:ok,event2:ok]::tmpl3="console.log('hello')"> 
					</div>`),
				},
			},
			want: fileInfo{
				name: "test.html",
				content: []byte(`<!DOCTYPE html> 
					<div key="1"
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok::tmpl2=""  
						 @fir:event:ok::tmpl2="" 
						 @fir:event:ok=""
						 @fir:[event1:ok,event2:ok]::tmpl3="console.log('hello')"> 
					</div>`),
				eventTemplates: eventTemplates{
					"event:ok": eventTemplate{
						"tmpl1": struct{}{},
						"tmpl2": struct{}{},
						"-":     struct{}{},
					},
					"event1:ok": eventTemplate{
						"tmpl3": struct{}{},
					},
					"event2:ok": eventTemplate{
						"tmpl3": struct{}{},
					},
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := query(tt.args.fi)
			if !reflect.DeepEqual(got.eventTemplates, tt.want.eventTemplates) {
				t.Errorf("eventTemplates = %v, want %v", got.eventTemplates, tt.want.eventTemplates)
			}
			if !areHTMLStringsEqual(t, got.content, tt.want.content) {
				t.Errorf("html \n %v, \n want \n %v", string(got.content), string(tt.want.content))
			}
		})
	}
}

func Test_transform(t *testing.T) {

	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		// TODO: Add test cases.
		{
			name: "add key attribute to children",
			input: []byte(`<!DOCTYPE html> 
					<div key="1"> 
						<div> <button @click="">  </button> </div>
						<div> <div> <div> <input @change="" > </div> </div> </div>
						<form @submit=""> </form>
					</div>
					`),

			want: []byte(`<!DOCTYPE html> 
					<div key="1"> 
						<div @fir:event:ok="" key="1" > </div>
						<div> <button @click="" key="1">  </button> </div>
						<div> <div> <div> <input @change="" key="1"> </div> </div> </div>
						<form @submit="" key="1"> </form>
					</div>
					`),
		},
		{
			name: "no filters in event string and key is not present",
			input: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>`),

			want: []byte(`<!DOCTYPE html> 
					<div class="fir-event-ok--tmpl1 fir-event-ok" 
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok="" 
						 </div>`),
		},
		{
			name: "same event multiple elements: no filters in event string and key is not present",
			input: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					<div
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					`),

			want: []byte(`<!DOCTYPE html> 
						<div class="fir-event-ok--tmpl1 fir-event-ok" 
							@fir:event:ok::tmpl1=""
							@fir:event:ok="" 
						 </div>
						 <div class="fir-event-ok--tmpl1 fir-event-ok" 
							@fir:event:ok::tmpl1=""
							@fir:event:ok="" 
						 </div>
						 `),
		},
		{
			name: "multiple events, multiple templates, multiple elements: no filters in event string and key is not present",
			input: []byte(`<!DOCTYPE html> 
					<div
						 @fir:event1:ok::tmpl1=""
						 @fir:event:ok=""> 
					</div>
					<div
						 @fir:event2:ok::tmpl2=""
						 @fir:event:ok=""> 
					</div>
					`),

			want: []byte(`<!DOCTYPE html> 
						<div class="fir-event1-ok--tmpl1 fir-event-ok" 
							@fir:event1:ok::tmpl1=""
							@fir:event:ok="" 
						 </div>
						 <div class="fir-event2-ok--tmpl2 fir-event-ok" 
							@fir:event2:ok::tmpl2=""
							@fir:event:ok="" 
						 </div>
						 `),
		},
		{
			name: "filters in event string and key is present",
			input: []byte(`<!DOCTYPE html> 
					<div key="1"
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok::tmpl2=""  
						 @fir:event:ok::tmpl2="" 
						 @fir:event:ok=""
						 @fir:[event1:ok,event2:ok]::tmpl3="console.log('hello')"> 
					</div>`),

			want: []byte(`<!DOCTYPE html> 
					<div key="1" 
						 class="fir-event-ok--tmpl1--1 fir-event-ok--tmpl2--1 fir-event-ok--1 fir-event1-ok--tmpl3--1 fir-event2-ok--tmpl3--1" 
						 @fir:event:ok::tmpl1=""
						 @fir:event:ok::tmpl2=""  
						 @fir:event:ok::tmpl2="" 
						 @fir:event:ok=""
						 @fir:event1:ok::tmpl3="console.log('hello')"
						 @fir:event2:ok::tmpl3="console.log('hello')"> 
						 </div>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transform(tt.input)

			if !areHTMLStringsEqual(t, got, tt.want) {
				t.Errorf("html \n %v, \n want \n %v", string(got), string(tt.want))
			}
		})
	}
}

func areHTMLStringsEqual(t *testing.T, html1, html2 []byte) bool {
	// Load the HTML strings into goquery documents
	doc1, err := goquery.NewDocumentFromReader(bytes.NewReader(html1))
	if err != nil {
		return false
	}

	doc2, err := goquery.NewDocumentFromReader(bytes.NewReader(html2))
	if err != nil {
		return false
	}

	doc1Attr := make(map[string]string)
	doc2Attr := make(map[string]string)

	doc1.Find("*").Each(func(_ int, node *goquery.Selection) {
		// fmt.Println("node attributes => ", node.Get(0).Attr)
		for _, attr := range node.Get(0).Attr {
			doc1Attr[attr.Key] = attr.Val
		}
	})

	doc2.Find("*").Each(func(_ int, node *goquery.Selection) {
		// fmt.Println("node attributes => ", node.Get(0).Attr)
		for _, attr := range node.Get(0).Attr {
			doc2Attr[attr.Key] = attr.Val
		}
	})

	for key, val := range doc1Attr {
		//fmt.Println("key => ", key)
		if key == "class" {
			var class1, class2 []string
			// remove whitespace
			for _, class := range strings.Split(val, " ") {
				if class == "" {
					continue
				}
				class1 = append(class1, strings.TrimSpace(class))
			}
			for _, class := range strings.Split(doc2Attr[key], " ") {
				if class == "" {
					continue
				}
				class2 = append(class2, strings.TrimSpace(class))
			}
			assert.ElementsMatch(t, class1, class2)
			continue
		}
		if doc2Attr[key] != val {
			return false
		}
	}

	return true
}

func TestDeepMergeEventTemplates(t *testing.T) {
	testCases := []struct {
		evt1           eventTemplates
		evt2           eventTemplates
		expectedResult eventTemplates
	}{
		// Test case 1: Merging two empty eventTemplates
		{
			evt1:           make(eventTemplates),
			evt2:           make(eventTemplates),
			expectedResult: make(eventTemplates),
		},
		// Test case 2: Merging two eventTemplates with non-empty values
		{
			evt1: eventTemplates{
				"event1": eventTemplate{"template1": struct{}{}, "template2": struct{}{}},
			},
			evt2: eventTemplates{
				"event1": eventTemplate{"template3": struct{}{}},
				"event2": eventTemplate{"template4": struct{}{}},
			},
			expectedResult: eventTemplates{
				"event1": eventTemplate{"template1": struct{}{}, "template2": struct{}{}, "template3": struct{}{}},
				"event2": eventTemplate{"template4": struct{}{}},
			},
		},
		// Test case 3: Merging two eventTemplates with duplicate values
		{
			evt1: eventTemplates{
				"event1": eventTemplate{"template1": struct{}{}, "template2": struct{}{}},
			},
			evt2: eventTemplates{
				"event1": eventTemplate{"template2": struct{}{}, "template3": struct{}{}},
			},
			expectedResult: eventTemplates{
				"event1": eventTemplate{"template1": struct{}{}, "template2": struct{}{}, "template3": struct{}{}},
			},
		},
		// Test case 4: Merging two eventTemplates with no common keys
		{
			evt1: eventTemplates{
				"event1": eventTemplate{"template1": struct{}{}},
			},
			evt2: eventTemplates{
				"event2": eventTemplate{"template2": struct{}{}},
			},
			expectedResult: eventTemplates{
				"event1": eventTemplate{"template1": struct{}{}},
				"event2": eventTemplate{"template2": struct{}{}},
			},
		},
	}

	for _, tc := range testCases {
		result := deepMergeEventTemplates(tc.evt1, tc.evt2)
		if !reflect.DeepEqual(result, tc.expectedResult) {
			t.Errorf("Expected %v, but got %v", tc.expectedResult, result)
		}
	}
}

func TestGetEventFilter(t *testing.T) {
	tests := []struct {
		input          string
		expectedBefore string
		expectedValues []string
		expectedAfter  string
		valid          bool
	}{
		{
			input:          "SomeText[value1:ok,value2:pending,value3:error]moreText",
			expectedBefore: "SomeText",
			expectedValues: []string{"value1:ok", "value2:pending", "value3:error"},
			expectedAfter:  "moreText",
			valid:          true,
		},
		{
			input:          "[value1:ok,value2:pending,value3:error]moreText",
			expectedBefore: "",
			expectedValues: []string{"value1:ok", "value2:pending", "value3:error"},
			expectedAfter:  "moreText",
			valid:          true,
		},
		{
			input:          "[value1:ok,value1-update:ok,value-update:pending,value3:error]",
			expectedBefore: "",
			expectedValues: []string{"value1:ok", "value1-update:ok", "value-update:pending", "value3:error"},
			expectedAfter:  "",
			valid:          true,
		},
		{
			input:          "SomeText:[value1:ok,value2:pending,value3:error]",
			expectedBefore: "SomeText:",
			expectedValues: []string{"value1:ok", "value2:pending", "value3:error"},
			expectedAfter:  "",
			valid:          true,
		},
		{
			input:          "SomeText:[value:ok]::moreText",
			expectedBefore: "SomeText:",
			expectedValues: []string{"value:ok"},
			expectedAfter:  "::moreText",
			valid:          true,
		},
		{
			input:          "SomeText[]moreText",
			expectedBefore: "SomeText",
			expectedValues: nil,
			expectedAfter:  "moreText",
			valid:          false,
		},
		{
			input:          "SomeText[invalidFormat]moreText",
			expectedBefore: "",
			expectedValues: nil,
			expectedAfter:  "",
			valid:          false,
		},
		{
			input:          "fir:event:ok::tmpl",
			expectedBefore: "",
			expectedValues: []string{"fir:event:ok::tmpl"},
			expectedAfter:  "",
			valid:          true,
		},
	}

	for _, test := range tests {
		ef, err := getEventFilter(test.input)
		if err != nil && test.valid {
			t.Fatalf("Failed to parse event filter for input: %s, error: = %v", test.input, err)
		}

		if err == nil && !test.valid {
			t.Fatalf("Expected error for input: %s, but got none", test.input)
		}

		if ef == nil && err != nil {
			continue
		}

		if ef.BeforeBracket != test.expectedBefore {
			t.Errorf("BeforeBracket mismatch for input: %s, expected: %s, got: %s", test.input, test.expectedBefore, ef.BeforeBracket)
		}

		if !reflect.DeepEqual(ef.Values, test.expectedValues) {
			t.Errorf("Values mismatch for input: %s, expected: %v, got: %v", test.input, test.expectedValues, ef.Values)
		}

		if ef.AfterBracket != test.expectedAfter {
			t.Errorf("AfterBracket mismatch for input: %s, expected: %s, got: %s", test.input, test.expectedAfter, ef.AfterBracket)
		}

	}
}
