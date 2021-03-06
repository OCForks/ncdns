package util_test

import "testing"
import "github.com/namecoin/ncdns/util"
import "gopkg.in/hlandau/madns.v1/merr"

type item struct {
	input            string
	expectedHead     string
	expectedRest     string
	expectedTail     string
	expectedTailRest string
}

var items = []item{
	item{"", "", "", "", ""},
	item{"a", "a", "", "a", ""},
	item{"alpha", "alpha", "", "alpha", ""},
	item{"alpha.beta", "beta", "alpha", "alpha", "beta"},
	item{"alpha.beta.gamma", "gamma", "alpha.beta", "alpha", "beta.gamma"},
	item{"alpha.beta.gamma.delta", "delta", "alpha.beta.gamma", "alpha", "beta.gamma.delta"},
	item{"alpha.beta.gamma.delta.", "delta", "alpha.beta.gamma", "alpha", "beta.gamma.delta."},
}

func TestSplitDomainHead(t *testing.T) {
	for i := range items {
		head, rest := util.SplitDomainHead(items[i].input)
		tail, trest := util.SplitDomainTail(items[i].input)
		if head != items[i].expectedHead {
			t.Errorf("Input \"%s\": head \"%s\" does not equal expected value \"%s\"", items[i].input, head, items[i].expectedHead)
		}
		if rest != items[i].expectedRest {
			t.Errorf("Input \"%s\": rest \"%s\" does not equal expected value \"%s\"", items[i].input, rest, items[i].expectedRest)
		}
		if tail != items[i].expectedTail {
			t.Errorf("Input \"%s\": tail \"%s\" does not equal expected value \"%s\"", items[i].input, tail, items[i].expectedTail)
		}
		if trest != items[i].expectedTailRest {
			t.Errorf("Input \"%s\": tail rest \"%s\" does not equal expected value \"%s\"", items[i].input, trest, items[i].expectedTailRest)
		}
	}
}

type aitem struct {
	input            string
	anchor           string
	expectedSubname  string
	expectedBasename string
	expectedRootname string
	expectedError    error
}

var aitems = []aitem{
	aitem{"", "bit", "", "", "", merr.ErrNotInZone},
	aitem{".", "bit", "", "", "", merr.ErrNotInZone},
	aitem{"d.", "bit", "", "", "", merr.ErrNotInZone},
	aitem{"a.b.c.d.", "bit", "", "", "", merr.ErrNotInZone},
	aitem{"a.b.c.d.bit.", "bit", "a.b.c", "d", "bit", nil},
	aitem{"d.bit.", "bit", "", "d", "bit", nil},
	aitem{"bit.", "bit", "", "", "bit", nil},
	aitem{"bit.x.y.z.", "bit", "", "", "bit.x.y.z", nil},
	aitem{"d.bit.x.y.z.", "bit", "", "d", "bit.x.y.z", nil},
	aitem{"c.d.bit.x.y.z.", "bit", "c", "d", "bit.x.y.z", nil},
	aitem{"a.b.c.d.bit.x.y.z.", "bit", "a.b.c", "d", "bit.x.y.z", nil},
}

func TestSplitDomainByFloatingAnchor(t *testing.T) {
	for i, it := range aitems {
		subname, basename, rootname, err := util.SplitDomainByFloatingAnchor(it.input, it.anchor)
		if subname != it.expectedSubname {
			t.Errorf("Item %d: subname \"%s\" does not equal expected value \"%s\"", i, subname, it.expectedSubname)
		}
		if basename != it.expectedBasename {
			t.Errorf("Item %d: basename \"%s\" does not equal expected value \"%s\"", i, basename, it.expectedBasename)
		}
		if rootname != it.expectedRootname {
			t.Errorf("Item %d: rootname \"%s\" does not equal expected value \"%s\"", i, basename, it.expectedRootname)
		}
		if err != it.expectedError {
			t.Errorf("Item %d: error \"%s\" does not equal expected error \"%s\"", i, err, it.expectedError)
		}
	}
}
