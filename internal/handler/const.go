package handler

var statusCodes = map[int]StatusCode{
	401: {"401",
		"Authorization required!",
		"You are not authorized, please go back to the home page and sign in"},
	403: {"404",
		"Forbidden!",
		"Request failed due to insufficient permissions"},
	404: {"404",
		"Page not found!",
		"The page that you've requested has not been found on this website, please try something else"},
}

var localization = []string{"homeButton", "addButton", "user", "register", "login", "logout", "titleMain", "titleEdit", "titleUploadArt", "titleSignUp", "titleLogin", "titleUserSettings", "titleAdmUsersList", "titleAdmManageUser", "listOfArticles", "lastModification", "uploadedBy", "backToTop", "preview", "save", "articleTitle", "upload", "by", "loginString", "passwordString", "fullNameString", "onlyLatin", "mustContain", "mustBeEmail", "submitButton", "loginButton", "settings", "security", "currentPass", "newPass", "reNewPass", "changePassButton", "misc", "language", "colorTheme", "changeSettingsButton", "adminPanel", "managePortalUsers", "listOfUsers", "thName", "thFullName", "thManage", "manageUser", "userStr", "isAdmin", "changeUserButton", "backToUsersList", "download"}
