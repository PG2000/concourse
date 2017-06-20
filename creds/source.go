package creds

import (
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
)

type Source struct {
	variablesResolver template.Variables
	rawSource         atc.Source
}

func NewSource(variables template.Variables, source atc.Source) *Source {
	return &Source{
		variablesResolver: variables,
		rawSource:         source,
	}
}

func (s *Source) Evaluate() (atc.Source, error) {
	var source atc.Source
	err := evaluate(s.variablesResolver, s.rawSource, &source)
	if err != nil {
		return nil, err
	}

	return source, nil
}
