package service

import (
	"testing"
)

func TestParseAccessToken(t *testing.T) {

	t.Run("normal en-us", func(tt *testing.T) {
		var cer = "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxMzAyNSwiTG9naW5OYW1lIjoiZGVtby0xNTk4NTg0NzczMTU1IiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNTk5NDQ4NzgwLCJSbmRJRCI6ImIyNzY0OTlmLTRhY2UtNGRhZC05MzZlLTA0ZWU0NjE4NTI3ZCIsIlNpZ25LZXkiOiIwMTJkOGJjNGYwYTZlNmNjYmYwMTdiOTVkODM5ZDBhNyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj5HW7ojskC2l3lkITLvePKl9ucZngaF/TqxtzI2t7Xq1y8TSQfeTJizscw8ExxUSX8lHo1JdZZM7Yi2KHRBHp044ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+wEJ9ZadpOsxjnq9uqGu4YdF7RbkBm++qsVfm36qJMXi0Ra9hTXDZLHYvw86hE1BKTYzCN22IldNnGFXM6MVHSiWIK/Brvij2VVYWwVJJUdATvwILBrjtNo9//x4BMat90Xbep64k5J0mL6DxVv3J8H2HgbfQhxre1NwDd387Oer9nFU3ciGsLvwkm6OOs/46vWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="
		accessToken := ParseAccessToken(cer)
		t.Logf("%+v", accessToken)
	})

	t.Run("normal zh-cn", func(tt *testing.T) {
		var cer = "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxMzAyNSwiTG9naW5OYW1lIjoiZGVtby0xNTk4NTg0NzczMTU1IiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiWkgtQ04iLCJUb2tlbkV4cGlyZWRUaW1lIjoxNTk5NDUxOTQ3LCJSbmRJRCI6IjQyMDM0YWRkLTg0ZmQtNGZkNS1hMWY4LTEyNTQ1MGNkNDk3YiIsIlNpZ25LZXkiOiI4NDA2NjhmNTFhNTBlM2E5OWE1NzY1MzVlYzQyMTk5YyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj5HW7ojskC2l3lkITLvePKl9ucZngaF/TqxtzI2t7Xq1y8TSQfeTJizscw8ExxUSX8lHo1JdZZM7Yi2KHRBHp044ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFbq99ORLYRW41n1rAu8Erk2jCFajdZ3/pI0VFhYvnD+xjb5tGhzAvpdg0EWeRIn3PM6KUHZVAWxkm44bbsHzTPFzfjp2kKPT77nqNzrrzvxr/AOwESiIf5tlcfsvtXmnDUfM/8B8muL+VjI+nFXB08FdfdM9rWJht2uNCUz22R/fAm51WtTnwkQrSIgaqM4acCXC+kg/wK9u0dSdxZ5jsPupbiuFiCqLSB3CseIUd+HjWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="
		accessToken := ParseAccessToken(cer)
		t.Logf("%+v", accessToken)
	})

	t.Run("empty", func(tt *testing.T) {
		var cer = ""
		accessToken := ParseAccessToken(cer)
		t.Logf("%+v", accessToken)
	})
}

func TestGetAccessToken(t *testing.T) {
	pid, cer := GetTestAccessToken()
	t.Log(pid, cer)
}
