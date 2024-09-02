//go:build ruleguard
// +build ruleguard

package ruleguard

import "github.com/quasilyte/go-ruleguard/dsl"

func junosDecode(m dsl.Matcher) { //nolint
	m.Match(`tfdata.JunosDecode($*_)`).
		Where(!m.File().PkgPath.Matches("internal/junos") && !m.File().PkgPath.Matches("internal/tfdata")).
		Suggest("use JunosDecode function of *junos.Session instead of function from tfdata package")
}

func junosdecodeDecode(m dsl.Matcher) { //nolint
	m.Match(`junosdecode.Decode($*_)`, `jdecode.Decode($*_)`).
		Where(!m.File().PkgPath.Matches("internal/tfdata")).
		Suggest("use JunosDecode function of *junos.Session instead of function from junosdecode package")
}
