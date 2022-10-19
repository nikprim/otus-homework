package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

var ErrEmpty = errors.New("domain is empty")

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if len(domain) == 0 {
		return nil, ErrEmpty
	}

	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	var user User
	domainLower := strings.ToLower(domain)

	for scanner.Scan() {
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("get user error: %w", err)
		}

		emailLower := strings.ToLower(user.Email)
		if strings.HasSuffix(emailLower, domainLower) {
			secondDomain := emailLower[strings.Index(user.Email, "@")+1:]
			result[secondDomain]++
		}
	}

	return result, nil
}
