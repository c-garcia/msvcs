package smtp

type EMailChecker interface {
	CheckEmail(dest, from, subject, body string) (bool, error)
}
