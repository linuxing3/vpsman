package notion

import (
	"github.com/kjk/u"
)

var srcFiles = u.MakeAllowedFileFilterForExts(".go")
var dirsToSkip = u.MakeExcludeDirsFilter("www")
var allFiles = u.MakeFilterAnd(srcFiles, dirsToSkip)

func doLineCount() int {
	stats := u.NewLineStats()
	err := stats.CalcInDir(".", allFiles, true)
	if err != nil {
		logf("doWordCount: stats.wcInDir() failed with '%s'\n", err)
		return 1
	}
	u.PrintLineStats(stats)
	return 0
}
