package handler

var toTheTop = []byte("\n<a href=\"#top\"><i>back to top</i></a>")

var statusCodes = map[int]StatusCode{
	401: {"401",
		"Authorization required!",
		"You are not authorized, please go back to the home page and sign in"},
	404: {"404",
		"Page not found!",
		"The page that you've requested has not been found on this website, please try something else"},
}
