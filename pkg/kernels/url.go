package kernels

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
)

type KernelURL interface {
	// returns the directory of the source code, or an error
	fetch(ctx context.Context, log *logrus.Logger, dir string, name string) error
}

func ParseURL(s string) (KernelURL, error) {

	kurl, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	switch kurl.Scheme {
	case "git":
		return NewGitURL(kurl)

	// NB: there are also git repos using http so we would need
	// some detection based on the suffix, e.g., .git vs .tgz
	case "http", "https":
		return nil, fmt.Errorf("%s support comming soon!", kurl.Scheme)

	default:
		return nil, fmt.Errorf("Unsupported URL: '%s'", kurl)
	}

}
