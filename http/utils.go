package http

import (
	"slices"
	"strings"

	u "github.com/adagit94/golib/utils"
)

func collectDynSegsIndices(segs []string) []int {
	dynSegsIndices := make([]int, 0)
	
	for i, seg := range segs {
		if strings.HasPrefix(seg, ":") {
			dynSegsIndices = append(dynSegsIndices, i)
		}
	}

	return dynSegsIndices
}

func genRouteSegs(segs []string) []*routeSeg {
	segsConfs := u.MapSlice(segs, func(seg string, i int) *routeSeg {
		isDyn := strings.HasPrefix(seg, ":")

		if isDyn {
			extractedSeg, _ := strings.CutPrefix(seg, ":")
			seg = extractedSeg
		}

		return &routeSeg{seg: seg, dyn: isDyn}
	})

	return segsConfs
}

func sortRoutes[H any](routes []*dynRoute[H]) {
	slices.SortFunc(routes, func(a, b *dynRoute[H]) int {
		for i := range a.segs {
			aSeg, bSeg := a.segs[i], b.segs[i]

			if !aSeg.dyn && bSeg.dyn {
				return -1
			}

			if aSeg.dyn && !bSeg.dyn {
				return 1
			}
		}

		return 0
	})
}
