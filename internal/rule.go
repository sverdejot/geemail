package internal

type Rule struct {
	Domain string `yaml:"domain"`
	Delete bool   `yaml:"delete"`
}

type RuleFile struct {
	Rules []Rule `yaml:"rules"`
}

func NewRuleFile() RuleFile {
	return RuleFile{}
}

func (f RuleFile) AddRule(domain string, del bool) RuleFile {
	return RuleFile{
		Rules: append(f.Rules, Rule{domain, del}),
	}
}
