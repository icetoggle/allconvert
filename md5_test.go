package main

import "testing"

func Test_md5file(t *testing.T) {
	md5num := fileConvertMd5("allconvert.sublime-project")
	if md5num != "a44ea86f68b347f3c99428e4bb9d5194" {
		t.Error("md5 check file")
	}
}
