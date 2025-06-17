package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	Unknown = iota
	InRule
	InUpdate
)

type Rule struct {
	former int
	latter int
}

func NewRule(former, latter string) (Rule, error) {
	fi, ferr := strconv.Atoi(former)
	if ferr != nil {
		return Rule{}, ferr
	}
	li, lerr := strconv.Atoi(latter)
	if lerr != nil {
		return Rule{}, lerr
	}
	return Rule{fi, li}, nil
}

type Update struct {
	pages []int
}

func NewUpdate() *Update {
	return &Update{make([]int, 0, 1024)}
}

func (u *Update) Add(s string) error {
	if u == nil {
		return fmt.Errorf("err null pointer to update")
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	u.pages = append(u.pages, i)
	return nil
}

func (u *Update) Middle() int {
	if u == nil {
		return 0
	}
	n := len(u.pages)
	return u.pages[n/2]
}

type Parser struct {
	rules   []Rule
	updates []Update
	src     []byte
	state   int
}

func NewParser(src []byte) *Parser {
	rules := make([]Rule, 0, 1024)
	updates := make([]Update, 0, 1024)
	return &Parser{rules, updates, src, Unknown}
}

func (p *Parser) Parse() error {
	b := make([]byte, 0, 1024)
	for len(p.src) > 0 {
		next := p.src[0]
		if next == '\n' {
			if e := p.handleNewline(b); e != nil {
				return e
			}
			b = b[:0]
			continue
		}

		if next == '|' {
			p.state = InRule
		}

		if next == ',' {
			p.state = InUpdate
		}

		b = append(b, next)
		p.src = p.src[1:]
	}
	return nil
}

func (p *Parser) handleNewline(b []byte) error {
	e := error(nil)
	if p.state == InRule {
		e = p.parseRule(b)
	}

	if p.state == InUpdate {
		e = p.parseUpdate(b)
	}

	p.src = p.src[1:]
	p.state = Unknown
	return e
}

func (p *Parser) parseRule(b []byte) error {
	i := bytes.IndexByte(b, '|')
	if i < 0 {
		return fmt.Errorf("parseRule() err rule is invalid")
	}
	former := string(b[0:i])
	latter := string(b[i+1:])
	rule, rerr := NewRule(former, latter)
	if rerr != nil {
		return fmt.Errorf("Err processing rule %e", rerr)
	}
	p.rules = append(p.rules, rule)
	return nil
}

func (p *Parser) parseUpdate(b []byte) error {
	s := make([]byte, 0, 10)
	update := NewUpdate()
	for _, ch := range b {
		if ch == ',' {
			err := update.Add(string(s))
			if err != nil {
				return err
			}
			s = s[:0]
			continue
		}
		s = append(s, ch)
	}
	err := update.Add(string(s))
	if err != nil {
		return err
	}
	p.updates = append(p.updates, *update)
	return nil
}

func main() {
	bytes, err := os.ReadFile("./cmd/day5/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	parser := NewParser(bytes)
	e := parser.Parse()
	if e != nil {
		fmt.Println(e)
		return
	}
	rules := make(map[int]map[int]bool) // map from a page# to the page#'s that MUST precede it.
	for _, rule := range parser.rules {
		r, exists := rules[rule.former]
		if !exists {
			rules[rule.former] = make(map[int]bool)
			r = rules[rule.former]
		}
		r[rule.latter] = true
	}

	midCount := 0
	for _, update := range parser.updates {
		updateOk := true
		for idx := range update.pages {
			page := update.pages[idx]
			rule := rules[page] // all of these numbers must occur after the current number
			if MapHas(rule, update.pages[:idx]) >= 0 {
				updateOk = false
				fixUpdate(update.pages, rules)
				break
			}
		}
		if !updateOk {
			midCount += update.Middle()
		}
	}

	fmt.Println(midCount)
}

func MapHas[T comparable](a map[T]bool, b []T) int {
	for i, n := range b {
		if _, hasKey := a[n]; hasKey {
			return i
		}
	}
	return -1
}

func fixUpdate(update []int, rules map[int]map[int]bool) {
	head := len(update) - 1 // start at the end and work toward the beginning
	for head > 0 {
		rule := rules[update[head]]
		if wrongIdx := MapHas(rule, update[:head]); wrongIdx >= 0 {
			tmp := update[head]
			update[head] = update[wrongIdx]
			update[wrongIdx] = tmp
		} else {
			head--
		}
	}
}
