package main

import (
	"flag"
	"log"

	"example.invalid/ServerTeam/Servers/internal/game/inrm/robot"
)

type user struct {
	uid int32
	cer string
}

var u1 = &user{uid: 1015452, cer: "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTQ1MiwiTG9naW5OYW1lIjoiZGVtby0xNjA2ODk5NTg4Mjk5IiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA3NzYzNjI2LCJSbmRJRCI6IjEzOGRlN2RhLTYxMDAtNDBlYy1hYzcwLWY5ZmIxZWUwMGZmNyIsIlNpZ25LZXkiOiI2OWY0MjJlNDU0YzQyNjBkNjJiNDk5NWY0MjNhNTJmMyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZjGLI8AjqnxOoHZgMlySrb1e7rMXexwUPG2xS4D75vNiaYbME9qhEDBxc5H5Q+blgh88napM+L/SQ14+jszIUed44ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+wKMskYDJ469zZUGfu/y7pE8sW6n4gp5mnmLkUrNi495xt5QNCnD6gYmpAQnFgWRVA0UY0rMveQzjHaQzK+Vk1W1MfEaQ/SYrAwH/D+4+lvpZIYj2Cu8NYVmzwqPVVgCauHcQsmSG5RIi+eW7IT7vTWS4T+WwJ/fEbdcJbTKwIRiUmXBcDv0cwPJ2qF/Fsap0DWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="}
var u2 = &user{uid: 1015453, cer: "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTQ1MywiTG9naW5OYW1lIjoiZGVtby0xNjA2ODk5NzE1MjcwIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA3NzYzNzE3LCJSbmRJRCI6IjZkMDFmNTRlLTBhN2YtNDAzYi05YTcwLWE2Njk0OWYwOGM1NCIsIlNpZ25LZXkiOiJlYTQzOTNiODBhNDBjODIwOGE1MWNhOTFmOGUxNjY2NyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZjGLI8AjqnxOoHZgMlySrb1edGipsSyFgpiUk8zrTdkzyYbME9qhEDBxc5H5Q+blghuY6QHUrwJqTxyvYbyOs6LI4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+wKMskYDJ469zZUGfu/y7pE3UeKHz1tIwwfqZM8bOAs2+1cMu4GphiX2hSJxq5+8kPzisPn7Fo+2L7aUmKL3pNoZ24klK7/zCXWVrk+ozAhyRpPfECjBXtQ+hGFYRN9pw/fkM8lI8mf8ANTX+nr3miuN3x2T52nJTa2fFP0rKWUs57jl7ueRtIPjqvCRQRsZLTWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="}
var u3 = &user{uid: 1015454, cer: "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTQ1NCwiTG9naW5OYW1lIjoiZGVtby0xNjA2ODk5Nzg0Nzc4IiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjA3NzYzNzg2LCJSbmRJRCI6ImQ0ODQyMTk0LWM5NTItNGUwNy04YmFiLWYwMTZiMTg0NDA3YyIsIlNpZ25LZXkiOiI1MzBmYTE1OWRjNzhiNDhmZDIzMGE4ODhlOGEyNjA4YSIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZjGLI8AjqnxOoHZgMlySrb1fLRkMfVhw6CvO0N/sJ0ctSYbME9qhEDBxc5H5Q+blghG+9T4JiwoscNL71sYtCYlY4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+wKMskYDJ469zZUGfu/y7pEb7B1r7CraTvtoivmCgqD2sbwR47sJprO5YsVL7TZKI82XS8/fOdIXKQlOqA0R3PzgF/xRuAlkRk7tVkSkSKEw/uo/6Nr2YcTIF9i9bRbEFwlY1krZjdBoiJrOCoAXiTbJ0prwENHjPs5IoZph6LA9Kvi6HucBUF2dNYQLpWv6+nWorM6mN5UFXk3SjvL4JlXfMij76EeCLCYYiW13ScpSg=="}
var pid int

func init() {
	flag.IntVar(&pid, "p", 0, "user1 or user2 or user3")
}

func main() {

	var u *user

	flag.Parse()

	switch pid {
	case 1:
		u = u1
	case 2:
		u = u2
	case 3:
		u = u3
	}

	if u == nil {
		log.Fatal("Not right user")
	}

	rbt := robot.NewMessagePureRobot(":3000", u.uid, u.cer)

	if err := rbt.Run(); err != nil {
		log.Fatal(err)
	}
}
